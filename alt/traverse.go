// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import (
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

type config struct {
	// values initialized directly from parameters
	start          graph.NI
	arcVisitor     func(n graph.NI, x int)
	iterateFrom    func(n graph.NI)
	levelVisitor   func(l int, n []graph.NI)
	nodeVisitor    func(n graph.NI)
	okArcVisitor   func(n graph.NI, x int) bool
	okNodeVisitor  func(n graph.NI) bool
	okLevelVisitor func(l int, n []graph.NI) bool
	rand           *rand.Rand
	visBits        *bits.Bits
	pathBits       *bits.Bits
	fromList       *graph.FromList

	// other stuff initialized in constructor
	rp   []graph.PathEnd     // fromList.Paths
	nvis func(graph.NI) bool // not-visited test
}

func newConfig(g graph.AdjacencyList, start graph.NI, opt []TraverseOption) *config {
	cf := &config{start: start}
	for _, o := range opt {
		o(cf)
	}
	// either visBits or fromList are suitable for recording visited nodes.
	// if neither is specified as an option, allocate bits.
	if cf.fromList == nil {
		if cf.visBits == nil {
			b := bits.New(len(g))
			cf.visBits = &b
		}
	} else {
		if cf.fromList.Paths == nil {
			*cf.fromList = graph.NewFromList(len(g))
		}
		cf.rp = cf.fromList.Paths
	}
	if cf.visBits != nil {
		if cf.visBits.Bit(int(cf.start)) == 1 {
			return nil
		}
		cf.visBits.SetBit(int(cf.start), 1)
		cf.nvis = func(n graph.NI) bool { return cf.visBits.Bit(int(n)) == 0 }
	} else {
		cf.nvis = func(n graph.NI) bool { return cf.rp[n].Len == 0 }
	}
	if cf.rp != nil {
		if cf.rp[cf.start].Len > 0 {
			return nil
		}
		cf.rp[cf.start] = graph.PathEnd{Len: 1, From: -1}
	}
	return cf
}

// A TraverseOption specifies an option for a breadth first or depth first
// traversal.
//
// Values of this type are returned by various TraverseOption constructor
// functions.  These constructors take optional values to be used
// in a traversal and wrap them in the TraverseOption type.  This type
// is actually a function.  The BreadthFirst and DepthFirst traversal
// methods call these functions in order, to initialize state that controls
// the traversal.
type TraverseOption func(*config)

// ArcVisitor specifies a visitor function to call at each arc.
//
// See also OkArcVisitor.
func ArcVisitor(v func(n graph.NI, x int)) TraverseOption {
	return func(c *config) {
		c.arcVisitor = v
	}
}

// From specifies a graph.FromList to populate.
func From(f *graph.FromList) TraverseOption {
	return func(c *config) {
		c.fromList = f
	}
}

// LevelVisitor specifies a visitor function to call at each level or depth.
//
// The level visitor function is called before any node or arc visitor
// functions.
//
// See also OkLevelVisitor.
func LevelVisitor(v func(level int, nodes []graph.NI)) TraverseOption {
	return func(c *config) {
		c.levelVisitor = v
	}
}

// NodeVisitor specifies a visitor function to call at each node.
//
// The node visitor function is called before any arc visitor functions.
//
// See also OkNodeVisitor.
func NodeVisitor(v func(graph.NI)) TraverseOption {
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
func OkArcVisitor(v func(n graph.NI, x int) bool) TraverseOption {
	return func(c *config) {
		c.okArcVisitor = v
	}
}

// OKLevelVisitor specifies a visitor function to call at each level or depth,
// returning a boolean result
//
// As long as v returns a result of true, the traverse progresses to traverse
// all nodes.  If v returns false, the traverse terminates immediately.
//
// The level visitor function is called before any node or arc visitor
// functions.
//
// See also LevelVisitor.
func OkLevelVisitor(v func(level int, nodes []graph.NI) bool) TraverseOption {
	return func(c *config) {
		c.okLevelVisitor = v
	}
}

// OkNodeVisitor specifies a visitor function to perform some test at each node
// and return a boolean result.
//
// As long as v returns a result of true, the traverse progresses to traverse
// all nodes.  If v returns false, the traverse terminates immediately.
//
// The node visitor function is called before any arc visitor functions.
//
// See also NodeVisitor.
func OkNodeVisitor(v func(graph.NI) bool) TraverseOption {
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
func PathBits(b *bits.Bits) TraverseOption {
	return func(c *config) { c.pathBits = b }
}

// Rand specifies to traverse edges from each visited node in random order.
func Rand(r *rand.Rand) TraverseOption {
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
func Visited(b *bits.Bits) TraverseOption {
	return func(c *config) { c.visBits = b }
}

// BreadthFirst traverses a directed or undirected graph in breadth first order.
//
// Argument start is the start node for the traversal.  Argument opt can be
// any number of values returned by a supported TraverseOption function.
//
// Supported:
//
//   From
//   ArcVisitor
//   LevelVisitor
//   NodeVisitor
//   OkArcVisitor
//   OkLevelVisitor
//   OkNodeVisitor
//   Rand
//   Visited
//
// Unsupported:
//
//   PathBits
//
// See also alt.BreadthFirst2, a direction optimizing breadth first algorithm.
func BreadthFirst(g graph.AdjacencyList, start graph.NI, options ...TraverseOption) {
	cf := newConfig(g, start, options)
	if cf == nil {
		return
	}
	// the frontier consists of nodes all at the same level
	frontier := []graph.NI{cf.start}
	level := 1
	var next []graph.NI
	visitNode := func(n graph.NI) bool {
		// visit nodes as they come off frontier
		if cf.nodeVisitor != nil {
			cf.nodeVisitor(n)
		}
		if cf.okNodeVisitor != nil {
			if !cf.okNodeVisitor(n) {
				return false
			}
		}
		for x, nb := range g[n] {
			if cf.arcVisitor != nil {
				cf.arcVisitor(n, x)
			}
			if cf.okArcVisitor != nil {
				if !cf.okArcVisitor(n, x) {
					return false
				}
			}
			if cf.nvis(nb) {
				next = append(next, nb)
				if cf.visBits != nil {
					cf.visBits.SetBit(int(nb), 1)
				}
				if cf.rp != nil {
					cf.rp[nb] = graph.PathEnd{From: n, Len: level}
				}
			}
		}
		return true
	}
	for {
		if cf.rand != nil {
			cf.rand.Shuffle(len(frontier), func(i, j int) {
				frontier[i], frontier[j] = frontier[j], frontier[i]
			})
		}
		if cf.levelVisitor != nil {
			cf.levelVisitor(level, frontier)
		}
		if cf.okLevelVisitor != nil && !cf.okLevelVisitor(level, frontier) {
			return
		}
		if cf.fromList != nil {
			cf.fromList.MaxLen = level
		}
		level++
		for _, n := range frontier {
			if !visitNode(n) {
				return
			}
		}
		if len(next) == 0 {
			break
		}
		frontier, next = next, nil
	}
}

// DepthFirst traverses a directed or undirected graph in depth first order.
//
// Argument start is the start node for the traversal.  Argument opt can be
// any number of values returned by a supported TraverseOption function.
//
// Supported:
//
//   From
//   ArcVisitor
//   NodeVisitor
//   OkArcVisitor
//   OkNodeVisitor
//   PathBits
//   Rand
//   Visited
//
// Unsupported:
//
//   LevelVisitor
//   OkLevelVisitor
func DepthFirst(g graph.AdjacencyList, start graph.NI, options ...TraverseOption) {
	cf := newConfig(g, start, options)
	if cf == nil {
		return
	}
	if cf.pathBits != nil {
		cf.pathBits.ClearAll()
	}
	var dfArc func(graph.NI, graph.NI, int, int) bool
	dfNode := func(n graph.NI, level int) bool {
		if cf.visBits != nil {
			cf.visBits.SetBit(int(n), 1)
		}
		if cf.pathBits != nil {
			cf.pathBits.SetBit(int(n), 1)
		}
		if cf.nodeVisitor != nil {
			cf.nodeVisitor(n)
		}
		if cf.okNodeVisitor != nil {
			if !cf.okNodeVisitor(n) {
				return false
			}
		}
		if cf.rand == nil {
			for x, to := range g[n] {
				if !dfArc(n, to, x, level) {
					return false
				}
			}
		} else {
			to := g[n]
			for _, x := range cf.rand.Perm(len(to)) {
				if !dfArc(n, to[x], x, level) {
					return false
				}
			}
		}
		if cf.pathBits != nil {
			cf.pathBits.SetBit(int(n), 0)
		}
		return true
	}
	dfArc = func(fr, to graph.NI, x, level int) bool {
		if cf.arcVisitor != nil {
			cf.arcVisitor(fr, x)
		}
		if cf.okArcVisitor != nil {
			if !cf.okArcVisitor(fr, x) {
				return false
			}
		}
		if !cf.nvis(to) {
			return true
		}
		if cf.rp != nil {
			cf.rp[fr] = graph.PathEnd{From: fr, Len: level}
		}
		return dfNode(to, level+1)
	}
	dfNode(cf.start, 1)
}
