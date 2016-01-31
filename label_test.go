// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleLabledAdjacencyList_DAGMaxLenPath() {
	// arcs directed right:
	//           (M)
	//    (W)  /------\
	//  3-----4  1-----0-----2
	//             (S)   (P)
	g := graph.LabeledAdjacencyList{
		3: {{To: 4, Label: 'W'}},
		4: {{To: 0, Label: 'M'}},
		1: {{To: 0, Label: 'S'}},
		0: {{To: 2, Label: 'P'}},
	}
	o, _ := g.Topological()
	fmt.Println("ordering:", o)
	n, p := g.DAGMaxLenPath(o)
	fmt.Printf("path from %d: %v\n", n, p)
	fmt.Print("label path: ")
	for _, h := range p {
		fmt.Print(string(h.Label))
	}
	fmt.Println()
	// Output:
	// ordering: [3 4 1 0 2]
	// path from 3: [{4 87} {0 77} {2 80}]
	// label path: WMP
}

func ExampleLabeledAdjacencyList_BoundsOk() {
	g := graph.LabeledAdjacencyList{
		0: {{0, -1}},
	}
	ok, _, _ := g.BoundsOk()
	fmt.Println(ok)
	g = graph.LabeledAdjacencyList{
		0: {{-1, -1}},
	}
	fmt.Println(g.BoundsOk())
	g = graph.LabeledAdjacencyList{
		0: {{9, -1}},
	}
	fmt.Println(g.BoundsOk())
	// Output:
	// true
	// false 0 {-1 -1}
	// false 0 {9 -1}
}

func ExampleLabeledAdjacencyList_IsUndirected() {
	// multigraph, edges with different labels
	//               ----0
	//  (Label: 0)  /   /  (Label: 1)
	//             1----
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 0}, {To: 1, Label: 1}},
		1: {{To: 0, Label: 0}, {To: 0, Label: 1}},
	}
	ud, _, _ := g.IsUndirected()
	fmt.Println(ud)
	// directed graph, arcs with different labels
	//               --->0
	//  (Label: 0)  /   /  (Label: 1)
	//             1<---
	g = graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 1}},
		1: {{To: 0, Label: 0}},
	}
	fmt.Println(g.IsUndirected())
	// Output:
	// true
	// false 0 {1 1}
}

// A directed graph with negative arc weights.
// Arc weights are encoded simply as label numbers.
func ExampleLabeledAdjacencyList_FloydWarshall_negative() {
	g := graph.LabeledAdjacencyList{
		0: {{To: 2, Label: -1}},
		1: {{To: 3, Label: -2}},
		2: {{To: 1, Label: 4}, {To: 3, Label: 3}},
		3: {{To: 0, Label: 2}},
	}
	d := g.FloydWarshall(func(l graph.LI) float64 { return float64(l) })
	for _, di := range d {
		fmt.Printf("%2.0f\n", di)
	}
	// Output:
	// [ 0  3 -1  1]
	// [ 0  0 -1 -2]
	// [ 4  4  0  2]
	// [ 2  5  1  0]
}

func ExampleLabeledAdjacencyList_NegativeArc() {
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 0}, {To: 1, Label: 1}},
	}
	arcWeights := []float64{0, .5}
	w := func(label graph.LI) float64 { return arcWeights[label] }
	fmt.Println(g.NegativeArc(w))
	g[0] = []graph.Half{{To: 1, Label: graph.LI(len(arcWeights))}}
	arcWeights = append(arcWeights, -2)
	fmt.Println(g.NegativeArc(w))
	// Output:
	// false
	// true
}

func ExampleLabeledAdjacencyList_Transpose() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}
	tr, m := g.Transpose()
	for fr, to := range tr {
		fmt.Printf("%d %#v\n", fr, to)
	}
	fmt.Println(m, "arcs")
	// Output:
	// 0 []graph.Half{graph.Half{To:2, Label:7}}
	// 1 []graph.Half{graph.Half{To:2, Label:9}}
	// 2 []graph.Half(nil)
	// 2 arcs
}

func ExampleLabeledAdjacencyList_UndirectedCopy() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}
	for fr, to := range g.UndirectedCopy() {
		fmt.Printf("%d %#v\n", fr, to)
	}
	// Output:
	// 0 []graph.Half{graph.Half{To:2, Label:7}}
	// 1 []graph.Half{graph.Half{To:2, Label:9}}
	// 2 []graph.Half{graph.Half{To:0, Label:7}, graph.Half{To:1, Label:9}}
}

func ExampleLabeledAdjacencyList_Unlabeled() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}
	fmt.Println("Input:")
	for fr, to := range g {
		fmt.Printf("%d, %#v\n", fr, to)
	}
	fmt.Println("\nUnlabeled:")
	for fr, to := range g.Unlabeled() {
		fmt.Printf("%d, %#v\n", fr, to)
	}
	// Output:
	// Input:
	// 0, []graph.Half(nil)
	// 1, []graph.Half(nil)
	// 2, []graph.Half{graph.Half{To:0, Label:7}, graph.Half{To:1, Label:9}}
	//
	// Unlabeled:
	// 0, []graph.NI{}
	// 1, []graph.NI{}
	// 2, []graph.NI{0, 1}
}

func ExampleLabeledAdjacencyList_UnlabeledTranspose() {
	// arcs directed down:
	//             2
	//  (label: 7)/ \(9)
	//           0   1
	g := graph.LabeledAdjacencyList{
		2: {{To: 0, Label: 7}, {To: 1, Label: 9}},
	}

	fmt.Println("two steps:")
	ut, m := g.Unlabeled().Transpose()
	for fr, to := range ut {
		fmt.Println(fr, to)
	}
	fmt.Println(m, "arcs")

	fmt.Println("direct:")
	ut, m = g.UnlabeledTranspose()
	for fr, to := range ut {
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
