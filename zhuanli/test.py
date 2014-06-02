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
# from scrapemark import scrape
import pymongo

import time
import string
import urllib
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
    return


def getFileUrlAndDown(br, fn):
    br.execute_script('detailTifDownload("gk")')
    elem = br.find_element_by_css_selector('td.ui_main div.ui_content div a')
    url = elem.get_attribute('href')
    downFile(url, fn)
    return 


def getContentAndParse(d):
    content = d.page_source


# display = Display(visible=0, size=(800,600))
# display.start()

init()
connection = pymongo.Connection('localhost',27017)
db = connection.zhuanliku
zl = db.zhuanli

browser = webdriver.Firefox() # Get local session of firefox
injectJQuery(browser)

# browser.implicitly_wait(10) # seconds
# browser = webdriver.PhantomJS() # or add to your PATH
# browser.set_window_size(1024, 768) # optional
browser.get('http://search.cnipr.com/')
# time.sleep(1)
# browser.save_screenshot('screen.png') # save a screenshot to disk

wait = WebDriverWait(browser, 10)
elem = wait.until(EC.presence_of_element_located((By.ID,'keywords')))
# elem = browser.find_element_by_css_selector('input#keywords')
elem.send_keys(u"人脸识别")
elem.send_keys(Keys.RETURN)

time.sleep(1)

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

curPageNo = 1
while True:
    print("Current Page NO: %d" % curPageNo)
    # fetch page content and parse it
    # html = browser.page_source
    items = browser.find_elements_by_css_selector("div.g_item div.g_tit")
    # browser.execute_script("document.getElementsByClassName('classname')[0].style.display='block'")
    # browser.execute_script('return jQuery("div.g_item")')
    # for x in items:
    #     k = x.text.split()
    #     print(string.join(k[:3], ','))

    oldtab = browser.current_window_handle

    for x in range(0, len(items)):
        vt = items[x].text.split()
        zs = {"bt": vt[0], "lx": vt[1], "zt": vt[2]}
        # print zs
        browser.execute_script("viewDetail(%d)" % x)
        # print browser.title
        # newtab = browser.current_window_handle
        browser.switch_to_window(browser.window_handles[-1])
        # ht = browser.execute_script('return $("div#nct1.n_cont1 div.nc_right").html()')
        elem = browser.find_element_by_css_selector('div#nct1.n_cont1 div.nc_right')
        content = elem.text
        # print "-----------------------------"
        # print content
        # print "*****************************"
        
        lines = content.split('\n')
        zs['zlh'] = getField(lines, 0)
        zs['sqr'] = getField(lines, 1)
        zs['gbh'] = getField(lines, 2)
        zs['gkr'] = getField(lines, 3)
        zs['zflh'] = getField(lines, 4)
        zs['flh'] = getField(lines, 5)
        zs['qlr'] = getField(lines, 6)
        zs['sjr'] = getField(lines, 7)
        zs['dz'] = getField(lines, 8)
        zs['gsdm'] = getField(lines, 9)
        zl.insert(zs)

        getFileUrlAndDown(browser, zs['zlh'])
        browser.close()
        browser.switch_to_window(oldtab)

    # contitems = browser.find_elements_by_css_selector("div.g_item div.g_cont div.g_cont_left")
    # for x in contitems:
    #     k = x.text.split("\n")
    #     print(string.join(k, '|'))

    try:
        lnkNext = browser.find_element_by_css_selector("#pagination_top a.next")
        clickAndWait(lnkNext)
        curPageNo += 1
    except NoSuchElementException:
        break

print("Current Page NO: %d" % curPageNo)

# elem.click()
# /html/body/div[2]/table/tbody/tr[2]/td[2]/div/table/tbody/tr[3]/td/div/input

# browser = webdriver.Firefox()
# browser.get("http://somedomain/url_that_delays_loading")
# try:
#     element = WebDriverWait(browser, 10).until(
#         EC.presence_of_element_located((By.ID, "myDynamicElement"))
#     )
# finally:
#     browser.quit()


# browser = webdriver.Firefox() # Get local session of firefox

# injectJQuery(browser)

# #now you can write some jquery code then execute_script them
# js = """
#     var str = "div#myPager table a:[href=\\"javascript:__doPostBack('myPager','%s')\\"]"
#     console.log(str)
#     var $next_anchor = $(str);
#     if ($next_anchor.length) {
#         return $next_anchor.get(0).click(); //do click and redirect
#     } else {
#         return false;
#     }""" % str(25) 

# success = browser.execute_script(js)
# if success == False:
#     break




