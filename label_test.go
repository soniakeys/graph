// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleLabeledAdjacencyList_DAGMaxLenPath() {
	// arcs directed right:
	//            (M)
	//    (W)  /---------\
	//  3-----4   1-------0-----2
	//         \    (S)  /  (P)
	//          \       /
	//           \-----/ (Q)
	g := graph.LabeledAdjacencyList{
		3: {{To: 0, Label: 'Q'}, {4, 'W'}},
		4: {{0, 'M'}},
		1: {{0, 'S'}},
		0: {{2, 'P'}},
	}
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

func ExampleLabeledAdjacencyList_HasArcLabel() {
	g := graph.LabeledAdjacencyList{
		2: {{0, 10}, {2, 20}, {0, 30}},
	}
	fmt.Println(g.HasArcLabel(2, 0, 30))
	// Output:
	// true 2
}

func ExampleLabeledAdjacencyList_HasParallelSort_parallelArcs() {
	g := graph.LabeledAdjacencyList{
		1: {{0, 0}, {0, 0}},
	}
	// result true 1 0 means parallel arc from node 1 to node 0
	fmt.Println(g.HasParallelSort())
	// Output:
	// true 1 0
}

func ExampleLabeledAdjacencyList_HasParallelSort_noParallelArcs() {
	g := graph.LabeledAdjacencyList{
		1: {{0, 0}},
	}
	// result false -1 -1 means no parallel arc
	fmt.Println(g.HasParallelSort())
	// Output:
	// false -1 -1
}

func ExampleLabeledAdjacencyList_IsUndirected() {
	//             0<--
	// (label 'A')  \  \ (matching label 'A' on reciprocal)
	//               -->1
	// 2<--\
	// |   | (label 'B' on loop)
	// \---/
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 'A'}},
		1: {{0, 'A'}},
		2: {{2, 'B'}},
	}
	ok, _, _ := g.IsUndirected()
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleLabeledAdjacencyList_IsUndirected_undirectedMultigraph() {
	// lines shown here are edges (arcs with reciprocals.)
	//               0---
	//  (Label: 'A')  \  \  (Label: 'B')
	//                 ---1
	g := &graph.LabeledAdjacencyList{}
	g.AddEdge(graph.Edge{0, 1}, 'A')
	g.AddEdge(graph.Edge{0, 1}, 'B')
	ok, _, _ := g.IsUndirected()
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleLabeledAdjacencyList_IsUndirected_labelMismatch() {
	// directed graph, arcs with different labels
	//               0<--
	//  (Label: 'A')  \  \  (Label: 'B')
	//                 -->1
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 'A'}},
		1: {{To: 0, Label: 'B'}},
	}
	ok, fr, to := g.IsUndirected()
	fmt.Printf("%t %d {%d %c}\n", ok, fr, to.To, to.Label)
	// Output:
	// false 0 {1 A}
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

func ExampleLabeledAdjacencyList_TarjanBiconnectedComponents() {
	// undirected edges:
	// 3---2---1---7---9
	//  \ / \ / \   \ /
	//   4   5---6   8
	g := graph.LabeledAdjacencyList{}
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{2, 5}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{6, 1}, 0)
	g.AddEdge(graph.Edge{6, 5}, 0)
	g.AddEdge(graph.Edge{7, 1}, 0)
	g.AddEdge(graph.Edge{7, 9}, 0)
	g.AddEdge(graph.Edge{7, 8}, 0)
	g.AddEdge(graph.Edge{9, 8}, 0)
	g.TarjanBiconnectedComponents(func(bcc []graph.LabeledEdge) bool {
		fmt.Println("Edges:")
		for _, e := range bcc {
			fmt.Println(e.Edge)
		}
		return true
	})
	// Output:
	// Edges:
	// {4 2}
	// {3 4}
	// {2 3}
	// Edges:
	// {6 1}
	// {5 6}
	// {5 1}
	// {2 5}
	// {1 2}
	// Edges:
	// {8 7}
	// {9 8}
	// {7 9}
	// Edges:
	// {1 7}
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

func ExampleLabeledAdjacencyList_WeightedInDegree() {
	//  0
	//  | (weight = label: 3)
	//  v
	//  1
	//  | (4)
	//  v
	//  2<-\
	//  \--/ (5)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 3}},
		1: {{2, 4}},
		2: {{2, 5}},
	}
	w := func(l graph.LI) float64 { return float64(l) }
	fmt.Println(g.WeightedInDegree(w))
	// Output:
	// [0 3 9]
}

func ExampleLabeledAdjacencyList_WeightedOutDegree() {
	//  0
	//  | (weight = label: 3)
	//  v
	//  1
	//  | (4)
	//  v
	//  2<-\
	//  \--/ (5)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 3}},
		1: {{2, 4}},
		2: {{2, 5}},
	}
	w := func(l graph.LI) float64 { return float64(l) }
	fmt.Println("node  weighted out degree")
	for n := range g {
		fmt.Println(n, "   ", g.WeightedOutDegree(graph.NI(n), w))
	}
	// Output:
	// node  weighted out degree
	// 0     3
	// 1     4
	// 2     5
}

func ExampleLabeledAdjacencyList_WeightedOutDegree_undirected() {
	//  0
	//  | (weight = label: 3)
	//  |
	//  1
	//  | (4)
	//  |
	//  2--\
	//  \--/ (5)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 3}},
		1: {{0, 3}, {2, 4}},
		2: {{1, 4}, {2, 5}},
	}
	w := func(l graph.LI) float64 { return float64(l) }
	ok, _, _ := g.IsUndirected()
	fmt.Println("undirected:", ok)
	fmt.Println()
	fmt.Println("node  weighted out-degree  weighted in-degree")
	ind := g.WeightedInDegree(w)
	for n := range g {
		fmt.Println(n, "   ", g.WeightedOutDegree(graph.NI(n), w),
			"                  ", ind[n])
	}
	// Output:
	// undirected: true
	//
	// node  weighted out-degree  weighted in-degree
	// 0     3                    3
	// 1     7                    7
	// 2     9                    9
}
