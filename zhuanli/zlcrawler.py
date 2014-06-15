# coding: utf8

'''
detail page:

detailTifDownload('gk')

jx = {"downloadabs":null,"downloadbean":{"fileURLPath":"http:\/\/search.cnipr.com\/tempfile\/11ae7fbd-6c1f-4210-9848-807a1b88a0f3.zip","message":"","requestCount":1,"responseCount":1,"returncode":1},"downloadenab":null,"filterChannel":null,"patInfos":null,"strSortMethod":null,"strSynonymous":null,"success":false,"tifInfos":"BOOKS\/XX\/2010\/20100825\/200920258513.7,5,CN200920258513.7"}

j = json.loads(jx)
j['downloadbean']['fileURLPath']

driver.findElement(By.cssSelector(".x-window input[name='name']"));

'''

from pyvirtualdisplay import Display
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.keys import Keys
from selenium.common.exceptions import NoSuchElementException
from selenium.common.exceptions import TimeoutException
from selenium.common.exceptions import StaleElementReferenceException
# from scrapemark import scrape
import pymongo

import time
import string
import urllib
import zipfile
import datetime
import os, sys

def init():
    if not os.path.exists('zlwj'):
        os.mkdir('zlwj')


def clickAndWait(elem):
    elem.click()
    time.sleep(1)


def injectJQuery(browser):
    #read the jquery from a file
    with open('jquery-1.7.1.min.js', 'r') as jquery_js: 
        jquery = jquery_js.read()
        browser.execute_script(jquery)  #active the jquery lib

def getField(lines, no):
    # print lines[no]
    return lines[no].split(u'：')[1].strip()


def downFile(url, fn):
    filename = os.path.join(os.getcwd(), 'zlwj', fn+'.zip')
    urllib.urlretrieve(url, filename)
    return filename


def getFileUrlAndDown(br, fn):
    wait = WebDriverWait(br, 30)
    br.execute_script('detailTifDownload("gk")')
    # html body div table.ui_border tbody tr td.ui_c div.ui_inner table.ui_dialog tbody tr td.ui_main div.ui_content div a
    # elem = br.find_element_by_css_selector('div.ui_inner table.ui_dialog td.ui_main div.ui_content div a')
    elem = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR,'div.ui_inner table.ui_dialog td.ui_main div.ui_content div a')))
    url = elem.get_attribute('href')
    rfn = downFile(url, fn)
    return rfn 


def unzip_file(zipfilename, unziptodir):
    if not os.path.exists(unziptodir): os.mkdir(unziptodir, 0777)
    zfobj = zipfile.ZipFile(zipfilename)
    for name in zfobj.namelist():
        name = name.replace('\\','/')
        if name.endswith('/'):
            os.mkdir(os.path.join(unziptodir, name))
        else:
            ext_filename = os.path.join(unziptodir, name)
            ext_dir= os.path.dirname(ext_filename)
            if not os.path.exists(ext_dir) : os.mkdir(ext_dir,0777)
            outfile = open(ext_filename, 'wb')
            outfile.write(zfobj.read(name))
            outfile.close()


def getContentAndParse(d):
    content = d.page_source


lxmap = {}
lxmap[u"发明专利"] = 0
lxmap[u"实用新型"] = 1
lxmap[u"外观设计"] = 2
def getLeixing(lx):
    if lx in lxmap.keys():
        return lxmap[lx]
    else:
        lxmap[lx] = lxmap.values()[-1] + 1
        return lxmap[lx]


ztmap = {}
ztmap[u"有效"] = 0
ztmap[u"无效"] = 1
ztmap[u"在审"] = 2
def getZhuangtai(zt):
    if zt in ztmap.keys():
        return ztmap[zt]
    else:
        ztmap[zt] = ztmap.values()[-1] + 1
        return ztmap[zt]


def getDate(ds):
    dl = ds.split('.')
    return datetime.datetime(int(dl[0]), int(dl[1]), int(dl[2]))


def getZlwj(zlh):
    browser = webdriver.Firefox() # Get local session of firefox
    try:
        browser.get('http://search.cnipr.com/')
        wait = WebDriverWait(browser, 10)
        elem = wait.until(EC.presence_of_element_located((By.ID,'keywords')))
        # elem = browser.find_element_by_css_selector('input#keywords')
        elem.send_keys(zlh)
        elem.send_keys(Keys.RETURN)

        try:
            # elem = browser.find_element_by_id("rememberme")
            elem = wait.until(EC.presence_of_element_located((By.ID,'rememberme')))
            clickAndWait(elem)

            btn = browser.find_element_by_css_selector("html body div table.ui_border tbody tr td.ui_c div.ui_inner table.ui_dialog tbody tr td div.ui_buttons input.ui_state_highlight")
            clickAndWait(btn)
        except NoSuchElementException:
            pass

        oldtab = browser.current_window_handle
        browser.execute_script("viewDetail(0)")
        browser.switch_to_window(browser.window_handles[-1])
        rfn = getFileUrlAndDown(browser, zlh)
        print "Downlaod %s ok" % rfn
        browser.close()
        browser.switch_to_window(oldtab)
        # unzip_file(rfn, os.path.dirname(rfn))
    finally:
        browser.close()


def updateZlTags(zl, oldzl, tags):
    zlh = oldzl['_id']
    # ts = oldzl['tags']
    # for x in tags:
    #     if x not in ts:
    #         ts.append(x)
    # if len(ts) > 0:
    #     zl.update({'_id': zlh}, {'$set': {'tags': ts}})
    zl.update({"_id", zlh}, {"$addToSet": 
        {"tags": {"$each": tags}}})
    return 


def updateZlDate(db):
    it = db.zhuanli.find().skip(0).limit(2000)
    for x in range(0, it.count()-1):
        a = it.next()
        # replace _id with zlh and del zlh field
        oid = a['_id']
        zlh = a['zlh']
        a['_id'] = zlh
        del a['zlh']
        # replace date field
        a['gkr'] = getDate(a['gkr'])
        a['sqr'] = getDate(a['sqr'])
        db.zhuanli.remove({'_id': oid})
        db.zhuanli.insert(a)
    return it.count()


def updateZl(zl, zs):
    zl.insert(zs)


# just and only update tags or zt field 
def insertOrUpdateZl(zl, zs):
    # find 
    it = zl.find({'_id': zs['_id']})
    if it.count() > 0:
        # update it 
        oldzl = it.next()
        if not (oldzl['tags'] == zs['tags']) and len(zs['tags']) > 0:
            updateZlTags(zl, oldzl, zs['tags'])
    else:
        zl.insert(zs)


def zlExists(zl, zlh):
    it = zl.find({'_id', zlh})
    if it.count() > 0:
        return True
    return False


###########################################################

connection = pymongo.Connection('localhost',27017)
db = connection.zhuanliku3
zl = db.zhuanli
db.zhuanli.create_index('zlh', 1, kwargs={'unique':True, 'dropDups':True})
# db.zhuanli.ensureIndex( { zlh: 1 }, { unique: true }, { dropDups: true })


def main(keyword=u"人脸识别", skippage=0):
    # display = Display(visible=0, size=(800,600))
    # display.start()
    init()
    browser = webdriver.Firefox() # Get local session of firefox
    injectJQuery(browser)
    browser.get('http://search.cnipr.com/')
    wait = WebDriverWait(browser, 10)
    elem = wait.until(EC.presence_of_element_located((By.ID,'keywords')))
    elem.clear()
    # elem = browser.find_element_by_css_selector('input#keywords')
    elem.send_keys(keyword)
    elem.send_keys(Keys.RETURN)

    try:
        # elem = browser.find_element_by_id("rememberme")
        elem = wait.until(EC.presence_of_element_located((By.ID,'rememberme')))
        clickAndWait(elem)

        btn = browser.find_element_by_css_selector("html body div table.ui_border tbody tr td.ui_c div.ui_inner table.ui_dialog tbody tr td div.ui_buttons input.ui_state_highlight")
        clickAndWait(btn)
    except NoSuchElementException:
        pass

    # lnk = browser.find_element_by_partial_link_text(u"中国发明专利")
    # clickAndWait(lnk)

    oldtab = browser.current_window_handle

    curPageNo = 1 
    skipped = False
    while True:
        browser.switch_to_window(oldtab)
        if not skipped and skippage > 0:
            lnkGobtn = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, '#pagination_top a.goBtn')))
            inputGo = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, '#gopagePagination')))
            inputGo.send_keys(str(skippage))
            # lnkNext = browser.find_element_by_css_selector("#pagination_top a.next")
            clickAndWait(lnkGobtn)
            curPageNo = skippage
            skipped = True

        print("Current Page NO: %d" % curPageNo)
        # fetch page content and parse it
        # html = browser.page_source
        items = browser.find_elements_by_css_selector("div.g_item div.g_tit")
        # browser.execute_script("document.getElementsByClassName('classname')[0].style.display='block'")
        # browser.execute_script('return jQuery("div.g_item")')
        # for x in items:
        #     k = x.text.split()
        #     print(string.join(k[:3], ','))

        for x in range(0, len(items)):
            browser.switch_to_window(oldtab)
            vt = items[x].text.split()
            zs = {"bt": vt[0], "lx": getLeixing(vt[1]), "zt": getZhuangtai(vt[2])}
            # print zs
            browser.execute_script("viewDetail(%d)" % x)
            # print browser.title
            # newtab = browser.current_window_handle
            browser.switch_to_window(browser.window_handles[-1])
            injectJQuery(browser)
            # ht = browser.execute_script('return $("div#nct1.n_cont1 div.nc_right").html()')

            content = ''
            lines = []
            try:
                if zs['lx'] == 2: # u'外观设计':
                    # print "------ Enter waiguangsheji ---------------"
                    # elem = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, 'div.info_con div.x_warp1 div.x_table table')))
                    content = browser.execute_script('return $("div.info_con div.x_warp1 div.x_table table").text()')
                    # content = elem.text
                    lines = content.split("\n")
                    # need to process it, remove empty lines
                    nls=filter(lambda x: len(x.strip())>0, lines) 

                    if len(nls) > 0:
                        lines = nls
                else:
                    elem = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, 'div#nct1.n_cont1 div.nc_right')))
                    content = elem.text
                    # elem = browser.find_element_by_css_selector('div#nct1.n_cont1 div.nc_right')
                    elem2 = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, 'div.info_con div#nct1.n_cont1 div.nc_left')))
                    xt = elem2.text
                    xl = xt.split()
                    # 摘要
                    zs['zhaiyao'] = xl[3].strip()
                    # 主权项
                    zs['zhuquanxiang'] = xl[5].strip()
                    lines = content.split('\n')
            except TimeoutException:
                # waiguansheji 
                # print "--------- TimeoutException -----------"
                continue
            
            # 申请(专利)号
            zs['_id'] = getField(lines, 0)
            # 申请日
            zs['sqr'] = getDate(getField(lines, 1))
            # 授权公告号
            zs['gbh'] = getField(lines, 2)
            # 授权公告日
            zs['gkr'] = getDate(getField(lines, 3))
            # 主分类号
            zs['zflh'] = getField(lines, 4)
            # 分类号
            zs['flh'] = getField(lines, 5)
            # 申请权利人
            zs['qlr'] = getField(lines, 6)
            # 发明设计人
            zs['sjr'] = getField(lines, 7)
            # 地址
            zs['dz'] = getField(lines, 8)
            # 国省代码
            zs['gsdm'] = getField(lines, 9)
            zs['tags'] = []
            zs['tags'].append(keyword)

            insertOrUpdateZl(zl, zs)
            # zl.insert(zs)

            # have a limit for count of download tiff 
            # rfn = getFileUrlAndDown(browser, zs['zlh'])
            # unzip_file(rfn, os.path.dirname(rfn))

            browser.close()
            browser.switch_to_window(oldtab)

        try:
            # lnkNext = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, '#pagination_top a.next')))
            lnkNext = browser.find_element_by_css_selector("#pagination_top a.next")
            clickAndWait(lnkNext)
            curPageNo += 1
        except NoSuchElementException:
            browser.close()
            break

    print("Current Page NO: %d" % curPageNo)
    return



if __name__ == '__main__':
    if len(sys.argv) < 2:
        print "Need one argument."
        sys.exit()
    kw = unicode(sys.argv[1], 'utf8')
    skipno = 0
    if len(sys.argv) >= 3:
        skipno = int(sys.argv[2])
    main(kw, skipno)
    # getZlwj("CN201120292255.1")  {"tags": ["人脸识别"]}




