// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package treevis_test

import (
	"os"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/treevis"
)

func ExampleWriteDirected() {
	g := graph.Directed{graph.AdjacencyList{
		0: {1, 2, 3}, // root
		1: {4, 5},    // ei
		5: {6},       // y
		6: {},        // eff
	}}
	labels := []string{
		0: "root",
		1: "ei",
		2: "bee",
		3: "si",
		4: "dee",
		5: "y",
		6: "eff",
	}
	treevis.WriteDirected(g,
		func(n graph.NI) string { return labels[n] },
		os.Stdout)
	// Output:
	// ┐ root
	// ├─┐ ei
	// │ ├─╴ dee
	// │ └─┐ y
	// │   └─╴ eff
	// ├─╴ bee
	// └─╴ si
}
