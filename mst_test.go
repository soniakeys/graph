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
	g := graph.LabeledAdjacencyList{
		0: {{1, 3}, {2, 5}},
		1: {{0, 3}, {2, 4}},
		2: {{0, 5}, {1, 4}},
		3: {{4, 2}},
		4: {{3, 2}},
	}
	w := func(arcLabel int) float64 { return float64(arcLabel) }

	// construct unweighted equivalent for some demonstrations
	ul := g.Unlabeled()

	// demonstration 1:  show that the graph is undirected
	ud, _, _ := ul.Undirected()
	fmt.Println("Undirected:", ud)

	// demonstration 2:  show connected components
	rep, nNodes := ul.ConnectedComponents()
	fmt.Println("Connected components:")
	fmt.Println("representative node - number of nodes in component")
	for i, r := range rep {
		fmt.Printf("%d %21d\n", r, nNodes[i])
	}

	// construct prim object on the original weighted graph
	p := graph.NewPrim(g, w)

	// construct spanning tree for each component
	for _, r := range rep {
		ns := p.Span(r)
		fmt.Printf("From node %d, %d nodes spanned.\n", r, ns)
	}

	fmt.Println("Spanning Forest:")
	fmt.Println("Node  From  Weight")
	for n, pe := range p.Tree.Paths {
		fmt.Printf("%d %8d %7.1f\n",
			n, pe.From, p.Wt[n])
	}
	// Output:
	// Undirected: true
	// Connected components:
	// representative node - number of nodes in component
	// 0                     3
	// 3                     2
	// From node 0, 3 nodes spanned.
	// From node 3, 2 nodes spanned.
	// Spanning Forest:
	// Node  From  Weight
	// 0       -1     0.0
	// 1        0     3.0
	// 2        1     4.0
	// 3       -1     0.0
	// 4        3     2.0
}
