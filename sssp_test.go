// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/soniakeys/graph"
)

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
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	h := func(from graph.NI) float64 { return h4[from] }
	p, d := g.AStarAPath(0, 4, h, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path distance:", d)
	// Output:
	// Shortest path: [0 2 3 4]
	// Path distance: 26
}

func ExampleLabeledAdjacencyList_AStarMPath() {
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
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	h := func(from graph.NI) float64 { return h4[from] }
	p, d := g.AStarMPath(0, 4, h, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path distance:", d)
	// Output:
	// Shortest path: [0 2 3 4]
	// Path distance: 26
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
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	var h graph.Heuristic = func(from graph.NI) float64 { return h4[from] }
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
	w := func(label graph.LI) float64 { return float64(label) }
	h4 := []float64{19, 20, 10, 6, 0, 9}
	var h graph.Heuristic = func(from graph.NI) float64 { return h4[from] }
	fmt.Println(h.Monotonic(g, w))
	// Output:
	// true
}

func ExampleBellmanFord() {
	//              /--------3        4<-------9
	//              |        ^        |   (6)  ^
	//              |(1)     |        |        |
	//              |        |(-2)    |(3)     |
	//    (wt: 10)  v   (2)  |        v        |
	//  1---------->2------->6<-------5--------/
	//  |           ^        ^   (-1)    (-10)
	//  |(8)        |        |
	//  |           |(-4)    |(-1)
	//  v     (1)   |        |
	//  8---------->7--------/
	g := graph.LabeledAdjacencyList{
		1: {{2, 10}, {8, 8}},
		2: {{6, 2}},
		3: {{2, 1}},
		4: {{5, 3}},
		5: {{6, -1}, {9, -10}},
		6: {{3, -2}},
		7: {{6, -1}, {2, -4}},
		8: {{7, 1}},
		9: {{4, 6}},
	}
	w := func(label graph.LI) float64 { return float64(label) }
	b := graph.NewBellmanFord(g, w)
	// graph contains negative cycle somewhere
	fmt.Println("negative cycle:", b.NegativeCycle())

	// but negative cycle not reached starting at node 1
	start := graph.NI(1)
	fmt.Println("start:", start)
	if !b.Start(start) {
		fmt.Println("negative cycle")
		return
	}
	fmt.Println("end   path  path")
	fmt.Println("node  len   dist   path")
	p := make([]graph.NI, b.Forest.MaxLen)
	for n, e := range b.Forest.Paths {
		fmt.Printf("%d       %d   %4.0f   %d\n",
			n, e.Len, b.Dist[n], b.Forest.PathTo(graph.NI(n), p))
	}
	// Output:
	// negative cycle: true
	// start: 1
	// end   path  path
	// node  len   dist   path
	// 0       0   +Inf   []
	// 1       1      0   [1]
	// 2       4      5   [1 8 7 2]
	// 3       6      5   [1 8 7 2 6 3]
	// 4       0   +Inf   []
	// 5       0   +Inf   []
	// 6       5      7   [1 8 7 2 6]
	// 7       3      9   [1 8 7]
	// 8       2      8   [1 8]
	// 9       0   +Inf   []
}
func ExampleAdjacencyList_BreadthFirstPath() {
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
	var start, end graph.NI = 1, 3
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
	p := make([]graph.NI, b.Result.MaxLen)
	for n := range b.Graph {
		fmt.Println(n, b.Result.PathTo(graph.NI(n), p))
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

func ExampleBreadthFirst_Traverse() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	b := graph.NewBreadthFirst(graph.AdjacencyList{
		0: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
		8: {},
	})
	b.Traverse(0, func(n graph.NI) bool {
		fmt.Println("visit", n, "level", b.Result.Paths[n].Len)
		return true
	})
	// Output:
	// visit 0 level 1
	// visit 1 level 2
	// visit 2 level 2
	// visit 3 level 2
	// visit 4 level 3
	// visit 5 level 3
	// visit 6 level 3
	// visit 7 level 3
	// visit 8 level 3
}

func ExampleBreadthFirst_Traverse_random() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	b := graph.NewBreadthFirst(graph.AdjacencyList{
		0: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
		8: {},
	})

	// only difference from non-random example
	b.Rand = rand.New(rand.NewSource(8))

	b.Traverse(0, func(n graph.NI) bool {
		fmt.Println("visit", n, "level", b.Result.Paths[n].Len)
		return true
	})
	// Output:
	// visit 0 level 1
	// visit 1 level 2
	// visit 3 level 2
	// visit 2 level 2
	// visit 8 level 3
	// visit 5 level 3
	// visit 6 level 3
	// visit 4 level 3
	// visit 7 level 3
}

func ExampleAdjacencyList_BreadthFirst2Path() {
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
	g := graph.Directed{graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}}
	t, m := g.Transpose()
	b := graph.NewBreadthFirst2(g.AdjacencyList, t.AdjacencyList, m)
	var start, end graph.NI = 1, 3
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
	g := graph.Directed{graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}}
	t, m := g.Transpose()
	b := graph.NewBreadthFirst2(g.AdjacencyList, t.AdjacencyList, m)
	b.AllPaths(1)
	fmt.Println("Max path length:", b.Result.MaxLen)
	p := make([]graph.NI, b.Result.MaxLen)
	for n := range b.To {
		fmt.Println(n, b.Result.PathTo(graph.NI(n), p))
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

func ExampleDAGPath_AllPaths() {
	// arcs are directed right:
	//   (11)
	// 0------2         4
	//                   \(11)
	//      (11)    (10)  \   (30)   (10)
	//    1-------3--------5-------7-------9
	//                      \     /
	//                   (10)\   /(20)
	//                        \ /
	//                         6------8
	//                           (10)
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 2, Label: 11}},
		1: {{3, 11}},
		3: {{5, 10}},
		4: {{5, 11}},
		5: {{6, 10}, {7, 30}},
		6: {{7, 20}, {8, 10}},
		7: {{9, 10}},
		9: {},
	}}
	o := []graph.NI{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	d := graph.NewDAGPath(g, o,
		func(l graph.LI) float64 { return float64(l) },
		false)
	reached := d.AllPaths(3)
	fmt.Println("node  path dist  path len  leaf")
	for n, pd := range d.Dist {
		fmt.Printf("%d  %9.0f  %9d %7d\n",
			n, pd, d.Forest.Paths[n].Len, d.Forest.Leaves.Bit(n))
	}
	fmt.Println()
	fmt.Println("Nodes reached:       ", reached)
	fmt.Println("Max path len:        ", d.Forest.MaxLen)
	// Output:
	// node  path dist  path len  leaf
	// 0          0          0       0
	// 1          0          0       0
	// 2          0          0       0
	// 3          0          1       0
	// 4          0          0       0
	// 5         10          2       0
	// 6         20          3       0
	// 7         40          3       0
	// 8         30          4       1
	// 9         50          4       1
	//
	// Nodes reached:        6
	// Max path len:         4
}

func ExampleDAGPath_Path_shortest() {
	// arcs are directed right:
	//             4
	//        (-3)/ \(-2)
	//           /   \
	//    (10)  /     \   (5)    (10)
	// 3-------1-------0-------6-------2
	//           (-10)  \     /
	//                   \   /
	//                 (2)\ /(3)
	//                     5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 5, Label: 2}, {6, 5}},
		1: {{0, -10}, {4, -3}},
		3: {{1, 10}},
		4: {{0, -2}},
		5: {{6, 3}},
		6: {{2, 10}},
	}}
	o, _ := g.Topological()
	fmt.Println("Ordering:", o)
	d := graph.NewDAGPath(g, o,
		func(l graph.LI) float64 { return float64(l) },
		false)
	var start, end graph.NI = 3, 2
	if !d.Path(start, end) {
		fmt.Println("not found")
	}
	fmt.Println("Path length:", d.Forest.Paths[end].Len)
	fmt.Println("Path:", d.Forest.PathTo(end, nil))
	fmt.Println("Distance:", d.Dist[end])
	// Output:
	// Ordering: [3 1 4 0 5 6 2]
	// Path length: 5
	// Path: [3 1 0 6 2]
	// Distance: 15
}

func ExampleDAGPath_Path_longest() {
	// arcs are directed right:
	//             4
	//        (-3)/ \(-2)
	//           /   \
	//    (10)  /     \   (5)    (10)
	// 3-------1-------0-------6-------2
	//           (-10)  \     /
	//                   \   /
	//                 (2)\ /(3)
	//                     5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 5, Label: 2}, {6, 5}},
		1: {{0, -10}, {4, -3}},
		3: {{1, 10}},
		4: {{0, -2}},
		5: {{6, 3}},
		6: {{2, 10}},
	}}
	o, _ := g.Topological()
	fmt.Println("Ordering:", o)
	d := graph.NewDAGPath(g, o,
		func(l graph.LI) float64 { return float64(l) },
		true)
	var start, end graph.NI = 3, 2
	if !d.Path(start, end) {
		fmt.Println("not found")
	}
	fmt.Println("Path length:", d.Forest.Paths[end].Len)
	fmt.Println("Path:", d.Forest.PathTo(end, nil))
	fmt.Println("Distance:", d.Dist[end])
	// Output:
	// Ordering: [3 1 4 0 5 6 2]
	// Path length: 7
	// Path: [3 1 4 0 5 6 2]
	// Distance: 20
}

func ExampleLabeledDirected_DAGMinDistPath() {
	// arcs are directed right:
	//             4
	//        (-3)/ \(-2)
	//           /   \
	//    (10)  /     \   (5)    (10)
	// 3-------1-------0-------6-------2
	//           (-10)  \     /
	//                   \   /
	//                 (2)\ /(3)
	//                     5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 5, Label: 2}, {6, 5}},
		1: {{0, -10}, {4, -3}},
		3: {{1, 10}},
		4: {{0, -2}},
		5: {{6, 3}},
		6: {{2, 10}},
	}}
	var start, end graph.NI = 3, 2
	w := func(l graph.LI) float64 { return float64(l) }
	p, dist, err := g.DAGMinDistPath(start, end, w)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Path:", p)
	fmt.Println("Distance:", dist)
	// Output:
	// Path: [3 1 0 6 2]
	// Distance: 15
}

func ExampleLabeledDirected_DAGMaxDistPath() {
	// arcs are directed right:
	//             4
	//        (-3)/ \(-2)
	//           /   \
	//    (10)  /     \   (5)    (10)
	// 3-------1-------0-------6-------2
	//           (-10)  \     /
	//                   \   /
	//                 (2)\ /(3)
	//                     5
	g := graph.LabeledDirected{graph.LabeledAdjacencyList{
		0: {{To: 5, Label: 2}, {6, 5}},
		1: {{0, -10}, {4, -3}},
		3: {{1, 10}},
		4: {{0, -2}},
		5: {{6, 3}},
		6: {{2, 10}},
	}}
	var start, end graph.NI = 3, 2
	w := func(l graph.LI) float64 { return float64(l) }
	p, dist, err := g.DAGMaxDistPath(start, end, w)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Path:", p)
	fmt.Println("Distance:", dist)
	// Output:
	// Path: [3 1 4 0 5 6 2]
	// Distance: 20
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
	w := func(label graph.LI) float64 { return float64(label) }
	p, d := g.DijkstraPath(1, 5, w)
	fmt.Println("Shortest path:", p)
	fmt.Println("Path distance:", d)
	// Output:
	// Shortest path: [1 6 5]
	// Path distance: 20
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
	w := func(label graph.LI) float64 { return float64(label) }
	d := graph.NewDijkstra(g, w)
	var start, end graph.NI = 1, 5
	if d.Path(start, end) {
		fmt.Println("Path", d.Forest.PathTo(end, nil), "distance", d.Dist[end])
	}
	// Output:
	// Path [1 6 5] distance 20
}

func ExampleNewDijkstra_consecutive() {
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
	w := func(label graph.LI) float64 { return float64(label) }
	d := graph.NewDijkstra(g, w)
	out := func(d *graph.Dijkstra, end graph.NI) {
		fmt.Println(
			"Path", d.Forest.PathTo(end, nil), "distance", d.Dist[end])
	}
	if d.Path(1, 5) {
		out(d, 5)
	}
	d.Reset()
	if d.Path(2, 5) {
		out(d, 5)
	}
	// Output:
	// Path [1 6 5] distance 20
	// Path [2 3 6 5] distance 21
}

func ExampleNewDijkstra_concurrent() {
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
	w := func(label graph.LI) float64 { return float64(label) }

	// result channel
	type result struct {
		d   *graph.Dijkstra
		end graph.NI
		ok  bool
	}
	ch := make(chan result)

	// Start two concurrent goroutines.  Each goroutine gets it's own
	// Dijkstra struct, but they search the same graph.
	f := func(d *graph.Dijkstra, start, end graph.NI) {
		ch <- result{d, end, d.Path(start, end)}
	}
	go f(graph.NewDijkstra(g, w), 1, 5)
	go f(graph.NewDijkstra(g, w), 2, 5)

	var out []string // to sort formatted output

	// format results from the two goroutines
	for i := 0; i < 2; i++ {
		if r := <-ch; r.ok {
			out = append(out, fmt.Sprintln(
				"Path", r.d.Forest.PathTo(r.end, nil),
				"distance", r.d.Dist[r.end]))
		}
	}

	// sort for determinism to pass go test
	sort.Strings(out)
	for _, l := range out {
		fmt.Print(l)
	}

	// Output:
	// Path [1 6 5] distance 20
	// Path [2 3 6 5] distance 21
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
	w := func(label graph.LI) float64 { return float64(label) }
	d := graph.NewDijkstra(g, w)
	fmt.Println(d.AllPaths(2), "paths found.")
	// column len is from Result, and will be equal to len(path).
	// column dist is from Result, and will be equal to sum.
	fmt.Println("node:  path                  len  dist")
	p := make([]graph.NI, d.Forest.MaxLen)
	for nd := range g {
		r := &d.Forest.Paths[nd]
		path := d.Forest.PathTo(graph.NI(nd), p)
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
