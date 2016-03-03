// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleDirected_Undirected() {
	// arcs directed down:
	//    0
	//   / \
	//  1   2
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {},
		2: {},
	}}
	u := g.Undirected()
	for fr, to := range u.AdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [1 2]
	// 1 [0]
	// 2 [0]
}

func ExampleDirected_Undirected_loopMultigraph() {
	//  0--\   /->1--\
	//  |  |   |  ^  |
	//  \--/   |  |  |
	//         \--2<-/
	g := graph.Directed{graph.AdjacencyList{
		0: {0},
		1: {2},
		2: {1, 1},
	}}
	u := g.Undirected()
	for fr, to := range u.AdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [0]
	// 1 [2 2]
	// 2 [1 1]
}

func ExampleAdjacencyList_IsUndirected() {
	// 0<--    2<--\
	//  \  \   |   |
	//   -->1  \---/
	g := graph.AdjacencyList{
		0: {1},
		1: {0},
		2: {2},
	}
	ud, _, _ := g.IsUndirected()
	fmt.Println(ud)
	// 0<--
	//  \  \
	//   -->1<--2
	g = graph.AdjacencyList{
		0: {1},
		1: {0},
		2: {1},
	}
	fmt.Println(g.IsUndirected())
	// Output:
	// true
	// false 2 1
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
