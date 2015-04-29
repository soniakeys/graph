// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_ConnectedComponents() {
	g := graph.AdjacencyList{
		0: {3, 4},
		1: {5},
		3: {0, 4},
		4: {0, 3},
		5: {1},
	}
	fmt.Println(g.ConnectedComponents())
	// Output:
	// [0 1 2] [3 2 1]
}

func ExampleTarjanBiconnectedComponents() {
	g := graph.AdjacencyList{}
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
	for _, bcc := range g.TarjanBiconnectedComponents() {
		fmt.Println("Edges:")
		for _, e := range bcc {
			fmt.Println(e)
		}
	}
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

func ExampleAdjacencyList_Degeneracy() {
	//
	//   /---2
	//  0--./|\
	//  |\ /\| \
	//  | .  6  3--5
	//  |/ \/| /
	//  1--'\|/
	//   \---4
	//
	// Same graph redrawn to show ordering:
	//
	//          /-----\
	//         /  /--\ \
	//      /-+--+-\  \|
	//  5--3--4--6--2--1--0
	//         \  \  \---/|\
	//          \  \-----/ |
	//           \--------/
	//
	g := graph.AdjacencyList{
		0: {1, 2, 4, 6},
		1: {0, 2, 4, 6},
		2: {0, 1, 3, 6},
		3: {2, 4, 5},
		4: {0, 1, 3, 6},
		5: {3},
		6: {0, 1, 2, 4},
	}
	k, ord, cores := g.Degeneracy()
	fmt.Println("Degeneracy:", k)
	fmt.Println("Ordering:", ord)
	fmt.Println("0-core:", ord[:cores[0]])
	for k := 1; k < len(cores); k++ {
		fmt.Printf("%d-core: %d\n", k, ord[cores[k-1]:cores[k]])
	}
	// Output:
	// Degeneracy: 3
	// Ordering: [5 3 4 6 2 1 0]
	// 0-core: []
	// 1-core: [5]
	// 2-core: [3]
	// 3-core: [4 6 2 1 0]
}
