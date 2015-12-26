package graph_test

import (
	"testing"

	"github.com/soniakeys/graph"
)

func TestKroneckerDir(t *testing.T) {
	g, _ := graph.KroneckerDir(10, 10)
	if s, n := g.Simple(); !s {
		t.Fatalf("KroneckerDir returned non-simple graph.  Node %d to: %v",
			n, g[n])
	}
}

func TestKroneckerUndir(t *testing.T) {
	g, _ := graph.KroneckerUndir(10, 10)
	if s, n := g.Simple(); !s {
		t.Fatalf("KroneckerUndir returned non-simple graph.  Node %d to: %v",
			n, g[n])
	}
	if u, from, to := g.IsUndirected(); !u {
		t.Fatalf("KroneckerUndir returned directed graph.  "+
			"Arc %d->%d has no reciprocal.", from, to)
	}
}
