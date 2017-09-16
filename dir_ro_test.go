// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// dir_ro_test.go -- tests on dir_RO.go
//
// These are tests on the code-generated unlabeled versions of methods.
//
// If testing prompts changes in the tested method, be sure to edit
// dir_cg.go, go generate to generate the dir_RO.go, and then retest.
// Do not edit dir_RO.go
//
// See also dir_cg_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleDirected_Balanced() {
	// 2
	// |
	// v
	// 0----->1
	g := graph.Directed{graph.AdjacencyList{
		2: {0},
		0: {1},
	}}
	fmt.Println(g.Balanced())

	// 0<--\
	// |    \
	// v     \
	// 1----->2
	g.AdjacencyList[1] = []graph.NI{2}
	fmt.Println(g.Balanced())
	// Output:
	// false
	// true
}

func ExampleDirected_Cyclic() {
	//   0
	//  / \
	// 1-->2-->3
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
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
	g.AdjacencyList[3] = []graph.NI{1}
	fmt.Println(g.Cyclic())

	// Output:
	// false
	// true 3 1
}

func ExampleDirected_Dominators() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {5},
	}}
	d := g.Dominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleDirected_Doms() {
	//   0   6
	//   |   |
	//   1   |
	//  / \  |
	// 2   3 |
	//  \ / \|
	//   4   5
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {5},
	}}
	// compute postorder with depth-first traversal
	var post []graph.NI
	var vis big.Int
	var f func(graph.NI)
	f = func(n graph.NI) {
		vis.SetBit(&vis, int(n), 1)
		for _, to := range g.AdjacencyList[n] {
			if vis.Bit(int(to)) == 0 {
				f(to)
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

func ExampleDirected_Eulerian() {
	//   /<--------\
	//  /   /<---\  \
	// 0-->1-->\ /  /
	//      \-->2--/
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.Eulerian())
	// Output:
	// -1 -1 <nil>
}

func ExampleDirected_EulerianCycle() {
	//   /<--------\
	//  /   /<---\  \
	// 0-->1-->\ /  /
	//      \-->2--/
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.EulerianCycle())
	// Output:
	// [0 1 2 1 2 2 0] <nil>
}

func ExampleDirected_EulerianCycleD() {
	//   /<--------\
	//  /   /<---\  \
	// 0-->1-->\ /  /
	//      \-->2--/
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println("ma =", g.ArcSize())
	fmt.Println(g.EulerianCycleD(6))
	// Output:
	// ma = 6
	// [0 1 2 1 2 2 0] <nil>
}

func TestEulerianCycle(t *testing.T) {
	same := func(a, b []graph.NI) bool {
		if len(a) != len(b) {
			return false
		}
		for i, x := range a {
			if b[i] != x {
				return false
			}
		}
		return true
	}
	var msg string
	for _, tc := range []struct {
		g     graph.AdjacencyList
		cycle []graph.NI
		ok    bool
	}{
		{nil, nil, true},
		{graph.AdjacencyList{nil}, []graph.NI{0}, true},    // 1 node, 0 arcs
		{graph.AdjacencyList{{0}}, []graph.NI{0, 0}, true}, // loop
		{graph.AdjacencyList{nil, nil}, nil, false},        // not connected
		{graph.AdjacencyList{{1}, nil}, nil, false},        // not balanced
		{graph.AdjacencyList{nil, {0}}, nil, false},        // not balanced
	} {
		got, err := graph.Directed{tc.g}.EulerianCycle()
		switch {
		case err != nil:
			if !tc.ok {
				continue
			}
			msg = "g.EulerianCycle() returned error" + err.Error()
		case !tc.ok:
			msg = fmt.Sprintf("g.EulerianCycle() = %v, want error", got)
		case !same(got, tc.cycle):
			msg = fmt.Sprintf("g.EulerianCycle() = %v, want %v", got, tc.cycle)
		default:
			continue
		}
		t.Log("g:", tc.g)
		t.Fatal(msg)
	}
}

func ExampleDirected_EulerianPath() {
	//      /<---\
	// 3-->1-->\ /
	//      \-->2-->0
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.EulerianPath())
	// Output:
	// [3 1 2 1 2 2 0] <nil>
}

func ExampleDirected_EulerianPathD() {
	//      /<---\
	// 3-->1-->\ /
	//      \-->2-->0
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println("ma =", g.ArcSize())
	fmt.Print("start: ")
	fmt.Println(g.EulerianStart())
	fmt.Println(g.EulerianPathD(6, 3))
	// Output:
	// ma = 6
	// start: 3 <nil>
	// [3 1 2 1 2 2 0] <nil>
}

func ExampleDirected_EulerianStart() {
	//      /<---\
	// 3-->1-->\ /
	//      \-->2-->0
	//         / \
	//        /<--\
	g := graph.Directed{graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
	fmt.Println(g.EulerianStart())
	// Output:
	// 3 <nil>
}

func TestEulerianPath(t *testing.T) {
	same := func(a, b []graph.NI) bool {
		if len(a) != len(b) {
			return false
		}
		for i, x := range a {
			if b[i] != x {
				return false
			}
		}
		return true
	}
	var msg string
	for _, tc := range []struct {
		g    graph.AdjacencyList
		path []graph.NI
		ok   bool
	}{
		{nil, nil, true},
		{graph.AdjacencyList{nil}, []graph.NI{0}, true},    // 1 node, 0 arcs
		{graph.AdjacencyList{{0}}, []graph.NI{0, 0}, true}, // loop
		{graph.AdjacencyList{{1}, nil}, []graph.NI{0, 1}, true},
		{graph.AdjacencyList{nil, {0}}, []graph.NI{1, 0}, true},
		{graph.AdjacencyList{nil, nil}, nil, false},         // not connected
		{graph.AdjacencyList{{1}, nil, {1}}, nil, false},    // two starts
		{graph.AdjacencyList{nil, nil, {0, 1}}, nil, false}, // two ends
	} {
		got, err := graph.Directed{tc.g}.EulerianPath()
		switch {
		case err != nil:
			if !tc.ok {
				continue
			}
			msg = "g.EulerianPath() returned error" + err.Error()
		case !tc.ok:
			msg = fmt.Sprintf("g.EulerianPath() = %v, want error", got)
		case !same(got, tc.path):
			msg = fmt.Sprintf("g.EulerianPath() = %v, want %v", got, tc.path)
		default:
			continue
		}
		t.Log("g:", tc.g)
		t.Fatal(msg)
	}
}

func ExampleDirected_InDegree() {
	// arcs directed down:
	//  0     2
	//  |
	//  1
	//  |\
	//  | \
	//  3  4<-\
	//     \--/
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {3, 4},
		4: {4},
	}}
	fmt.Println("node:    0 1 2 3 4")
	fmt.Println("in-deg:", g.InDegree())
	// Output:
	// node:    0 1 2 3 4
	// in-deg: [0 1 0 1 2]
}

func ExampleDirected_IsTree() {
	// Example graph
	// Arcs point down unless otherwise indicated
	//           1
	//          / \
	//         0   5
	//        /   / \
	//       2   3-->4
	g := graph.Directed{graph.AdjacencyList{
		1: {0, 5},
		0: {2},
		5: {3, 4},
		3: {4},
	}}
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	// Output:
	// true false
	// false false
}

func ExampleDirected_MaximalNonBranchingPaths() {
	// 0-->1-->2-->3
	//          \
	//   -->6    ->4
	//  /  /
	// 5<--
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {3, 4},
		5: {6},
		6: {5},
	}}
	g.MaximalNonBranchingPaths(func(p []graph.NI) bool {
		fmt.Println(p)
		return true
	})
	// Output:
	// [0 1 2]
	// [2 3]
	// [2 4]
	// [5 6 5]
}

func ExampleDirected_PostDominators() {
	// Example graph here is transpose of that in the Dominators example
	// to show result is the same.
	//   4   5
	//  / \ /|
	// 2   3 |
	//  \ /  |
	//   1   |
	//   |   |
	//   0   6
	g := graph.Directed{graph.AdjacencyList{
		4: {2, 3},
		5: {3, 6},
		2: {1},
		3: {1},
		1: {0},
		6: {},
	}}
	d := g.PostDominators(0)
	fmt.Println(d.Immediate)
	// Output:
	// [0 0 1 1 1 3 -1]
}

func ExampleDirected_StronglyConnectedComponents() {
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
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
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

func ExampleDirected_Condensation() {
	// input:          condensation:
	// /---0---\      <->  /---3
	// |   |\--/           |   |
	// |   v               |   v
	// |   5<=>4---\  <->  |   2--\
	// |   |   |   |       |   |  |
	// v   v   |   |       |   v  |
	// 7<=>6   |   |  <->  \-->1  |
	//     |   v   v           |  v
	//     \-->3<--2  <->      \->0
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
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
	// 0 [3 1 2]
	// 1 [7 6]
	// 2 [4 5]
	// 3 [0]
	// condensation:
	// 0 []
	// 1 [0]
	// 2 [0 1]
	// 3 [2 1]
}

func ExampleDirected_Topological() {
	g := graph.Directed{graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}}
	fmt.Println(g.Topological())
	g.AdjacencyList[2] = []graph.NI{3}
	fmt.Println(g.Topological())
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleDirected_TopologicalKahn() {
	g := graph.Directed{graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}}
	tr, _ := g.Transpose()
	fmt.Println(g.TopologicalKahn(tr))

	g.AdjacencyList[2] = []graph.NI{3}
	tr, _ = g.Transpose()
	fmt.Println(g.TopologicalKahn(tr))
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleDirected_TopologicalSubgraph() {
	// arcs directected down unless otherwise indicated
	// 0       1<-\
	//  \     / \ /
	//   2   3   4
	//    \ / \
	//     5   6
	g := graph.Directed{graph.AdjacencyList{
		0: {2},
		1: {3, 4},
		2: {5},
		3: {5, 6},
		4: {1},
		6: {},
	}}
	fmt.Println(g.TopologicalSubgraph([]graph.NI{0, 3}))
	// Output:
	// [3 6 0 2 5] []
}
