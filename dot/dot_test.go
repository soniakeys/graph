package dot_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func ExampleStringAdjacencyList() {
	// arcs directed down:
	// 0  4
	// | /|
	// |/ |
	// 2  3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	// (default indent is 2)
	s, _ := dot.StringAdjacencyList(g, dot.Indent(""))
	fmt.Println(s)
	// Output:
	// digraph {
	// 0 -> 2
	// 4 -> {2 3}
	// }
}

func ExampleStringAdjacencyList_parallelArcs() {
	// arcs directed down:
	// 0  4
	// | /|\
	// |/ |/
	// 2  3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3, 3},
	}
	s, _ := dot.StringAdjacencyList(g, dot.Indent(""))
	fmt.Println(s)
	// Output:
	// digraph {
	// 0 -> 2
	// 4 -> {2 3}
	// 4 -> {3}
	// }
}

func ExampleIndent() {
	// arcs directed down:
	// 0  4
	// | /|
	// |/ |
	// 2  3
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	// All other examples have Indent("") to avoid a quirk of go test
	// that it can't handle leading space in the output.  In this example a
	// nonbreaking space works around the quirk to show indented output that
	// looks like the default two space indent.
	// (But then if you render it with graphviz, graphviz picks up the nbsp
	// as a node statement...)
	s, _ := dot.StringAdjacencyList(g, dot.Indent("\u00a0 "))
	fmt.Println(s)
	// Output:
	// digraph {
	//   0 -> 2
	//   4 -> {2 3}
	// }
}

func ExampleDirected() {
	// arcs directed down:
	// 0  2
	// | /|
	// |/ |
	// 3  4
	g := graph.AdjacencyList{
		0: {3},
		2: {3, 4},
		4: {},
	}
	// default is directed
	s, _ := dot.StringAdjacencyList(g, dot.Indent(""))
	fmt.Println(s)
	fmt.Println()

	// Directed(false) generates error witout reciprocal arcs
	_, err := dot.StringAdjacencyList(g, dot.Directed(false))
	fmt.Println("Error:", err)
	fmt.Println()

	// undirected
	u := g.UndirectedCopy()
	s, _ = dot.StringAdjacencyList(u.AdjacencyList,
		dot.Directed(false), dot.Indent(""))
	fmt.Println(s)

	// Output:
	// digraph {
	// 0 -> 3
	// 2 -> {3 4}
	// }
	//
	// Error: directed graph
	//
	// graph {
	// 0 -- 3
	// 2 -- {3 4}
	// }
}

func ExampleNodeLabel() {
	// arcs directed down:
	// A  D
	// | /|
	// |/ |
	// B  C
	labels := []string{
		0: "A",
		4: "D",
		2: "B",
		3: "C",
	}
	lf := func(n graph.NI) string { return labels[n] }
	g := graph.AdjacencyList{
		0: {2},
		4: {2, 3},
	}
	s, _ := dot.StringAdjacencyList(g, dot.Indent(""), dot.NodeLabel(lf))
	fmt.Println(s)
	// Output:
	// digraph {
	// A -> B
	// D -> {B C}
	// }
}

func ExampleNodeLabel_construction() {
	// arcs directed down:
	// A  D
	// | /|
	// |/ |
	// B  C
	var g graph.AdjacencyList

	// example graph construction mechanism
	labels := []string{}
	nodes := map[string]graph.NI{}
	node := func(l string) graph.NI {
		if n, ok := nodes[l]; ok {
			return n
		}
		n := graph.NI(len(labels))
		labels = append(labels, l)
		g = append(g, nil)
		nodes[l] = n
		return n
	}
	addArc := func(fr, to string) {
		f := node(fr)
		g[f] = append(g[f], node(to))
	}

	// construct graph
	addArc("A", "B")
	addArc("D", "B")
	addArc("D", "C")

	// generate dot
	lf := func(n graph.NI) string { return labels[n] }
	s, _ := dot.StringAdjacencyList(g, dot.Indent(""), dot.NodeLabel(lf))
	fmt.Println(s)

	// Output:
	// digraph {
	// A -> B
	// D -> {B C}
	// }
}

func ExampleStringLabeledAdjacencyList() {
	// arcs directed down:
	//     0      4
	// (30)|     /|
	//     | (20) |
	//     |/     |(10)
	//     2      3
	g := graph.LabeledAdjacencyList{
		0: {{2, 30}},
		4: {{2, 20}, {3, 10}},
	}
	s, _ := dot.StringLabeledAdjacencyList(g, dot.Indent(""))
	fmt.Println(s)
	// Output:
	// digraph {
	// 0 -> 2 [label = 30]
	// 4 -> 2 [label = 20]
	// 4 -> 3 [label = 10]
	// }
}

func ExampleEdgeLabel() {
	// arcs directed down:
	//      0       4
	// (.33)|      /|
	//      | (1.7) |
	//      |/      |(2e117)
	//      2       3
	weights := map[int]float64{
		30: .33,
		20: 1.7,
		10: 2e117,
	}
	lf := func(l graph.LI) string {
		return fmt.Sprintf(`"%g"`, weights[int(l)])
	}
	g := graph.LabeledAdjacencyList{
		0: {{2, 30}},
		4: {{2, 20}, {3, 10}},
	}
	s, _ := dot.StringLabeledAdjacencyList(g, dot.EdgeLabel(lf), dot.Indent(""))
	fmt.Println(s)
	// Output:
	// digraph {
	// 0 -> 2 [label = "0.33"]
	// 4 -> 2 [label = "1.7"]
	// 4 -> 3 [label = "2e+117"]
	// }
}

func ExampleStringFromList() {
	//     0
	//    / \
	//   /   2
	//  /     \
	// 1       3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
	}}
	f.Leaves.SetBit(&f.Leaves, 1, 1)
	f.Leaves.SetBit(&f.Leaves, 3, 1)
	s, _ := dot.StringFromList(f, dot.Indent(""))
	fmt.Println(s)
	// Output:
	// digraph {
	// rankdir = BT
	// 1 -> 0
	// 2 -> 0
	// 3 -> 2
	// {rank = same 1 3}
	// }
}

func ExampleStringWeightedEdgeList() {
	//              (label 0, wt 1.6)
	//          0----------------------2
	// (label 1 |                     /
	//  wt .33) |  ------------------/
	//          | / (label 2, wt 1.7)
	//          |/
	//          1
	weights := []float64{
		0: 1.6,
		1: .33,
		2: 1.7,
	}
	g := graph.WeightedEdgeList{
		WeightFunc: func(l graph.LI) float64 { return weights[int(l)] },
		Order:      3,
		Edges: []graph.LabeledEdge{
			{graph.Edge{0, 1}, 1},
			{graph.Edge{1, 0}, 1},

			{graph.Edge{0, 2}, 0},
			{graph.Edge{2, 0}, 0},

			{graph.Edge{1, 2}, 2},
			{graph.Edge{2, 1}, 2},
		},
	}
	s, _ := dot.StringWeightedEdgeList(g, dot.Indent(""))
	fmt.Println(s)
	// Output:
	// graph {
	// 0 -- 1 [label = "0.33"]
	// 0 -- 2 [label = "1.6"]
	// 1 -- 2 [label = "1.7"]
	// }
}

func ExampleUndirectArcs() {
	//              (label 0, wt 1.6)
	//          0----------------------2
	// (label 1 |                     /
	//  wt .33) |  ------------------/
	//          | / (label 2, wt 1.7)
	//          |/
	//          1
	weights := []float64{
		0: 1.6,
		1: .33,
		2: 1.7,
	}
	g := graph.WeightedEdgeList{
		WeightFunc: func(l graph.LI) float64 { return weights[int(l)] },
		Order:      3,
		Edges: []graph.LabeledEdge{
			{graph.Edge{0, 1}, 1},
			{graph.Edge{0, 2}, 0},
			{graph.Edge{1, 2}, 2},
		},
	}
	s, _ := dot.StringWeightedEdgeList(g,
		dot.UndirectArcs(true),
		dot.Indent(""))
	fmt.Println(s)
	// Output:
	// graph {
	// 0 -- 1 [label = "0.33"]
	// 0 -- 2 [label = "1.6"]
	// 1 -- 2 [label = "1.7"]
	// }
}
