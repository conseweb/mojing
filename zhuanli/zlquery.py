# coding: utf8

import re
import pymongo

import time
import string
import urllib
import zipfile
import os, sys

connection = pymongo.Connection('localhost',27017)
db = connection.zhuanliku
zl = db.zhuanli
# db.zhuanli.create_index('zlh', 1, kwargs={'unique':True, 'dropDups':True})
# db.zhuanli.ensureIndex( { zlh: 1 }, { unique: true }, { dropDups: true })


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


def getCountByQlrKw(kw, lx=u'发明专利'):
    pat = re.compile(u".*%s.*" % kw)
    it = db.zhuanli.find({'qlr': pat, 'lx': lx, 'zt': {'$nin': [u'无效']}})
    return it.count()


def getFmCountByQlrKw(kw):
    pat = re.compile(u".*%s.*" % kw)
    # it = db.zhuanli.find({'qlr': pat, 'lx': u'发明专利', 'zt': {'$nin': [u'无效', u'在审']}})
    it = db.zhuanli.find({'qlr': pat, 'lx': u'发明专利', 'zt': u'有效'})
    return it.count()


def getQlrStats():
    xm = {}
    xl = getAllQlr()
    for x in xl:
        cnt = getFmCountByQlrKw(x)
        if cnt > 0:
            xm[x] = cnt
    sortDictByValues(xm)


def getAllQlr():
    ql = []
    it = db.zhuanli.find().skip(0).limit(3000)
    for x in it:
        qlr = x['qlr']
        if qlr not in ql:
            ql.append(qlr)
    return ql


def sortDictByValues(mydict):
    for key, value in sorted(mydict.iteritems(), key=lambda (k,v): (v,k), reverse=True):
        print(u"%s, %s" % (key, value)) 


if __name__ == '__main__':
    # if len(sys.argv) < 2:
    #     print "Need one argument."
    #     sys.exit()
    # kw = unicode(sys.argv[1], 'utf8')
    # lc = len(sys.argv)
    # cnt = 0
    # if lc == 2:
    #   lx = u'发明专利'
    #   cnt = getCountByQlrKw(kw, lx)
    # elif lc >= 3:
    #   lx = unicode(sys.argv[2], 'utf8')
    #   cnt = getCountByQlrKw(kw, lx)
    # print cnt
    getQlrStats()

