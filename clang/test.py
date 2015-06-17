#!/usr/bin/env python
""" Usage: call with <filename> <typename>

how to find libclang.dylib
mdfind -name libclang.dylib
or
find / -name 'libclang.dylib' 
find / -type f -name libclang.dylib -o -name libclang.so 2> /dev/null
"""

import sys
import clang.cindex

def find_typerefs(node, typename):
    """ Find all references to the type named 'typename'
    """
    if node.kind.is_reference():
        print node
        ref_node = clang.cindex.Cursor(node)
        if ref_node.spelling == typename:
            print 'Found %s [line=%s, col=%s]' % (
                typename, node.location.line, node.location.column)
    # Recurse for children of this node
    for c in node.get_children():
        find_typerefs(c, typename)


clang.cindex.Config.set_library_file('/Library/Developer/CommandLineTools/usr/lib/libclang.dylib')
index = clang.cindex.Index.create()
tu = index.parse(sys.argv[1])
print 'Translation unit:', tu.spelling
find_typerefs(tu.cursor, sys.argv[2])

