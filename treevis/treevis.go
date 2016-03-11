// Copyright 2016 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package treevis

import (
	"fmt"
	"io"
	"math/big"

	"github.com/soniakeys/graph"
)

func WriteDirected(g graph.Directed, root graph.NI, label func(graph.NI) string, w io.Writer) {
	if len(g.AdjacencyList) == 0 {
		fmt.Fprintln(w, "<empty>")
		return
	}
	var vis big.Int
	var f func(graph.NI, string) bool
	f = func(n graph.NI, pre string) bool {
		if vis.Bit(int(n)) != 0 {
			fmt.Fprintln(w, "<cycle detected>")
			return false
		}
		vis.SetBit(&vis, int(n), 1)
		to := g.AdjacencyList[n]
		if len(to) == 0 {
			fmt.Fprintln(w, "╴", label(n))
			return true
		}
		fmt.Fprintln(w, "┐", label(n))
		last := len(to) - 1
		for _, to := range to[:last] {
			fmt.Fprint(w, pre, "├─")
			f(to, pre+"│ ")
		}
		fmt.Fprint(w, pre, "└─")
		f(to[last], pre+"  ")
		return true
	}
	f(root, "")
}
