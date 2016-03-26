// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleUndirected_Bipartite() {
	// 0 1 2
	//  \|/|
	//   3 4
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	b, c1, c2, _ := g.Bipartite(0)
	if b {
		fmt.Println("n:  43210")
		fmt.Printf("c1: %05b\n", c1)
		fmt.Printf("c2: %05b\n", c2)
	}
	// Output:
	// n:  43210
	// c1: 00111
	// c2: 11000
}

func ExampleUndirected_Bipartite_oddCycle() {
	// 0 1  2
	//  \|/ |
	//   3--4
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	b, _, _, oc := g.Bipartite(0)
	if !b {
		fmt.Println("odd cycle:", oc)
	}
	// Output:
	// odd cycle: [3 4 2]
}

func ExampleUndirected_BronKerbosch1() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch1(func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BKPivotMaxDegree() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BKPivotMinP() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch2(g.BKPivotMinP, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BronKerbosch2() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleUndirected_BronKerbosch3() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.Undirected
	g.AddEdge(0, 4)
	g.AddEdge(4, 5)
	g.AddEdge(4, 3)
	g.AddEdge(3, 2)
	g.AddEdge(5, 2)
	g.AddEdge(5, 1)
	g.AddEdge(2, 1)
	g.BronKerbosch3(g.BKPivotMaxDegree, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [0 4]
	// [3 4]
	// [4 5]
	// [2 3]
	// [1 2 5]
}

func ExampleUndirected_ConnectedComponentBits() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	f := g.ConnectedComponentBits()
	fmt.Println("o  543210")
	fmt.Println("-  ------")
	for o, b := f(); o > 0; o, b = f() {
		fmt.Printf("%d  %0*b\n", o, len(g.AdjacencyList), &b)
	}
	// Output:
	// o  543210
	// -  ------
	// 3  011001
	// 2  100010
	// 1  000100
}

func ExampleUndirected_ConnectedComponentLists() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	f := g.ConnectedComponentLists()
	for l := f(); l != nil; l = f() {
		fmt.Println(l)
	}
	// Output:
	// [0 3 4]
	// [1 5]
	// [2]
}

func ExampleUndirected_ConnectedComponentReps() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	fmt.Println(g.ConnectedComponentReps())
	// Output:
	// [0 1 2] [3 2 1]
}

func ExampleUndirected_ConnectedComponentReps_collectingBits() {
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	rep, order := g.ConnectedComponentReps()
	fmt.Println("543210  rep  order")
	fmt.Println("------  ---  -----")
	for i, r := range rep {
		var bits graph.Bits
		g.DepthFirst(r, &bits, nil)
		fmt.Printf("%0*b   %d     %d\n",
			len(g.AdjacencyList), &bits, r, order[i])
	}
	// Output:
	// 543210  rep  order
	// ------  ---  -----
	// 011001   0     3
	// 100010   1     2
	// 000100   2     1
}

func ExampleUndirected_ConnectedComponentReps_collectingLists() {
	var g graph.Undirected
	g.AddEdge(0, 3)
	g.AddEdge(0, 4)
	g.AddEdge(3, 4)
	g.AddEdge(1, 5)
	rep, _ := g.ConnectedComponentReps()
	for _, r := range rep {
		var m []graph.NI
		g.DepthFirst(r, nil, func(n graph.NI) bool {
			m = append(m, n)
			return true
		})
		fmt.Println(m)
	}
	// Output:
	// [0 3 4]
	// [1 5]
	// [2]
}

func ExampleUndirected_Degeneracy() {
	//
	//   /---2
	//  0--./|\
	//  |\ /\| \
	//  | .  6  3--5
	//  |/ \/| /
	//  1--'\|/
	//   \---4
	//
	// Same graph redrawn to show ordering:
	//
	//          /-----\
	//         /  /--\ \
	//      /-+--+-\  \|
	//  5--3--4--6--2--1--0
	//         \  \  \---/|\
	//          \  \-----/ |
	//           \--------/
	//
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(0, 4)
	g.AddEdge(0, 6)
	g.AddEdge(1, 2)
	g.AddEdge(1, 4)
	g.AddEdge(1, 6)
	g.AddEdge(6, 2)
	g.AddEdge(6, 4)
	g.AddEdge(3, 2)
	g.AddEdge(3, 4)
	g.AddEdge(3, 5)
	k, ord, cores := g.Degeneracy()
	fmt.Println("Degeneracy:", k)
	fmt.Println("Ordering:", ord)
	fmt.Println("0-core:", ord[:cores[0]])
	for k := 1; k < len(cores); k++ {
		fmt.Printf("%d-core: %d\n", k, ord[cores[k-1]:cores[k]])
	}
	// Output:
	// Degeneracy: 3
	// Ordering: [5 3 4 6 2 1 0]
	// 0-core: []
	// 1-core: [5]
	// 2-core: [3]
	// 3-core: [4 6 2 1 0]
}

func ExampleUndirected_Degree() {
	// 0---1--\
	//      \-/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(1, 1)
	fmt.Println(g.Degree(0))
	fmt.Println(g.Degree(1))
	// Output:
	// 1
	// 3
}

func ExampleUndirected_FromList() {
	//    4   3
	//   / \
	//  2   1
	//       \
	//        0
	var g graph.Undirected
	g.AddEdge(2, 4)
	g.AddEdge(4, 1)
	g.AddEdge(1, 0)
	f, cycle := g.FromList(4)
	if cycle >= 0 {
		return
	}
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %3d\n", n, e.From, e.Len)
	}
	fmt.Println("MaxLen:  ", f.MaxLen)
	// Output:
	// n  from  path len
	// 0    1    3
	// 1    4    2
	// 2    4    2
	// 3   -1    0
	// 4   -1    1
	// MaxLen:   3
}

func ExampleUndirected_FromList_cycle() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3---4
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 2
}

func ExampleUndirected_FromList_loop() {
	//    0
	//   / \ /-\
	//  1   2--/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 2)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 2
}

func ExampleUndirected_FromList_loopDisconnected() {
	//    0
	//   /   /-\
	//  1   2--/
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(2, 2)
	f, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: -1
	fmt.Println("n  from  path len")
	for n, e := range f.Paths {
		fmt.Printf("%d  %3d  %3d\n", n, e.From, e.Len)
	}
	// Output:
	// cycle: -1
	// n  from  path len
	// 0   -1    1
	// 1    0    2
	// 2   -1    0
}

func ExampleUndirected_FromList_multigraph() {
	//    0
	//   / \
	//  1   2==3
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(2, 3)
	g.AddEdge(2, 3)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 3
}

func ExampleUndirected_IsConnected() {
	// undirected graph:
	//   0
	//  / \
	// 1   2
	var g graph.Undirected
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	u, _, _ := g.IsUndirected()
	fmt.Println("undirected:", u)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// undirected: true
	// connected:  true
}

func ExampleUndirected_IsConnected_notConnected() {
	// undirected graph:
	//   0   1
	//  / \
	// 2   3
	var g graph.Undirected
	g.AddEdge(0, 2)
	g.AddEdge(0, 3)
	u, _, _ := g.IsUndirected()
	fmt.Println("undirected:", u)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// undirected: true
	// connected:  false
}
func ExampleUndirected_IsTree() {
	//  0--\
	//  |  |
	//  \--/   1   3
	//        /   / \
	//       2   4---5
	var g graph.Undirected
	g.AddEdge(0, 0)
	g.AddEdge(1, 2)
	g.AddEdge(3, 4)
	g.AddEdge(3, 5)
	g.AddEdge(4, 5)
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	fmt.Println(g.IsTree(3))
	// Output:
	// false false
	// true false
	// false false
}

func ExampleUndirected_Size() {
	//   0--\
	//  / \-/
	// 1
	var g graph.Undirected
	g.AddEdge(0, 0)
	g.AddEdge(0, 1)
	fmt.Println("Size:", g.Size())
	fmt.Printf("(Arc size = %d)\n", g.ArcSize())
	// Output:
	// Size: 2
	// (Arc size = 3)
}
