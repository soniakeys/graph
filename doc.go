// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Graph is a simple and fast graph library.
//
// This is a graph library of the kind where you create graphs out of
// concrete types defined in this package, perhaps parallel to existing
// graph data structures in your application.  You call some function such
// as a graph search, then use the result to navigate your application data.
//
// The types of this package represent only data minimally neccessary for
// search functions.  This minimalism simplifies package code and allows faster
// searches.  Zero-based integer node IDs serve directly as slice indexes.
// Nodes and edges are structs rather than interfaces.  Maps are not needed
// to associate arbitrary IDs with node or edge types.  These graphs are memory
// efficient and large graphs can potentially be handled.
//
// Representation
//
// The only graph representation currently is an adjacency list, although
// there are separate types for arc-weighted graphs and unweighted graphs.
// The types AdjacencyList and WeightedAdjacencyList are simply slices
// of slices.  Construct with make; there is no special constructor.
// Directed and undirected graphs use the same types.  Construct an undirected
// graph by adding reciprocal edges.  Methods specific to either directed
// or undirected graphs will be documented as such.
//
// Terminology
//
// This package uses the term "node" rather than "vertex."  It uses "arc"
// to mean a directed edge, and uses "from" and "to" to refer to the ends
// of an arc.  It uses "start" and "end" to refer to endpoints of a search
// or traversal.
//
// A float64 value associated with an arc is "weight."  The sum of arc weights
// along a path is a "distance."  The number of nodes in a path, including
// start and end nodes, is the path's "length."
//
// A "half arc" represents just one end of an arc, perhaps assocating it with
// an arc weight.  The more common half to work with is the "to half" (the
// type name is simply "Half".)  A list of half arcs can represent a
// "neighbor list," neighbors of a single node.  A list of neighbor lists
// forms an "adjacency list" which represents a directed graph.
//
// Two arcs are "reciprocal" if they connect two distinct nodes n1 and n2,
// one arc leading from n1 to n2 and the other arc leading from n2 to n1.
// Undirected graphs are represented with reciprocal arcs.  A graph is
// undirected if a reciprocal arc exists for every arc connecting distinct
// nodes.
//
// A node that is a neighbor of itself represents a "loop."  Duplicate
// neighbors (when a node appears more than once in the same neighbor list)
// represent "parallel arcs."  A graph with no loops or parallel arcs
// is "simple."
//
// Finally, this package documentation takes back the word "object" to
// refer to a Go value, especially a value of a type with methods.
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
//
// The result structure contains a "from-tree" which is an efficient
// encoding of the path results.  Other terms for this data structure include
// "predecessor list", "in-tree", "inverse arborescence", and "spaghetti stack."
package graph
