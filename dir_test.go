// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_DAGMaxLenPath() {
	// arcs directed right:
	//      /---\
	//  3--4  1--0--2
	//   \------/
	g := graph.AdjacencyList{
		3: {0, 4},
		4: {0},
		1: {0},
		0: {2},
	}
	o, _ := g.Topological()
	fmt.Println(o)
	fmt.Println(g.DAGMaxLenPath(o))
	// Output:
	// [3 4 1 0 2]
	// [3 4 0 2]
}

func ExampleAdjacencyList_IsTreeDirected() {
	// Example graph
	// Arcs point down unless otherwise indicated
	//           1
	//          / \
	//         0   5
	//        /   / \
	//       2   3-->4
	g := graph.AdjacencyList{
		1: {0, 5},
		0: {2},
		5: {3, 4},
		3: {4},
	}
	fmt.Println(g.IsTreeDirected(0))
	fmt.Println(g.IsTreeDirected(1))
	// Output:
	// true
	// false
}

func ExampleAdjacencyList_EulerianCycle() {
	g := graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}
	fmt.Println(g.EulerianCycle())
	// Output:
	// [0 1 2 1 2 2 0] <nil>
}

func ExampleAdjacencyList_EulerianCycleUndirD() {
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {0, 2, 2, 2},
		2: {0, 1, 1, 1, 2, 2},
	}
	fmt.Println(g.EulerianCycleUndirD(6))
	// Output:
	// [0 1 2 2 1 2 0] <nil>
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
		got, err := tc.g.EulerianCycle()
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

func ExampleAdjacencyList_EulerianPath() {
	g := graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}
	fmt.Println(g.EulerianPath())
	// Output:
	// [3 1 2 1 2 2 0] <nil>
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
		got, err := tc.g.EulerianPath()
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

func ExampleAdjacencyList_FromList() {
	//    4   3
	//   / \
	//  2   1
	//       \
	//        0
	g := graph.AdjacencyList{
		4: {2, 1},
		1: {0},
	}
	f, _ := g.FromList()
	fmt.Println("Paths:")
	fmt.Println("N  From  Len")
	for n, e := range f.Paths {
		fmt.Printf("%d %4d %5d\n", n, e.From, e.Len)
	}
	fmt.Println("Leaves:")
	fmt.Println("43210")
	fmt.Println("-----")
	fmt.Printf("%05b\n", &f.Leaves)
	fmt.Println("MaxLen:", f.MaxLen)
	// Output:
	// Paths:
	// N  From  Len
	// 0    1     3
	// 1    4     2
	// 2    4     2
	// 3   -1     1
	// 4   -1     1
	// Leaves:
	// 43210
	// -----
	// 01101
	// MaxLen: 3
}

func ExampleAdjacencyList_FromList_nonTree() {
	//    0
	//   / \
	//  1   2
	//   \ /
	//    3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {3},
		2: {3},
		3: {},
	}
	fmt.Println(g.FromList())
	// Output:
	// <nil> 3
}

func ExampleAdjacencyList_MaximalNonBranchingPaths() {
	// 0-->1-->2-->3
	//          \
	//   -->6    ->4
	//  /  /
	// 5<--
	g := graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {3, 4},
		5: {6},
		6: {5},
	}
	for p := range g.MaximalNonBranchingPaths() {
		fmt.Println(p)
	}
	// Output:
	// [0 1 2]
	// [2 3]
	// [2 4]
	// [6 5 6]
}

func ExampleAdjacencyList_TarjanCondensation() {
	// input:          condensation:
	// /---0---\        /---0
	// |   |\--/        |   |
	// |   v            |   v
	// |   5<=>4---\    |   1--\
	// |   |   |   |    |   |  |
	// v   v   |   |    |   v  |
	// 7<=>6   |   |    \-->2  |
	//     |   v   v        |  v
	//     \-->3<--2        \->3
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}
	scc, cd := g.TarjanCondensation()
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
	// 3 [1 3 2]
	// condensation:
	// 0 [1 2]
	// 1 [3 2]
	// 2 [3]
	// 3 []
}

func ExampleAdjacencyList_Topological() {
	g := graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}
	fmt.Println(g.Topological())
	g[2] = []graph.NI{3}
	fmt.Println(g.Topological())
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleAdjacencyList_TopologicalKahn() {
	g := graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}
	tr, _ := g.Transpose()
	fmt.Println(g.TopologicalKahn(tr))

	g[2] = []graph.NI{3}
	tr, _ = g.Transpose()
	fmt.Println(g.TopologicalKahn(tr))
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleAdjacencyList_Transpose() {
	g := graph.AdjacencyList{
		2: {0, 1},
	}
	t, m := g.Transpose()
	for n, nbs := range t {
		fmt.Printf("%d: %v\n", n, nbs)
	}
	fmt.Println(m)
	// Output:
	// 0: [2]
	// 1: [2]
	// 2: []
	// 2
}
