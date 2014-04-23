// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// cross_test.go, tests across multiple search algorithms

package graph_test

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/soniakeys/graph"
)

// duplicate code in instr_test.go
type testCase struct {
	g          graph.WeightedAdjacencyList
	start, end int
	h          graph.Heuristic
}

var s = rand.New(rand.NewSource(59))
var r100 = r(100, 200, 62)
var (
	r1k, r10k, r100k testCase
	once             sync.Once
	bigger           = func() {
		r1k = r(1e3, 3e3, 66)   // (15x as many arcs as r100)
		r10k = r(1e4, 5e4, 59)  // (17x as many arcs as r1k)
		r100k = r(1e5, 1e6, 59) // (20x as many arcs as r10k)
	}
)

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
	tc.g = make(graph.WeightedAdjacencyList, nNodes)
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
		tc.g[n1] = append(tc.g[n1], graph.Half{n2, dist})
		i++
	}
	return tc
}

// end duplicate code

func TestR(t *testing.T) {
	tcs := []testCase{r100, r1k, r10k, r100k}
	if testing.Short() {
		tcs = tcs[:1]
	}
	for _, tc := range tcs {
		if s, cx := tc.g.Unweighted().Simple(); !s {
			t.Fatal(len(tc.g), "not simple at node", cx)
		}
	}
}

func TestSSSP(t *testing.T) {
	tx := func(tc testCase) {
		d := graph.NewDijkstra(tc.g)
		pathD, distD := d.Path(tc.start, tc.end)
		// test that repeating same search on same d gives same result
		path2, dist2 := d.Path(tc.start, tc.end)
		if len(pathD) != len(path2) || distD != dist2 {
			t.Fatal(len(tc.g), "D, D2 len or dist mismatch")
		}
		for i, half := range pathD {
			if path2[i] != half {
				t.Fatal(len(tc.g), "D, D2 path mismatch")
			}
		}
		// A*
		a := graph.NewAStar(tc.g)
		pathA, distA := a.AStarAPath(tc.start, tc.end, tc.h)
		// test that a* path is same distance and length as dijkstra path
		if len(pathA) != len(pathD) {
			t.Log("pathA:", pathA)
			t.Log("pathD:", pathD)
			t.Fatal(len(tc.g), "A, D len mismatch")
		}
		//fudge coded when math was a little different. not needed currently.
		//if math.Abs((distA - distD)/distA) > 1e-15 {
		if distA != distD {
			t.Log("distA:", distA)
			t.Log("distD:", distD)
			t.Log("delta:", math.Abs(distA-distD))
			t.Fatal(len(tc.g), "A, D dist mismatch")
		}
	}
	tx(r100)
	if testing.Short() {
		t.Skip()
	}
	once.Do(bigger)
	tx(r1k)
	tx(r10k)
	tx(r100k)
}
