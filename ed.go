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

// FromHalf is a half arc, representing a directed weighted arc and the
// "neighbor" node that the arc originates from.
type FromHalf struct {
	From      int
	ArcWeight float64
}
