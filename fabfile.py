#!/usr/bin/python
# -*- coding: utf-8 -*-

from __future__ import with_statement
from fabric.api import *
from fabric.contrib.console import confirm
from fabric.context_managers import *
# from fabric.api import env, run
import string
from random import choice
import socket
import paramiko
from fabric.contrib.project import rsync_project


# env.use_ssh_config = True
env.user = 'ubuntu'
# env.password = 'rtmap2013'
env.hosts = ['54.200.138.17', ]
env.key_filename = 'aws.pem'
env.roledefs = {
    'test': ['115.28.171.71']
}

remote_home = '/home/ubuntu'
remote_working_dir = '/home/ubuntu/workspace/mojing'

def host_type():
    run('uname -s')


# @task
# @parallel
# def passwd(user, passwd=False):
#     with settings(hide('running', 'stdout', 'stderr'), warn_only=True):
#         if isup(env.host):
#             if not passwd:
#                 passwd = genpass()
#             sudo("echo -e '%s\n%s' | passwd %s" % (passwd, passwd, user))

# def genpass(length=10):
#     return ''.join(choice(string.ascii_letters + string.digits) for _ in range(length))

# def isup(host):
#     print 'connecting host: %s' % host
#     timeout = socket.getdefaulttimeout()
#     socket.setdefaulttimeout(1)
#     up = True
#     try:
#         paramiko.Transport((host, 22))
#     except Exception, e:
#         up = False
#         print '%s down, %s' % (host, e)
#     finally:
#         socket.setdefaulttimeout(timeout)
#         return up

def test():
    with settings(warn_only=True):
        result = local('./manage.py test my_app', capture=True)
    if result.failed and not confirm("Tests failed. Continue anyway?"):
        abort("Aborting at user request.")

def commit():
    local("git add -p && git commit")

def push():
    local("git push origin master")

def pull():
    local("git pull origin master")

def clone(user, host, repos):
    local("git clone ssh://%s@%s/%s.git" % (user, host, repos))

def prepare_deploy():
    test()
    commit()
    push()

def deploy():
    code_dir = '/srv/django/myproject'
    with settings(warn_only=True):
        if run("test -d %s" % code_dir).failed:
            run("git clone user@vcshost:/path/to/repo/.git %s" % code_dir)
    with cd(code_dir):
        run("git pull")
        run("touch app.wsgi")

def deploy_proj(loc):
    rsync_project(local_dir=loc, remote_dir=remote_working_dir, exclude=['.git','bin','data','logs'])
    stop()
    build()
    start()

def start_circusd():
    with cd(remote_working_dir):
        run("/usr/local/bin/circusd --daemon circus.ini")

def quit_circusd():
    run("/usr/local/bin/circusctl quit")

def start():
    run("/usr/local/bin/circusctl start mojing")

def stop():
    run("/usr/local/bin/circusctl stop mojing")

def build():
    with shell_env(GOROOT="/home/ubuntu/go", GOPATH="/home/ubuntu/.go"):
        with cd(remote_working_dir):
            run('/home/ubuntu/go/bin/go build -o bin/mojing src/main.go', shell=False)



