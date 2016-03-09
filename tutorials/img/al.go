package main

import (
	"log"
	"os/exec"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/dot"
)

func main() {
	for _, f := range []func() error{al} {
		if err := f(); err != nil {
			log.Fatal(err)
		}
	}
}

func al() error {
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	c := exec.Command("dot", "-Tsvg", "-o", "al.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	dot.WriteAdjacencyList(g, w, dot.GraphAttr("rankdir", "LR"))
	w.Close()
	return c.Wait()
}

func al0() error {
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	c := exec.Command("dot", "-Tsvg", "-o", "al0.svg")
	w, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Start()
	dot.WriteAdjacencyList(g, w, dot.GraphAttr("rankdir", "LR"))
	w.Close()
	return c.Wait()
}
