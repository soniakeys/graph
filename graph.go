// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import "math/big"

//go:generate cp cg_label.go cg_adj.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w cg_adj.go
//go:generate gofmt -r "n.To -> n" -w cg_adj.go
//go:generate gofmt -r "Half -> NI" -w cg_adj.go

//go:generate cp cg_undir.go cg_undir_al.go
//go:generate gofmt -r "UndirectedLabeled -> Undirected" -w cg_undir_al.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w cg_undir_al.go
//go:generate gofmt -r "n.To -> n" -w cg_undir_al.go
//go:generate gofmt -r "Half -> NI" -w cg_undir_al.go

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

// NodeList satisfies sort.Interface.
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
