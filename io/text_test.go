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
	g, err := io.ReadAdjacencyList(r)
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

/* func ExampleReadAdjacencyListOrder() {
	r := bytes.NewBufferString("2 1 1\n\n1")
	g, err := io.ReadAdjacencyListOrder(r, 3)
	for n, to := range g {
		fmt.Println(n, to)
	}
	fmt.Println("err: ", err)
	// Output:
	// 0 [2 1 1]
	// 1 []
	// 2 [1]
	// err:  <nil>
} */

func ExampleReadAdjacencyListNIs() {
	r := bytes.NewBufferString(`
0: 2 1 1
2: 1`)
	g, err := io.ReadAdjacencyListNIs(r, "")
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

func ExampleReadAdjacencyListNIsBase() {
	r := bytes.NewBufferString(`
1 2 3 // no 0
4 5
`)
	g, err := io.ReadAdjacencyListNIsBase(r, "//", 10)
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

func ExampleReadAdjacencyListNames() {
	r := bytes.NewBufferString(`
a b c  # source target target
d e 
`)
	g, names, m, err := io.ReadAdjacencyListNames(r, "", "", "#")
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

func ExampleReadArcNIs() {
	r := bytes.NewBufferString(`
0 2  
0 1
0 1 // parallel
2 1
`)
	g, err := io.ReadArcNIs(r, "//")
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

func ExampleReadArcNames() {
	r := bytes.NewBufferString(`
a b
a b // parallel
a c  
c b
`)
	g, names, m, err := io.ReadArcNames(r, "", "//")
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

func ExampleWriteAdjacencyListNIs() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.WriteAdjacencyListNIs(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0: 2 1 1
	// 2: 1
	// bytes: 14, err: <nil>
}

func ExampleWriteAdjacencyListNames() {
	//   a   d
	//  / \   \
	// b   c   e
	g := graph.AdjacencyList{
		0: {1, 2},
		3: {4},
		4: {},
	}
	names := []string{"a", "b", "c", "d", "e"}
	n, err := io.WriteAdjacencyListNames(g, os.Stdout, ": ", " ",
		func(n graph.NI) string { return names[n] })
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a: b c
	// d: e
	// bytes: 12, err: <nil>
}

func ExampleWriteArcNIs() {
	//   0
	//  / \\
	// 2-->1
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	n, err := io.WriteArcNIs(g, os.Stdout)
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// 0 2
	// 0 1
	// 0 1
	// 2 1
	// bytes: 16, err: <nil>
}

func ExampleWriteArcNames() {
	//   a
	//  / \\
	// c-->b
	g := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	names := []string{"a", "b", "c"}
	n, err := io.WriteArcNames(g, os.Stdout, " ",
		func(n graph.NI) string { return names[n] })
	fmt.Printf("bytes: %d, err: %v\n", n, err)
	// Output:
	// a c
	// a b
	// a b
	// c b
	// bytes: 16, err: <nil>
}
