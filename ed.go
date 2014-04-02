// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

// Half is a half arc, representing a "neighbor" of a node.
//
// Halfs can be composed to form an adjacency list.
type Half struct {
	To        int // node ID, usable as a slice index
	ArcWeight float64
}
