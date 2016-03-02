// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/soniakeys/graph"
)

// duplicate code in instr_test.go
type testCase struct {
	l graph.DirectedLabeled // generated labeled directed graph
	w []float64             // arc weights for l
	// variants
	g graph.Directed // unlabeled
	t graph.Directed // transpose

	h graph.Heuristic

	start, end graph.NI
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
	tc := testCase{start: graph.NI(s.Intn(nNodes))}
	// end is point at distance nearest target distance
	const target = .3
	nearest := 2.
	c1 := coords[tc.start]
	for i, c2 := range coords {
		d := math.Abs(target - math.Hypot(c2.x-c1.x, c2.y-c1.y))
		if d < nearest {
			tc.end = graph.NI(i)
			nearest = d
		}
	}
	// with end chosen, define heuristic
	ce := coords[tc.end]
	tc.h = func(n graph.NI) float64 {
		cn := &coords[n]
		return math.Hypot(ce.x-cn.x, ce.y-cn.y)
	}
	// graph
	tc.l = graph.DirectedLabeled{make(graph.LabeledAdjacencyList, nNodes)}
	tc.w = make([]float64, nArcs)
	// arcs
	var tooFar, dup int
arc:
	for i := 0; i < nArcs; {
		if tooFar == nArcs || dup == nArcs {
			panic(fmt.Sprint("tooFar", tooFar, "dup", dup, "nArcs", nArcs,
				"nNodes", nNodes, "seed", seed))
		}
		n1 := graph.NI(s.Intn(nNodes))
		n2 := n1
		for n2 == n1 {
			n2 = graph.NI(s.Intn(nNodes)) // no graph loops
		}
		c1 := &coords[n1]
		c2 := &coords[n2]
		dist := math.Hypot(c2.x-c1.x, c2.y-c1.y)
		if dist > s.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		for _, nb := range tc.l.LabeledAdjacencyList[n1] {
			if nb.To == n2 {
				dup++
				continue arc
			}
		}
		tc.w[i] = dist
		tc.l.LabeledAdjacencyList[n1] = append(tc.l.LabeledAdjacencyList[n1],
			graph.Half{To: n2, Label: graph.LI(i)})
		i++
	}
	// variants
	tc.g = tc.l.Unlabeled()
	tc.t, tc.m = tc.g.Transpose()
	return tc
}

// end duplicate code

func TestR(t *testing.T) {
	tc := r100
	if s, cx := tc.g.IsSimple(); !s {
		t.Fatal(len(tc.w), "not simple at node", cx)
	}
}

// kronecker test case
type kronTest struct {
	// parameters:
	scale      uint
	edgeFactor float64
	starts     []graph.NI // the parameter here is len(starts)
	// generated:
	g graph.Undirected
	m int // number of arcs in g
	// also generated are values for starts[]
}

var k7 = k(7, 8, 4)

// generate kronecker graph and start points
func k(scale uint, ef float64, nStarts int) (kt kronTest) {
	kt.g, kt.m = graph.KroneckerUndir(scale, ef)
	// extract giant connected component
	rep, nc := kt.g.ConnectedComponentReps()
	var x, max int
	for i, n := range nc {
		if n > max {
			x = i
			max = n
		}
	}
	gcc := new(big.Int)
	kt.g.DepthFirst(rep[x], gcc, nil)
	kt.starts = make([]graph.NI, nStarts)
	for i := 0; i < nStarts; {
		if s := rand.Intn(len(kt.g.AdjacencyList)); gcc.Bit(s) == 1 {
			kt.starts[i] = graph.NI(s)
			i++
		}
	}
	return
}

func testSSSP(tc testCase, t *testing.T) {
	w := func(label graph.LI) float64 { return tc.w[label] }
	d := graph.NewDijkstra(tc.l.LabeledAdjacencyList, w)
	d.Path(tc.start, tc.end)
	pathD := d.Tree.PathTo(tc.end, nil)
	distD := d.Dist[tc.end]
	// test that repeating same search on same d gives same result
	d.Path(tc.start, tc.end)
	path2 := d.Tree.PathTo(tc.end, nil)
	dist2 := d.Dist[tc.end]
	if len(pathD) != len(path2) || distD != dist2 {
		t.Fatal(len(tc.w), "D, D2 len or dist mismatch")
	}
	for i, half := range pathD {
		if path2[i] != half {
			t.Fatal(len(tc.w), "D, D2 path mismatch")
		}
	}
	// A*
	pathA, distA := tc.l.AStarAPath(tc.start, tc.end, tc.h, w)
	// test that a* path is same distance and length as dijkstra path
	if len(pathA) != len(pathD) {
		t.Log("pathA:", pathA)
		t.Log("pathD:", pathD)
		t.Fatal(len(tc.w), "A, D len mismatch")
	}
	if distA != distD {
		t.Log("distA:", distA)
		t.Log("distD:", distD)
		t.Log("delta:", math.Abs(distA-distD))
		t.Fatal(len(tc.w), "A, D dist mismatch")
	}
	// test Bellman Ford against Dijkstra all paths
	d.Reset()
	d.AllPaths(tc.start)
	b := graph.NewBellmanFord(tc.l.LabeledAdjacencyList, w)
	b.Start(tc.start)
	// result objects should be identical
	dr := d.Tree
	br := b.Tree
	if len(dr.Paths) != len(br.Paths) {
		t.Fatal("len(dr.Paths), len(br.Paths)",
			len(dr.Paths), len(br.Paths))
	}
	/* this test not working, possibly not a valid test.
	t.Log(dr.Paths)
	t.Log(br.Paths)
	for i, de := range dr.Paths {
		t.Log(de, br.Paths[i])
		if de != br.Paths[i] {
			t.Fatal("dr.Paths ne br.Paths")
		}
	}
	*/
	// breadth first, compare to dijkstra with unit weights
	d.Weight = func(graph.LI) float64 { return 1 }
	d.Reset()
	d.AllPaths(tc.start)
	ur := d.Tree
	bfs := graph.NewBreadthFirst(tc.g.AdjacencyList)
	np := bfs.AllPaths(tc.start)
	bfsr := bfs.Result
	var ml, npf int
	for i, ue := range ur.Paths {
		bl := bfsr.Paths[i].Len
		if bl != ue.Len {
			t.Fatal("ue.From.Len, bfsr.Paths[i].Len", ue.Len, bl)
		}
		if bl > ml {
			ml = bl
		}
		if bl > 0 {
			npf++
		}
	}
	if ml != bfsr.MaxLen {
		t.Fatal("bfsr.MaxLen, recomputed", bfsr.MaxLen, ml)
	}
	if npf != np {
		t.Fatal("bfs all paths returned", np, "recount:", npf)
	}
	// breadth first 2
	bfs2 := graph.NewBreadthFirst2(tc.g.AdjacencyList, tc.t.AdjacencyList, tc.m)
	np2 := bfs2.AllPaths(tc.start)
	bfs2r := bfs2.Result
	var ml2, npf2 int
	for i, e := range bfsr.Paths {
		bl2 := bfs2r.Paths[i].Len
		if bl2 != e.Len {
			t.Fatal("bfsr.Paths[i].Len, bfs2r", e.Len, bl2)
		}
		if bl2 > ml2 {
			ml2 = bl2
		}
		if bl2 > 0 {
			npf2++
		}
	}
	if ml2 != bfs2r.MaxLen {
		t.Fatal("bfs2r.MaxLen, recomputed", bfs2r.MaxLen, ml)
	}
	if npf2 != np2 {
		t.Fatal("bfs2 all paths returned", np2, "recount:", npf2)
	}
	if ml2 != ml {
		t.Fatal("bfs max len, bfs2", ml, ml2)
	}
	if npf2 != npf {
		t.Fatal("bfs return, bfs2", npf, npf2)
	}
}

func TestSSSP(t *testing.T) {
	testSSSP(r100, t)
}

func BenchmarkDijkstra100(b *testing.B) {
	// 100 nodes, 200 edges
	tc := r100
	w := func(label graph.LI) float64 { return tc.w[label] }
	d := graph.NewDijkstra(tc.l.LabeledAdjacencyList, w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkBFS_K07(b *testing.B) {
	tc := k7
	bf := graph.NewBreadthFirst(tc.g.AdjacencyList)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.starts[x])
		x = (x + 1) % len(tc.starts)
	}
}

func BenchmarkBFS2_K07(b *testing.B) {
	tc := k7
	bf := graph.NewBreadthFirst2(tc.g.AdjacencyList, tc.g.AdjacencyList, tc.m)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.starts[x])
		x = (x + 1) % len(tc.starts)
	}
}
