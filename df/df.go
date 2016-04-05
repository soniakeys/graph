// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// Package df provides a paramertized depth-first search.
//
// A single variadic function, Search, takes options in the form of configuration functions.
package df

import (
	"errors"
	"math/rand"

	"github.com/soniakeys/graph"
)

// Search performs a depth-first search or traversal of graph g starting at
// node start.
//
// Options controlling the search are specified with configuration functions
// defined in this package.
//
// A non-nil error indicates some problem initializing the search, such as
// an invalid graph type or options.
func Search(g interface{}, start graph.NI, options ...func(*config)) (err error) {
	cf := &config{}
	for _, o := range options {
		o(cf)
	}
	if cf.nodeVisitor != nil && cf.okNodeVisitor != nil {
		return errors.New("NodeVisitor and OkNodeVisitor cannot both be specified")
	}
	if cf.bits == nil { // for now, bits required as visit marks
		cf.bits = &graph.Bits{}
	}
	var f func(start graph.NI)
	switch t := g.(type) {
	case graph.AdjacencyList:
		f, err = cf.adjFunc(t)
	case graph.LabeledAdjacencyList:
		f, err = cf.labFunc(t)
	default:
		return errors.New("invalid graph type")
	}
	if err == nil {
		f(start)
	}
	return
}

// skeleton for df traversal involves three functions, traverse, visited, and recurse.
// traverse and recurse are mutually recursive.  traverse is a method so you start by taking a
// traverse "method value" t then creating recurse as a a closure that uses t.
// visited can be created independently.
type dfTraverseNodes struct {
	visited func(graph.NI) bool
	recurse func(graph.NI)
}

func (f *dfTraverseNodes) traverse(n graph.NI) {
	if !f.visited(n) {
		f.recurse(n)
	}
}

type dfSearchNodes struct {
	visited func(graph.NI) bool
	recurse func(graph.NI) bool
}

func (f *dfSearchNodes) search(n graph.NI) bool {
	return f.visited(n) || f.recurse(n)
}

func (cf *config) adjFunc(g graph.AdjacencyList) (func(graph.NI), error) {
	if v := cf.okNodeVisitor; v != nil {
		f := dfSearchNodes{visited: cf.visitedFunc()}
		// take method value
		search := f.search
		// define recurse using the method value
		f.recurse = cf.composeOkVisitor(v, cf.adjRecurseSearch(g, search))
		// closure to drop return value
		return func(start graph.NI) { search(start) }, nil
	}
	// simpler case of full traversal
	f := dfTraverseNodes{visited: cf.visitedFunc()}
	traverse := f.traverse
	f.recurse = cf.composeVisitor(cf.adjRecurseTraverse(g, traverse))
	return traverse, nil
}

func (cf *config) visitedFunc() func(graph.NI) bool {
	b := cf.bits
	return func(n graph.NI) (t bool) {
		if b.Bit(n) != 0 {
			return true
		}
		b.SetBit(n, 1)
		return false
	}
}

func (cf *config) composeOkVisitor(v graph.OkNodeVisitor, f func(graph.NI) bool) func(graph.NI) bool {
	if cf.visitOk != nil { // TODO move this to a pre-func
		*cf.visitOk = true
	}
	return func(n graph.NI) bool {
		if !v(n) {
			if cf.visitOk != nil { // and post-func?
				*cf.visitOk = false
			}
			return false
		}
		return f(n)
	}
}

func (cf *config) composeVisitor(f func(graph.NI)) func(graph.NI) {
	if v := cf.nodeVisitor; v != nil {
		return func(n graph.NI) {
			v(n)
			f(n)
		}
	}
	return f
}

func (cf *config) adjRecurseSearch(g graph.AdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
	if r := cf.rand; r != nil {
		return cf.adjRandSearch(g, search, r)
	}
	return cf.adjToSearch(g, search)
}

func (cf *config) adjRandSearch(g graph.AdjacencyList, search func(graph.NI) bool, r *rand.Rand) func(graph.NI) bool {
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

func (cf *config) adjToSearch(g graph.AdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
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
	return func(n graph.NI) {
		to := g[n]
		for _, i := range r.Perm(len(to)) {
			traverse(to[i])
		}
	}
}

func (cf *config) adjToTraverse(g graph.AdjacencyList, traverse func(graph.NI)) func(graph.NI) {
	return func(n graph.NI) {
		for _, to := range g[n] {
			traverse(to)
		}
	}
}

func (cf *config) labFunc(g graph.LabeledAdjacencyList) (func(graph.NI), error) {
	if v := cf.okNodeVisitor; v != nil {
		f := dfSearchNodes{visited: cf.visitedFunc()}
		search := f.search
		f.recurse = cf.composeOkVisitor(v, cf.labRecurseSearch(g, search))
		return func(start graph.NI) { search(start) }, nil
	}
	f := dfTraverseNodes{visited: cf.visitedFunc()}
	traverse := f.traverse
	f.recurse = cf.composeVisitor(cf.labRecurseTraverse(g, traverse))
	return traverse, nil
}

func (cf *config) labRecurseSearch(g graph.LabeledAdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
	if r := cf.rand; r != nil {
		return cf.labRandSearch(g, search, r)
	}
	return cf.labToSearch(g, search)
}

func (cf *config) labRandSearch(g graph.LabeledAdjacencyList, search func(graph.NI) bool, r *rand.Rand) func(graph.NI) bool {
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

func (cf *config) labToSearch(g graph.LabeledAdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
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
	return func(n graph.NI) {
		to := g[n]
		for _, i := range r.Perm(len(to)) {
			traverse(to[i].To)
		}
	}
}

func (cf *config) labToTraverse(g graph.LabeledAdjacencyList, traverse func(graph.NI)) func(graph.NI) {
	return func(n graph.NI) {
		for _, to := range g[n] {
			traverse(to.To)
		}
	}
}
