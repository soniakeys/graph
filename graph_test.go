// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph_test

import (
	"fmt"
	"math/big"

	"github.com/soniakeys/graph"
)

func ExampleOneBits() {
	g := make(graph.AdjacencyList, 5)
	var b big.Int
	fmt.Printf("%b\n", graph.OneBits(&b, len(g)))
	// Output:
	// 11111
}
