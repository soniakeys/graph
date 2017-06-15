// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// A single variadic function, DepthFirst, takes options in the form of configuration functions.
package searchx

import (
	"errors"
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// DepthFirst performs a depth-first search or traversal of graph g starting at
// node start.
//
// Options controlling the search are specified with configuration functions
// defined in this package.
//
// A non-nil error indicates some problem initializing the search, such as
// an invalid graph type or options.
func DepthFirst(g interface{}, start graph.NI, options ...func(*config)) error {
	cf := &config{}
	for _, o := range options {
		o(cf)
	}
	if cf.nodeVisitor != nil && cf.okNodeVisitor != nil {
		return errors.New("NodeVisitor and OkNodeVisitor cannot both be specified")
	}
	if cf.arcVisitor != nil && cf.okArcVisitor != nil {
		return errors.New("ArcVisitor and OkArcVisitor cannot both be specified")
	}
	if cf.visited == nil { // for now, visited required internally
		cf.visited = &bits.Bits{}
	}
	var f func(start graph.NI)
	n := 0
	switch t := g.(type) {
	case graph.AdjacencyList:
		n = len(t)
		f = cf.adjFunc(t)
	case graph.LabeledAdjacencyList:
		n = len(t)
		f = cf.labFunc(t)
	default:
		return errors.New("invalid graph type")
	}
	if cf.visited.Num != n {
		*cf.visited = bits.New(n)
	}
	f(start)
	return nil
}

// skeleton for df traversal involves three functions, traverse, visited, and
// recurse.  traverse and recurse are mutually recursive.  traverse is a method
// so you start by taking a traverse "method value" t then creating recurse as
// a a closure that uses t.  visited can be created independently.
type dfTraverseNodes struct {
	visited func(graph.NI) bool
	recurse func(graph.NI)
}

func (f *dfTraverseNodes) traverse(n graph.NI) {
	if !f.visited(n) {
		f.recurse(n)
	}
}

// skeleton fo df search.  similar to traverse but boolean result is propagated
// back through search.
type dfDepthFirstNodes struct {
	visited func(graph.NI) bool
	recurse func(graph.NI) bool
}

func (f *dfDepthFirstNodes) search(n graph.NI) bool {
	return f.visited(n) || f.recurse(n)
}

func (cf *config) adjFunc(g graph.AdjacencyList) func(graph.NI) {
	if cf.okNodeVisitor == nil && cf.okArcVisitor == nil {
		// simpler case of full traversal
		f := dfTraverseNodes{visited: cf.visitedFunc()}
		// take method value
		traverse := f.traverse
		// define recurse using the method value
		f.recurse = cf.composeTraverseVisitor(cf.adjRecurseTraverse(g, traverse))
		return traverse
	}
	f := dfDepthFirstNodes{visited: cf.visitedFunc()}
	search := f.search
	f.recurse = cf.composeDepthFirstVisitor(cf.adjRecurseDepthFirst(g, search))
	// closure to drop final return value
	return func(start graph.NI) { search(start) }
}

func (cf *config) visitedFunc() func(graph.NI) bool {
	// only option for now is to use bits
	b := cf.visited
	return func(n graph.NI) (t bool) {
		if b.Bit(int(n)) != 0 {
			return true
		}
		b.SetBit(int(n), 1)
		return false
	}
}

func (cf *config) composeTraverseVisitor(f func(graph.NI)) func(graph.NI) {
	if v := cf.nodeVisitor; v != nil {
		return func(n graph.NI) {
			v(n)
			f(n)
		}
	}
	return f
}

func (cf *config) composeDepthFirstVisitor(f func(graph.NI) bool) func(graph.NI) bool {
	if v := cf.okNodeVisitor; v != nil {
		return func(n graph.NI) bool {
			return v(n) && f(n)
		}
	}
	if v := cf.nodeVisitor; v != nil {
		return func(n graph.NI) bool {
			v(n)
			return f(n)
		}
	}
	return f
}

func (cf *config) adjRecurseDepthFirst(g graph.AdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
	if r := cf.rand; r != nil {
		return cf.adjRandDepthFirst(g, search, r)
	}
	return cf.adjToDepthFirst(g, search)
}

func (cf *config) adjRandDepthFirst(g graph.AdjacencyList, search func(graph.NI) bool, r *rand.Rand) func(graph.NI) bool {
	if v := cf.okArcVisitor; v != nil {
		return func(n graph.NI) bool {
			to := g[n]
			for _, x := range r.Perm(len(to)) {
				if !v(n, x) || !search(to[x]) {
					return false
				}
			}
			return true
		}
	}
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) bool {
			to := g[n]
			for _, x := range r.Perm(len(to)) {
				v(n, x)
				if !search(to[x]) {
					return false
				}
			}
			return true
		}
	}
	return func(n graph.NI) bool {
		to := g[n]
		for _, i := range r.Perm(len(to)) {
			if !search(to[i]) {
				return false
			}
		}
		return true
	}
}

func (cf *config) adjToDepthFirst(g graph.AdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
	if v := cf.okArcVisitor; v != nil {
		return func(n graph.NI) bool {
			for x, to := range g[n] {
				if !v(n, x) || !search(to) {
					return false
				}
			}
			return true
		}
	}
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) bool {
			for x, to := range g[n] {
				v(n, x)
				if !search(to) {
					return false
				}
			}
			return true
		}
	}
	return func(n graph.NI) bool {
		for _, to := range g[n] {
			if !search(to) {
				return false
			}
		}
		return true
	}
}

func (cf *config) adjRecurseTraverse(g graph.AdjacencyList, traverse func(graph.NI)) func(graph.NI) {
	if r := cf.rand; r != nil {
		return cf.adjRandTraverse(g, traverse, r)
	}
	return cf.adjToTraverse(g, traverse)
}

func (cf *config) adjRandTraverse(g graph.AdjacencyList, traverse func(graph.NI), r *rand.Rand) func(graph.NI) {
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) {
			to := g[n]
			for _, x := range r.Perm(len(to)) {
				v(n, x)
				traverse(to[x])
			}
		}
	}
	return func(n graph.NI) {
		to := g[n]
		for _, x := range r.Perm(len(to)) {
			traverse(to[x])
		}
	}
}

func (cf *config) adjToTraverse(g graph.AdjacencyList, traverse func(graph.NI)) func(graph.NI) {
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) {
			for x, to := range g[n] {
				v(n, x)
				traverse(to)
			}
		}
	}
	return func(n graph.NI) {
		for _, to := range g[n] {
			traverse(to)
		}
	}
}

func (cf *config) labFunc(g graph.LabeledAdjacencyList) func(graph.NI) {
	if cf.okNodeVisitor == nil && cf.okArcVisitor == nil {
		f := dfTraverseNodes{visited: cf.visitedFunc()}
		traverse := f.traverse
		f.recurse = cf.composeTraverseVisitor(cf.labRecurseTraverse(g, traverse))
		return traverse
	}
	f := dfDepthFirstNodes{visited: cf.visitedFunc()}
	search := f.search
	f.recurse = cf.composeDepthFirstVisitor(cf.labRecurseDepthFirst(g, search))
	return func(start graph.NI) { search(start) }
}

func (cf *config) labRecurseDepthFirst(g graph.LabeledAdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
	if r := cf.rand; r != nil {
		return cf.labRandDepthFirst(g, search, r)
	}
	return cf.labToDepthFirst(g, search)
}

func (cf *config) labRandDepthFirst(g graph.LabeledAdjacencyList, search func(graph.NI) bool, r *rand.Rand) func(graph.NI) bool {
	if v := cf.okArcVisitor; v != nil {
		return func(n graph.NI) bool {
			to := g[n]
			for _, x := range r.Perm(len(to)) {
				if !v(n, x) || !search(to[x].To) {
					return false
				}
			}
			return true
		}
	}
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) bool {
			to := g[n]
			for _, x := range r.Perm(len(to)) {
				v(n, x)
				if !search(to[x].To) {
					return false
				}
			}
			return true
		}
	}
	return func(n graph.NI) bool {
		to := g[n]
		for _, i := range r.Perm(len(to)) {
			if !search(to[i].To) {
				return false
			}
		}
		return true
	}
}

func (cf *config) labToDepthFirst(g graph.LabeledAdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
	if v := cf.okArcVisitor; v != nil {
		return func(n graph.NI) bool {
			for x, to := range g[n] {
				if !v(n, x) || !search(to.To) {
					return false
				}
			}
			return true
		}
	}
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) bool {
			for x, to := range g[n] {
				v(n, x)
				if !search(to.To) {
					return false
				}
			}
			return true
		}
	}
	return func(n graph.NI) bool {
		for _, to := range g[n] {
			if !search(to.To) {
				return false
			}
		}
		return true
	}
}

func (cf *config) labRecurseTraverse(g graph.LabeledAdjacencyList, traverse func(graph.NI)) func(graph.NI) {
	if r := cf.rand; r != nil {
		return cf.labRandTraverse(g, traverse, r)
	}
	return cf.labToTraverse(g, traverse)
}

func (cf *config) labRandTraverse(g graph.LabeledAdjacencyList, traverse func(graph.NI), r *rand.Rand) func(graph.NI) {
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) {
			to := g[n]
			for _, x := range r.Perm(len(to)) {
				v(n, x)
				traverse(to[x].To)
			}
		}
	}
	return func(n graph.NI) {
		to := g[n]
		for _, i := range r.Perm(len(to)) {
			traverse(to[i].To)
		}
	}
}

func (cf *config) labToTraverse(g graph.LabeledAdjacencyList, traverse func(graph.NI)) func(graph.NI) {
	if v := cf.arcVisitor; v != nil {
		return func(n graph.NI) {
			for x, to := range g[n] {
				v(n, x)
				traverse(to.To)
			}
		}
	}
	return func(n graph.NI) {
		for _, to := range g[n] {
			traverse(to.To)
		}
	}
}
