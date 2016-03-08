// Package dot writes graphs from package graph in the Graphviz dot format.
//
// This package provides a minimal capability to output graphs simply and
// efficiently.
//
// There is no goal to provide a rich API to the many capabilities of the
// dot format.  Someday, maybe, another package.  Not now.
//
// The scheme
//
// The dot package is a separate package from graph.  It includes graph;
// graph knows nothing of dot.  This keeps the graph package uncluttered by
// file format specific code.  Dot functions are functions then, not methods
// of graph representations.
//
// For each graph representation there is a Write function that takes a graph,
// an io.Writer, and optional arguments.  For convenience, there is also a
// String function that does not require an io.Writer and simply returns the
// dot format as a string.
//
// Optional arguments are variadic and consist of calls to configuration
// functions defined in this package.  Not all configuration functions are
// meaningful for all graph types.  When a Write or String function is called
// it (1) initializes a Config struct from the package variable Defaults,
// then (2) in some cases initializes some members according to the graph type,
// then (3) calls the config functions in order.  Each config function can
// modify the Config struct.  After processing options, the funcion generates
// a dot file using the options specified in the Config struct.
package dot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/soniakeys/graph"
)

// AttrVal represents the dot format concept of an attribute-value pair.
type AttrVal struct {
	Attr string
	Val  string
}

// Config holds options that control the dot output.
//
// See Overview/Scheme for an overview of how this works.  Generally you will
// not set members of a Config struct directly.  There is an option function
// for each member.  To set a member, pass the option function as an optional
// argument to a Write or String function.
type Config struct {
	Directed     bool
	EdgeLabel    func(graph.LI) string
	GraphAttr    []AttrVal
	Indent       string
	NodeLabel    func(graph.NI) string
	UndirectArcs bool
}

// Defaults holds a package default Config struct.
//
// Defaults is copied as the first configuration step.  See Overview/Scheme.
var Defaults = Config{
	Directed:  true,
	EdgeLabel: func(l graph.LI) string { return strconv.Itoa(int(l)) },
	Indent:    "  ",
	NodeLabel: func(n graph.NI) string { return strconv.Itoa(int(n)) },
}

// Directed specifies whether to write a dot format directected or undirected
// graph.
//
// The default, Directed(true), specifies a dot format directected graph,
// or "digraph."  In this case the Write or String function outputs each arc
// of the graph as a dot format directed edge.
//
// Directed(false) specifies a dot format undirected graph or simply a "graph."
// In this case the Write or String function requires that all arcs between
// distinct nodes occur in reciprocal pairs.  For each pair the function
// outputs a single edge in dot format.
func Directed(d bool) func(*Config) {
	return func(c *Config) { c.Directed = d }
}

// EdgeLabel specifies a function to generate edge label strings for the
// dot format given the arc label integers of graph package.
//
// The default function is simply strconv.Itoa of the graph package arc label.
func EdgeLabel(f func(graph.LI) string) func(*Config) {
	return func(c *Config) { c.EdgeLabel = f }
}

// GraphAttr adds a dot format graph attribute.
//
// Graph attributes are held in a slice, and so are ordered.  This function
// updates the value of the last matching attribute if it exists, or adds a
// new attribute to the end of the list.
func GraphAttr(attr, val string) func(*Config) {
	return func(c *Config) {
		for i := len(c.GraphAttr) - 1; i >= 0; i-- {
			if c.GraphAttr[i].Attr == attr {
				c.GraphAttr[i].Val = val
				return
			}
		}
		c.GraphAttr = append(c.GraphAttr, AttrVal{attr, val})
	}
}

// Indent specifies an indent string for the body of the dot format.
//
// The default is two spaces.
func Indent(i string) func(*Config) {
	return func(c *Config) { c.Indent = i }
}

// NodeLabel specifies a function to generate node label strings for the
// dot format given the node integers of graph package.
//
// The default function is simply strconv.Itoa of the graph package node
// integer.
func NodeLabel(f func(graph.NI) string) func(*Config) {
	return func(c *Config) { c.NodeLabel = f }
}

// UndirectArcs, for the WeightedEdgeList graph type, specifies to write
// each element of the edge list as a dot file undirected edge.
//
// Note that while Directed(false) requires arcs to occur in reciprocal pairs,
// UndirectArcs(true) does not, and does not collapse reciprocal arc pairs to
// single dot format edges.
//
// See WriteWeightedEdgeList for more detail.
func UndirectArcs(u bool) func(*Config) {
	return func(c *Config) { c.UndirectArcs = u }
}

// StringAdjacencyList generates a dot format string for an AdjacencyList.
//
// See WriteAdjacencyList for options.
func StringAdjacencyList(g graph.AdjacencyList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteAdjacencyList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

// WriteAdjacencyList writes dot format text for an AdjacencyList to an
// io.Writer.
//
// Supported options:
//   Directed
//   GraphAttr
//   Indent
//   NodeLabel
func WriteAdjacencyList(g graph.AdjacencyList, w io.Writer, options ...func(*Config)) error {
	cf := Defaults
	for _, o := range options {
		o(&cf)
	}
	b := bufio.NewWriter(w)
	if err := writeHead(&cf, b); err != nil {
		return err
	}
	wf := writeALUndirected
	if cf.Directed {
		wf = writeALDirected
	}
	if err := wf(g, &cf, b); err != nil {
		return err
	}
	return writeTail(b)
}

func writeHead(cf *Config, b *bufio.Writer) error {
	t := "graph"
	if cf.Directed {
		t = "digraph"
	}
	if _, err := fmt.Fprintf(b, "%s {\n", t); err != nil {
		return err
	}
	for _, av := range cf.GraphAttr {
		_, err := fmt.Fprintf(b, "%s%s = %s\n", cf.Indent, av.Attr, av.Val)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeTail(b *bufio.Writer) error {
	if err := b.WriteByte('}'); err != nil {
		return err
	}
	return b.Flush()
}

func writeALDirected(g graph.AdjacencyList, cf *Config, b *bufio.Writer) error {
	for fr, to := range g {
		if err := writeALEdgeStmt(fr, to, "->", cf, b); err != nil {
			return err
		}
	}
	return nil
}

func writeALEdgeStmt(fr int, to []graph.NI, op string, cf *Config, b *bufio.Writer) error {
	// fast paths
	if len(to) == 0 {
		return nil
	}
	if len(to) == 1 {
		_, err := fmt.Fprintf(b, "%s%s %s %s\n",
			cf.Indent, cf.NodeLabel(graph.NI(fr)), op, cf.NodeLabel(to[0]))
		return err
	}
	// otherwise it's complicated.  we like to use a subgraph rhs to keep
	// output compact, but graphviz (some version) won't separate parallel
	// arcs in a subgraph, so in that case we write multiple edge statments.
	_, err := fmt.Fprintf(b, "%s%s %s ",
		cf.Indent, cf.NodeLabel(graph.NI(fr)), op)
	if err != nil {
		return err
	}
	var s1 big.Int
	m := map[graph.NI]int{} // multiset of defered duplicates
	c := "{"
	// first pass is over the to-list, the slice
	for _, to := range to {
		if s1.Bit(int(to)) == 0 {
			if _, err := b.WriteString(c + cf.NodeLabel(to)); err != nil {
				return err
			}
			c = " "
			s1.SetBit(&s1, int(to), 1)
		} else {
			m[to]++
		}
	}
	if _, err := b.WriteString("}\n"); err != nil {
		return err
	}
	// make additional passes over the map until it's fully consumed
	for len(m) > 0 {
		_, err := fmt.Fprintf(b, "%s%s %s ",
			cf.Indent, cf.NodeLabel(graph.NI(fr)), op)
		if err != nil {
			return err
		}
		c1 := "{"
		for n, c := range m {
			if _, err := b.WriteString(c1 + cf.NodeLabel(n)); err != nil {
				return err
			}
			if c == 1 {
				delete(m, n)
			} else {
				m[n]--
			}
			c1 = " "
		}
		if _, err := b.WriteString("}\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeALUndirected(g graph.AdjacencyList, cf *Config, b *bufio.Writer) error {
	// Similar code in undir.go at IsUndirected
	unpaired := make(graph.AdjacencyList, len(g))
	for fr, to := range g {
		// first collect unpaired subset of to
		var uto []graph.NI
	arc: // for each arc in g
		for _, to := range to {
			if to == graph.NI(fr) {
				uto = append(uto, to) // loop
				continue
			}
			// search unpaired arcs
			ut := unpaired[to]
			for i, u := range ut {
				if u == graph.NI(fr) { // found reciprocal
					last := len(ut) - 1
					ut[i] = ut[last]
					unpaired[to] = ut[:last]
					continue arc
				}
			}
			// reciprocal not found
			uto = append(uto, to)
			unpaired[fr] = append(unpaired[fr], to)
		}
		if err := writeALEdgeStmt(fr, uto, "--", cf, b); err != nil {
			return err
		}
	}
	for _, to := range unpaired {
		if len(to) > 0 {
			return fmt.Errorf("directed graph")
		}
	}
	return nil
}

// StringLabeledAdjacencyList generates a dot format string for a
// LabeledAdjacencyList.
//
// See WriteLabeledAdjacencyList for options.
func StringLabeledAdjacencyList(g graph.LabeledAdjacencyList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteLabeledAdjacencyList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

// WriteLabeledAdjacencyList writes dot format text for a LabeledAdjacencyList
// to an io.Writer.
//
// Supported options:
//   Directed
//   GraphAttr
//   Indent
//   NodeLabel
//   EdgeLabel
func WriteLabeledAdjacencyList(g graph.LabeledAdjacencyList, w io.Writer, options ...func(*Config)) error {
	cf := Defaults
	for _, o := range options {
		o(&cf)
	}
	b := bufio.NewWriter(w)
	if err := writeHead(&cf, b); err != nil {
		return err
	}
	wf := writeLALUndirected
	if cf.Directed {
		wf = writeLALDirected
	}
	if err := wf(g, &cf, b); err != nil {
		return err
	}
	return writeTail(b)
}

func writeLALDirected(g graph.LabeledAdjacencyList, cf *Config, b *bufio.Writer) error {
	for fr, to := range g {
		if err := writeLALEdgeStmt(fr, to, "->", cf, b); err != nil {
			return err
		}
	}
	return nil
}

func writeLALEdgeStmt(fr int, to []graph.Half, op string, cf *Config, b *bufio.Writer) error {
	for _, to := range to {
		_, err := fmt.Fprintf(b, "%s%s %s %s [label = %s]\n",
			cf.Indent, cf.NodeLabel(graph.NI(fr)), op, cf.NodeLabel(to.To),
			cf.EdgeLabel(to.Label))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLALUndirected(g graph.LabeledAdjacencyList, cf *Config, b *bufio.Writer) error {
	// Similar code in undir.go at IsUndirected
	unpaired := make(graph.LabeledAdjacencyList, len(g))
	for fr, to := range g {
		// first collect unpaired subset of to
		var uto []graph.Half
	arc: // for each arc in g
		for _, to := range to {
			if to.To == graph.NI(fr) {
				uto = append(uto, to) // loop
				continue
			}
			// search unpaired arcs
			ut := unpaired[to.To]
			for i, u := range ut {
				if u.To == graph.NI(fr) && u.Label == to.Label { // found reciprocal
					last := len(ut) - 1
					ut[i] = ut[last]
					unpaired[to.To] = ut[:last]
					continue arc
				}
			}
			// reciprocal not found
			uto = append(uto, to)
			unpaired[fr] = append(unpaired[fr], to)
		}
		if err := writeLALEdgeStmt(fr, uto, "--", cf, b); err != nil {
			return err
		}
	}
	for _, to := range unpaired {
		if len(to) > 0 {
			return fmt.Errorf("directed graph")
		}
	}
	return nil
}

// StringFromList (that's "String" "FromList", not "String" from "List")
// generates a dot format string for a graph.FromList.
//
// See WriteFromList for options.
func StringFromList(g graph.FromList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteFromList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

// WriteFromList writes dot format text for a graph.FromList
// to an io.Writer.
//
// Supported options:
//   Indent
//   GraphAttr
//   NodeLabel
func WriteFromList(g graph.FromList, w io.Writer, options ...func(*Config)) error {
	cf := Defaults
	GraphAttr("rankdir", "BT")(&cf)
	for _, o := range options {
		o(&cf)
	}
	b := bufio.NewWriter(w)
	if err := writeHead(&cf, b); err != nil {
		return err
	}
	for n, e := range g.Paths {
		fr := e.From
		if fr < 0 {
			continue
		}
		_, err := fmt.Fprintf(b, "%s%s -> %s\n",
			cf.Indent, cf.NodeLabel(graph.NI(n)), cf.NodeLabel(fr))
		if err != nil {
			return err
		}
	}
	if g.Leaves.BitLen() > 0 {
		if _, err := b.WriteString(cf.Indent + "{rank = same"); err != nil {
			return err
		}
		for n := range g.Paths {
			if g.Leaves.Bit(n) == 0 {
				continue
			}
			_, err := b.WriteString(" " + cf.NodeLabel(graph.NI(n)))
			if err != nil {
				return err
			}
		}
		if _, err := b.WriteString("}\n"); err != nil {
			return err
		}
	}
	return writeTail(b)
}

// StringWeightedEdgeList generates a dot format string for a
// graph.WeightedEdgeList.
//
// See WriteWeightedEdgeList for options.
func StringWeightedEdgeList(g graph.WeightedEdgeList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteWeightedEdgeList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

// WriteWeightedEdgeList writes dot format text for a graph.WeightedEdgeList
// to an io.Writer.
//
// The WeightedEdgeList, as used by the Kruskal methods, is a bit strange
// in that Kruskal interprets it as an undirected graph, but does not require
// that reciprocal edges be present.  Depending on how you construct a
// WeightedEdgeList, you may or may not have reciprocal edges.  If you do
// have reciprocal edges, the Directed(false) option is appropriate for
// collapsing reciprocals as usual and writing an undirected dot file.
// If, for Kruskal, for example, you constructed a WeightedEdgeList without
// reciprocals, then the UndirectArcs(true) is appropriate for writing an
// undirected dot file.  Specifying neither option and using the default of
// Directed(true) will produce a directed dot file.
//
// Supported options:
//   Directed
//   EdgeLabel
//   GraphAttr
//   Indent
//   NodeLabel
//   UndirectArcs
func WriteWeightedEdgeList(g graph.WeightedEdgeList, w io.Writer, options ...func(*Config)) error {
	cf := Defaults
	cf.Directed = false
	cf.EdgeLabel = func(l graph.LI) string {
		return fmt.Sprintf(`"%g"`, g.WeightFunc(l))
	}
	for _, o := range options {
		o(&cf)
	}
	if cf.UndirectArcs {
		cf.Directed = false
	}
	b := bufio.NewWriter(w)
	if err := writeHead(&cf, b); err != nil {
		return err
	}
	wf := writeWELNoRecip
	if cf.UndirectArcs || cf.Directed {
		wf = writeWELAllArcs
	}
	if err := wf(g, &cf, b); err != nil {
		return err
	}
	return writeTail(b)
}

func writeWELNoRecip(g graph.WeightedEdgeList, cf *Config, b *bufio.Writer) error {
	unpaired := make(graph.LabeledAdjacencyList, g.Order)
edge:
	for _, e := range g.Edges {
		// search unpaired arcs
		u2 := unpaired[e.N2]
		for i, u := range u2 {
			if u.To == e.N1 && u.Label == e.LI { // found reciprocal
				// write the edge
				_, err := fmt.Fprintf(b, "%s%s -- %s [label = %s]\n",
					cf.Indent, cf.NodeLabel(e.N2), cf.NodeLabel(e.N1),
					cf.EdgeLabel(e.LI))
				if err != nil {
					return err
				}
				// delete reciprocal
				last := len(u2) - 1
				u2[i] = u2[last]
				unpaired[e.N2] = u2[:last]
				continue edge
			}
		}
		// reciprocal not found
		unpaired[e.N1] = append(unpaired[e.N1], graph.Half{e.N2, e.LI})
	}
	for _, to := range unpaired {
		if len(to) > 0 {
			return fmt.Errorf("directed graph")
		}
	}
	return nil
}

func writeWELAllArcs(g graph.WeightedEdgeList, cf *Config, b *bufio.Writer) error {
	op := "--"
	if cf.Directed {
		op = "->"
	}
	for _, e := range g.Edges {
		_, err := fmt.Fprintf(b, "%s%s %s %s [label = %s]\n",
			cf.Indent, cf.NodeLabel(e.N1), op, cf.NodeLabel(e.N2),
			cf.EdgeLabel(e.LI))
		if err != nil {
			return err
		}
	}
	return nil
}
