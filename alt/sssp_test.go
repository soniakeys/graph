// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func ExampleBreadthFirst_allPaths() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	var f graph.FromList
	alt.BreadthFirst(g, nil, 0, 1, &f, func(n graph.NI) bool {
		return true
	})
	fmt.Println("Max path length:", f.MaxLen)
	p := make([]graph.NI, f.MaxLen)
	for n := range g {
		fmt.Println(n, f.PathTo(graph.NI(n), p))
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

/*
func TestSSSP(t *testing.T) {
	r100 := r(100, 200, 62)
	testSSSP(r100, t)
}

func testSSSP(tc testCase, t *testing.T) {
	w := func(label graph.LI) float64 { return tc.w[label] }
	f, dist, _ := tc.l.LabeledAdjacencyList.Dijkstra(tc.start, tc.end, w)
	pathD := f.PathTo(tc.end, nil)
	distD := dist[tc.end]
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
	dr, _, _ := tc.l.LabeledAdjacencyList.Dijkstra(tc.start, -1, w)
	br, _, _ := tc.l.BellmanFord(w, tc.start)
	// result objects should be identical
	if len(dr.Paths) != len(br.Paths) {
		t.Fatal("len(dr.Paths), len(br.Paths)",
			len(dr.Paths), len(br.Paths))
	}
	// breadth first, compare to dijkstra with unit weights
	w = func(graph.LI) float64 { return 1 }
	ur, _, _ := tc.l.LabeledAdjacencyList.Dijkstra(tc.start, -1, w)
	var bfsr graph.FromList
	np, _ := tc.g.AdjacencyList.BreadthFirst(tc.start, nil, &bfsr,
		func(n graph.NI) bool { return true })
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
	var bfs2r graph.FromList
	np2 := graph.BreadthFirst2(tc.g.AdjacencyList, tc.t.AdjacencyList, tc.m,
		tc.start, &bfs2r, func(n graph.NI) bool { return true })
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
*/
