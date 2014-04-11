package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

func ExampleBreadthFirst_AllPaths() {
	b := ed.NewBreadthFirst(ed.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	})
	b.AllPaths(1)
	fmt.Println("Max path length:", b.Result.MaxLen)
	for n := range b.Graph {
		fmt.Println(n, b.Result.PathTo(n))
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

func ExampleBreadthFirst2_AllPaths() {
	g := ed.AdjacencyList{
		1: {4},
		2: {1},
		3: {5},
		4: {3, 6},
		6: {5, 6},
	}
	from, m := g.Inverse()
	b := ed.NewBreadthFirst2(g, from, m)
	b.AllPaths(1)
	fmt.Println("Max path length:", b.Result.MaxLen)
	for n := range b.To {
		fmt.Println(n, b.Result.PathTo(n))
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
