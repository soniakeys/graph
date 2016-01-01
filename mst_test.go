package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

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
	g := &graph.LabeledAdjacencyList{}
	g.AddEdge(graph.Edge{0, 1}, 3)
	g.AddEdge(graph.Edge{1, 2}, 4)
	g.AddEdge(graph.Edge{2, 0}, 5)
	g.AddEdge(graph.Edge{3, 4}, 2)
	w := func(arcLabel int) float64 { return float64(arcLabel) }

	// get connected components
	reps, orders := g.ConnectedComponentReps()
	fmt.Println("Connected components:")
	fmt.Println("representative node - order (number of nodes) in component")
	for i, r := range reps {
		fmt.Printf("%d %21d\n", r, orders[i])
	}

	// construct prim object
	p := graph.NewPrim(*g, w)

	// construct spanning tree for each component
	for _, r := range reps {
		ns := p.Span(r)
		fmt.Printf("From node %d, %d nodes spanned.\n", r, ns)
	}

	fmt.Println("Spanning Forest:")
	fmt.Println("Node  From")
	for n, pe := range p.Tree.Paths {
		fmt.Printf("%d %8d\n", n, pe.From)
	}
	// Output:
	// Connected components:
	// representative node - order (number of nodes) in component
	// 0                     3
	// 3                     2
	// From node 0, 3 nodes spanned.
	// From node 3, 2 nodes spanned.
	// Spanning Forest:
	// Node  From
	// 0       -1
	// 1        0
	// 2        1
	// 3       -1
	// 4        3
}
