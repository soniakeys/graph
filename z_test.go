// +build big

// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"testing"

	"github.com/soniakeys/graph"
)

var (
	r1k   = r(1e3, 3e3, 66) // (15x as many arcs as r100)
	r10k  = r(1e4, 5e4, 59) // (17x as many arcs as r1k)
	r100k = r(1e5, 1e6, 59) // (20x as many arcs as r10k)
	//      r1m = r(1e6, 16e6, 59) // (16x as many arcs as r100k)
)

func TestRBig(t *testing.T) {
	for _, tc := range []testCase{r1k, r10k, r100k /*, r1m*/} {
		if s, cx := tc.g.Simple(); !s {
			t.Fatal(len(tc.g), "not simple at node", cx)
		}
	}
}

func TestSSSPBig(t *testing.T) {
	testSSSP(r1k, t)
	testSSSP(r10k, t)
	testSSSP(r100k, t)
	//    testSSSP(r1m, t)
}

func BenchmarkDijkstra1e3(b *testing.B) {
	// 1000 nodes, 3000 edges
	tc := r1k
	w := func(label int) float64 { return tc.w[label] }
	d := graph.NewDijkstra(tc.l, w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e4(b *testing.B) {
	// 10k nodes, 50k edges
	tc := r10k
	w := func(label int) float64 { return tc.w[label] }
	d := graph.NewDijkstra(tc.l, w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e5(b *testing.B) {
	// 100k nodes, 1m edges
	tc := r100k
	w := func(label int) float64 { return tc.w[label] }
	d := graph.NewDijkstra(tc.l, w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

var (
	k10 = k(10, 10, 8)
	k13 = k(13, 12, 16)
	k16 = k(16, 14, 32)

//	k20 = k(20, 16, 64)
)

func BenchmarkBFS_K10(b *testing.B) {
	kt := k10
	bf := graph.NewBreadthFirst(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}

func BenchmarkBFS_K13(b *testing.B) {
	kt := k13
	bf := graph.NewBreadthFirst(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}

func BenchmarkBFS_K16(b *testing.B) {
	kt := k16
	bf := graph.NewBreadthFirst(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}

/*
func BenchmarkBFS_K20(b *testing.B) {
	kt := k20
	bf := graph.NewBreadthFirst(kt.g)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(kt.starts[x])
		x = (x + 1) % len(kt.starts)
	}
}
*/

func BenchmarkBFS2_K10(b *testing.B) {
	tc := k10
	bf := graph.NewBreadthFirst2(tc.g, tc.g, tc.m)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.starts[x])
		x = (x + 1) % len(tc.starts)
	}
}

func BenchmarkBFS2_K13(b *testing.B) {
	tc := k13
	bf := graph.NewBreadthFirst2(tc.g, tc.g, tc.m)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.starts[x])
		x = (x + 1) % len(tc.starts)
	}
}

func BenchmarkBFS2_K16(b *testing.B) {
	tc := k16
	bf := graph.NewBreadthFirst2(tc.g, tc.g, tc.m)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.starts[x])
		x = (x + 1) % len(tc.starts)
	}
}

/*
func BenchmarkBFS2_K20(b *testing.B) {
    tc := k20
    bf := graph.NewBreadthFirst2(tc.g, tc.g, tc.m)
    x := 0
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bf.AllPaths(tc.starts[x])
        x = (x + 1) % len(tc.starts)
    }
}
*/
