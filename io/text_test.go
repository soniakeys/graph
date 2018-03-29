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

func ExampleText_ReadAdjacencyList_dense() {
	r := bytes.NewBufferString(`2 1 1

1`)
	g, _, _, err := io.Text{Format: io.Dense}.ReadAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [2 1 1]
	// 1 []
	// 2 [1]
	// err:  <nil>
}

func ExampleText_ReadAdjacencyList() {
	r := bytes.NewBufferString(`
0: 2 1 1
2: 1`)
	g, _, _, err := io.Text{}.ReadAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [2 1 1]
	// 1 []
	// 2 [1]
	// err:  <nil>
}

func ExampleText_ReadAdjacencyList_names() {
	r := bytes.NewBufferString(`
a b c  # source target target
d e
`)
	t := io.Text{ReadNames: true, Comment: "#"}
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
	// map[a:0 b:1 c:2 d:3 e:4 ]
	// err: <nil>
}

func ExampleText_ReadAdjacencyList_arcs() {
	r := bytes.NewBufferString(`
0 2  
0 1
0 1 // parallel
2 1
`)
	t := io.Text{Format: io.Arcs, Comment: "//"}
	g, _, _, err := t.ReadAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [2 1 1]
	// 1 []
	// 2 [1]
	// err:  <nil>
}

func ExampleText_ReadAdjacencyList_arcNames() {
	r := bytes.NewBufferString(`
a b
a b // parallel
a c  
c b
`)
	t := io.Text{Format: io.Arcs, ReadNames: true, Comment: "//"}
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
	fmt.Println("err: ", err)
	// Output:
	// names:
	// 0 a
	// 1 b
	// 2 c
	// graph:
	// 0 [1 1 2]
	// 1 []
	// 2 [1]
	// map[a:0 b:1 c:2 ]
	// err:  <nil>
}

func ExampleText_ReadLabeledAdjacencyList() {
	r := bytes.NewBufferString(`2 101 1 102 1 102

1 103`)
	g, err := io.Text{}.ReadLabeledAdjacencyList(r)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [{2 101} {1 102} {1 102}]
	// 1 []
	// 2 [{1 103}]
	// err:  <nil>
}

func ExampleText_WriteAdjacencyList_dense() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.Text{Format: io.Dense}.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 2 1 1
	//
	// 1
	// bytes: 9, err: <nil>
}

func ExampleText_WriteAdjacencyList() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.Text{}.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0: 2 1 1
	// 2: 1
	// bytes: 14, err: <nil>
}

func ExampleText_WriteAdjacencyList_names() {
	//   a   d
	//  / \   \
	// b   c   e
	g := graph.AdjacencyList{
		0: {1, 2},
		3: {4},
		4: {},
	}
	names := []string{"a", "b", "c", "d", "e"}
	t := io.Text{WriteName: func(n graph.NI) string { return names[n] }}
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a: b c
	// d: e
	// bytes: 12, err: <nil>
}

func ExampleText_WriteAdjacencyList_arcs() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.Text{Format: io.Arcs}.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0 2
	// 0 1
	// 0 1
	// 2 1
	// bytes: 16, err: <nil>
}

func ExampleText_WriteAdjacencyList_arcNames() {
	//   a
	//  / \\
	// c-->b
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	names := []string{"a", "b", "c"}
	t := io.Text{
		Format:    io.Arcs,
		WriteName: func(n graph.NI) string { return names[n] },
	}
	n, err := t.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a c
	// a b
	// a b
	// c b
	// bytes: 16, err: <nil>
}

func ExampleText_WriteLabeledAdjacencyList() {
	//        0
	// (101) / \\ (102)
	//      2-->1
	//      (103)
	g := graph.LabeledAdjacencyList{
		0: {{2, 101}, {1, 102}, {1, 102}},
		2: {{1, 103}},
	}
	n, err := io.Text{}.WriteLabeledAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 2 101 1 102 1 102
	//
	// 1 103
	// bytes: 25, err: <nil>
}

func ExampleText_WriteAdjacencyList_undirectedDense() {
	//   0
	//  / \\
	// 1---2--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(0, 2)
	g.AddEdge(2, 2)
	t := io.Text{Format: io.Dense, WriteArcs: io.Upper}
	n, err := t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 1 2 2
	//
	// 2
	// bytes: 9, err: <nil>
}

func ExampleText_WriteAdjacencyList_undirectedNames() {
	//   a
	//  / \\
	// b---c--\
	//      \-/
	var g graph.Undirected
	names := []string{"a", "b", "c"}
	ni := map[string]graph.NI{}
	for i, s := range names {
		ni[s] = graph.NI(i)
	}
	g.AddEdge(ni["a"], ni["b"])
	g.AddEdge(ni["a"], ni["c"])
	g.AddEdge(ni["a"], ni["c"])
	g.AddEdge(ni["c"], ni["c"])
	t := io.Text{
		WriteArcs: io.Upper,
		WriteName: func(n graph.NI) string { return names[n] },
	}
	n, err := t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a: b c c
	// c: c
	// bytes: 14, err: <nil>
}

func ExampleText_WriteAdjacencyList_undirected() {
	//   0
	//  / \\
	// 1---2--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(0, 2)
	g.AddEdge(2, 2)
	t := io.Text{WriteArcs: io.Upper}
	n, err := t.WriteAdjacencyList(g.AdjacencyList, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0: 1 2 2
	// 2: 2
	// bytes: 14, err: <nil>
}
