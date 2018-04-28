// Copyright 2018 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package io_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/io"
)

func ExampleArcDir() {
	//   0
	//  / \\
	// 1   2--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(0, 2)
	g.AddEdge(2, 2)

	fmt.Println("Default WriteArcs All:  writes all arcs:")
	t := io.Text{}
	n, err := t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)

	fmt.Println("Upper triangle:")
	t.WriteArcs = io.Upper
	n, err = t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)

	fmt.Println("Lower triangle:")
	t.WriteArcs = io.Lower
	n, err = t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// Default WriteArcs All:  writes all arcs:
	// 0: 1 2 2
	// 1: 0
	// 2: 0 0 2
	// bytes: 23, err: <nil>
	//
	// Upper triangle:
	// 0: 1 2 2
	// 2: 2
	//bytes: 14, err: <nil>
	//
	// Lower triangle:
	// 1: 0
	// 2: 0 0 2
	// bytes: 14, err: <nil>
}

func ExampleFormat() {
	//   0
	//  / \\
	// 2-->3
	g := graph.AdjacencyList{
		0: {2, 3, 3},
		2: {3},
		3: {},
	}
	fmt.Println("Default Format Sparse:")
	var t io.Text
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)

	fmt.Println("Format Dense:")
	t.Format = io.Dense
	n, err = t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)

	fmt.Println("Format Arcs:")
	t.Format = io.Arcs
	n, err = t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)
	// Output:
	// Default Format Sparse:
	// 0: 2 3 3
	// 2: 3
	// 3:
	// bytes: 17, err: <nil>
	//
	// Format Dense:
	// 2 3 3
	//
	// 3
	//
	// bytes: 10, err: <nil>
	//
	// Format Arcs:
	// 0 2
	// 0 3
	// 0 3
	// 2 3
	// bytes: 16, err: <nil>
}

func ExampleNewText() {
	r := bytes.NewBufferString(`
0: 1 // node 0 comment
`)
	g, _, _, err := io.NewText().ReadAdjacencyList(r)
	for fr, to := range g {
		if len(to) > 0 {
			fmt.Println(fr, to)
		}
	}
	fmt.Println("err:", err)
	// Output:
	// 0 [1]
	// err: <nil>
}

// Example with zero value Text.  For options see examples under Text.
func ExampleText_ReadAdjacencyList() {
	// set r to output of example Text.WriteAdjacencyList:
	r := bytes.NewBufferString(`
0: 2 3 3
2: 3
3:
`)
	g, _, _, err := io.Text{}.ReadAdjacencyList(r)
	// result g should be:
	//   0
	//  / \\
	// 2-->3
	for fr, to := range g {
		if len(to) > 0 {
			fmt.Println(fr, to)
		}
	}
	fmt.Println("err:", err)
	// Output:
	// 0 [2 3 3]
	// 2 [3]
	// err: <nil>
}

// Example with zero value Text.  For options see examples under Text.
func ExampleText_WriteAdjacencyList() {
	//   0
	//  / \\
	// 2-->3
	g := graph.AdjacencyList{
		0: {2, 3, 3},
		2: {3},
		3: {},
	}
	n, err := io.Text{}.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)
	// Output:
	// 0: 2 3 3
	// 2: 3
	// 3:
	// bytes: 17, err: <nil>
}

func ExampleText_base() {
	const zz = 36*36 - 1
	g := graph.AdjacencyList{
		zz: {0},
	}
	n, err := io.Text{Base: 36}.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)
	// Output:
	// zz: 0
	// bytes: 6, err: <nil>
}

func ExampleText_delimiters() {
	//   a   d
	//  / \   \
	// b   c   e
	name := []string{"a", "b", "c", "d", "e"}
	g := graph.AdjacencyList{
		0: {1, 2},
		3: {4},
		4: {},
	}
	t := io.Text{
		NodeName: func(n graph.NI) string { return name[n] },
		FrDelim:  "->",
		ToDelim:  ":",
	}
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)
	// Output:
	// a->b:c
	// d->e
	// bytes: 12, err: <nil>
}

func ExampleText_mapNames() {
	//   a   d
	//  / \   \
	// b   c   e
	r := bytes.NewBufferString(`
a b c  # source target target
d e   
`)
	// For reading, default blank delimiter fields enable
	// delimiting by whitespace.
	t := io.Text{MapNames: true, Comment: "#"}
	g, names, m, err := t.ReadAdjacencyList(r)

	fmt.Println("names:")
	for i, n := range names {
		fmt.Println(i, n)
	}
	fmt.Println("graph:")
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println(graph.OrderMap(m))
	fmt.Println("err:", err)
	// Output:
	// names:
	// 0 a
	// 1 b
	// 2 c
	// 3 d
	// 4 e
	// graph:
	// 0 [1 2]
	// 1 []
	// 2 []
	// 3 [4]
	// 4 []
	// map[a:0 b:1 c:2 d:3 e:4]
	// err: <nil>
}

func ExampleText_nodeName() {
	//   a   d
	//  / \   \
	// b   c   e
	name := []string{"a", "b", "c", "d", "e"}
	g := graph.AdjacencyList{
		0: {1, 2},
		3: {4},
		4: {},
	}
	t := io.Text{NodeName: func(n graph.NI) string { return name[n] }}
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)
	// Output:
	// a: b c
	// d: e
	// bytes: 12, err: <nil>
}

/*
func ExampleText_WriteLabeledAdjacencyList() {
	//     0
	//    / \\
	// a /  b\\c
	//  /     \\
	// 2------->3
	//     d
	g := graph.LabeledAdjacencyList{
		0: {{2, 'a'}, {3, 'b'}, {3, 'c'}},
		2: {{3, 'd'}},
		3: {},
	}
	n, err := io.Text{}.WriteLabeledAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n\n", n, err)
	// Output:
	// 0: (2 96) (3 97) (3 98)
	// 2: (3 99)
	// 3:
	// bytes: 17, err: <nil>
}
*/
