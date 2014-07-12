// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleBreadthFirstPath() {
	g := graph.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	}
	fmt.Println(graph.BreadthFirstPath(g, 1, 3))
	// Output:
	// [1 4 3]
}

func ExampleBreadthFirst_Path() {
	b := graph.NewBreadthFirst(graph.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	})
	start, end := 1, 3
	if !b.Path(start, end) {
		return
	}
	fmt.Println("Path length:", b.Result.Paths[end].Len)
	fmt.Print("Backtrack to start: ", end)
	rp := b.Result.Paths
	for n := end; n != start; {
		n = rp[n].From
		fmt.Print(" ", n)
	}
	fmt.Println()
	// Output:
	// Path length: 3
	// Backtrack to start: 3 4 1
}

func ExampleBreadthFirst_AllPaths() {
	b := graph.NewBreadthFirst(graph.AdjacencyList{
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

func ExampleBreadthFirst2Path() {
	g := graph.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	}
	fmt.Println(graph.BreadthFirst2Path(g, 1, 3))
	// Output:
	// [1 4 3]
}

func ExampleBreadthFirst2_Path() {
	g := graph.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	}
	t, m := g.Transpose()
	b := graph.NewBreadthFirst2(g, t, m)
	start, end := 1, 3
	if !b.Path(start, end) {
		return
	}
	fmt.Println("Path length:", b.Result.Paths[end].Len)
	fmt.Print("Backtrack to start: ", end)
	rp := b.Result.Paths
	for n := end; n != start; {
		n = rp[n].From
		fmt.Print(" ", n)
	}
	fmt.Println()
	// Output:
	// Path length: 3
	// Backtrack to start: 3 4 1
}

func ExampleBreadthFirst2_AllPaths() {
	g := graph.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	}
	t, m := g.Transpose()
	b := graph.NewBreadthFirst2(g, t, m)
	b.AllPaths(1)
	fmt.Println("Max path length:", b.Result.MaxLen)
	for n := range b.To {
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

func ExampleDijkstraPath() {
	g := graph.WeightedAdjacencyList{
		1: {{2, 7}, {3, 9}, {6, 11}},
		2: {{3, 10}, {4, 15}},
		3: {{4, 11}, {6, 2}},
		4: {{5, 7}},
		6: {{5, 9}},
	}
	p, l := graph.DijkstraPath(g, 1, 5)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{1 +Inf} {6 11} {5 9}]
	// Path length: 20
}

func ExampleDijkstra_Path() {
	d := graph.NewDijkstra(graph.WeightedAdjacencyList{
		1: {{2, 7}, {3, 9}, {6, 11}},
		2: {{3, 10}, {4, 15}},
		3: {{4, 11}, {6, 2}},
		4: {{5, 7}},
		6: {{5, 9}},
	})
	start, end := 1, 5
	if !d.Path(start, end) {
		return
	}
	fmt.Println("Path distance:", d.Result.Paths[end].Dist)
	fmt.Print("Backtrack to start: ", end)
	rp := d.Result.Paths
	for n := end; n != start; {
		n = rp[n].From.From
		fmt.Print(" ", n)
	}
	fmt.Println()
	// Output:
	// Path distance: 20
	// Backtrack to start: 5 6 1
}

func ExampleDijkstra_AllPaths() {
	g := [][]graph.Half{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1.0}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: nil,
	}
	d := graph.NewDijkstra(g)
	fmt.Println(d.AllPaths(2), "paths found.")
	// column len is from Result, and will be equal to len(path).
	// column dist is from Result, and will be equal to sum.
	fmt.Println("node:  path                       len  dist   sum")
	for nd := range g {
		r := &d.Result.Paths[nd]
		path, dist := d.Result.PathTo(nd)
		fmt.Printf("%d:     %-27s %d   %4.1f  %4.1f\n",
			nd, fmt.Sprint(path), r.Len, r.Dist, dist)
	}
	// Output:
	// 4 paths found.
	// node:  path                       len  dist   sum
	// 0:     []                          0   +Inf  +Inf
	// 1:     []                          0   +Inf  +Inf
	// 2:     [{2 +Inf}]                  1    0.0   0.0
	// 3:     [{2 +Inf} {3 1.1}]          2    1.1   1.1
	// 4:     [{2 +Inf} {3 1.1} {4 0.6}]  3    1.7   1.7
	// 5:     [{2 +Inf} {5 0.2}]          2    0.2   0.2
}

func ExampleAStarAPath() {
	g := graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	}
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	h := func(from int) float64 { return h4[from] }
	p, l := graph.AStarAPath(g, 0, 4, h)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{0 +Inf} {2 0.9} {3 1.1} {4 0.6}]
	// Path length: 2.6
}

func ExampleAStarMPath() {
	g := graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	}
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	h := func(from int) float64 { return h4[from] }
	p, l := graph.AStarMPath(g, 0, 4, h)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [{0 +Inf} {2 0.9} {3 1.1} {4 0.6}]
	// Path length: 2.6
}

func ExampleHeuristic_Admissable() {
	g := graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	}
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	var h graph.Heuristic = func(from int) float64 { return h4[from] }
	fmt.Println(h.Admissable(g, 4))
	// Output:
	// true
}

func ExampleHeuristic_Monotonic() {
	g := graph.WeightedAdjacencyList{
		0: {{1, .7}, {2, .9}, {5, 1.4}},
		1: {{2, 1}, {3, 1.5}},
		2: {{3, 1.1}, {5, .2}},
		3: {{4, .6}},
		4: {{5, .9}},
		5: {},
	}
	h4 := []float64{1.9, 2, 1, .6, 0, .9}
	var h graph.Heuristic = func(from int) float64 { return h4[from] }
	fmt.Println(h.Monotonic(g))
	// Output:
	// true
}

func ExampleBellmanFord() {
	g := graph.WeightedAdjacencyList{
		1: {{2, 10}, {8, 8}},
		2: {{6, 2}},
		3: {{2, 1}, {4, 1}},
		4: {{5, 3}},
		5: {{6, -1}},
		6: {{3, -2}},
		7: {{6, -1}, {2, -4}},
		8: {{7, 1}},
		9: {{5, -10}, {4, 7}},
	}
	b := graph.NewBellmanFord(g)
	if !b.Run(1) {
		return
	}
	fmt.Println("end    path  path")
	fmt.Println("node:  len   dist")
	for n, p := range b.Result.Paths {
		fmt.Printf("%d:       %d   %4.0f\n", n, p.Len, p.Dist)
	}
	// Output:
	// end    path  path
	// node:  len   dist
	// 0:       0   +Inf
	// 1:       1      0
	// 2:       4      5
	// 3:       6      5
	// 4:       7      6
	// 5:       8      9
	// 6:       5      7
	// 7:       3      9
	// 8:       2      8
	// 9:       0   +Inf
}

// duplicate code in instr_test.go
type testCase struct {
	w graph.WeightedAdjacencyList // generated directed graph
	// variants
	unit graph.WeightedAdjacencyList // unit arc weights
	g    graph.AdjacencyList         // unweighted
	t    graph.AdjacencyList         // transpose

	h graph.Heuristic

	start, end, m int

	// kronecker
	kg      graph.AdjacencyList // undirected
	km      int                 // number of arcs
	kStarts []int               // a list of random start points
}

var s = rand.New(rand.NewSource(59))
var r100 = r(100, 200, 62)
var (
	r1k, r10k, r100k/*, r1m*/ testCase
	once   sync.Once
	bigger = func() {
		r1k = r(1e3, 3e3, 66)   // (15x as many arcs as r100)
		r10k = r(1e4, 5e4, 59)  // (17x as many arcs as r1k)
		r100k = r(1e5, 1e6, 59) // (20x as many arcs as r10k)
		//		r1m = r(1e6, 16e6, 59) // (16x as many arcs as r100k)
	}
)

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
	tc.w = make(graph.WeightedAdjacencyList, nNodes)
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
		for _, nb := range tc.w[n1] {
			if nb.To == n2 {
				dup++
				continue arc
			}
		}
		tc.w[n1] = append(tc.w[n1], graph.Half{To: n2, ArcWeight: dist})
		i++
	}
	// variants
	tc.unit = make(graph.WeightedAdjacencyList, len(tc.w))
	for n, nbs := range tc.w {
		u := make([]graph.Half, len(nbs))
		for i, nb := range nbs {
			u[i] = graph.Half{To: nb.To, ArcWeight: 1}
		}
		tc.unit[n] = u
	}
	tc.g = tc.w.Unweighted()
	tc.t, tc.m = tc.g.Transpose()

	// kronecker
	scale := uint(math.Log2(float64(nNodes)) + .5)
	arcFactor := float64(nArcs) / float64(nNodes)
	tc.kg, tc.km = graph.KroneckerUndir(scale, arcFactor)
	// extract giant connected component
	rep, nc := tc.kg.ConnectedComponents()
	var x, max int
	for i, n := range nc {
		if n > max {
			x = i
			max = n
		}
	}
	gcc := new(big.Int)
	tc.kg.DepthFirst(rep[x], gcc, nil)
	// 16 start points, a little arbitrary but that's what Graph 500 uses.
	tc.kStarts = make([]int, 16)
	for i := 0; i < len(tc.kStarts); {
		if s := rand.Intn(len(tc.kg)); gcc.Bit(s) == 1 {
			tc.kStarts[i] = s
			i++
		}
	}
	return tc
}

// end duplicate code

func TestR(t *testing.T) {
	tcs := []testCase{r100, r1k, r10k, r100k} //, r1m}
	if testing.Short() {
		tcs = tcs[:1]
	}
	for _, tc := range tcs {
		if s, cx := tc.w.Unweighted().Simple(); !s {
			t.Fatal(len(tc.w), "not simple at node", cx)
		}
	}
}

func TestSSSP(t *testing.T) {
	tx := func(tc testCase) {
		d := graph.NewDijkstra(tc.w)
		d.Path(tc.start, tc.end)
		pathD, distD := d.Result.PathTo(tc.end)
		// test that repeating same search on same d gives same result
		d.Path(tc.start, tc.end)
		path2, dist2 := d.Result.PathTo(tc.end)
		if len(pathD) != len(path2) || distD != dist2 {
			t.Fatal(len(tc.w), "D, D2 len or dist mismatch")
		}
		for i, half := range pathD {
			if path2[i] != half {
				t.Fatal(len(tc.w), "D, D2 path mismatch")
			}
		}
		// A*
		pathA, distA := graph.AStarAPath(tc.w, tc.start, tc.end, tc.h)
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
		d.AllPaths(tc.start)
		b := graph.NewBellmanFord(tc.w)
		b.Run(tc.start)
		// result objects should be identical
		dr := d.Result
		br := b.Result
		if dr.NoPath != br.NoPath {
			t.Fatal("dr.NoPath, br.NoPath", dr.NoPath, br.NoPath)
		}
		if len(dr.Paths) != len(br.Paths) {
			t.Fatal("len(dr.Paths), len(br.Paths)",
				len(dr.Paths), len(br.Paths))
		}
		for i, de := range dr.Paths {
			if de != br.Paths[i] {
				t.Fatal("dr.Paths ne br.Paths")
			}
		}
		// breadth first, compare to dijkstra with unit weights
		u := graph.NewDijkstra(tc.unit)
		u.AllPaths(tc.start)
		ur := u.Result
		bfs := graph.NewBreadthFirst(tc.g)
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
		bfs2 := graph.NewBreadthFirst2(tc.g, tc.t, tc.m)
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
	tx(r100)
	if testing.Short() {
		t.Skip()
	}
	once.Do(bigger)
	tx(r1k)
	tx(r10k)
	tx(r100k)
	//	tx(r1m)
}

func BenchmarkDijkstra100(b *testing.B) {
	// 100 nodes, 200 edges
	tc := r100
	d := graph.NewDijkstra(tc.w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e3(b *testing.B) {
	// 1000 nodes, 3000 edges
	once.Do(bigger)
	tc := r1k
	d := graph.NewDijkstra(tc.w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e4(b *testing.B) {
	// 10k nodes, 50k edges
	once.Do(bigger)
	tc := r10k
	d := graph.NewDijkstra(tc.w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkDijkstra1e5(b *testing.B) {
	// 100k nodes, 1m edges
	once.Do(bigger)
	tc := r100k
	d := graph.NewDijkstra(tc.w)
	for i := 0; i < b.N; i++ {
		d.AllPaths(tc.start)
	}
}

func BenchmarkBFS_K100(b *testing.B) {
	tc := r100
	bf := graph.NewBreadthFirst(tc.kg)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

func BenchmarkBFS_K1e3(b *testing.B) {
	once.Do(bigger)
	tc := r1k
	bf := graph.NewBreadthFirst(tc.kg)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

func BenchmarkBFS_K1e4(b *testing.B) {
	once.Do(bigger)
	tc := r10k
	bf := graph.NewBreadthFirst(tc.kg)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

func BenchmarkBFS_K1e5(b *testing.B) {
	once.Do(bigger)
	tc := r100k
	bf := graph.NewBreadthFirst(tc.kg)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

/*
func BenchmarkBFS_K1e6(b *testing.B) {
	once.Do(bigger)
	tc := r1m
	bf := graph.NewBreadthFirst(tc.kg)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}
*/
func BenchmarkBFS2_K100(b *testing.B) {
	tc := r100
	bf := graph.NewBreadthFirst2(tc.kg, tc.kg, tc.km)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

func BenchmarkBFS2_K1e3(b *testing.B) {
	once.Do(bigger)
	tc := r1k
	bf := graph.NewBreadthFirst2(tc.kg, tc.kg, tc.km)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

func BenchmarkBFS2_K1e4(b *testing.B) {
	once.Do(bigger)
	tc := r10k
	bf := graph.NewBreadthFirst2(tc.kg, tc.kg, tc.km)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

func BenchmarkBFS2_K1e5(b *testing.B) {
	once.Do(bigger)
	tc := r100k
	bf := graph.NewBreadthFirst2(tc.kg, tc.kg, tc.km)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}

/*
func BenchmarkBFS2_K1e6(b *testing.B) {
	once.Do(bigger)
	tc := r1m
	bf := graph.NewBreadthFirst2(tc.kg, tc.kg, tc.km)
	x := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.AllPaths(tc.kStarts[x])
		x = (x + 1) % len(tc.kStarts)
	}
}
*/
