package search_test

import (
	"fmt"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/search"
)

// Cyclic reimplements AdjacencyList.Cyclic using search.DepthFirst
func cyclic(g graph.AdjacencyList) (cyclic bool, fr, to graph.NI) {
	vis := bits.New(len(g))
	path := bits.New(len(g))
	for n := range g {
		if vis.Bit(n) == 1 {
			continue
		}
		search.DepthFirst(g, graph.NI(n),
			search.Visited(&vis), search.PathBits(&path),
			search.OkArcVisitor(func(n graph.NI, x int) bool {
				to = g[n][x]
				cyclic, fr = path.Bit(int(to)) == 1, n
				return !cyclic
			}))
		if cyclic {
			return
		}
	}
	return false, -1, -1
}

// The example here repeats the example from AdjacencyList.Cyclic.
func Example_cyclic() {
	//   0
	//  / \
	// 1-->2-->3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {},
	}
	c, _, _ := cyclic(g)
	fmt.Println(c)

	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g[3] = []graph.NI{1}
	fmt.Println(cyclic(g))

	// Output:
	// false
	// true 3 1
}
