// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package treevis

import (
	"fmt"
	"io"

	"github.com/soniakeys/graph"
)

func WriteDirected(g graph.Directed, label func(graph.NI) string, w io.Writer) {
	if len(g.AdjacencyList) == 0 {
		fmt.Fprintln(w, "<empty>")
		return
	}
	var f func(graph.NI, string)
	f = func(n graph.NI, pre string) {
		to := g.AdjacencyList[n]
		if len(to) == 0 {
			fmt.Fprintln(w, "╴", label(n))
			return
		}
		fmt.Fprintln(w, "┐", label(n))
		last := len(to) - 1
		for _, to := range to[:last] {
			fmt.Fprint(w, pre, "├─")
			f(to, pre+"│ ")
		}
		fmt.Fprint(w, pre, "└─")
		f(to[last], pre+"  ")
	}
	f(0, "")
}
