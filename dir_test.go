// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleAdjacencyList_Cyclic() {
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {},
	}
	fmt.Println(g.Cyclic())
	g[3] = []int{2}
	fmt.Println(g.Cyclic())
	// Output:
	// false
	// true
}

func ExampleAdjacencyList_Topological() {
	g := graph.AdjacencyList{
		1: {2},
		3: {1, 2},
		4: {3, 2},
	}
	fmt.Println(g.Topological())
	g[2] = []int{3}
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

	g[2] = []int{3}
	tr, _ = g.Transpose()
	fmt.Println(g.TopologicalKahn(tr))
	// Output:
	// [4 3 1 2 0] []
	// [] [1 2 3]
}

func ExampleAdjacencyList_Tarjan() {
	g := graph.AdjacencyList{
		0: {1},
		1: {4, 2, 5},
		2: {3, 6},
		3: {2, 7},
		4: {5, 0},
		5: {6},
		6: {5},
		7: {3, 6},
	}
	scc := g.Tarjan()
	for _, c := range scc {
		fmt.Println(c)
	}
	fmt.Println(len(scc))
	// Output:
	// [6 5]
	// [7 3 2]
	// [4 1 0]
	// 3
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

func ExampleAdjacencyList_EulerianCycleUndir() {
	g := graph.AdjacencyList{
		0: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}
	g, m := g.CopyUndir()
	fmt.Println(g.EulerianCycleUndirD(m))
	// Output:
	// [0 1 2 2 1 2 0] <nil>
}

func TestEulerianCycle(t *testing.T) {
	same := func(a, b []int) bool {
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
		cycle []int
		ok    bool
	}{
		{nil, nil, true},
		{graph.AdjacencyList{nil}, []int{0}, true},    // 1 node, 0 arcs
		{graph.AdjacencyList{{0}}, []int{0, 0}, true}, // loop
		{graph.AdjacencyList{nil, nil}, nil, false},   // not connected
		{graph.AdjacencyList{{1}, nil}, nil, false},   // not balanced
		{graph.AdjacencyList{nil, {0}}, nil, false},   // not balanced
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
	same := func(a, b []int) bool {
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
		path []int
		ok   bool
	}{
		{nil, nil, true},
		{graph.AdjacencyList{nil}, []int{0}, true},    // 1 node, 0 arcs
		{graph.AdjacencyList{{0}}, []int{0, 0}, true}, // loop
		{graph.AdjacencyList{{1}, nil}, []int{0, 1}, true},
		{graph.AdjacencyList{nil, {0}}, []int{1, 0}, true},
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
