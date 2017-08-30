// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

// undir_cg_test.go -- tests for code in undir_cg.go.
//
// These are tests on the labeled versions of methods.
//
// See also undir_ro_test.go when editing this file.  Try to keep the tests
// in the two files as similar as possible.

import (
	"fmt"
	"log"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

func ExampleLabeledUndirected_Bipartite() {
	// 0 1 2
	//  \|/|
	//   3 4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	b, c1, c2, _ := g.Bipartite(0)
	if b {
		fmt.Println("n:  43210")
		fmt.Println("c1:", c1)
		fmt.Println("c2:", c2)
	}
	// Output:
	// n:  43210
	// c1: 00111
	// c2: 11000
}

func ExampleLabeledUndirected_Bipartite_oddCycle() {
	// 0 1  2
	//  \|/ |
	//   3--4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	b, _, _, oc := g.Bipartite(0)
	if !b {
		fmt.Println("odd cycle:", oc)
	}
	// Output:
	// odd cycle: [3 4 2]
}

func ExampleLabeledUndirected_BronKerbosch1() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch1(func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BKPivotMaxDegree() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BKPivotMinP() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch2(g.BKPivotMinP, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [1 2 5]
	// [2 3]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BronKerbosch2() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch2(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [0 4]
	// [2 3]
	// [1 2 5]
	// [3 4]
	// [4 5]
}

func ExampleLabeledUndirected_BronKerbosch3() {
	// 0--4--5-
	//    |  | \
	//    3--2--1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 3}, 0)
	g.AddEdge(graph.Edge{3, 2}, 0)
	g.AddEdge(graph.Edge{5, 2}, 0)
	g.AddEdge(graph.Edge{5, 1}, 0)
	g.AddEdge(graph.Edge{2, 1}, 0)
	g.BronKerbosch3(g.BKPivotMaxDegree, func(c bits.Bits) bool {
		fmt.Println(c.Slice())
		return true
	})
	// Output:
	// [1 2 5]
	// [4 5]
	// [2 3]
	// [3 4]
	// [0 4]
}

func ExampleLabeledUndirected_ConnectedComponentBits() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	f := g.ConnectedComponentBits()
	fmt.Println("n  ma  543210")
	fmt.Println("-  --  ------")
	for {
		n, ma, b := f()
		if n == 0 {
			break
		}
		fmt.Printf("%d  %2d  %s\n", n, ma, b)
	}
	// Output:
	// n  ma  543210
	// -  --  ------
	// 3   6  011001
	// 2   2  100010
	// 1   0  000100
}

func ExampleLabeledUndirected_ConnectedComponentLists() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	f := g.ConnectedComponentLists()
	for {
		l, ma := f()
		if l == nil {
			break
		}
		fmt.Println(l, ma)
	}
	// Output:
	// [0 3 4] 6
	// [1 5] 2
	// [2] 0
}

func ExampleLabeledUndirected_ConnectedComponentReps() {
	//    0   1   2
	//   / \   \
	//  3---4   5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{0, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{1, 5}, 0)
	reps, orders, arcSizes := g.ConnectedComponentReps()
	fmt.Println("reps:    ", reps)
	fmt.Println("orders:  ", orders)
	fmt.Println("arcSizes:", arcSizes)
	// Output:
	// reps:     [0 1 2]
	// orders:   [3 2 1]
	// arcSizes: [6 2 0]
}

func ExampleLabeledUndirected_Degeneracy() {
	//   1   ----5
	//  / \ /   / \
	// 0---2---4  |
	//      \   \ /
	//   3   ----6
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{2, 5}, 0)
	g.AddEdge(graph.Edge{2, 6}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	g.AddEdge(graph.Edge{4, 6}, 0)
	g.AddEdge(graph.Edge{5, 6}, 0)
	ord, breaks := g.DegeneracyOrdering()
	fmt.Println("Degeneracy:", len(breaks)-1)
	fmt.Println("k-breaks:", breaks)
	fmt.Println("Ordering:", ord)
	for k, x := range breaks {
		fmt.Printf("nodes of %d-core(s): %d\n", k, ord[:x])
	}
	// Output:
	// Degeneracy: 3
	// k-breaks: [7 6 6 4]
	// Ordering: [4 5 6 2 0 1 3]
	// nodes of 0-core(s): [4 5 6 2 0 1 3]
	// nodes of 1-core(s): [4 5 6 2 0 1]
	// nodes of 2-core(s): [4 5 6 2 0 1]
	// nodes of 3-core(s): [4 5 6 2]

}

func ExampleLabeledUndirected_Degree() {
	// 0---1--\
	//      \-/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{1, 1}, 0)
	fmt.Println(g.Degree(0))
	fmt.Println(g.Degree(1))
	// Output:
	// 1
	// 3
}

func ExampleLabeledUndirected_Density() {
	// 0---1
	// |
	// 2---3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	fmt.Println(g.Density())
	// Output:
	// 0.5
}

func ExampleLabeledUndirected_Eulerian() {
	//   0--
	//  /   \
	//  \    \
	//   1----3
	//  / \  /
	//  \ / /
	//   2--
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	fmt.Println(g.Eulerian())
	// Output:
	// 2 3 <nil>
}

func ExampleLabeledUndirected_EulerianCycle() {
	//    0---
	// a /    \ b
	//   \  c  \
	//    1-----3
	// d / \  f/ \ g
	//   \ e\ /  /
	//    ---2---
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 3}, 'b')
	g.AddEdge(graph.Edge{1, 2}, 'd')
	g.AddEdge(graph.Edge{1, 2}, 'e')
	g.AddEdge(graph.Edge{1, 3}, 'c')
	g.AddEdge(graph.Edge{2, 3}, 'f')
	g.AddEdge(graph.Edge{2, 3}, 'g')
	c, err := g.EulerianCycle()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{0 -1} {1 97} {3 99} {2 103} {1 100} {2 101} {3 102} {0 98}]
	//
	// 0 --a-- 1 --c-- 3 --g-- 2 --d-- 1 --e-- 2 --f-- 3 --b-- 0
}

func ExampleLabeledUndirected_EulerianCycleD() {
	//    a
	// 0-------1
	//  \    c/|\
	//   \   /d| \
	//  b \  \ | /e
	//     \  \|/
	//      ---2--\
	//          \-/f
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 2}, 'b')
	g.AddEdge(graph.Edge{1, 2}, 'c')
	g.AddEdge(graph.Edge{1, 2}, 'd')
	g.AddEdge(graph.Edge{1, 2}, 'e')
	g.AddEdge(graph.Edge{2, 2}, 'f') // 6 edges total
	c, err := g.EulerianCycleD(6)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{0 -1} {1 97} {2 101} {1 99} {2 100} {2 102} {0 98}]
	//
	// 0 --a-- 1 --e-- 2 --c-- 1 --d-- 2 --f-- 2 --b-- 0
}

func ExampleLabeledUndirected_EulerianPath() {
	//    0
	//  a/|\
	//  /b| \
	//  \ | /c
	//   \|/
	//    1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 1}, 'b')
	g.AddEdge(graph.Edge{0, 1}, 'c')
	c, err := g.EulerianPath()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{0 -1} {1 97} {0 99} {1 98}]
	//
	// 0 --a-- 1 --c-- 0 --b-- 1
}

func ExampleLabeledUndirected_EulerianPathD() {
	//    0
	//  a/|\
	//  /b| \
	//  \ | /c
	//   \|/
	//    1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 'a')
	g.AddEdge(graph.Edge{0, 1}, 'b')
	g.AddEdge(graph.Edge{0, 1}, 'c')
	c, err := g.EulerianPathD(3, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	// prettier,
	fmt.Print("\n", c[0].To)
	for _, to := range c[1:] {
		fmt.Printf(" --%c-- %d", to.Label, to.To)
	}
	fmt.Println()
	// Output:
	// [{0 -1} {1 97} {0 99} {1 98}]
	//
	// 0 --a-- 1 --c-- 0 --b-- 1
}

func ExampleLabeledUndirected_EulerianStart() {
	//   0--
	//  /   \
	//  \    \
	//   1----3
	//  / \  /
	//  \ / /
	//   2--
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{0, 3}, 0)
	g.AddEdge(graph.Edge{1, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	fmt.Println(g.EulerianStart())
	// Output:
	// 2
}

func ExampleLabeledUndirected_FromList() {
	//    4   3
	//   / \
	//  2   1
	//       \
	//        0
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{4, 1}, 0)
	g.AddEdge(graph.Edge{1, 0}, 0)
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

func ExampleLabeledUndirected_FromList_cycle() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3---4
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 4}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 2
}

func ExampleLabeledUndirected_FromList_loop() {
	//    0
	//   / \ /-\
	//  1   2--/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 2}, 0)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 2
}

func ExampleLabeledUndirected_FromList_loopDisconnected() {
	//    0
	//   /   /-\
	//  1   2--/
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{2, 2}, 0)
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

func ExampleLabeledUndirected_FromList_multigraph() {
	//    0
	//   / \
	//  1   2==3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	g.AddEdge(graph.Edge{2, 3}, 0)
	_, cycle := g.FromList(0)
	fmt.Println("cycle:", cycle)
	// Output:
	// cycle: 3
}

func ExampleLabeledUndirected_IsConnected() {
	// undirected graph:
	//   0
	//  / \
	// 1   2
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 1}, 0)
	g.AddEdge(graph.Edge{0, 2}, 0)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  true
}

func ExampleLabeledUndirected_IsConnected_notConnected() {
	// undirected graph:
	//   0   1
	//  / \
	// 2   3
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 2}, 0)
	g.AddEdge(graph.Edge{0, 3}, 0)
	fmt.Println("connected: ", g.IsConnected())
	// Output:
	// connected:  false
}
func ExampleLabeledUndirected_IsTree() {
	//  0--\
	//  |  |
	//  \--/   1   3
	//        /   / \
	//       2   4---5
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 0)
	g.AddEdge(graph.Edge{1, 2}, 0)
	g.AddEdge(graph.Edge{3, 4}, 0)
	g.AddEdge(graph.Edge{3, 5}, 0)
	g.AddEdge(graph.Edge{4, 5}, 0)
	fmt.Println(g.IsTree(0))
	fmt.Println(g.IsTree(1))
	fmt.Println(g.IsTree(3))
	// Output:
	// false false
	// true false
	// false false
}

func ExampleLabeledUndirected_Size() {
	//   0--\
	//  / \-/
	// 1
	var g graph.LabeledUndirected
	g.AddEdge(graph.Edge{0, 0}, 0)
	g.AddEdge(graph.Edge{0, 1}, 0)
	fmt.Println("Size:", g.Size())
	fmt.Printf("(Arc size = %d)\n", g.ArcSize())
	// Output:
	// Size: 2
	// (Arc size = 3)
}
