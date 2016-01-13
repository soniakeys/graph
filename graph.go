// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// graph.go
//
// Definitions for unlabeled graphs, and methods not specific to directed
// or undirected graphs.  Method docs need not mention that they work on both.

package graph

import (
	"math/big"
	"sort"
)

//go:generate cp cg_label.go cg_adj.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w cg_adj.go
//go:generate gofmt -r "n.To -> n" -w cg_adj.go
//go:generate gofmt -r "Half -> NI" -w cg_adj.go

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
// It is a node number.  It is used extensively as a slice index.
// Node numbers also account for a significant fraction of the memory
// required to represent a graph.
type NI int32

type NodeList []NI

func (l NodeList) Len() int           { return len(l) }
func (l NodeList) Less(i, j int) bool { return l[i] < l[j] }
func (l NodeList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
//
// Adjacency lists are inherently directed. To represent an undirected graph,
// create reciprocal neighbors.
type AdjacencyList [][]NI

// Simple checks for loops and parallel arcs.
//
// A graph is "simple" if it has no loops or parallel arcs.
//
// Simple returns true, -1 for simple graphs.  If a loop or parallel arc is
// found, simple returns false and and a node that represents a counterexample
// to the graph being simple.
func (g AdjacencyList) Simple() (s bool, n NI) {
	var t NodeList
	for n, nbs := range g {
		if len(nbs) == 0 {
			continue
		}
		t = append(t[:0], nbs...)
		sort.Sort(t)
		if t[0] == NI(n) {
			return false, NI(n)
		}
		for i, nb := range t[1:] {
			if nb == NI(n) || nb == t[i] {
				return false, NI(n)
			}
		}
	}
	return true, -1
}
