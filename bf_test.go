// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleBreadthFirstBit_AllPaths() {
	b := graph.NewBreadthFirstBit(graph.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	})
	b.AllPaths(1)
	fmt.Println("Max path length:", b.Result.MaxLen)
	for n := range b.Graph {
		fmt.Println(n, b.Result.PathTo(n))
	}
	// Output:
	// Max path length: 4
	// 0 []
	// 1 [1]
	// 2 []
	// 3 [1 4 3]
	// 4 [1 4]
	// 5 [1 4 3 5]
	// 6 [1 4 6]
}

// kronecker test case
type kronTest struct {
	// parameters:
	scale      uint
	edgeFactor float64
	starts     []int // the parameter here is len(starts)
	// generated:
	g graph.AdjacencyList // an undirected graph
	m int
	// also generated are values for starts[]
}

var (
	k7  = k(7, 8, 4)
	k10 = k(10, 10, 8)
	k13 = k(13, 12, 16)
	k16 = k(16, 14, 32)
	k20 = k(20, 16, 64)
)

// generate kronecker graph and start points
func k(scale uint, ef float64, nStarts int) (kt kronTest) {
	kt.g, kt.m = graph.KroneckerUndir(scale, ef)
	// extract giant connected component
	rep, nc := kt.g.ConnectedComponents()
	var x, max int
	for i, n := range nc {
		if n > max {
			x = i
			max = n
		}
	}
	gcc := new(big.Int)
	kt.g.DepthFirst(rep[x], gcc, nil)
	kt.starts = make([]int, nStarts)
	for i := 0; i < nStarts; {
		if s := rand.Intn(len(kt.g)); gcc.Bit(s) == 1 {
			kt.starts[i] = s
			i++
		}
	}
	return
}

func BenchmarkBFBit_K7(b *testing.B) {
	kt := k7
	bf := graph.NewBreadthFirstBit(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}

func BenchmarkBFBit_K10(b *testing.B) {
	kt := k10
	bf := graph.NewBreadthFirstBit(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}

func BenchmarkBFBit_K13(b *testing.B) {
	kt := k13
	bf := graph.NewBreadthFirstBit(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}

func BenchmarkBFBit_K16(b *testing.B) {
	kt := k16
	bf := graph.NewBreadthFirstBit(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}

func BenchmarkBFBit_K20(b *testing.B) {
	kt := k20
	bf := graph.NewBreadthFirstBit(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}
