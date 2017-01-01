// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// adj_cg_test.go -- tests for code in adj_cg.go.
//
// These are tests on the labeled versions of methods.
//
// See also adj_ro_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleLabeledAdjacencyList_ArcDensity() {
	// 0-->1
	// |
	// v
	// 2-->3
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		2: {{To: 3}},
		3: {},
	}
	fmt.Println(g.ArcDensity())
	// Output:
	// 0.25
}

func ExampleLabeledAdjacencyList_ArcSize() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{ // simple graph
		2: {{To: 0}, {To: 1}},
	}}
	fmt.Println(g.ArcSize())

	// with reciprocals now.
	u := g.Undirected()
	// common term "size" for undirected graph is number of undirected edges.
	// size, m = ArcSize() / 2 here, but only because there are no loops.
	fmt.Println(u.ArcSize())

	g.LabeledAdjacencyList[1] = []graph.Half{{To: 1}} // add a loop
	//   2
	//  / \
	// 0   1---\
	//      \--/
	fmt.Println(g.ArcSize())

	// loops have no reciprocals.  ArcSize() / 2 no longer meaningful.
	fmt.Println(g.Undirected().ArcSize())
	// Output:
	// 2
	// 4
	// 3
	// 5
}

func ExampleLabeledAdjacencyList_BoundsOk() {
	var g graph.LabeledAdjacencyList
	ok, _, _ := g.BoundsOk() // zero value adjacency list is valid
	fmt.Println(ok)
	g = graph.LabeledAdjacencyList{
		0: {{To: 9}},
	}
	fmt.Println(g.BoundsOk()) // arc 0 to 9 invalid with only one node
	// Output:
	// true
	// false 0 {9 0}
}

func ExampleLabeledAdjacencyList_BreadthFirstPath() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.LabeledAdjacencyList{
		2: {{To: 1}},
		1: {{To: 4}},
		4: {{To: 3}, {To: 6}},
		3: {{To: 5}},
		6: {{To: 5}, {To: 6}},
	}
	fmt.Println(g.BreadthFirstPath(1, 3))
	// Output:
	// [1 4 3]
}

func ExampleLabeledAdjacencyList_BreadthFirst_singlePath() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.LabeledAdjacencyList{
		2: {{To: 1}},
		1: {{To: 4}},
		4: {{To: 3}, {To: 6}},
		3: {{To: 5}},
		6: {{To: 5}, {To: 6}},
	}
	var start, end graph.NI = 1, 6
	var f graph.FromList
	n, _ := g.BreadthFirst(start, nil, &f, func(n graph.NI) bool {
		return n != end
	})
	fmt.Println(n, "nodes visited")
	fmt.Println("path:", f.PathTo(end, nil))
	// Output:
	// 4 nodes visited
	// path: [1 4 6]
}

func ExampleLabeledAdjacencyList_BreadthFirst_allPaths() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.LabeledAdjacencyList{
		2: {{To: 1}},
		1: {{To: 4}},
		4: {{To: 3}, {To: 6}},
		3: {{To: 5}},
		6: {{To: 5}, {To: 6}},
	}
	start := graph.NI(1)
	var f graph.FromList
	g.BreadthFirst(start, nil, &f, func(n graph.NI) bool {
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

func ExampleLabeledAdjacencyList_BreadthFirst_traverse() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}, {To: 3}},
		2: {{To: 4}, {To: 5}, {To: 6}},
		3: {{To: 7}, {To: 8}},
		8: {},
	}
	var f graph.FromList
	g.BreadthFirst(0, nil, &f, func(n graph.NI) bool {
		fmt.Println("visit", n, "level", f.Paths[n].Len)
		return true
	})
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

func ExampleLabeledAdjacencyList_BreadthFirst_traverseRandom() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}, {To: 3}},
		2: {{To: 4}, {To: 5}, {To: 6}},
		3: {{To: 7}, {To: 8}},
		8: {},
	}

	// only difference from non-random example
	r := rand.New(rand.NewSource(8))

	var f graph.FromList
	g.BreadthFirst(0, r, &f, func(n graph.NI) bool {
		fmt.Println("visit", n, "level", f.Paths[n].Len)
		return true
	})
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

func ExampleLabeledAdjacencyList_DepthFirst() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		1: {{To: 2}},
		2: {{To: 3}},
		3: {{To: 1}},
	}
	ok := g.DepthFirst(0, nil, func(n graph.NI) (ok bool) {
		fmt.Println("visit", n)
		return true
	})
	fmt.Println(ok)
	// Output:
	// visit 0
	// visit 1
	// visit 2
	// visit 3
	// true
}

func ExampleLabeledAdjacencyList_DepthFirst_earlyTermination() {
	//   0-->3
	//  / \
	// 1-->2
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}, {To: 3}},
		1: {{To: 2}},
		3: {},
	}}
	ok := g.DepthFirst(0, nil, func(n graph.NI) bool {
		fmt.Println("visit", n)
		return n != 2
	})
	fmt.Println(ok)
	// Output:
	// visit 0
	// visit 1
	// visit 2
	// false
}

func ExampleLabeledAdjacencyList_DepthFirst_bitmap() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		1: {{To: 2}},
		2: {{To: 3}},
		3: {{To: 1}},
	}
	var vis graph.Bits
	fmt.Println("3210")
	fmt.Println("----")
	g.DepthFirst(0, &vis, func(graph.NI) bool {
		fmt.Printf("%04b\n", &vis)
		return true
	})
	// Output:
	// 3210
	// ----
	// 0001
	// 0011
	// 0111
	// 1111
}

func TestLabeledAdjacencyList_DepthFirst_bothNil(t *testing.T) {
	// for coverage
	var g graph.LabeledAdjacencyList
	if g.DepthFirst(0, nil, nil) {
		t.Fatal("DepthFirst both nil must return false")
	}
}

func ExampleLabeledAdjacencyList_DepthFirstRandom() {
	//     ----0-----
	//    /    |     \
	//   1     2      3
	//  /|\   /|\   / | \
	// 4 5 6 7 8 9 10 11 12
	g := graph.LabeledAdjacencyList{
		0:  {{To: 1}, {To: 2}, {To: 3}},
		1:  {{To: 4}, {To: 5}, {To: 6}},
		2:  {{To: 7}, {To: 8}, {To: 9}},
		3:  {{To: 10}, {To: 11}, {To: 12}},
		12: {},
	}
	r := rand.New(rand.NewSource(12))
	f := func(n graph.NI) (ok bool) {
		fmt.Println("visit", n)
		return true
	}
	g.DepthFirstRandom(0, nil, f, r)
	// Output:
	// visit 0
	// visit 1
	// visit 6
	// visit 4
	// visit 5
	// visit 3
	// visit 12
	// visit 11
	// visit 10
	// visit 2
	// visit 9
	// visit 7
	// visit 8
}

func ExampleLabeledAdjacencyList_HasArc() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 2}, {To: 0}, {To: 1}, {To: 1}},
	}
	fmt.Println(g.HasArc(2, 1))
	fmt.Println(g.HasArc(2, 2)) // test for loop
	// Output:
	// true 3
	// true 1
}

func ExampleLabeledAdjacencyList_AnyLoop_loop() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 2}},
	}
	fmt.Println(g.AnyLoop())
	// Output:
	// true 2
}

func ExampleLabeledAdjacencyList_AnyLoop_noLoop() {
	g := graph.LabeledAdjacencyList{
		1: {{To: 0}},
	}
	lp, _ := g.AnyLoop()
	fmt.Println("has loop:", lp)
	// Output:
	// has loop: false
}

func ExampleLabeledAdjacencyList_AnyParallelMap_parallelArcs() {
	g := graph.LabeledAdjacencyList{
		1: {{To: 0}, {To: 0}},
	}
	// result true 1 0 means parallel arcs from node 1 to node 0
	fmt.Println(g.AnyParallelMap())
	// Output:
	// true 1 0
}

func ExampleLabeledAdjacencyList_AnyParallelMap_noParallelArcs() {
	g := graph.LabeledAdjacencyList{
		1: {{To: 0}},
	}
	fmt.Println(g.AnyParallelMap()) // result false -1 -1 means no parallel arc
	// Output:
	// false -1 -1
}

func ExampleLabeledAdjacencyList_IsolatedNodes() {
	//   0  1
	//  / \
	// 2   3  4
	g := graph.LabeledAdjacencyList{
		0: {{To: 2}, {To: 3}},
		4: {},
	}
	fmt.Println(g.IsolatedNodes().Slice())
	// Output:
	// [1 4]
}

func ExampleLabeledAdjacencyList_IsSimple() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 1}},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// true -1
}

func ExampleLabeledAdjacencyList_IsSimple_loop() {
	// arcs directed down
	//   2
	//  / \
	// 0   1---\
	//      \--/
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 1}},
		1: {{To: 1}}, // loop
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 1
}

func ExampleLabeledAdjacencyList_IsSimple_parallelArc() {
	// arcs directed down
	//   2
	//  /|\
	//  |/ \
	//  0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 1}, {To: 0}},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 2
}

func ExampleLabeledAdjacencyList_Order() {
	g := graph.LabeledAdjacencyList{ // maybe you think of node 0 as "unused",
		1: {{To: 2}},
		2: {{To: 1}},
	}
	fmt.Println(g.Order())                    // but graph still has 3 nodes.
	g = make(graph.LabeledAdjacencyList, 300) // empty graph,
	fmt.Println(g.Order())                    // with 300 nodes.
	fmt.Println(len(g))                       // equivalent.

	var u graph.LabeledUndirected
	u.AddEdge(graph.Edge{0, 99}, 0)
	fmt.Println(len(u.LabeledAdjacencyList)) // explicit.
	fmt.Println(u.Order())                   // a little more concise.
	// Output:
	// 3
	// 300
	// 300
	// 100
	// 100
}

func ExampleLabeledAdjacencyList_ParallelArcs() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 0}, {To: 2}, {To: 0}, {To: 1}, {To: 1}},
	}
	fmt.Println(g.ParallelArcs(0, 2))
	fmt.Println(g.ParallelArcs(2, 0))
	fmt.Println(g.ParallelArcs(2, 1))
	fmt.Println(g.ParallelArcs(2, 2)) // returns loops on 2
	// Output:
	// []
	// [0 2]
	// [3 4]
	// [1]
}
