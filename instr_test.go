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

// duplicate code in random_test.go (execpt for package qualification)

type testCase struct {
	l LabeledAdjacencyList // generated labeled directed graph
	w []float64            // arc weights for l
	// variants
	g AdjacencyList // unlabeled
	t AdjacencyList // transpose

	h Heuristic

	start, end NI
	m          int
}

var s = rand.New(rand.NewSource(59))
var r100 = r(100, 200, 62)

// generate random graphs and end points to test
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
	tc := testCase{start: NI(s.Intn(nNodes))}
	// end is point at distance nearest target distance
	const target = .3
	nearest := 2.
	c1 := coords[tc.start]
	for i, c2 := range coords {
		d := math.Abs(target - math.Hypot(c2.x-c1.x, c2.y-c1.y))
		if d < nearest {
			tc.end = NI(i)
			nearest = d
		}
	}
	// with end chosen, define heuristic
	ce := coords[tc.end]
	tc.h = func(n NI) float64 {
		cn := &coords[n]
		return math.Hypot(ce.x-cn.x, ce.y-cn.y)
	}
	// graph
	tc.l = make(LabeledAdjacencyList, nNodes)
	tc.w = make([]float64, nArcs)
	// arcs
	var tooFar, dup int
arc:
	for i := 0; i < nArcs; {
		if tooFar == nArcs || dup == nArcs {
			panic(fmt.Sprint("tooFar", tooFar, "dup", dup, "nArcs", nArcs,
				"nNodes", nNodes, "seed", seed))
		}
		n1 := NI(s.Intn(nNodes))
		n2 := n1
		for n2 == n1 {
			n2 = NI(s.Intn(nNodes)) // no graph loops
		}
		c1 := &coords[n1]
		c2 := &coords[n2]
		dist := math.Hypot(c2.x-c1.x, c2.y-c1.y)
		if dist > s.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		for _, nb := range tc.l[n1] {
			if nb.To == n2 {
				dup++
				continue arc
			}
		}
		tc.w[i] = dist
		tc.l[n1] = append(tc.l[n1], Half{To: n2, Label: i})
		i++
	}
	// variants
	tc.g = tc.l.Unlabeled()
	tc.t, tc.m = tc.g.Transpose()
	return tc
}

// end duplicate code

func TestInstr(t *testing.T) {
	ti := func(tc testCase) {
		w := func(label int) float64 { return tc.w[label] }
		d := NewDijkstra(tc.l, w)
		d.Path(tc.start, tc.end)
		t.Log("NV AV:", d.ndVis, d.arcVis)
		ndVis1 := d.ndVis
		arcVis1 := d.arcVis
		// test that repeating same search on same d has same signature
		d.Reset()
		d.Path(tc.start, tc.end)
		if d.ndVis != ndVis1 || d.arcVis != arcVis1 {
			t.Fatal(len(tc.g), ndVis1, d.ndVis, arcVis1, d.arcVis)
		}
	}
	ti(r100)
}
