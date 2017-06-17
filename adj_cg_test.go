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
	"reflect"
	"sort"
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

func ExampleLabeledAdjacencyList_Permute() {
	//    0                              2
	//  x/ \y  Permute([2 0 1]) gives  x/ \y
	//  1-->2                          0-->1
	//    z                              z
	g := graph.LabeledAdjacencyList{
		0: {{1, 'x'}, {2, 'y'}},
		1: {{2, 'z'}},
		2: {},
	}
	g.Permute([]int{2, 0, 1})
	for fr, to := range g {
		fmt.Print(fr, ":")
		for _, to := range to {
			fmt.Printf(" (%d %c)", to.To, to.Label)
		}
		fmt.Println()
	}
	// Output:
	// 0: (1 z)
	// 1:
	// 2: (0 x) (1 y)
}

// not much of a test.  doesn't actually test that shuffle did anything,
// does a bunch of unrelated stuff.  It at least tests that shuffle doesn't
// corrupt the graph.
func TestShuffleArcListsLabeled(t *testing.T) {
	testCase := func(rad float64, r *rand.Rand) {
		g, _, _ := graph.LabeledGeometric(10, rad, r)
		c, _ := g.LabeledAdjacencyList.Copy()
		c.ShuffleArcLists(r)
		for fr, to := range g.LabeledAdjacencyList {
			sort.Slice(to, func(i, j int) bool {
				return to[i].To < to[j].To ||
					to[i].To == to[j].To && to[i].Label < to[j].Label
			})
			sh := c[fr]
			sort.Slice(sh, func(i, j int) bool {
				return sh[i].To < sh[j].To ||
					sh[i].To == sh[j].To && sh[i].Label < sh[j].Label
			})
			if (len(to) > 0 || len(sh) > 0) && !reflect.DeepEqual(to, sh) {
				t.Fatal(rad, r, fr, to, sh)
			}
		}
	}
	testCase(.2, nil)
	testCase(.6, rand.New(rand.NewSource(3)))
}
