// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
)

type testCase struct {
	g          WeightedAdjacencyList
	start, end int
	h          Heuristic
}

var s = rand.New(rand.NewSource(59))
var r100 = r(100, 200, 62)
var r1k, r10k, r100k testCase

func bigger() {
	r1k = r(1e3, 3e3, 66)   // (15x as many arcs as r100)
	r10k = r(1e4, 5e4, 59)  // (17x as many arcs as r1k)
	r100k = r(1e5, 1e6, 59) // (20x as many arcs as r10k)
}

var once sync.Once // for bigger

// generate random directed graph and end points to test
func r(nNodes, nArcs int, seed int64) testCase {
	s.Seed(seed)
	// generate random coordinates
	type xy struct{ x, y float64 }
	coords := make([]xy, nNodes)
	for i := range coords {
		coords[i].x = s.Float64()
		coords[i].y = s.Float64()
	}
	// random start
	tc := testCase{start: s.Intn(nNodes)}
	// end is point at distance nearest target distance
	const target = .3
	nearest := 2.
	c1 := coords[tc.start]
	for i, c2 := range coords {
		d := math.Abs(target - math.Hypot(c2.x-c1.x, c2.y-c1.y))
		if d < nearest {
			tc.end = i
			nearest = d
		}
	}
	// with end chosen, define heuristic
	ce := coords[tc.end]
	tc.h = func(n int) float64 {
		cn := &coords[n]
		return math.Hypot(ce.x-cn.x, ce.y-cn.y)
	}
	// graph
	tc.g = make(WeightedAdjacencyList, nNodes)
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
		for _, nb := range tc.g[n1] {
			if nb.To == n2 {
				dup++
				continue arc
			}
		}
		tc.g[n1] = append(tc.g[n1], Half{n2, dist})
		i++
	}
	return tc
}

func Test100(t *testing.T) {
	tc := r100
	d := NewDijkstra(tc.g)
	pathD, distD := d.Path(tc.start, tc.end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
	// test that repeating same search on same d gives same result
	path2, dist2 := d.Path(tc.start, tc.end)
	if len(pathD) != len(path2) || distD != dist2 {
		t.Fatal("D, D2 len or dist mismatch")
	}
	for i, half := range pathD {
		if path2[i] != half {
			t.Fatal("D, D2 path mismatch")
		}
	}
	// A*
	a := NewAStar(tc.g)
	pathA, distA := a.AStarA(tc.start, tc.end, tc.h)
	// test that a* path is same distance and length a dijkstra path
	if len(pathA) != len(pathD) || distA != distD {
		t.Fatal("A, D len or dist mismatch")
	}
}

func BenchmarkDijkstra100(b *testing.B) {
	// 100 nodes, 200 edges
	tc := r100
	d := NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func TestDijkstra1e3(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	once.Do(bigger)
	tc := r1k
	d := NewDijkstra(tc.g)
	d.Path(tc.start, tc.end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
}

func BenchmarkDijkstra1e3(b *testing.B) {
	// 1000 nodes, 3000 edges
	once.Do(bigger)
	tc := r1k
	d := NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func TestDijkstra1e4(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	once.Do(bigger)
	tc := r10k
	d := NewDijkstra(tc.g)
	d.Path(tc.start, tc.end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
}

func BenchmarkDijkstra1e4(b *testing.B) {
	// 10k nodes, 50k edges
	once.Do(bigger)
	tc := r10k
	d := NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func TestDijkstra1e5(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	once.Do(bigger)
	tc := r100k
	d := NewDijkstra(tc.g)
	d.Path(tc.start, tc.end)
	t.Log("NV AV:", d.ndVis, d.arcVis)
}

func BenchmarkDijkstra1e5(b *testing.B) {
	// 100k nodes, 1m edges
	once.Do(bigger)
	tc := r100k
	d := NewDijkstra(tc.g)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}
