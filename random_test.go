// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package graph_test

import (
	"testing"

	"github.com/soniakeys/graph"
)

func TestEuclidean(t *testing.T) {
	var g graph.Directed
	var err error
	for {
		if g, _, err = graph.Euclidean(10, 30, 2, 10, nil); err == nil {
			break
		}
	}
	if s, n := g.IsSimple(); !s {
		t.Fatalf("Euclidean returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
}

func TestKroneckerDir(t *testing.T) {
	g, _ := graph.KroneckerDir(10, 10, nil)
	if s, n := g.IsSimple(); !s {
		t.Fatalf("KroneckerDir returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
}

func TestKroneckerUndir(t *testing.T) {
	g, _ := graph.KroneckerUndir(10, 10, nil)
	if s, n := g.IsSimple(); !s {
		t.Fatalf("KroneckerUndir returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
	if u, from, to := g.IsUndirected(); !u {
		t.Fatalf("KroneckerUndir returned directed graph.  "+
			"Arc %d->%d has no reciprocal.", from, to)
	}
}
