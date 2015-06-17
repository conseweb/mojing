import sys
import clang.cindex

def callexpr_visitor(node, parent, userdata):
    if node.kind == clang.cindex.CursorKind.CALL_EXPR:
        print 'Found %s [line=%s, col=%s]' % (
                node.spelling, node.location.line, node.location.column)
    return 2 # means continue visiting recursively


clang.cindex.Config.set_library_file('/Library/Developer/CommandLineTools/usr/lib/libclang.dylib')
index = clang.cindex.Index.create()
tu = index.parse(sys.argv[1])
# clang.cindex.Cursor(
#         tu.cursor,
#         clang.cindex.Cursor_visit_callback(callexpr_visitor),
#         None)


