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

func ExampleReadAdjacencyList() {
	r := bytes.NewBufferString("2 1 1\n\n1")
	g, err := io.ReadAdjacencyList(r, 3)
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

func ExampleReadAdjacencyListKeyed() {
	r := bytes.NewBufferString(`
0: 2 1 1
2: 1`)
	g, err := io.ReadAdjacencyListKeyed(r, 0, "")
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

func ExampleReadAdjacencyListKeyedBase() {
	r := bytes.NewBufferString(`
1 2 3 // no 0
4 5
`)
	g, err := io.ReadAdjacencyListKeyedBase(r, 0, "//", 10)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 []
	// 1 [2 3]
	// 2 []
	// 3 []
	// 4 [5]
	// err:  <nil>
}

func ExampleReadAdjacencyListNamed() {
	r := bytes.NewBufferString(`
a b c  # source target target
d e 
`)
	g, names, m, err := io.ReadAdjacencyListNamed(r, "", "", "#")
	fmt.Println("names:")
	for i, n := range names {
		fmt.Println(i, n)
	}
	fmt.Println("graph:")
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("map:")
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
	// map:
	// map[a:0 b:1 c:2 d:3 e:4 ]
	// err: <nil>
}

func ExampleWriteAdjacencyList() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.WriteAdjacencyList(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 2 1 1
	//
	// 1
	// bytes: 9, err: <nil>
}

func ExampleWriteAdjacencyListKeyed() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.WriteAdjacencyListKeyed(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0: 2 1 1
	// 2: 1
	// bytes: 14, err: <nil>
}

func ExampleWriteAdjacencyListNamed() {
	//   a   d
	//  / \   \
	// b   c   e
	g := graph.AdjacencyList{
		0: {1, 2},
		3: {4},
		4: {},
	}
	names := []string{"a", "b", "c", "d", "e"}
	n, err := io.WriteAdjacencyListNamed(g, os.Stdout, ": ", " ",
		func(n graph.NI) string { return names[n] })
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a: b c
	// d: e
	// bytes: 12, err: <nil>
}
