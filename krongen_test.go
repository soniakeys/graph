package graph_test

import (
	"testing"

	"github.com/soniakeys/graph"
)

func TestNewDir(t *testing.T) {
	g, _ := graph.NewDir(10, 10)
	if s, n := g.Simple(); !s {
		t.Fatalf("NewDir returned non-simple graph.  Node %d to: %v", n, g[n])
	}
}

func TestKrongen(t *testing.T) {
	g, _ := graph.NewUnDir(10, 10)
	if s, n := g.Simple(); !s {
		t.Fatalf("NewDir returned non-simple graph.  Node %d to: %v", n, g[n])
	}
	if u, from, to := g.Undirected(); !u {
		t.Fatalf("NewUnDir returned directed graph.  Arc %d->%d has no reciprocal.", from, to)
	}
}
