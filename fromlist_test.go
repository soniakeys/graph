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

func ExampleFromList_CommonAncestor() {
	//       4
	//      /
	//     1
	//    / \
	//   0   2
	//  /
	// 3
	t := &graph.FromList{Paths: []graph.PathEnd{
		4: {From: -1, Len: 1},
		1: {From: 4, Len: 2},
		0: {From: 1, Len: 3},
		2: {From: 1, Len: 3},
		3: {From: 0, Len: 4},
	}}
	fmt.Println(t.CommonAncestor(2, 3))
	// Output:
	// 1
}

func ExampleFromList_Labeled() {
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
	for fr, to := range f.Labeled(nil) {
		fmt.Println(fr, to)
	}
	// Output:
	// 0 [{1 1} {2 2}]
	// 1 []
	// 2 [{3 3}]
	// 3 []
}

func ExampleFromList_Labeled_indexed() {
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
	for fr, to := range f.Labeled(labels) {
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
	t.Leaves.SetBit(&t.Leaves, 0, 1)
	t.Leaves.SetBit(&t.Leaves, 2, 1)
	t.Leaves.SetBit(&t.Leaves, 3, 1)
	// end at non-leaf, let PathEnd allocate result
	fmt.Println(t.PathTo(1, nil))
	fmt.Println()
	// preallocate buffer, enumerate paths to all leaves
	p := make([]graph.NI, t.MaxLen)
	for n := range t.Paths {
		if t.Leaves.Bit(n) == 1 {
			fmt.Println(t.PathTo(graph.NI(n), p))
		}
	}
	// Output:
	// [4 1]
	//
	// [4 1 0]
	// [4 1 2]
	// [3]
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
		fmt.Println(n, "      ", f.Leaves.Bit(n))
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
	for n, fr := range g {
		fmt.Println(n, fr)
	}
	// Output:
	// 0 [1 2]
	// 1 []
	// 2 []
	// 3 []
}

func ExampleFromList_Undirected() {
	//    0   3
	//   / \
	//  1   2
	t := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: -1},
	}}
	g := t.Undirected()
	for n, fr := range g {
		fmt.Println(n, fr)
	}
	ud, _, _ := g.IsUndirected()
	fmt.Println(ud)
	// Output:
	// 0 [1 2]
	// 1 [0]
	// 2 [0]
	// 3 []
	// true
}

func ExampleFromList_UndirectedLabeled() {
	//    0   3
	//   / \
	//  1   2
	t := graph.FromList{Paths: []graph.PathEnd{
		0: {From: -1},
		1: {From: 0},
		2: {From: 0},
		3: {From: -1},
	}}
	g := t.UndirectedLabeled(nil)
	for n, fr := range g {
		fmt.Printf("%d %#v\n", n, fr)
	}
	// Output:
	// 0 []graph.Half{graph.Half{To:1, Label:1}, graph.Half{To:2, Label:2}}
	// 1 []graph.Half{graph.Half{To:0, Label:1}}
	// 2 []graph.Half{graph.Half{To:0, Label:2}}
	// 3 []graph.Half(nil)
}
