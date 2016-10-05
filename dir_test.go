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
	// [5 6 5]
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

func ExampleLabeledDirected_DAGMaxLenPath() {
	// arcs directed right:
	//            (M)
	//    (W)  /---------\
	//  3-----4   1-------0-----2
	//         \    (S)  /  (P)
	//          \       /
	//           \-----/ (Q)
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		3: {{To: 0, Label: 'Q'}, {4, 'W'}},
		4: {{0, 'M'}},
		1: {{0, 'S'}},
		0: {{2, 'P'}},
	}}
	o, _ := g.Topological()
	fmt.Println("ordering:", o)
	n, p := g.DAGMaxLenPath(o)
	fmt.Printf("path from %d:", n)
	for _, e := range p {
		fmt.Printf(" {%d, '%c'}", e.To, e.Label)
	}
	fmt.Println()
	fmt.Print("label path: ")
	for _, h := range p {
		fmt.Print(string(h.Label))
	}
	fmt.Println()
	// Output:
	// ordering: [3 4 1 0 2]
	// path from 3: {4, 'W'} {0, 'M'} {2, 'P'}
	// label path: WMP
}

func ExampleLabeledDirected_Transpose() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}}
	tr, m := g.Transpose()
	for fr, to := range tr.LabeledAdjacencyList {
		fmt.Printf("%d %#v\n", fr, to)
	}
	fmt.Println(m, "arcs")
	// Output:
	// 0 []graph.Half{graph.Half{To:2, Label:7}}
	// 1 []graph.Half{graph.Half{To:2, Label:9}}
	// 2 []graph.Half(nil)
	// 2 arcs
}

func ExampleLabeledDirected_Undirected() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}}
	for fr, to := range g.Undirected().LabeledAdjacencyList {
		fmt.Printf("%d %#v\n", fr, to)
	}
	// Output:
	// 0 []graph.Half{graph.Half{To:2, Label:7}}
	// 1 []graph.Half{graph.Half{To:2, Label:9}}
	// 2 []graph.Half{graph.Half{To:0, Label:7}, graph.Half{To:1, Label:9}}
}

func ExampleLabeledDirected_UnlabeledTranspose() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}}

	fmt.Println("two steps:")
	ut, m := g.Unlabeled().Transpose()
	for fr, to := range ut.AdjacencyList {
		fmt.Println(fr, to)
	}
	fmt.Println(m, "arcs")

	fmt.Println("direct:")
	ut, m = g.UnlabeledTranspose()
	for fr, to := range ut.AdjacencyList {
		fmt.Println(fr, to)
	}
	fmt.Println(m, "arcs")
	// Output:
	// two steps:
	// 0 [2]
	// 1 [2]
	// 2 []
	// 2 arcs
	// direct:
	// 0 [2]
	// 1 [2]
	// 2 []
	// 2 arcs
}

func ExampleDominators_Frontiers() {
	//   0
	//   |
	//   1
	//  / \
	// 2   3
	//  \ / \
	//   4   5   6
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {},
	}}
	for n, f := range g.Dominators(0).Frontiers() {
		fmt.Print(n, ":")
		if f == nil {
			fmt.Println(" nil")
			continue
		}
		for n := range f {
			fmt.Print(" ", n)
		}
		fmt.Println()
	}
	// Output:
	// 0:
	// 1:
	// 2: 4
	// 3: 4
	// 4:
	// 5:
	// 6: nil
}

func ExampleDominators_Set() {
	//   0
	//   |
	//   1
	//  / \
	// 2   3
	//  \ / \
	//   4   5   6
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2, 3},
		2: {4},
		3: {4, 5},
		6: {},
	}}
	d := g.Dominators(0)
	for n := range g.AdjacencyList {
		fmt.Println(n, d.Set(graph.NI(n)))
	}
	// Output:
	// 0 [0]
	// 1 [1 0]
	// 2 [2 1 0]
	// 3 [3 1 0]
	// 4 [4 1 0]
	// 5 [5 3 1 0]
	// 6 []
}
