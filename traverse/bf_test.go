// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package traverse_test

import (
	"fmt"
	"math/rand"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/traverse"
)

func ExampleBreadthFirst_singlePath() {
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
	var start, end graph.NI = 1, 6
	var f graph.FromList
	var visited int
	traverse.BreadthFirst(g, start, traverse.FromList(&f),
		traverse.OkNodeVisitor(func(n graph.NI) bool {
			visited++
			return n != end
		}))
	fmt.Println(visited, "nodes visited")
	fmt.Println("path:", f.PathTo(end, nil))
	// Output:
	// 4 nodes visited
	// path: [1 4 6]
}

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
	start := graph.NI(1)
	var f graph.FromList
	traverse.BreadthFirst(g, start, traverse.FromList(&f))
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

func ExampleBreadthFirst_traverse() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	g := graph.AdjacencyList{
		0: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
		8: {},
	}
	var f graph.FromList
	traverse.BreadthFirst(g, 0, traverse.FromList(&f),
		traverse.NodeVisitor(func(n graph.NI) {
			fmt.Println("visit", n, "level", f.Paths[n].Len)
		}))
	// Output:
	// visit 0 level 1
	// visit 1 level 2
	// visit 2 level 2
	// visit 3 level 2
	// visit 4 level 3
	// visit 5 level 3
	// visit 6 level 3
	// visit 7 level 3
	// visit 8 level 3
}

func ExampleBreadthFirst_traverseRandom() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	g := graph.AdjacencyList{
		0: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
		8: {},
	}

	// only difference from non-random example
	r := rand.New(rand.NewSource(8))

	var f graph.FromList
	traverse.BreadthFirst(g, 0, traverse.Rand(r), traverse.FromList(&f),
		traverse.NodeVisitor(func(n graph.NI) {
			fmt.Println("visit", n, "level", f.Paths[n].Len)
		}))
	// Output:
	// visit 0 level 1
	// visit 1 level 2
	// visit 3 level 2
	// visit 2 level 2
	// visit 8 level 3
	// visit 5 level 3
	// visit 6 level 3
	// visit 4 level 3
	// visit 7 level 3
}
