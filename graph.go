// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math/big"
	"sort"
)

//go:generate cp adj_cg.go adj_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w adj_RO.go
//go:generate gofmt -r "n.To -> n" -w adj_RO.go
//go:generate gofmt -r "Half -> NI" -w adj_RO.go

//go:generate cp dir_cg.go dir_RO.go
//go:generate gofmt -r "DirectedLabeled -> Directed" -w dir_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w dir_RO.go
//go:generate gofmt -r "n.To -> n" -w dir_RO.go
//go:generate gofmt -r "Half -> NI" -w dir_RO.go

//go:generate cp undir_cg.go undir_RO.go
//go:generate gofmt -r "UndirectedLabeled -> Undirected" -w undir_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w undir_RO.go
//go:generate gofmt -r "n.To -> n" -w undir_RO.go
//go:generate gofmt -r "Half -> NI" -w undir_RO.go

var one = big.NewInt(1)

// OneBits sets a big.Int to a number that is all 1s in binary.
//
// It's a utility function useful for initializing a bitmap of a graph
// to all ones; that is, with a bit set to 1 for each node of the graph.
//
// OneBits modifies b, then returns b for convenience.
func OneBits(b *big.Int, n int) *big.Int {
	return b.Sub(b.Lsh(one, uint(n)), one)
}

// NI is a "node int"
//
// It is a node number or node ID.  NIs are used extensively as slice indexes.
// NIs typically account for a significant fraction of the memory footprint of
// a graph.
type NI int32

// NodeList satisfies sort.Interface.
type NodeList []NI

func (l NodeList) Len() int           { return len(l) }
func (l NodeList) Less(i, j int) bool { return l[i] < l[j] }
func (l NodeList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
//
// Adjacency lists are inherently directed but can be used to represent
// directed or undirected graphs.  See types Directed and Undirected.
type AdjacencyList [][]NI

// HasParallelSort identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the results fr and to represent an example
// where there are parallel arcs from node fr to node to.
//
// If there are no parallel arcs, the method returns false -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Sort" in the method name indicates that sorting is used to detect parallel
// arcs.  Compared to method HasParallelMap, this may give better performance
// for small or sparse graphs but will have asymtotically worse performance for
// large dense graphs.
func (g AdjacencyList) HasParallelSort() (has bool, fr, to NI) {
	var t NodeList
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		// different code in the labeled version, so no code gen.
		t = append(t[:0], to...)
		sort.Sort(t)
		t0 := t[0]
		for _, to := range t[1:] {
			if to == t0 {
				return true, NI(n), t0
			}
			t0 = to
		}
	}
	return false, -1, -1
}
