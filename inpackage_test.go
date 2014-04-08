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
func r(nNodes, nArcs int, seed int64) (g [][]Half, start, end int) {
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
		d := math.Abs(target - math.Hypot(c2.x-c1.x, c2.y-c1.y))
		if d < nearest {
			end = i
			nearest = d
		}
	}
	// graph
	g = make([][]Half, nNodes)
	// arcs
	var tooFar, dup int
arc:
	for i := 0; i < nArcs; {
		if tooFar == nArcs || dup == nArcs {
			panic(fmt.Sprint("tooFar", tooFar, "dup", dup, "nArcs", nArcs,
				"nNodes", nNodes, "seed", seed))
		}
		n1 := s.Intn(nNodes)
		n2 := n1
		for n2 == n1 {
			n2 = s.Intn(nNodes) // no graph loops
		}
		c1 := &coords[n1]
		c2 := &coords[n2]
		dist := math.Hypot(c2.x-c1.x, c2.y-c1.y)
		if dist > s.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		for _, nb := range g[n1] {
			if nb.To == n2 {
				dup++
				continue arc
			}
		}
		g[n1] = append(g[n1], Half{n2, dist})
		i++
	}
	return
}

func Test100(t *testing.T) {
	g, start, end := r(100, 200, 62)
	d := NewDijkstra(g)
	p1, l1 := d.Path(start, end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
	// test that repeating same search on same d gives same result
	p2, l2 := d.Path(start, end)
	if len(p1) != len(p2) || l1 != l2 {
		t.Fatal("len")
	}
	for i, h := range p1 {
		if p2[i] != h {
			t.Fatal("path")
		}
	}
}

func Benchmark100(b *testing.B) {
	// 100 nodes, 200 edges
	g, start, end := r(100, 200, 62)
	d := NewDijkstra(g)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Path(start, end)
	}
}

func Test1e3(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	g, start, end := r(1000, 3000, 66)
	d := NewDijkstra(g)
	d.Path(start, end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
}

func Benchmark1e3(b *testing.B) {
	// 1000 nodes, 3000 edges
	g, start, end := r(1000, 3000, 66)
	d := NewDijkstra(g)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Path(start, end)
	}
}

func Test1e4(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	g, start, end := r(1e4, 5e4, 59)
	d := NewDijkstra(g)
	d.Path(start, end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
}

func Benchmark1e4(b *testing.B) {
	// 10k nodes, 100k edges
	g, start, end := r(1e4, 5e4, 59)
	d := NewDijkstra(g)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Path(start, end)
	}
}

func Test1e5(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	g, start, end := r(1e5, 1e6, 59)
	d := NewDijkstra(g)
	d.Path(start, end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
}

func Benchmark1e5(b *testing.B) {
	// 100k nodes, 1m edges
	g, start, end := r(1e5, 1e6, 59)
	d := NewDijkstra(g)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Path(start, end)
	}
}
