******************
  SCons Go Tools
******************

SCons Go Tools is a collection of builders that makes it easy to compile Go_
projects in SCons_.

.. _Go: http://golang.org/
.. _SCons: http://www.scons.org/

================
  Installation
================

Download and extract the latest release of SCons Go Tools from the project
homepage, then copy the go.py script into the site_tools subdirectory of your
site_scons directory.

If you don't have a site_scons directory yet, you will need to create one.  By
default, SCons looks for the site_scons directory in the current directory.  If
you only need to use this for one project, just create a site_scons at the root
of your project source and that will work fine.  For other setups, consult the
SCons documentation.

=============
  Embedding
=============

If you want to use SCons Go Tools in your own project, but you aren't sure
whether your end-users will have the tools installed, you can embed the tools
into your project's source tree.  Simply copy the SCons Go Tools directory into
site_scons/site_tools at your project's root.  Make sure that the SCons Go
Tools directory is named ``go``, or your SConstruct file won't be able to find
it.

=========
  Usage
=========

Once the tools are installed, using them in your SConstruct file is easy::

   # SConstruct
   env = Environment(TOOLS=['default', 'go'])
   
   # A simple program
   env.GoProgram('foo', 'src/foo.go')
   
   # A multi-package program
   bar = env.Go('bar', 'src/bar.go')
   env.Go('baz', ['src/baz1.go', 'src/baz2.go'])
   env.GoProgram('bar', bar)
   
   # Cross-compiling
   windowsEnv = env.Clone()
   windowsEnv.GoTarget('windows', '386')
   windowsEnv.GoProgram('test.exe', 'src/testwindows.go')
   
.. Note::
   You don't specify all of the object files when you go to link a program;
   just the one that contains the main function.  The Go linker does this for
   you automatically; however, the SCons Go Tools are smart enough to
   determine the dependencies as well, so the program will always be rebuilt
   when one of the packages changes.

===========
  Testing
===========

If you want to use gotest-style unit tests in your project, SCons Go Tools
allows you to easily collect them and produce a test program::

    # SConstruct
    env = Environment(TOOLS=['default', 'go'])
    env['GO_STRIPTESTS'] = False
    mypackage = env.Go('mypackage.go', 'mypackage_test.go')
    env.GoProgram('runtests', env.GoTest('tests.go', mypackage))

This will collect all of the functions whose name starts with ``Test`` in the
mypackage package and make a source file called ``tests.go``.  ``tests.go`` is
compiled and linked into a program called ``runtests``, which you can then run
from the command line to run your unit tests.

===============
  Environment
===============

The SCons Go Tools use the following parameters set in the Environment object:

GO_GC
   The path to the `gc`_ program for this platform.

GO_GCFLAGS
    Flags for `gc`_.

GO_LD
   The path to the `ld`_ program for this platform.

GO_LDFLAGS
   Flags for `ld`_.

GO_A
   The path to the assembler program.

GO_LIBPATH
   A list of paths that will be searched for imports (this is used for both
   compiling and linking, since you will usually be using the same place for
   both)

GO_STRIPTESTS
    Whether to ignore Go source files that end in ``_test.go``.  Defaults to
    False.

.. _gc: http://golang.org/cmd/gc/
.. _ld: http://golang.org/cmd/ld/

===========
  License
===========

Copyright (c) 2010, Ross Light
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

   Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

   Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

   Neither the name of the SCons Go Tools nor the names of its contributors may
   be used to endorse or promote products derived from this software without
   specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

.. vim: ft=rst sw=3 sts=3 et
