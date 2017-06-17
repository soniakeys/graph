// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// adj_ro_test.go -- tests on adj_RO.go
//
// These are tests on the code-generated unlabeled versions of methods.
//
// If testing prompts changes in the tested method, be sure to edit
// adj_cg.go, go generate to generate the adj_RO.go, and then retest.
// Do not edit adj_RO.go
//
// See also adj_cg_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_ArcDensity() {
	// 0-->1
	// |
	// v
	// 2-->3
	g := graph.AdjacencyList{
		0: {1, 2},
		2: {3},
		3: {},
	}
	fmt.Println(g.ArcDensity())
	// Output:
	// 0.25
}

func ExampleAdjacencyList_ArcSize() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.Directed{graph.AdjacencyList{ // simple graph
		2: {0, 1},
	}}
	fmt.Println(g.ArcSize())

	// with reciprocals now.
	u := g.Undirected()
	// common term "size" for undirected graph is number of undirected edges.
	// size, m = ArcSize() / 2 here, but only because there are no loops.
	fmt.Println(u.ArcSize())

	g.AdjacencyList[1] = []graph.NI{1} // add a loop
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

func ExampleAdjacencyList_BoundsOk() {
	var g graph.AdjacencyList
	ok, _, _ := g.BoundsOk() // zero value adjacency list is valid
	fmt.Println(ok)
	g = graph.AdjacencyList{
		0: {9},
	}
	fmt.Println(g.BoundsOk()) // arc 0 to 9 invalid with only one node
	// Output:
	// true
	// false 0 9
}

func ExampleAdjacencyList_HasArc() {
	g := graph.AdjacencyList{
		2: {0, 2, 0, 1, 1},
	}
	fmt.Println(g.HasArc(2, 1))
	fmt.Println(g.HasArc(2, 2)) // test for loop
	// Output:
	// true 3
	// true 1
}

func ExampleAdjacencyList_AnyLoop_loop() {
	g := graph.AdjacencyList{
		2: {2},
	}
	fmt.Println(g.AnyLoop())
	// Output:
	// true 2
}

func ExampleAdjacencyList_AnyLoop_noLoop() {
	g := graph.AdjacencyList{
		1: {0},
	}
	lp, _ := g.AnyLoop()
	fmt.Println("has loop:", lp)
	// Output:
	// has loop: false
}

func ExampleAdjacencyList_AnyParallelMap_parallelArcs() {
	g := graph.AdjacencyList{
		1: {0, 0},
	}
	// result true 1 0 means parallel arcs from node 1 to node 0
	fmt.Println(g.AnyParallelMap())
	// Output:
	// true 1 0
}

func ExampleAdjacencyList_AnyParallelMap_noParallelArcs() {
	g := graph.AdjacencyList{
		1: {0},
	}
	fmt.Println(g.AnyParallelMap()) // result false -1 -1 means no parallel arc
	// Output:
	// false -1 -1
}

func ExampleAdjacencyList_IsolatedNodes() {
	//   0  1
	//  / \
	// 2   3  4
	g := graph.AdjacencyList{
		0: {2, 3},
		4: {},
	}
	fmt.Println(g.IsolatedNodes().Slice())
	// Output:
	// [1 4]
}

func ExampleAdjacencyList_IsSimple() {
	// arcs directed down
	//   2
	//  / \
	// 0   1
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// true -1
}

func ExampleAdjacencyList_IsSimple_loop() {
	// arcs directed down
	//   2
	//  / \
	// 0   1---\
	//      \--/
	g := graph.AdjacencyList{
		2: {0, 1},
		1: {1}, // loop
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 1
}

func ExampleAdjacencyList_IsSimple_parallelArc() {
	// arcs directed down
	//   2
	//  /|\
	//  |/ \
	//  0   1
	g := graph.AdjacencyList{
		2: {0, 1, 0},
	}
	fmt.Println(g.IsSimple())
	// Output:
	// false 2
}

func ExampleAdjacencyList_Order() {
	g := graph.AdjacencyList{ // maybe you think of node 0 as "unused",
		1: {2},
		2: {1},
	}
	fmt.Println(g.Order())             // but graph still has 3 nodes.
	g = make(graph.AdjacencyList, 300) // empty graph,
	fmt.Println(g.Order())             // with 300 nodes.
	fmt.Println(len(g))                // equivalent.

	var u graph.Undirected
	u.AddEdge(0, 99)
	fmt.Println(len(u.AdjacencyList)) // explicit.
	fmt.Println(u.Order())            // a little more concise.
	// Output:
	// 3
	// 300
	// 300
	// 100
	// 100
}

func ExampleAdjacencyList_ParallelArcs() {
	g := graph.AdjacencyList{
		2: {0, 2, 0, 1, 1},
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

func ExampleAdjacencyList_Permute() {
	//    0                             2
	//   / \   Permute([2 0 1]) gives  / \
	//  1-->2                         0-->1
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {},
	}
	g.Permute([]int{2, 0, 1})
	for fr, to := range g {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [1]
	// 1 []
	// 2 [0 1]
}

// not much of a test.  doesn't actually test that shuffle did anything,
// does a bunch of unrelated stuff.  It at least tests that shuffle doesn't
// corrupt the graph.
func TestShuffleArcLists(t *testing.T) {
	testCase := func(p float64, r *rand.Rand) {
		g, _ := graph.GnpUndirected(10, p, r)
		c, _ := g.AdjacencyList.Copy()
		c.ShuffleArcLists(r)
		for fr, to := range g.AdjacencyList {
			sort.Slice(to, func(i, j int) bool { return to[i] < to[j] })
			sh := c[fr]
			sort.Slice(sh, func(i, j int) bool { return sh[i] < sh[j] })
			if (len(to) > 0 || len(sh) > 0) && !reflect.DeepEqual(to, sh) {
				t.Fatal(p, r, fr, to, sh)
			}
		}
	}
	testCase(.1, nil)
	testCase(.9, rand.New(rand.NewSource(3)))
}
