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

type AttrVal struct {
	Attr string
	Val  string
}

type Config struct {
	Directed     bool
	EdgeLabel    func(int) string
	GraphAttr    []AttrVal
	Indent       string
	NodeLabel    func(graph.NI) string
	PathAttr     func([]graph.NI) (string, string)
	RankLeaves   bool
	UndirectArcs bool
}

var Defaults = Config{
	Directed:  true,
	EdgeLabel: func(l int) string { return strconv.Itoa(l) },
	Indent:    "  ",
	NodeLabel: func(n graph.NI) string { return strconv.Itoa(int(n)) },
}

func Directed(d bool) func(*Config) {
	return func(c *Config) { c.Directed = d }
}

func EdgeLabel(f func(int) string) func(*Config) {
	return func(c *Config) { c.EdgeLabel = f }
}

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

func Indent(i string) func(*Config) {
	return func(c *Config) { c.Indent = i }
}

func NodeLabel(f func(graph.NI) string) func(*Config) {
	return func(c *Config) { c.NodeLabel = f }
}

func PathAttr(f func([]graph.NI) (string, string)) func(*Config) {
	return func(c *Config) { c.PathAttr = f }
}

func RankLeaves(r bool) func(*Config) {
	return func(c *Config) { c.RankLeaves = r }
}

func UndirectArcs(u bool) func(*Config) {
	return func(c *Config) { c.UndirectArcs = u }
}

func StringAdjacencyList(g graph.AdjacencyList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteAdjacencyList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

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
	if _, err := b.WriteString("}\n"); err != nil {
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

func StringLabeledAdjacencyList(g graph.LabeledAdjacencyList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteLabeledAdjacencyList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

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

func StringFromList(g graph.FromList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteFromList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

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
	if cf.RankLeaves {
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

func StringWeightedEdgeList(g graph.WeightedEdgeList, options ...func(*Config)) (string, error) {
	var b bytes.Buffer
	if err := WriteWeightedEdgeList(g, &b, options...); err != nil {
		return "", err
	}
	return b.String(), nil
}

func WriteWeightedEdgeList(g graph.WeightedEdgeList, w io.Writer, options ...func(*Config)) error {
	cf := Defaults
	cf.Directed = false
	cf.EdgeLabel = func(l int) string {
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
			if u.To == e.N1 && u.Label == e.Label { // found reciprocal
				// write the edge
				_, err := fmt.Fprintf(b, "%s%s -- %s [label = %s]\n",
					cf.Indent, cf.NodeLabel(e.N2), cf.NodeLabel(e.N1),
					cf.EdgeLabel(e.Label))
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
		unpaired[e.N1] = append(unpaired[e.N1], graph.Half{e.N2, e.Label})
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
			cf.EdgeLabel(e.Label))
		if err != nil {
			return err
		}
	}
	return nil
}
