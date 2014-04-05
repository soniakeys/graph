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
package ed

// file ed.go contains definitions common to different search functions

// Half is a half arc, representing a directed weighted arc and the "neighbor"
// node that the arc leads to.
//
// Halfs can be composed to form an adjacency list.
type Half struct {
	To        int // node ID, usable as a slice index
	ArcWeight float64
}

// An adjacency list represents a graph as a list of neighbors for each node.
// Node IDs correspond to slice indexes of the AdjacencyList. Each neighbor
// list is a []Half where Half.To fields contain a node ID greater than or
// equal to zero and strictly less than len() of the AdjacencyList.
type AdjacencyList [][]Half

// FromHalf is a half arc, representing a directed weighted arc and the
// "neighbor" node that the arc originates from.
type FromHalf struct {
	From      int
	ArcWeight float64
}
