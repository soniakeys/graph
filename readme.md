Graph
=====

A graph library.  For goals of speed and simplicity, It uses zero-based
integer node IDs and omits maps and interfaces that would accomodate
user data or user implemented behavior.

To use functions of this package, you may have to create a data structure
parallel to an existing application data structure.  After calling a search
method, you then use the result to navigate the original application
data structure.

For a different approach, where graphs are defined with interfaces you
can implement on your own types, see github.com/soniakeys/graph2.
