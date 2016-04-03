// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package df

import (
	"errors"

	"github.com/soniakeys/graph"
)

func Search(g interface{}, start graph.NI, options ...func(*config)) (err error) {
	cf := &config{}
	for _, o := range options {
		o(cf)
	}
	if cf.bits == nil {
		cf.bits = &graph.Bits{}
	}
	var f func(start graph.NI)
	switch t := g.(type) {
	case graph.AdjacencyList:
		f, err = cf.adjSearchFunc(t)
	case graph.LabeledAdjacencyList:
		f, err = cf.labSearchFunc(t)
	default:
		return errors.New("unsupported graph type")
	}
	if err == nil {
		f(start)
	}
	return
}

type dff struct {
	term    func(graph.NI) bool
	recurse func(graph.NI)
}

func (f *dff) search(n graph.NI) {
	if !f.term(n) {
		f.recurse(n)
	}
}

func (cf *config) adjSearchFunc(g graph.AdjacencyList) (func(graph.NI), error) {
	f := &dff{}
	search := f.search
	f.term = cf.termFunc()
	f.recurse = cf.adjRecurseFunc(g, search)
	return search, nil
}

func (cf *config) termFunc() func(graph.NI) bool {
	b := cf.bits
	t := func(n graph.NI) (t bool) {
		if t = b.Bit(n) != 0; !t {
			b.SetBit(n, 1)
		}
		return
	}
	if v := cf.visitor; v != nil {
		return func(n graph.NI) bool {
			return t(n) || !v(n)
		}
	}
	return t
}

func (cf *config) adjRecurseFunc(g graph.AdjacencyList, search func(graph.NI)) func(graph.NI) {
	if r := cf.rand; r != nil {
		return func(n graph.NI) {
			to := g[n]
			for _, i := range r.Perm(len(to)) {
				search(to[i])
			}

		}
	}
	return func(n graph.NI) {
		for _, to := range g[n] {
			search(to)
		}
	}
}

func (cf *config) labSearchFunc(g graph.LabeledAdjacencyList) (func(start graph.NI), error) {
	return nil, nil
}
