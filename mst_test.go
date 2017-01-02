// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleLabeledUndirected_Kruskal() {
	//       (10)
	//     0------4----\
	//     |     /|     \(70)
	// (30)| (40) |(60)  \
	//     |/     |      |
	//     1------2------3
	//       (50)   (20)
	w := func(l graph.LI) float64 { return float64(l) }
	// undirected graph
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 30)
	g.AddEdge(graph.Edge{0, 4}, 10)
	g.AddEdge(graph.Edge{1, 2}, 50)
	g.AddEdge(graph.Edge{1, 4}, 40)
	g.AddEdge(graph.Edge{2, 3}, 20)
	g.AddEdge(graph.Edge{2, 4}, 60)
	g.AddEdge(graph.Edge{3, 4}, 70)

	t, dist := g.Kruskal(w)

	fmt.Println("spanning tree as undirected graph:")
	for n, to := range t.LabeledAdjacencyList {
		fmt.Println(n, to)
	}
	fmt.Println("total distance: ", dist)
	// Output:
	// spanning tree as undirected graph:
	// 0 [{4 10} {1 30}]
	// 1 [{0 30} {2 50}]
	// 2 [{3 20} {1 50}]
	// 3 [{2 20}]
	// 4 [{0 10}]
	// total distance:  110
}

func ExampleWeightedEdgeList_Kruskal() {
	//       (10)
	//     0------4----\
	//     |     /|     \(70)
	// (30)| (40) |(60)  \
	//     |/     |      |
	//     1------2------3
	//       (50)   (20)
	w := func(l graph.LI) float64 { return float64(l) }
	// construct WeightedEdgeList directly
	l := graph.WeightedEdgeList{5, w, []graph.LabeledEdge{
		{graph.Edge{0, 1}, 30},
		{graph.Edge{0, 4}, 10},
		{graph.Edge{1, 2}, 50},
		{graph.Edge{1, 4}, 40},
		{graph.Edge{2, 3}, 20},
		{graph.Edge{2, 4}, 60},
		{graph.Edge{3, 4}, 70},
	}}

	t, dist := l.Kruskal()

	fmt.Println("spanning tree as undirected graph:")
	for n, to := range t.LabeledAdjacencyList {
		fmt.Println(n, to)
	}
	fmt.Println("total distance: ", dist)
	// Output:
	// spanning tree as undirected graph:
	// 0 [{4 10} {1 30}]
	// 1 [{0 30} {2 50}]
	// 2 [{3 20} {1 50}]
	// 3 [{2 20}]
	// 4 [{0 10}]
	// total distance:  110
}

func ExampleWeightedEdgeList_Kruskal_fromUndirected() {
	//       (10)
	//     0------4----\
	//     |     /|     \(70)
	// (30)| (40) |(60)  \
	//     |/     |      |
	//     1------2------3
	//       (50)   (20)
	w := func(l graph.LI) float64 { return float64(l) }
	// undirected graph
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 30)
	g.AddEdge(graph.Edge{0, 4}, 10)
	g.AddEdge(graph.Edge{1, 2}, 50)
	g.AddEdge(graph.Edge{1, 4}, 40)
	g.AddEdge(graph.Edge{2, 3}, 20)
	g.AddEdge(graph.Edge{2, 4}, 60)
	g.AddEdge(graph.Edge{3, 4}, 70)
	// convert to edge list for Kruskal.
	l := g.WeightedArcsAsEdges(w)

	t, dist := l.Kruskal()

	fmt.Println("spanning tree as undirected graph:")
	for n, to := range t.LabeledAdjacencyList {
		fmt.Println(n, to)
	}
	fmt.Println("total distance: ", dist)
	// Output:
	// spanning tree as undirected graph:
	// 0 [{4 10} {1 30}]
	// 1 [{0 30} {2 50}]
	// 2 [{3 20} {1 50}]
	// 3 [{2 20}]
	// 4 [{0 10}]
	// total distance:  110
}

func ExampleWeightedEdgeList_KruskalSorted() {
	//       (10)
	//     0------4----\
	//     |     /|     \(70)
	// (30)| (40) |(60)  \
	//     |/     |      |
	//     1------2------3
	//       (50)   (20)
	w := func(l graph.LI) float64 { return float64(l) }
	// Bypass construction of an undirected graph if you can, by directly
	// constructing an edge list.  No need for reciprocal arcs.  Also if
	// you can, construct it already sorted by weight.
	l := graph.WeightedEdgeList{5, w, []graph.LabeledEdge{
		{graph.Edge{0, 4}, 10},
		{graph.Edge{2, 3}, 20},
		{graph.Edge{0, 1}, 30},
		{graph.Edge{1, 4}, 40},
		{graph.Edge{1, 2}, 50},
		{graph.Edge{2, 4}, 60},
		{graph.Edge{3, 4}, 70},
	}}

	t, dist := l.KruskalSorted()

	fmt.Println("spanning tree as undirected graph:")
	for n, to := range t.LabeledAdjacencyList {
		fmt.Println(n, to)
	}
	fmt.Println("total distance: ", dist)
	// Output:
	// spanning tree as undirected graph:
	// 0 [{4 10} {1 30}]
	// 1 [{0 30} {2 50}]
	// 2 [{3 20} {1 50}]
	// 3 [{2 20}]
	// 4 [{0 10}]
	// total distance:  110
}

func ExampleLabeledUndirected_Prim() {
	// graph:
	//
	//  (2)     (3)
	//   |\       \
	//   | \       \ 2
	//   |  \       \
	// 4 |   \ 5    (4)
	//   |    \
	//   |     \
	//   |      \
	//  (1)-----(0)
	//       3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 3)
	g.AddEdge(graph.Edge{1, 2}, 4)
	g.AddEdge(graph.Edge{2, 0}, 5)
	g.AddEdge(graph.Edge{3, 4}, 2)
	// weight function
	w := func(arcLabel graph.LI) float64 { return float64(arcLabel) }

	// get connected components
	reps, orders := g.ConnectedComponentReps()
	fmt.Println(len(reps), "connected components:")
	fmt.Println("Representative node  Order (number of nodes in component)")
	for i, r := range reps {
		fmt.Printf("%d %20d\n", r, orders[i])
	}

	a := g.LabeledAdjacencyList
	f := graph.NewFromList(len(a))
	labels := make([]graph.LI, len(a))

	// construct spanning tree for each component
	fmt.Println("Span results:")
	fmt.Println("Root  Nodes spanned  Total tree distance  Leaves")
	for _, r := range reps {
		var leaves graph.Bits
		ns, dist := g.Prim(r, w, &f, labels, &leaves)
		fmt.Printf("%d %17d %20.0f  %d\n", r, ns, dist, leaves.Slice())
	}

	// show final forest
	fmt.Println("Spanning forest:")
	fmt.Println("Node  From  Arc distance  Path length  Leaf")
	for n, pe := range f.Paths {
		fmt.Printf("%d %8d %13.0f %12d %5d\n",
			n, pe.From, w(labels[n]), pe.Len, f.Leaves.Bit(graph.NI(n)))
	}

	// optionally, convert to undirected graph
	u := f.TransposeLabeled(labels).Undirected()
	fmt.Println("Equivalent undirected graph:")
	for fr, to := range u.LabeledAdjacencyList {
		fmt.Printf("%d:  %#v\n", fr, to)
	}

	// Output:
	// 2 connected components:
	// Representative node  Order (number of nodes in component)
	// 0                    3
	// 3                    2
	// Span results:
	// Root  Nodes spanned  Total tree distance  Leaves
	// 0                 3                    7  [2]
	// 3                 2                    2  [4]
	// Spanning forest:
	// Node  From  Arc distance  Path length  Leaf
	// 0       -1             0            1     0
	// 1        0             3            2     0
	// 2        1             4            3     1
	// 3       -1             0            1     0
	// 4        3             2            2     1
	// Equivalent undirected graph:
	// 0:  []graph.Half{graph.Half{To:1, Label:3}}
	// 1:  []graph.Half{graph.Half{To:2, Label:4}, graph.Half{To:0, Label:3}}
	// 2:  []graph.Half{graph.Half{To:1, Label:4}}
	// 3:  []graph.Half{graph.Half{To:4, Label:2}}
	// 4:  []graph.Half{graph.Half{To:3, Label:2}}
}

func TestPrim100(t *testing.T) {
	r100 := r(100, 200, 62)
	u100 := r100.l.Undirected()
	reps, orders := u100.ConnectedComponentReps()
	w := func(l graph.LI) float64 { return r100.w[l] }
	var f graph.FromList
	// construct spanning tree for each component
	for i, r := range reps {
		ns, _ := u100.Prim(r, w, &f, nil, nil)
		if ns != orders[i] {
			t.Fatal("Not all nodes spanned within a connected component.")
		}
	}
}

func BenchmarkMST100(b *testing.B) {
	r100 := r(100, 200, 62)
	u100 := r100.l.Undirected()
	reps, _ := u100.ConnectedComponentReps()
	b.Log(len(reps), "connected components")
	w := func(l graph.LI) float64 { return r100.w[l] }
	b.Run("Krus", func(b *testing.B) { benchKruskal(u100, w, b) })
	b.Run("WKrus", func(b *testing.B) { benchWKrus(u100, w, b) })
	b.Run("WKrusD", func(b *testing.B) { benchWKrusD(r100.l, w, b) })
	b.Run("WKrusS", func(b *testing.B) { benchWKrusS(r100.l, w, b) })
	b.Run("CCPrim", func(b *testing.B) { benchCCPrim(u100, w, b) })
	b.Run("PrAll", func(b *testing.B) { benchPrimAll(u100, w, b) })
	b.Run("PrMin", func(b *testing.B) { benchPrimMin(u100, w, b) })
	b.Run("PrUnd", func(b *testing.B) { benchPrimUnd(u100, w, b) })
}

func benchKruskal(u graph.LabeledUndirected, w graph.WeightFunc, b *testing.B) {
	for i := 0; i < b.N; i++ {
		u.Kruskal(w)
	}
}
func benchWKrus(u graph.LabeledUndirected, w graph.WeightFunc, b *testing.B) {
	wl := u.WeightedArcsAsEdges(w)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wl.Kruskal()
	}
}
func benchWKrusD(d graph.LabeledDirected, w graph.WeightFunc, b *testing.B) {
	wl := d.WeightedArcsAsEdges(w)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wl.Kruskal()
	}
}
func benchWKrusS(d graph.LabeledDirected, w graph.WeightFunc, b *testing.B) {
	wl := d.WeightedArcsAsEdges(w)
	sort.Sort(wl)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wl.KruskalSorted()
	}
}
func benchCCPrim(u graph.LabeledUndirected, w graph.WeightFunc, b *testing.B) {
	f := graph.NewFromList(u.Order())
	lab := make([]graph.LI, u.Order())
	var lvs graph.Bits
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reps, _ := u.ConnectedComponentReps()
		for _, r := range reps {
			u.Prim(r, w, &f, lab, &lvs)
		}
	}
}
func benchPrimAll(u graph.LabeledUndirected, w graph.WeightFunc, b *testing.B) {
	reps, _ := u.ConnectedComponentReps()
	f := graph.NewFromList(u.Order())
	lab := make([]graph.LI, u.Order())
	var lvs graph.Bits
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range reps {
			u.Prim(r, w, &f, lab, &lvs)
		}
	}
}
func benchPrimMin(u graph.LabeledUndirected, w graph.WeightFunc, b *testing.B) {
	reps, _ := u.ConnectedComponentReps()
	f := graph.NewFromList(u.Order())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range reps {
			u.Prim(r, w, &f, nil, nil)
		}
	}
}
func benchPrimUnd(u graph.LabeledUndirected, w graph.WeightFunc, b *testing.B) {
	reps, _ := u.ConnectedComponentReps()
	f := graph.NewFromList(u.Order())
	lab := make([]graph.LI, u.Order())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range reps {
			u.Prim(r, w, &f, lab, nil)
			f.TransposeLabeled(lab).Undirected()
		}
	}
}
