// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt_test

import (
	"fmt"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func ExampleSCCPathBased() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}}
	alt.SCCPathBased(g, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [1 3 2]
	// [7 6]
	// [4 5]
	// [0]
}

func ExampleSCCTarjan() {
	// /---0---\
	// |   |\--/
	// |   v
	// |   5<=>4---\
	// |   |   |   |
	// v   v   |   |
	// 7<=>6   |   |
	//     |   v   v
	//     \-->3<--2
	//         |   ^
	//         |   |
	//         \-->1
	g := graph.Directed{graph.AdjacencyList{
		0: {0, 5, 7},
		5: {4, 6},
		4: {5, 2, 3},
		7: {6},
		6: {7, 3},
		3: {1},
		1: {2},
		2: {3},
	}}
	alt.SCCTarjan(g, func(c []graph.NI) bool {
		fmt.Println(c)
		return true
	})
	// Output:
	// [1 3 2]
	// [7 6]
	// [4 5]
	// [0]
}
