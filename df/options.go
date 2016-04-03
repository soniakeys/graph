// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package df

import (
	"math/rand"

	"github.com/soniakeys/graph"
)

type config struct {
	iterateFrom func(n graph.NI)
	bits        *graph.Bits
	visitor     graph.Visitor
	visitOk     *bool
	rand        *rand.Rand
}

func Bits(b *graph.Bits) func(*config) {
	return func(c *config) { c.bits = b }
}

func Visitor(v graph.Visitor, ok *bool) func(*config) {
	return func(c *config) {
		c.visitor = v
		c.visitOk = ok
	}
}

func Rand(r *rand.Rand) func(*config) {
	return func(c *config) { c.rand = r }
}
