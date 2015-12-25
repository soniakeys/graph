// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleFromTree_CommonAncestor() {
	// tree:
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

func ExampleFromList_Undirected() {
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
	// Example:
	// 0 [1 2]
	// 1 [0]
	// 2 [0]
	// 3 []
}

func ExampleFromList_Transpose() {
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
	// Example:
	// 0 [1 2]
	// 1 []
	// 2 []
	// 3 []
}
