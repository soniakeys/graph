// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleFromList_BoundsOk() {
	//   0
	//  / \
	// 1   2
	t := &graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
	}}
	ok, _ := t.BoundsOk()
	fmt.Println(ok)
	t = &graph.FromList{Paths: []graph.PathEnd{
		0: {From: 9},
		1: {From: 0},
		2: {From: 0},
	}}
	fmt.Println(t.BoundsOk())
	// Output:
	// true
	// false 0
}

func ExampleFromList_CommonStart() {
	//   4   5
	//  /   /
	// 6   1
	//    / \
	//   0   2
	//  /
	// 3
	t := &graph.FromList{Paths: []graph.PathEnd{
		4: {From: -1, Len: 1},
		6: {From: 4, Len: 2},
		5: {From: -1, Len: 1},
		1: {From: 5, Len: 2},
		0: {From: 1, Len: 3},
		2: {From: 1, Len: 3},
		3: {From: 0, Len: 4},
	}}
	fmt.Println(t.CommonStart(2, 3))
	fmt.Println(t.CommonStart(6, 3))
	// Output:
	// 1
	// -1
}

func ExampleFromList_Cyclic_acyclic() {
	//   0
	//  / \
	// 1   2
	//    /
	//   3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
	}}
	fmt.Println(f.Cyclic())
	// Output:
	// false -1
}

func ExampleFromList_Cyclic_cyclic() {
	//   0
	//  /
	// 1   2
	//    /| (other arcs directed up)
	//   3<-
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 3},
		3: {From: 2},
	}}
	fmt.Println(f.Cyclic())
	// Output:
	// true 2
}

func ExampleFromList_IsolatedNodes() {
	//   0  1
	//  / \
	// 2   3  4
	g := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: -1},
		2: {From: 0},
		3: {From: 0},
		4: {From: -1},
	}}
	fmt.Println(g.IsolatedNodes().Slice())
	// Output:
	// [1 4]
}

func ExamplePathTo() {
	//       4  3
	//      /
	//     1
	//    / \
	//   0   2
	p := []graph.PathEnd{
		4: {From: -1, Len: 1},
		3: {From: -1, Len: 1},
		1: {From: 4, Len: 2},
		0: {From: 1, Len: 3},
		2: {From: 1, Len: 3},
	}
	// end at non-leaf, let PathEnd allocate result
	fmt.Println(graph.PathTo(p, 1, nil))
	// Output:
	// [4 1]
}

func ExampleFromList_PathTo() {
	//       4  3
	//      /
	//     1
	//    / \
	//   0   2
	t := &graph.FromList{
		Paths: []graph.PathEnd{
			4: {From: -1, Len: 1},
			3: {From: -1, Len: 1},
			1: {From: 4, Len: 2},
			0: {From: 1, Len: 3},
			2: {From: 1, Len: 3},
		},
		MaxLen: 3,
	}
	lv := &t.Leaves
	lv.SetBit(0, 1)
	lv.SetBit(2, 1)
	lv.SetBit(3, 1)
	// end at non-leaf, let PathEnd allocate result
	fmt.Println(t.PathTo(1, nil))
	fmt.Println()
	// preallocate buffer, enumerate paths to all leaves
	p := make([]graph.NI, t.MaxLen)
	lv.Iterate(func(n graph.NI) bool {
		fmt.Println(t.PathTo(graph.NI(n), p))
		return true
	})
	// Output:
	// [4 1]
	//
	// [4 1 0]
	// [4 1 2]
	// [3]
}

func ExampleFromList_Preorder() {
	//     2
	//    / \
	//   0   3
	//  /
	// 1
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: 2},
		1: {From: 0},
		2: {From: -1},
		3: {From: 2},
	}}
	f.RecalcLeaves()
	f.Preorder(func(n graph.NI) bool {
		fmt.Println(n)
		return true
	})
	// Output:
	// 2
	// 0
	// 1
	// 3
}

func ExampleFromList_Preorder_forest() {
	//     4      0
	//    / \    / \
	//   5   3  2   6
	//  /
	// 1
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 5},
		2: {From: 0},
		3: {From: 4},
		4: {From: -1},
		5: {From: 4},
		6: {From: 0},
	}}
	f.RecalcLeaves()
	f.Preorder(func(n graph.NI) bool {
		fmt.Println(n)
		return true
	})
	// Output:
	// 4
	// 5
	// 1
	// 0
	// 2
	// 3
	// 6
}

func ExampleFromList_RecalcLeaves() {
	//   0
	//  / \
	// 1   2
	//      \
	//       3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
	}}
	f.RecalcLeaves()
	fmt.Println("Node  Leaf")
	for n := range f.Paths {
		fmt.Println(n, "      ", f.Leaves.Bit(graph.NI(n)))
	}
	// Output:
	// Node  Leaf
	// 0        0
	// 1        1
	// 2        0
	// 3        1
}

func ExampleFromList_RecalcLen() {
	//   0
	//  / \
	// 1   2
	//      \
	//       3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
	}}
	f.RecalcLeaves() // Leaves required for RecalcLen
	f.RecalcLen()
	fmt.Println("Node  Path len")
	for n := range f.Paths {
		fmt.Println(n, "          ", f.Paths[n].Len)
	}
	fmt.Println("MaxLen:", f.MaxLen)
	// Output:
	// Node  Path len
	// 0            1
	// 1            2
	// 2            2
	// 3            3
	// MaxLen: 3
}

func ExampleFromList_ReRoot() {
	//   0
	//  / \
	// 1   2
	//    / \
	//   3   4
	//  /
	// 5
	t := &graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
		4: {From: 2},
		5: {From: 3},
	}}
	t.ReRoot(3)
	for n, e := range t.Paths {
		fmt.Println(n, e)
	}
	// Output:
	// 0 {2 0}
	// 1 {0 0}
	// 2 {3 0}
	// 3 {-1 0}
	// 4 {2 0}
	// 5 {3 0}
}

func ExampleFromList_Root() {
	//       4
	//      /
	//     1
	//    / \
	//   2   0
	//  /
	// 3
	t := &graph.FromList{Paths: []graph.PathEnd{
		4: {From: -1},
		1: {From: 4},
		2: {From: 1},
		0: {From: 1},
		3: {From: 2},
	}}
	fmt.Println(t.Root(2))
	// Output:
	// 4
}

func ExampleFromList_Transpose() {
	//    0   3
	//   / \
	//  1   2
	t := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: -1},
	}}
	g := t.Transpose()
	fmt.Println("fr to")
	for n, fr := range g.AdjacencyList {
		fmt.Println(n, "", fr)
	}
	// Output:
	// fr to
	// 0  [1 2]
	// 1  []
	// 2  []
	// 3  []
}

func ExampleFromList_TransposeRoots() {
	//    0   3
	//   / \
	//  1   2
	t := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: -1},
	}}
	g, nr, roots := t.TransposeRoots()
	fmt.Println("fr to")
	for n, fr := range g.AdjacencyList {
		fmt.Println(n, "", fr)
	}
	fmt.Println(nr, "roots:")
	fmt.Println("3210")
	fmt.Printf("%04b\n", &roots)
	// Output:
	// fr to
	// 0  [1 2]
	// 1  []
	// 2  []
	// 3  []
	// 2 roots:
	// 3210
	// 1001
}

func ExampleFromList_TransposeLabeled() {
	//   0
	//  / \
	// 1   2
	//      \
	//       3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
	}}
	g := f.TransposeLabeled(nil)
	for fr, to := range g.LabeledAdjacencyList {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [{1 1} {2 2}]
	// 1 []
	// 2 [{3 3}]
	// 3 []
}

func ExampleFromList_TransposeLabeled_indexed() {
	//      0
	// 'A' / \ 'B'
	//    1   2
	//         \ 'C'
	//          3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
	}}
	labels := []graph.LI{
		1: 'A',
		2: 'B',
		3: 'C',
	}
	g := f.TransposeLabeled(labels)
	for fr, to := range g.LabeledAdjacencyList {
		fmt.Print(fr)
		for _, to := range to {
			fmt.Printf(" {%d %c}", to.To, to.Label)
		}
		fmt.Println()
	}
	// Output:
	// 0 {1 A} {2 B}
	// 1
	// 2 {3 C}
	// 3
}

func ExampleFromList_TransposeLabeledRoots() {
	//      0        4
	// 'A' / \ 'B'
	//    1   2
	//         \ 'C'
	//          3
	f := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: 2},
		4: {From: -1},
	}}
	labels := []graph.LI{
		1: 'A',
		2: 'B',
		3: 'C',
	}
	g, n, r := f.TransposeLabeledRoots(labels)
	for fr, to := range g.LabeledAdjacencyList {
		fmt.Print(fr)
		for _, to := range to {
			fmt.Printf(" {%d %c}", to.To, to.Label)
		}
		fmt.Println()
	}
	fmt.Println(n, "roots:", r.Slice())
	// Output:
	// 0 {1 A} {2 B}
	// 1
	// 2 {3 C}
	// 3
	// 4
	// 2 roots: [0 4]
}
