// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

var s = rand.New(rand.NewSource(59))

// generate random directed graph and end points to test
func r(nNodes, nArcs int, seed int64) (g Graph, start, end int) {
	s.Seed(seed)
	// generate random coordinates
	type xy struct{ x, y float64 }
	coords := make([]xy, nNodes)
	for i := range coords {
		coords[i].x = s.Float64()
		coords[i].y = s.Float64()
	}
	// random start
	start = s.Intn(nNodes)
	// end is point at distance nearest target distance
	const target = .3
	nearest := 2.
	c1 := coords[start]
	for i, c2 := range coords {
		if d := math.Abs(target - math.Hypot(c2.x-c1.x, c2.y-c1.y)); d < nearest {
			end = i
			nearest = d
		}
	}
	// graph
	g = New(nNodes)
	// arcs
	var tooFar, dup int
	for i := 0; i < nArcs; {
		if tooFar == nArcs || dup == nArcs {
			panic(fmt.Sprint("tooFar", tooFar, "dup", dup, "nArcs", nArcs,
				"nNodes", nNodes, "seed", seed))
		}
		n1 := s.Intn(nNodes)
		n2 := s.Intn(nNodes)
		c1 := &coords[n1]
		c2 := &coords[n2]
		dist := math.Hypot(c2.x-c1.x, c2.y-c1.y)
		if dist > s.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		if !g.AddArcSimple(n1, Half{n2, dist}) { // only count additions
			dup++
			continue
		}
		i++
	}
	return
}

func Test100(t *testing.T) {
	g, start, end := r(100, 200, 62)
	t.Log(g.ShortestPath(start, end))
	t.Log("NV AV:", ndVis, arcVis)
	t.Log(g.ShortestPath(start, end))
	t.Log("NV AV:", ndVis, arcVis)
}

func Benchmark100(b *testing.B) {
	// 100 nodes, 200 edges
	g, start, end := r(100, 200, 62)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.ShortestPath(start, end)
	}
}

func Test1e3(t *testing.T) {
	g, start, end := r(1000, 3000, 66)
	t.Log(g.ShortestPath(start, end))
	t.Log("NV AV:", ndVis, arcVis)
}

func Benchmark1e3(b *testing.B) {
	// 1000 nodes, 3000 edges
	g, start, end := r(1000, 3000, 66)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.ShortestPath(start, end)
	}
}

func Test1e4(t *testing.T) {
	g, start, end := r(1e4, 5e4, 59)
	t.Log(g.ShortestPath(start, end))
	t.Log("NV AV:", ndVis, arcVis)
}

func Benchmark1e4(b *testing.B) {
	// 10k nodes, 100k edges
	g, start, end := r(1e4, 5e4, 59)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.ShortestPath(start, end)
	}
}

func Test1e5(t *testing.T) {
	g, start, end := r(1e5, 1e6, 59)
	t.Log(g.ShortestPath(start, end))
	t.Log("NV AV:", ndVis, arcVis)
}

func Benchmark1e5(b *testing.B) {
	// 100k nodes, 1m edges
	g, start, end := r(1e5, 1e6, 59)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.ShortestPath(start, end)
	}
}
