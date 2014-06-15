# coding: utf8

'''
{"tags": {"$and": {"$in": ["虹膜识别"]} , {"$in": ["指纹识别"]}}}
{"tags": {"$in": ["虹膜识别", "指纹识别"]}}



'''

import re
import pymongo

import time
import string
import urllib
import zipfile
import os, sys

connection = pymongo.Connection('localhost',27017)
db = connection.zhuanliku3
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


def getCountByQlrKw(kw, lx=0):
    pat = re.compile(u".*%s.*" % kw)
    it = db.zhuanli.find({'qlr': pat, 'lx': lx, 'zt': {'$nin': [1]}})
    return it.count()


def getFmCountByQlrKw(kw):
    pat = re.compile(u".*%s.*" % kw)
    # it = db.zhuanli.find({'qlr': pat, 'lx': u'发明专利', 'zt': {'$nin': [u'无效', u'在审']}})
    it = db.zhuanli.find({'qlr': pat, 'lx': 0, 'zt': 0})
    return it.count()


def getQlrStats(tag):
    xm = {}
    xl = getAllQlr(tag)
    for x in xl:
        cnt = getFmCountByQlrKw(x)
        if cnt > 0:
            xm[x] = cnt
    sortDictByValues(xm)


def getAllQlr(tag):
    ql = []
    it = db.zhuanli.find({"tags": {"$in": [tag]}}).skip(0).limit(5000)
    for x in it:
        qlr = x['qlr']
        if qlr not in ql:
            ql.append(qlr)
    return ql


# def getAllQlr():
#     ql = []
#     it = db.zhuanli.find().skip(0).limit(5000)
#     for x in it:
#         qlr = x['qlr']
#         if qlr not in ql:
#             ql.append(qlr)
#     return ql


def sortDictByValues(mydict):
    for key, value in sorted(mydict.iteritems(), key=lambda (k,v): (v,k), reverse=True):
        print("%s, %s" % (key.encode('utf8'), value)) 


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print "Need one argument."
        sys.exit()
    kw = unicode(sys.argv[1], 'utf8')
    # lc = len(sys.argv)
    # cnt = 0
    # if lc == 2:
    #   lx = u'发明专利'
    #   cnt = getCountByQlrKw(kw, lx)
    # elif lc >= 3:
    #   lx = unicode(sys.argv[2], 'utf8')
    #   cnt = getCountByQlrKw(kw, lx)
    # print cnt
    getQlrStats(kw)

