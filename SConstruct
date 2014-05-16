
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

# SConstruct
env = Environment(TOOLS=['default', 'go'])

# inc = ['/usr/local/include/', '/usr/local/lib/erlang/usr/include/']
# lib_path = ['/usr/lib', '/usr/local/lib']
# g++ -I /usr/local/boost_1_42_0 -lboost_system -lboost_thread a.cpp
# g++ -I /usr/local/boost_1_42_0 -lboost_system-mt -lboost_thread-mt a.cpp
# libs = ['websock', 'protobuf-c', 'erl_interface', 'ei']
# lib_path = ['.'] # LIBPATH = ['/usr/lib', '/usr/local/lib']

# common = Glob('../common/*.cpp')
# srclist = ['LocNetPacketParser.cpp', 'rtmap_nif.cpp', 'niftools.cpp']
srclist = ['src/server.go']

# env = Environment(TARGET_ARCH='x86_64')

# exec external command to compile device.proto to c 
# device.pb-c.c depends on device.proto 
# env.Command('device.pb.go', 'device.proto', 'protoc --go_out=. device.proto')
# A simple program
env.GoProgram('bin/mojing', srclist)

# A multi-package program
# bar = env.Go('bar', srclist)
# env.Go('baz', ['src/baz1.go', 'src/baz2.go'])
# env.GoProgram('bar', bar)

# Decider('MD5-timestamp')
env.Decider('MD5')


# A simple program
