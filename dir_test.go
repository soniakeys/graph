// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

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

func ExampleDominanceFrontiers_Closure() {
	//     0
	//     |
	//     1
	//     |
	// --->2
	// |  / \
	// --3   4
	//      / \
	//     5   6
	//      \ /
	//       7
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {3, 4},
		3: {2},
		4: {5, 6},
		5: {7},
		6: {7},
		7: {},
	}}
	f := g.Dominators(0).Frontiers()
	type ns map[graph.NI]struct{}
	fmt.Println(f.Closure(ns{
		0: struct{}{},
		1: struct{}{},
		3: struct{}{},
	}))
	// Output:
	// map[2:{}]
}

func ExampleDominanceFrontiers_Frontier() {
	//     0
	//     |
	//     1
	//     |
	// --->2
	// |  / \
	// --3   4
	//      / \
	//     5   6
	//      \ /
	//       7
	g := graph.Directed{graph.AdjacencyList{
		0: {1},
		1: {2},
		2: {3, 4},
		3: {2},
		4: {5, 6},
		5: {7},
		6: {7},
		7: {},
	}}
	f := g.Dominators(0).Frontiers()
	type ns map[graph.NI]struct{}
	fmt.Println(f.Frontier(ns{
		0: struct{}{},
		1: struct{}{},
		3: struct{}{},
	}))
	// Output:
	// map[2:{}]
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

// ------- Labeled examples -------

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
