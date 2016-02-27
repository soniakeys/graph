// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"testing"

	"github.com/soniakeys/graph"
)

func ExampleDirected_DAGMaxLenPath() {
	// arcs directed right:
	//      /---\
	//  3--4  1--0--2
	//   \------/
	g := graph.Directed{graph.AdjacencyList{
		3: {0, 4},
		4: {0},
		1: {0},
		0: {2},
	}}
	o, _ := g.Topological()
	fmt.Println(o)
	fmt.Println(g.DAGMaxLenPath(o))
	// Output:
	// [3 4 1 0 2]
	// [3 4 0 2]
}

func ExampleDirected_EulerianCycle() {
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
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(1, 2)
	g.AddEdge(2, 2)
	fmt.Println(g.EulerianCycleD(6))
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
	g := graph.Directed{graph.AdjacencyList{
		3: {1},
		1: {2, 2},
		2: {0, 1, 2},
	}}
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
	// [6 5 6]
}

func ExampleDirected_Transpose() {
	g := graph.Directed{graph.AdjacencyList{
		2: {0, 1},
	}}
	t, m := g.Transpose()
	for n, nbs := range t.AdjacencyList {
		fmt.Printf("%d: %v\n", n, nbs)
	}
	fmt.Println(m)
	// Output:
	// 0: [2]
	// 1: [2]
	// 2: []
	// 2
}
