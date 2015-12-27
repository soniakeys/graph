// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph is a simple and fast graph library.
//
// This is a graph library of integer indexes.  To use it with application
// data, associate data with integer indexes, perform searches or other
// operations with the library, and then use the integer index results to refer
// back to the application data.
//
// Thus it does not store application data, pointers to application data,
// or require application data to implement an interface.  The idea is to
// keep the library methods fast and lean.
//
// Representation overview
//
// The package currently implements three graph representations through types
// AdjacencyList, LabeledAdjacencyList, and FromList.
//
// AdjacencyList is the common "list of lists" representation.  It is a list
// with one element for each node of the graph.  Each element is a list
// itself, a list of neighbor nodes.
//
// LabeledAdjacencyList is similar, but each node-to-neighbor "arc" has an
// associated label.
//
// FromList is a compact rooted tree respresentation.  Like AdjacencyList and
// LabeledAdjacencyList, it is a list with one element for each node of the
// graph.  Each element contains only a single neighbor however, its parent
// in the tree, the "from" node.
//
// Terminology
//
// This package uses the term "node" rather than "vertex."  It uses "arc"
// to mean a directed edge, and uses "from" and "to" to refer to the ends
// of an arc.  It uses "start" and "end" to refer to endpoints of a search
// or traversal.
//
// The usage of "to" and "from" is perhaps most strange.  Throughout the
// the package they are used as adjectives, for example to refer to the
// "from node" of an arc or the "to node".  The type "FromList" is named
// to indicate it stores a list of "from" values.
//
// A "half arc" refers to just one end of an arc, either the to or from end.
//
// Two arcs are "reciprocal" if they connect two distinct nodes n1 and n2,
// one arc leading from n1 to n2 and the other arc leading from n2 to n1.
// "Undirected graphs" are represented with reciprocal arcs.
//
// A node that is a neighbor of itself represents a "loop."  Duplicate
// neighbors (when a node appears more than once in the same neighbor list)
// represent "parallel arcs."  A graph with no loops or parallel arcs
// is "simple."  A graph that allows parallel arcs is a "multigraph"
//
// A number of graph search algorithms use a concept of arc "weights."
// The sum of arc weights along a path is a "distance."  In contrast, the
// number of nodes in a path, including start and end nodes, is the path's
// "length."  (Yes, mixing weights and lengths would be nonsense physically,
// but the terms used here are just distinct terms for abstract values.
// The actual meaning to an application is not relevant within this package.)
//
// Finally, this package documentation takes back the word "object" in some
// places to refer to a Go value, especially a value of a type with methods.
//
// Single source shortest path searches on weighted graphs
//
// This package implements a number of single source shortest path searches.
// These all work with graphs that are directed or undirected, and with graphs
// that may have loops or parallel arcs.  For weighted graphs, "Shortest"
// is defined as the path distance (sum of arc weights) with path length
// (number of nodes) breaking ties.  If multiple paths have the same minimum
// distance with the same minumum length, search methods are free to return
// any of them.
//
//  Type name      Description, methods
//  BestFirst      Unweigted arcs, traversal, single path search or all paths.
//  BestFirst2     Direction-optimizing variant of BestFirst.
//  Dijkstra       Non-negative arc weights, single or all paths.
//  AStar          Non-negative arc weights, heuristic guided, single path.
//  BellmanFord    Negative arc weights allowed, no negative cycles, all paths.
//
// These searches are all done in a similar way that involves creating a
// search object, running a search method on the object, and decoding a result
// structure.  Convenience functions are provided that perform these
// steps for single path searches.
package graph
