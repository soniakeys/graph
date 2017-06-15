// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package gen_test

import "github.com/soniakeys/graph/gen"

func ExampleFuncData_okNodeVisitor() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := `
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}`
	data := &gen.Data{
		Pkg: "main",
		Funcs: []gen.FuncData{{
			Func:          "dfZ",
			OkNodeVisitor: true,
		}},
	}
	call := `
	dfZ(g, 0, func(n graph.NI) bool {
		fmt.Println("visit", n)
		return n != 2
	})`
	run(g, data, call)
	// Output:
	// visit 0
	// visit 1
	// visit 2
}
