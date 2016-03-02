package graph_test

import (
	"fmt"
	"math/big"

	"github.com/soniakeys/graph"
)

func ExampleWeightedEdgeList_Kruskal() {
	//       (10)
	//     0------4----\
	//     |     /|     \(70)
	// (30)| (40) |(60)  \
	//     |/     |      |
	//     1------2------3
	//       (50)   (20)
	w := func(l graph.LI) float64 { return float64(l) }
	// undirected graph
	var g graph.UndirectedLabeled
	g.AddEdge(graph.Edge{0, 1}, 30)
	g.AddEdge(graph.Edge{0, 4}, 10)
	g.AddEdge(graph.Edge{1, 2}, 50)
	g.AddEdge(graph.Edge{1, 4}, 40)
	g.AddEdge(graph.Edge{2, 3}, 20)
	g.AddEdge(graph.Edge{2, 4}, 60)
	g.AddEdge(graph.Edge{3, 4}, 70)
	// convert to edge list for Kruskal, but no need to sort it.
	// Kruskal will sort it.
	l := g.WeightedEdgeList(w)

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
	// constructing an edge list.  No need for reciprocal edges.  Also if
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

func ExamplePrim_Span() {
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
	var g graph.UndirectedLabeled
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

	// construct prim object
	p := graph.NewPrim(g, w)

	// construct spanning tree for each component
	fmt.Println("Span results:")
	fmt.Println("Root  Nodes spanned  Total tree distance  Leaves")
	for _, r := range reps {
		var leaves big.Int
		ns, dist := p.Span(r, &leaves)
		// collect leaf node ints from bitmap
		var ll []int
		for n := range g.LabeledAdjacencyList {
			if leaves.Bit(n) == 1 {
				ll = append(ll, n)
			}
		}
		fmt.Printf("%d %17d %20.0f  %d\n", r, ns, dist, ll)
	}

	// show final forest
	fmt.Println("Spanning forest:")
	fmt.Println("Node  From  Arc distance  Path length  Leaf")
	for n, pe := range p.Forest.Paths {
		fmt.Printf("%d %8d %13.0f %12d %5d\n",
			n, pe.From, w(p.Labels[n]), pe.Len, p.Forest.Leaves.Bit(n))
	}

	// optionally, convert to undirected graph
	u := p.Forest.TransposeLabeled(p.Labels).Undirected()
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

/* invalid.  needs to work on undirected and r100 is directed
func TestPrim100(t *testing.T) {
	reps, orders := r100.g.ConnectedComponentReps()
	p := graph.NewPrim(r100.l, func(l graph.LI) float64 { return r100.w[l] })

	// construct spanning tree for each component
	for i, r := range reps {
		ns, _ := p.Span(r, nil)
		if ns != orders[i] {
			t.Fatal("Not all nodes spanned within a connected component.")
		}
	}
}

func BenchmarkPrim100(b *testing.B) {
	reps, _ := r100.g.ConnectedComponentReps()
	p := graph.NewPrim(r100.l, func(l graph.LI) float64 { return r100.w[l] })
	for i := 0; i < b.N; i++ {
		p.Reset()
		for _, r := range reps {
			p.Span(r, nil)
		}
	}
}
*/
