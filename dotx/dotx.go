package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func main() {
	cg_adj_dots()
	dir_dots()
	graph_dots()
}

func runDot(fn string, f func(w io.WriteCloser)) {
	cmd := exec.Command("dot", "-Tsvg") // dot uses stdin and stdout
	// stdin we'll pipe from graph.Write-something
	dotIn, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout, err = os.Create(fn + ".svg") // stdout goes to a file
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()

	f(dotIn)      // run the graph code
	dotIn.Close() // dot doesn't exit until stdin closes

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

// graph_test.go
func graph_dots() {
	fn := "AdjacencyList_BoundsOk"
	runDot(fn, func(w io.WriteCloser) {
		var g graph.AdjacencyList
		ok, _, _ := g.BoundsOk()
		dot.WriteAdjacencyList(g, w,
			dot.GraphAttr("label", fmt.Sprintf(`"%s: %t
Empty graph is ok."`, fn, ok)))
	})

	// No dot support for ExampleOneBits really.
	// Subgraph from bitmap somehow?

	fn = "AdjacencyList_Simple"
	runDot(fn, func(w io.WriteCloser) {
		g := graph.AdjacencyList{
			2: {0, 1},
		}
		s, _ := g.Simple()
		dot.WriteAdjacencyList(g, w,
			dot.GraphAttr("label", fmt.Sprintf(`"%s: %t
No loops or parallel arcs."`, fn, s)))
	})

	fn = "AdjacencyList_Simple(loop)"
	runDot(fn, func(w io.WriteCloser) {
		g := graph.AdjacencyList{
			1: {1},
			2: {0, 1},
		}
		s, _ := g.Simple()
		dot.WriteAdjacencyList(g, w,
			dot.GraphAttr("label", fmt.Sprintf(`"%s: %t
Loop not simple."`, fn, s)))
	})

	fn = "AdjacencyList_Simple(parallel)"
	runDot(fn, func(w io.WriteCloser) {
		g := graph.AdjacencyList{
			2: {0, 1, 0},
		}
		s, _ := g.Simple()
		dot.WriteAdjacencyList(g, w,
			dot.GraphAttr("label", fmt.Sprintf(`"%s: %t
Parallel arcs not simple."`, fn, s)))
	})
}

// dir_test.go
func dir_dots() {
}

// cg_adj_test.go
func cg_adj_dots() {
	fn := "AdjacencyList_Cyclic(acyclic)"
	runDot(fn, func(w io.WriteCloser) {
		g := graph.AdjacencyList{
			0: {1, 2},
			1: {2},
			2: {3},
			3: {},
		}
		cyc, _, _ := g.Cyclic()
		dot.WriteAdjacencyList(g, w,
			dot.GraphAttr("label", fmt.Sprintf(`"%s: %t
Acyclic."`, fn, cyc)))
	})

	fn = "AdjacencyList_Cyclic"
	runDot(fn, func(w io.WriteCloser) {
		g := graph.AdjacencyList{
			0: {1, 2},
			1: {2},
			2: {3},
			3: {1},
		}
		cyc, fr, to := g.Cyclic()
		dot.WriteAdjacencyList(g, w,
			dot.GraphAttr("label", fmt.Sprintf(`"%s: %t
Acyclic:  Edge %d -> %d"`, fn, cyc, fr, to)))
	})

	{
		g := graph.AdjacencyList{
			0: {1, 2},
			1: {2},
			2: {3},
			3: {1},
		}
		cyc, fr, to := g.Cyclic()
		dot.WriteAdjacencyList(g, os.Stdout,
			dot.GraphAttr("label", fmt.Sprintf(`"%s: %t
Acyclic:  Edge %d -> %d"`, fn, cyc, fr, to)))
	}
}
