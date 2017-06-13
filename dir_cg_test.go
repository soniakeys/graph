// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// dir_cg_test.go -- tests for code in dir_cg.go.
//
// These are tests on the labeled versions of methods.
//
// See also dir_ro_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"math/big"

	"github.com/soniakeys/graph"
)

func ExampleLabeledDirected_Balanced() {
	// 2
	// |
	// v
	// 0----->1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0}},
		0: {{To: 1}},
	}}
	fmt.Println(g.Balanced())

	// 0<--\
	// |    \
	// v     \
	// 1----->2
	g.LabeledAdjacencyList[1] = []graph.Half{{To: 2}}
	fmt.Println(g.Balanced())
	// Output:
	// false
	// true
}

func ExampleLabeledDirected_Cyclic() {
	//   0
	//  / \
	// 1-->2-->3
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		1: {{To: 2}},
		2: {{To: 3}},
		3: {},
	}}
	cyclic, _, _ := g.Cyclic()
	fmt.Println(cyclic)

	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g.LabeledAdjacencyList[3] = []graph.Half{{To: 1}}
	fmt.Println(g.Cyclic())

	// Output:
	// false
	// true 3 {1 0}
}

func ExampleLabeledDirected_Dominators() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 2}, {To: 3}},
		2: {{To: 4}},
		3: {{To: 4}, {To: 5}},
		6: {{To: 5}},
	}}
	d := g.Dominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleLabeledDirected_Doms() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 2}, {To: 3}},
		2: {{To: 4}},
		3: {{To: 4}, {To: 5}},
		6: {{To: 5}},
	}}
	// compute postorder with depth-first traversal
	var post []graph.NI
	var vis big.Int
	var f func(graph.NI)
	f = func(n graph.NI) {
		vis.SetBit(&vis, int(n), 1)
		for _, to := range g.LabeledAdjacencyList[n] {
			if vis.Bit(int(to.To)) == 0 {
				f(to.To)
			}
		}
		post = append(post, n)
	}
	f(0)
	fmt.Println("post:", post)
	tr, _ := g.Transpose()
	d := g.Doms(tr, post)
	fmt.Println("doms:", d.Immediate)
	// Output:
	// post: [4 2 5 3 1 0]
	// doms: [0 0 1 1 1 3 -1]
}

func ExampleLabeledDirected_FromList() {
	//    4   3
	//   / \
	//  2   1
	//       \
	//        0
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		4: {{To: 2}, {To: 1}},
		1: {{To: 0}},
	}}
	f, n := g.FromList()
	fmt.Println("n:", n)
	fmt.Println("N  From")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d\n", n, e.From)
	}
	// Output:
	// n: -1
	// N  From
	// 0    1
	// 1    4
	// 2    4
	// 3   -1
	// 4   -1
}

func ExampleLabeledDirected_FromList_nonTree() {
	//    0
	//   / \
	//  1   2
	//   \ /
	//    3
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}},
		1: {{To: 3}},
		2: {{To: 3}},
		3: {},
	}}
	fmt.Println(g.FromList())
	// Output:
	// <nil> 3
}

func ExampleLabeledDirected_FromList_multigraphTree() {
	//    0
	//   / \\
	//  1   2
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}, {To: 2}, {To: 2}},
		2: {},
	}}
	fmt.Println(g.FromList())
	// Output:
	// <nil> 2
}

func ExampleLabeledDirected_FromList_rootLoops() {
	//     /-\
	//    0--/  3--\
	//   / \     \-/
	//  1   2
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 0}, {To: 1}, {To: 2}},
		3: {{To: 3}},
	}}
	f, n := g.FromList()
	fmt.Println("n:", n)
	fmt.Println("N  From")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d\n", n, e.From)
	}
	// Output:
	// n: -1
	// N  From
	// 0    0
	// 1    0
	// 2    0
	// 3    3
}

func ExampleLabeledDirected_InDegree() {
	// arcs directed down:
	//  0     2
	//  |
	//  1
	//  |\
	//  | \
	//  3  4<-\
	//     \--/
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 1}},
		1: {{To: 3}, {To: 4}},
		4: {{To: 4}},
	}}
	fmt.Println("node:    0 1 2 3 4")
	fmt.Println("in-deg:", g.InDegree())
	// Output:
	// node:    0 1 2 3 4
	// in-deg: [0 1 0 1 2]
}

func ExampleLabeledDirected_IsTree() {
	// Example graph
	// Arcs point down unless otherwise indicated
	//           1
	//          / \
	//         0   5
	//        /   / \
	//       2   3-->4
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{To: 0}, {To: 5}},
		0: {{To: 2}},
		5: {{To: 3}, {To: 4}},
		3: {{To: 4}},
	}}
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	// Output:
	// true false
	// false false
}

func ExampleLabeledDirected_PostDominators() {
	// Example graph here is transpose of that in the Dominators example
	// to show result is the same.
	//   4   5
	//  / \ /|
	// 2   3 |
	//  \ /  |
	//   1   |
	//   |   |
	//   0   6
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		4: {{To: 2}, {To: 3}},
		5: {{To: 3}, {To: 6}},
		2: {{To: 1}},
		3: {{To: 1}},
		1: {{To: 0}},
		6: {},
	}}
	d := g.PostDominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleLabeledDirected_StronglyConnectedComponents() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 0}, {To: 5}, {To: 7}},
		5: {{To: 4}, {To: 6}},
		4: {{To: 5}, {To: 2}, {To: 3}},
		7: {{To: 6}},
		6: {{To: 7}, {To: 3}},
		3: {{To: 1}},
		1: {{To: 2}},
		2: {{To: 3}},
	}}
	g.StronglyConnectedComponents(func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [3 1 2]
	// [7 6]
	// [4 5]
	// [0]
}

func ExampleLabeledDirected_Condensation() {
	// input:          condensation:
	// /---0---\      <->  /---0
	// |   |\--/           |   |
	// |   v               |   v
	// |   5<=>4---\  <->  |   1--\
	// |   |   |   |       |   |  |
	// v   v   |   |       |   v  |
	// 7<=>6   |   |  <->  \-->2  |
	//     |   v   v           |  v
	//     \-->3<--2  <->      \->3
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 0}, {To: 5}, {To: 7}},
		5: {{To: 4}, {To: 6}},
		4: {{To: 5}, {To: 2}, {To: 3}},
		7: {{To: 6}},
		6: {{To: 7}, {To: 3}},
		3: {{To: 1}},
		1: {{To: 2}},
		2: {{To: 3}},
	}}
	scc, cd := g.Condensation()
	fmt.Println(len(scc), "components:")
	for cn, c := range scc {
		fmt.Println(cn, c)
	}
	fmt.Println("condensation:")
	for cn, to := range cd {
		fmt.Println(cn, to)
	}
	// Output:
	// 4 components:
	// 0 [0]
	// 1 [4 5]
	// 2 [7 6]
	// 3 [3 1 2]
	// condensation:
	// 0 [1 2]
	// 1 [3 2]
	// 2 [3]
	// 3 []
}

func ExampleLabeledDirected_Topological() {
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{To: 2}},
		3: {{To: 1}, {To: 2}},
		4: {{To: 3}, {To: 2}},
	}}
	fmt.Println(g.Topological())
	g.LabeledAdjacencyList[2] = []graph.Half{{To: 3}}
	fmt.Println(g.Topological())
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleLabeledDirected_TopologicalKahn() {
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		1: {{To: 2}},
		3: {{To: 1}, {To: 2}},
		4: {{To: 3}, {To: 2}},
	}}
	tr, _ := g.UnlabeledTranspose()
	fmt.Println(g.TopologicalKahn(tr))

	g.LabeledAdjacencyList[2] = []graph.Half{{To: 3}}
	tr, _ = g.UnlabeledTranspose()
	fmt.Println(g.TopologicalKahn(tr))
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleLabeledDirected_TopologicalSubgraph() {
	// arcs directected down unless otherwise indicated
	// 0       1<-\
	//  \     / \ /
	//   2   3   4
	//    \ / \
	//     5   6
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 2}},
		1: {{To: 3}, {To: 4}},
		2: {{To: 5}},
		3: {{To: 5}, {To: 6}},
		4: {{To: 1}},
		6: {},
	}}
	fmt.Println(g.TopologicalSubgraph([]graph.NI{0, 3}))
	// Output:
	// [3 6 0 2 5] []
}
