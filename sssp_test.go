// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExampleBreadthFirstPath() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	fmt.Println(g.BreadthFirstPath(1, 3))
	// Output:
	// [1 4 3]
}

func ExampleBreadthFirst_Path() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	b := graph.NewBreadthFirst(graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	})
	start, end := 1, 3
	if !b.Path(start, end) {
		return
	}
	fmt.Println("Path length:", b.Result.Paths[end].Len)
	fmt.Print("Backtrack to start: ", end)
	rp := b.Result.Paths
	for n := end; n != start; {
		n = rp[n].From
		fmt.Print(" ", n)
	}
	fmt.Println()
	// Output:
	// Path length: 3
	// Backtrack to start: 3 4 1
}

func ExampleBreadthFirst_AllPaths() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	b := graph.NewBreadthFirst(graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	})
	b.AllPaths(1)
	fmt.Println("Max path length:", b.Result.MaxLen)
	p := make([]int, b.Result.MaxLen)
	for n := range b.Graph {
		fmt.Println(n, b.Result.PathTo(n, p))
	}
	// Output:
	// Max path length: 4
	// 0 []
	// 1 [1]
	// 2 []
	// 3 [1 4 3]
	// 4 [1 4]
	// 5 [1 4 3 5]
	// 6 [1 4 6]
}

func ExampleBreadthFirst2Path() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	fmt.Println(g.BreadthFirst2Path(1, 3))
	// Output:
	// [1 4 3]
}

func ExampleBreadthFirst2_Path() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	t, m := g.Transpose()
	b := graph.NewBreadthFirst2(g, t, m)
	start, end := 1, 3
	if !b.Path(start, end) {
		return
	}
	fmt.Println("Path length:", b.Result.Paths[end].Len)
	fmt.Print("Backtrack to start: ", end)
	rp := b.Result.Paths
	for n := end; n != start; {
		n = rp[n].From
		fmt.Print(" ", n)
	}
	fmt.Println()
	// Output:
	// Path length: 3
	// Backtrack to start: 3 4 1
}

func ExampleBreadthFirst2_AllPaths() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	t, m := g.Transpose()
	b := graph.NewBreadthFirst2(g, t, m)
	b.AllPaths(1)
	fmt.Println("Max path length:", b.Result.MaxLen)
	p := make([]int, b.Result.MaxLen)
	for n := range b.To {
		fmt.Println(n, b.Result.PathTo(n, p))
	}
	// Output:
	// Max path length: 4
	// 0 []
	// 1 [1]
	// 2 []
	// 3 [1 4 3]
	// 4 [1 4]
	// 5 [1 4 3 5]
	// 6 [1 4 6]
}

func ExampleLabeledAdjacencyList_DijkstraPath() {
	// arcs are directed right:
	//          (wt: 11)
	//       --------------6----
	//      /             /     \
	//     /             /(2)    \(9)
	//    /     (9)     /         \
	//   1-------------3----       5
	//    \           /     \     /
	//     \     (10)/   (11)\   /(7)
	//   (7)\       /         \ /
	//       ------2-----------4
	//                 (15)
	g := graph.LabeledAdjacencyList{
		1: {{To: 2, Label: 7}, {To: 3, Label: 9}, {To: 6, Label: 11}},
		2: {{To: 3, Label: 10}, {To: 4, Label: 15}},
		3: {{To: 4, Label: 11}, {To: 6, Label: 2}},
		4: {{To: 5, Label: 7}},
		6: {{To: 5, Label: 9}},
	}
	w := func(label int) float64 { return float64(label) }
	p, l := g.DijkstraPath(1, 5, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [1 6 5]
	// Path length: 20
}

func ExampleDijkstra_Path() {
	// arcs are directed right:
	//          (wt: 11)
	//       --------------6----
	//      /             /     \
	//     /             /(2)    \(9)
	//    /     (9)     /         \
	//   1-------------3----       5
	//    \           /     \     /
	//     \     (10)/   (11)\   /(7)
	//   (7)\       /         \ /
	//       ------2-----------4
	//                 (15)
	g := graph.LabeledAdjacencyList{
		1: {{To: 2, Label: 7}, {To: 3, Label: 9}, {To: 6, Label: 11}},
		2: {{To: 3, Label: 10}, {To: 4, Label: 15}},
		3: {{To: 4, Label: 11}, {To: 6, Label: 2}},
		4: {{To: 5, Label: 7}},
		6: {{To: 5, Label: 9}},
	}
	w := func(label int) float64 { return float64(label) }
	d := graph.NewDijkstra(g, w)
	start, end := 1, 5
	if !d.Path(start, end) {
		return
	}
	fmt.Print("Backtrack to start: ", end)
	rp := d.Tree.Paths
	for n := end; n != start; {
		n = rp[n].From
		fmt.Print(" ", n)
	}
	fmt.Println()
	fmt.Println("Path distance:", d.Dist[end])
	// Output:
	// Backtrack to start: 5 6 1
	// Path distance: 20
}

func ExampleDijkstra_AllPaths() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label int) float64 { return float64(label) }
	d := graph.NewDijkstra(g, w)
	fmt.Println(d.AllPaths(2), "paths found.")
	// column len is from Result, and will be equal to len(path).
	// column dist is from Result, and will be equal to sum.
	fmt.Println("node:  path                  len  dist")
	p := make([]int, d.Tree.MaxLen)
	for nd := range g {
		r := &d.Tree.Paths[nd]
		path := d.Tree.PathTo(nd, p)
		if r.Len > 0 {
			fmt.Printf("%d:     %-23s %d    %2.0f\n",
				nd, fmt.Sprint(path), r.Len, d.Dist[nd])
		}
	}

	// Output:
	// 4 paths found.
	// node:  path                  len  dist
	// 2:     [2]                     1     0
	// 3:     [2 3]                   2    11
	// 4:     [2 3 4]                 3    17
	// 5:     [2 5]                   2     2
}

func ExampleLabeledAdjacencyList_AStarAPath() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label int) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	h := func(from int) float64 { return h4[from] }
	p, l := g.AStarAPath(0, 4, h, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [0 2 3 4]
	// Path length: 26
}

func ExampleAStarMPath() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label int) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	h := func(from int) float64 { return h4[from] }
	p, l := g.AStarMPath(0, 4, h, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Shortest path: [0 2 3 4]
	// Path length: 26
}

func ExampleHeuristic_Admissable() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label int) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	var h graph.Heuristic = func(from int) float64 { return h4[from] }
	fmt.Println(h.Admissable(g, w, 4))
	// Output:
	// true
}

func ExampleHeuristic_Monotonic() {
	// arcs are directed right:
	//       -----------------------
	//      /      (wt: 14)         \
	//     /                         \
	//    /     (9)           (2)     \
	//   0-------------2---------------5
	//    \           / \             /
	//     \     (10)/   \(11)    (9)/
	//   (7)\       /     \         /
	//       ------1-------3-------4
	//               (15)     (6)
	g := graph.LabeledAdjacencyList{
		0: {{To: 1, Label: 7}, {To: 2, Label: 9}, {To: 5, Label: 14}},
		1: {{To: 2, Label: 10}, {To: 3, Label: 15}},
		2: {{To: 3, Label: 11}, {To: 5, Label: 2}},
		3: {{To: 4, Label: 6}},
		4: {{To: 5, Label: 9}},
		5: {},
	}
	w := func(label int) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	var h graph.Heuristic = func(from int) float64 { return h4[from] }
	fmt.Println(h.Monotonic(g, w))
	// Output:
	// true
}

func ExampleBellmanFord() {
	//              /--------3------->4<-------9
	//              |        ^   (1)  |   (7)  |
	//              |(1)     |        |        |
	//              |        |(-2)    |(3)     |
	//    (wt: 10)  v   (2)  |        v        |
	//  1---------->2------->6<-------5<-------/
	//  |           ^        ^   (-1)    (-10)
	//  |(8)        |        |
	//  |           |(-4)    |(-1)
	//  v     (1)   |        |
	//  8---------->7--------/
	g := graph.LabeledAdjacencyList{
		1: {{2, 10}, {8, 8}},
		2: {{6, 2}},
		3: {{2, 1}, {4, 1}},
		4: {{5, 3}},
		5: {{6, -1}},
		6: {{3, -2}},
		7: {{6, -1}, {2, -4}},
		8: {{7, 1}},
		9: {{5, -10}, {4, 7}},
	}
	w := func(label int) float64 { return float64(label) }
	b := graph.NewBellmanFord(g, w)
	start := 1
	fmt.Println("start:", start)
	if !b.Run(start) {
		fmt.Println("negative cycle")
		return
	}
	fmt.Println("end   path  path")
	fmt.Println("node  len   dist   path")
	p := make([]int, b.Tree.MaxLen)
	for n, e := range b.Tree.Paths {
		fmt.Printf("%d       %d   %4.0f   %d\n",
			n, e.Len, b.Dist[n], b.Tree.PathTo(n, p))
	}
	// start: 1
	// end   path  path
	// node  len   dist   path
	// 0       0   +Inf   []
	// 1       1      0   [1]
	// 2       4      5   [1 8 7 2]
	// 3       6      5   [1 8 7 2 6 3]
	// 4       7      6   [1 8 7 2 6 3 4]
	// 5       8      9   [1 8 7 2 6 3 4 5]
	// 6       5      7   [1 8 7 2 6]
	// 7       3      9   [1 8 7]
	// 8       2      8   [1 8]
	// 9       0   +Inf   []
}
