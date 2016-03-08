package dot_test

import (
	"fmt"
	"os"

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

func ExampleWriteAdjacencyList() {
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
	dot.WriteAdjacencyList(g, os.Stdout, dot.Indent(""))
	// Output:
	// digraph {
	// 0 -> 2
	// 4 -> {2 3}
	// }
}
