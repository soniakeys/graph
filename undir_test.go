// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleArcDensity() {
	fmt.Println(graph.ArcDensity(4, 3))
	// Output:
	// 0.25
}

func ExampleDensity() {
	fmt.Println(graph.Density(4, 3))
	// Output:
	// 0.5
}

func ExampleUndirected_Edges() {
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 0) // parallel
	g.AddEdge(0, 0) // loop
	g.Edges(func(e graph.Edge) {
		fmt.Println(e)
	})
	// Output:
	// {0 0}
	// {1 0}
	// {2 1}
	// {2 0}
	// {2 0}
}

func ExampleUndirected_EulerianCycleD() {
	var g graph.Undirected
	// add 6 edges
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 2) // loop
	m := g.Size()
	fmt.Println("m =", m)
	fmt.Println(g.EulerianCycleD(m))
	// Output:
	// m = 6
	// [0 1 2 2 1 2 0] <nil>
}

func ExampleUndirected_SimpleEdges() {
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(2, 0) // parallel
	g.AddEdge(0, 0) // loop
	g.SimpleEdges(func(e graph.Edge) {
		fmt.Println(e)
	})
	// Output:
	// {0 1}
	// {0 2}
	// {1 2}
}

func ExampleUndirected_TarjanBiconnectedComponents() {
	// undirected edges:
	// 3---2---1---7---9
	//  \ / \ / \   \ /
	//   4   5---6   8
	var g graph.Undirected
	g.AddEdge(3, 4)
	g.AddEdge(3, 2)
	g.AddEdge(2, 4)
	g.AddEdge(2, 5)
	g.AddEdge(2, 1)
	g.AddEdge(5, 1)
	g.AddEdge(6, 1)
	g.AddEdge(6, 5)
	g.AddEdge(7, 1)
	g.AddEdge(7, 9)
	g.AddEdge(7, 8)
	g.AddEdge(9, 8)
	g.TarjanBiconnectedComponents(func(bcc []graph.Edge) bool {
		fmt.Println("Edges:")
		for _, e := range bcc {
			fmt.Println(e)
		}
		return true
	})
	// Output:
	// Edges:
	// {4 2}
	// {3 4}
	// {2 3}
	// Edges:
	// {6 1}
	// {5 6}
	// {5 1}
	// {2 5}
	// {1 2}
	// Edges:
	// {8 7}
	// {9 8}
	// {7 9}
	// Edges:
	// {1 7}
}

func ExampleLabeledUndirected_Edges() {
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 'A')
	g.AddEdge(graph.Edge{0, 1}, 'B')
	g.AddEdge(graph.Edge{1, 2}, 'C')
	g.AddEdge(graph.Edge{1, 2}, 'D')
	g.Edges(func(e graph.LabeledEdge) {
		fmt.Printf("%d %c\n", e.Edge, e.LI)
	})
	// Output:
	// {0 0} A
	// {1 0} B
	// {2 1} C
	// {2 1} D
}

func ExampleLabeledUndirected_TarjanBiconnectedComponents() {
	// undirected edges:
	// 3---2---1---7---9
	//  \ / \ / \   \ /
	//   4   5---6   8
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{2, 5}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{6, 1}, 0)
	g.AddEdge(graph.Edge{6, 5}, 0)
	g.AddEdge(graph.Edge{7, 1}, 0)
	g.AddEdge(graph.Edge{7, 9}, 0)
	g.AddEdge(graph.Edge{7, 8}, 0)
	g.AddEdge(graph.Edge{9, 8}, 0)
	g.TarjanBiconnectedComponents(func(bcc []graph.LabeledEdge) bool {
		fmt.Println("Edges:")
		for _, e := range bcc {
			fmt.Println(e.Edge)
		}
		return true
	})
	// Output:
	// Edges:
	// {4 2}
	// {3 4}
	// {2 3}
	// Edges:
	// {6 1}
	// {5 6}
	// {5 1}
	// {2 5}
	// {1 2}
	// Edges:
	// {8 7}
	// {9 8}
	// {7 9}
	// Edges:
	// {1 7}
}

/* shelved
func ExampleBiconnectedComponents_Find() {
	g := graph.AdjacencyList{
		0:  {1, 7},
		1:  {2, 4, 0},
		2:  {3, 1},
		3:  {2, 4},
		4:  {3, 1},
		5:  {6, 12},
		6:  {5, 12, 8},
		7:  {8, 0},
		8:  {6, 7, 9, 10},
		9:  {8},
		10: {8, 13, 11},
		11: {10},
		12: {5, 6, 13},
		13: {12, 10},
	}
	b := graph.NewBiconnectedComponents(g)
	b.Find(0)
	fmt.Println("n: cut from")
	for n, f := range b.From {
		fmt.Printf("%d: %d %d\n",
			n, b.Cuts.Bit(n), f)
	}
	fmt.Println("Leaves:", b.Leaves)
	// Output:
	// n: cut from
	// 0: 1 -1
	// 1: 1 0
	// 2: 0 1
	// 3: 0 2
	// 4: 0 3
	// 5: 0 6
	// 6: 0 8
	// 7: 1 0
	// 8: 1 7
	// 9: 0 8
	// 10: 1 13
	// 11: 0 10
	// 12: 0 5
	// 13: 0 12
	// Leaves: [4 11 9]
}
*/
