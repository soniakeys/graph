package graph_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/soniakeys/graph"
)

func TestEuclidean(t *testing.T) {
	var g graph.Directed
	var err error
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		if g, _, err = graph.Euclidean(10, 30, 2, 10, r); err == nil {
			break
		}
	}
	if s, n := g.IsSimple(); !s {
		t.Fatalf("Euclidean returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
}

func TestKroneckerDir(t *testing.T) {
	g, _ := graph.KroneckerDir(10, 10)
	if s, n := g.IsSimple(); !s {
		t.Fatalf("KroneckerDir returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
}

func TestKroneckerUndir(t *testing.T) {
	g, _ := graph.KroneckerUndir(10, 10)
	if s, n := g.IsSimple(); !s {
		t.Fatalf("KroneckerUndir returned non-simple graph.  Node %d to: %v",
			n, g.AdjacencyList[n])
	}
	if u, from, to := g.IsUndirected(); !u {
		t.Fatalf("KroneckerUndir returned directed graph.  "+
			"Arc %d->%d has no reciprocal.", from, to)
	}
}
