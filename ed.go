// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Ed is a simple and fast graph library.
//
// For goals of speed and simplicity, Ed uses zero-based integer node IDs
// and omits interfaces that would accomodate user data or user implemented
// behavior.
//
// To use Ed functions, you typically create a data structure parallel
// to your application data, call an Ed function, and use the result to
// navigate your application data.
//
// Terminology
//
// This package uses the term "node" rather than "vertex."  It uses "arc"
// to mean a directed edge, and uses "from" and "to" to refer to the ends
// of an arc.
//
// A node that is a neighbor of itself is termed a "loop."  Duplicate
// neighbors, or when a node appears more than once in the same neighbor list
// are termed "parallel arcs."
package ed

// file ed.go contains definitions common to different search functions

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
//
// Adjacency lists are inherently directed. To represent an undirected graph,
// create reciprocal neighbors.
type AdjacencyList [][]int

// A WeightedAdjacencyList represents a graph as a list of neighbors for each
// node, connected by weighted arcs.
type WeightedAdjacencyList [][]Half

// Half is a half arc, representing a weighted arc and the "neighbor" node
// that the arc leads to.
//
// Halfs can be composed to form a weighted adjacency list.
type Half struct {
	To        int // node ID, usable as a slice index
	ArcWeight float64
}

// HalfFrom is a half arc, representing a weighted arc and the "neighbor" node
// that the arc originates from.
type HalfFrom struct {
	From      int
	ArcWeight float64
}
