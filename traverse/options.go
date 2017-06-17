// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package traverse

import (
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

type config struct {
	start         graph.NI
	arcVisitor    func(n graph.NI, x int)
	iterateFrom   func(n graph.NI)
	nodeVisitor   func(n graph.NI)
	okArcVisitor  func(n graph.NI, x int) bool
	okNodeVisitor func(n graph.NI) bool
	rand          *rand.Rand
	visBits       *bits.Bits
	pathBits      *bits.Bits
	fromList      *graph.FromList
}

type Option func(*config)

// ArcVisitor specifies a visitor function to call at each arc.
//
// See also OkArcVisitor.
func ArcVisitor(v func(n graph.NI, x int)) Option {
	return func(c *config) {
		c.arcVisitor = v
	}
}

// FromList specifies a graph.FromList to populate.
func FromList(f *graph.FromList) Option {
	return func(c *config) {
		c.fromList = f
	}
}

// NodeVisitor specifies a visitor function to call at each node.
//
// See also OkNodeVisitor.
func NodeVisitor(v func(graph.NI)) Option {
	return func(c *config) {
		c.nodeVisitor = v
	}
}

// OkArcVisitor specifies a visitor function to perform some test at each arc
// and return a boolean result.
//
// As long as v return a result of true, the traverse progresses to traverse all
// arcs.
//
// If v returns false, the traverse terminates immediately.
//
// See also ArcVisitor.
func OkArcVisitor(v func(n graph.NI, x int) bool) Option {
	return func(c *config) {
		c.okArcVisitor = v
	}
}

// OkNodeVisitor specifies a visitor function to perform some test at each node
// and return a boolean result.
//
// As long as v return a result of true, the traverse progresses to traverse all
// nodes.
//
// If v returns false, the traverse terminates immediately.
//
// See also NodeVisitor.
func OkNodeVisitor(v func(graph.NI) bool) Option {
	return func(c *config) {
		c.okNodeVisitor = v
	}
}

// PathBits specifies a bits.Bits value for nodes of the path to the
// currently visited node.
//
// A use for PathBits is identifying back arcs in a traverse.
//
// Unlike Visited, PathBits are zeroed at the start of a traverse.
func PathBits(b *bits.Bits) Option {
	return func(c *config) { c.pathBits = b }
}

// Rand specifies to traverse edges from each visited node in random order.
func Rand(r *rand.Rand) Option {
	return func(c *config) { c.rand = r }
}

// Visited specifies a bits.Bits value to record visited nodes.
//
// For each node visited, the corresponding bit is set to 1.  Other bits
// are not modified.
//
// The traverse algorithm controls the traverse using a bits.Bits.  If this
// function is used, argument b will be used as the controlling value.
//
// Bits are not zeroed at the start of a traverse, so the initial Bits value
// passed in should generally be zero.  Non-zero bits will limit the traverse.
func Visited(b *bits.Bits) Option {
	return func(c *config) { c.visBits = b }
}
