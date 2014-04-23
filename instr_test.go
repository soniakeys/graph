// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// file instr_test.go
//
// Tests on unexported instrumentation.

package graph

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// duplicate code in cross_test
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

// end duplicate code

func TestInstr(t *testing.T) {
	ti := func(tc testCase) {
		d := NewDijkstra(tc.g)
		d.Path(tc.start, tc.end)
		t.Log("NV AV:", d.ndVis, d.arcVis)
		ndVis1 := d.ndVis
		arcVis1 := d.arcVis
		// test that repeating same search on same d has same signature
		d.Path(tc.start, tc.end)
		if d.ndVis != ndVis1 || d.arcVis != arcVis1 {
			t.Fatal(len(tc.g), ndVis1, d.ndVis, arcVis1, d.arcVis)
		}
	}
	ti(r100)
	if testing.Short() {
		t.Skip()
	}
	bigger()
	ti(r1k)
	ti(r10k)
	ti(r100k)
}
