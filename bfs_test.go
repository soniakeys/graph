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
	fmt.Println(b.MaxLevel)
	fmt.Println(b.Level)
	// Output:
	// 3
	// [-1 0 -1 2 1 3 2]
}
