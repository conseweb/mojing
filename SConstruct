
# Library('foo', ['f1.c', 'f2.c', 'f3.c'])
# StaticLibrary('foo', ['f1.c', 'f2.c', 'f3.c'])
# SharedLibrary('foo', ['f1.c', 'f2.c', 'f3.c'])
#
# 
# Library('foo', ['f1.c', 'f2.c', 'f3.c'])
# Program('prog.c', LIBS=['foo', 'bar'], LIBPATH='.')
# 
# Object('hello.c', CCFLAGS='-DHELLO')
# Object('goodbye.c', CCFLAGS='-DGOODBYE')
# Program(['hello.o', 'goodbye.o'])

# File and Dir
# 
# hello_c = File('hello.c')
# Program(hello_c)
# 
# classes = Dir('classes')
# Java(classes, 'src')

# Program('hello.c', CPPPATH = ['include', '/home/project/inc'])

# env = Environment( MSVC_USE_SCRIPT = "c:\\Program Files (x86)\\Microsoft Visual Studio 11.0\\VC\\bin\\vcvars32.bat")

# env.Append(CCFLAGS = ['-g','-O3'])
# env.Append(CPPDEFINES=['BIG_ENDIAN']) 
# env.Append(CPPDEFINES={'RELEASE_BUILD' : '1'})
# env.Append(LIBPATH = ['/usr/local/lib/'])
# env.Append(LIBS = ['SDL_image','GL'])
# env.Append(LINKFLAGS = ['-Wl,--rpath,/usr/local/lib/'])

import sys
import os

# SConstruct
env = Environment(TOOLS=['default', 'go'], ENV = {'PATH' : os.environ['PATH']+':/home/ubuntu/go/bin'})
#env.Append(GO_LIBPATH = [''])
env.Decider('MD5')

srclist = ['src/main.go']

# A simple program
env.GoProgram('bin/mojing', srclist)

# Decider('MD5-timestamp')



