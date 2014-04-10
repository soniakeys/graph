Ed
==

A graph library.  For goals of speed and simplicity, Ed uses zero-based
integer node IDs and omits maps and interfaces that would accomodate
user data or user implemented behavior.

To use Ed functions, you may have to create a data structure parallel
to an existing application data structure.  After calling an Ed search
method, you then use the result to navigate the original application
data structure.

This library is a reaction to other Go graph libraries where you adapt your
application data to support library functions, either by implementing
interfaces or by storing application data in library data structures.

The name Ed means nothing.  You could think of it as short for someone’s
name or standing for something but really it’s just short an easy to type.
