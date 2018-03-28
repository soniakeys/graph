// Copyright 2018 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Package graph/io provides graph readers and writers.
//
// High-level "Names" functions:
//
//   ReadAdjacencyListNames
//   ReadArcNames
//   WriteAdjacencyListNames
//   WriteArcNames
//   WriteUpperNames
//
// These functions read and write graphs of named nodes, providing translation
// to and from the integer NIs used by the graph algorithms.  The translation
// might also be useful for data that has integer node IDs, but sparsely
// assigned.  The functions have flexible delimiter options.  The readers
// allow blank lines, whitespace, and have end-of-line comment options.
//
// High-level "NIs" functions:
//
//   ReadAdjacencyListNIs
//   ReadArcNIs
//   WriteAdjacencyListNIs
//   WriteArcNIs
//   WriteUpperNIs
//
// These functions read and write formats similar to the "Names" functions
// but with node NIs rather that arbitrary names.
//
// Lower-level "NIsBase" functions:
//
//   ReadAdjacencyListNIsBase
//   ReadArcNIsBase
//   WriteAdjacencyListNIsBase
//   WriteArcNIsBase
//   WriteUpperNIsBase
//
// These functions are like the "NIs" functions but read and write NIs in a
// specified base.  (Mostly just because the strconv functions do it, but
// higher bases will represent data somewhat more compactly.)
//
// Lower-level node-implicit functions:
//
//   ReadAdjacencyList
//   ReadAdjacencyListBase
//   WriteAdjacencyList
//   WriteAdjacencyListBase
//   WriteUpper
//   WriteUpperBase
//
// These functions read and write a more primitive format of to-lists, with
// the from-NI implied, derived by counting input lines.
package io

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/soniakeys/graph"
)

type Text struct {
	Comment string
	FrDelim string
	ToDelim string
	Base    int
}

func NewText() *Text {
	return &Text{Comment: "//", Base: 10}
}

// ReadAdjacencyList reads a graph from a simple text format.
//
// The format has a line for each node of the graph, each line consisting of
// NIs of a node's to-list.  The from-node is implied by the (0-based) line
// number of input text.  To-node IDs read from the file are interpreted
// directly as graph.NIs.
//
// ReadAdjacencyList reads to EOF.
//
// ReadAdjacencyList will read text written by WriteAdjacencyList.
func (tr *Text) ReadAdjacencyList(r io.Reader) (
	g graph.AdjacencyList, err error) {
	sep := bx(tr.Base)
	for b := bufio.NewReader(r); ; {
		to, err := tr.readToBase(b, sep)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return g, nil
		}
		g = append(g, to)
	}
}

// read and parse a single line of a to-list.  if the last line is missing
// the newline, it returns the data and err == nil.  a subsequent call will
// return io.EOF.
func (tr *Text) readToBase(r *bufio.Reader, sep *regexp.Regexp) (
	to []graph.NI, err error) {
	s, err := r.ReadString('\n') // but allow last line without \n
	if err != nil {
		if err != io.EOF || s == "" {
			return
		}
	}
	if s == "\n" { // fast path for blank line
		return
	}
	// s is non-empty at this point
	f := sep.Split(s, -1)
	// some minimal possibilities:
	// single sep: "\n" -> ["", ""], the empty strings before and after sep
	// single non-sep: "0" -> ["0"], a non-empty string
	if f[0] == "" {
		f = f[1:] // toss "" before leading any separator
	}
	last := len(f) - 1
	if f[last] == "" {
		f = f[:last] // toss "" after trailing separator, usually the \n
	}
	to = make([]graph.NI, len(f))
	for x, s := range f {
		i, err := strconv.ParseInt(s, tr.Base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		to[x] = graph.NI(i)
	}
	return to, nil
}

// ReadAdjacencyListNIs reads a graph from a simple text format keyed by
// node NI.
//
// The format is similar to Go keyed composite literals in that each line has
// a from-NI as a "key" followed by a list of to-NIs.  NIs are delimited by
// any non-empty string of characters not in "+-0123456789".  A final newline
// is not required.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
//
// ReadAdjacencyListNIs will read text written by WriteAdjacencyListNIs.
func ReadAdjacencyListNIs(r io.Reader, comment string) (
	graph.AdjacencyList, error) {
	return ReadAdjacencyListNIsBase(r, comment, 10)
}

var bx10 = regexp.MustCompile("[^0-9-+]+")

// return regular expression delimiting numbers of given base.
func bx(base int) *regexp.Regexp {
	if base == 10 || base == 0 {
		return bx10
	}
	expr := ""
	if base <= 10 {
		expr = fmt.Sprintf("[^0-%d-+]+", base-1)
	} else {
		expr = fmt.Sprintf("[^0-9a-%c-+]+", 'a'+base-11)
	}
	return regexp.MustCompile(expr)
}

// ReadAdjacencyListNIsBase reads a graph from the simple text format keyed by
// node NI with NIs in an arbitrary base.
//
// Like ReadAdjacencyListNIs but parsing numbers in the given base. Argument
// base is passed to strconv.ParseInt.
func ReadAdjacencyListNIsBase(r io.Reader, comment string, base int) (
	g graph.AdjacencyList, err error) {
	sep := bx(base)
	b := bufio.NewReader(r)
	for {
		s, err := b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if s == "" {
				return g, nil
			}
			// allow last line without \n
		}
		if comment > "" {
			if i := strings.Index(s, comment); i >= 0 {
				s = s[:i]
			}
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue // allow and ignore blank line
		}
		f := sep.Split(s, -1)
		// non-blank line must start with valid from-NI, not delimiters
		if f[0] == "" {
			return nil, fmt.Errorf("invalid: %q", s)
		}
		fr, err := strconv.ParseInt(f[0], base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		for int(fr) >= len(g) {
			g = append(g, nil)
		}
		if len(f) == 1 {
			continue // from-NI with no to-list is allowed.
		}
		f = f[1:]
		last := len(f) - 1
		if f[last] == "" {
			f = f[:last] // allow line to end with delimiters
		}
		to := make([]graph.NI, len(f))
		for x, s := range f {
			i, err := strconv.ParseInt(s, base, graph.NIBits)
			if err != nil {
				return nil, err
			}
			to[x] = graph.NI(i)
		}
		if g[fr] == nil {
			g[fr] = to
		} else {
			g[fr] = append(g[fr], to...)
		}
	}
}

// ReadAdjacencyListNames reads a graph from a simple text format using
// arbitrary node names.
//
// The format is similar to that of ReadAdjacencyListNIs but represents an
// input graph of named nodes rather than NIs.
//
// As data is read, unique node names are accumulatedin a list parallel to the
// returned graph, naming each node.  A map, with the reverse mapping from name
// to NI is also accumulated.
//
// Argument frDelim specifies a delimiter following the from-node.  If an
// empty string is passed, then any string of whitespace will delimit the
// from-node.
//
// Argument toDelim specifies a delimiter to separate to-nodes.  If an empty
// string is passed, then any string of whitespace will separate to-nodes.
//
// Node names are also stripped of leading and trailing space.  An empty
// string is not allowed as a node name. A final newline is not required.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
//
// ReadAdjacencyListNames will read text written by WriteAdjacencyListNames.
func ReadAdjacencyListNames(r io.Reader, frDelim, toDelim, comment string) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	ni = map[string]graph.NI{}
	getNI := func(s string) graph.NI {
		n, ok := ni[s]
		if !ok {
			n = graph.NI(len(g))
			g = append(g, nil)
			name = append(name, s)
			ni[s] = n
		}
		return n
	}
	frSpace := strings.TrimSpace(frDelim) == ""
	toSpace := strings.TrimSpace(toDelim) == ""
	frIndex := func(s string) int { return strings.Index(s, frDelim) }
	if frSpace {
		frIndex = func(s string) int {
			return strings.IndexFunc(s, unicode.IsSpace)
		}
	}
	b := bufio.NewReader(r)
	s := ""
	for line := 1; ; line++ {
		s, err = b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, fmt.Errorf("line %d: %v", line, err)
			}
			if s == "" {
				err = nil
				return
			}
			// allow last line without \n
		}
		if comment > "" {
			if i := strings.Index(s, comment); i >= 0 {
				s = s[:i]
			}
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue // allow and ignore blank line
		}
		f0 := s
		switch i := frIndex(s); {
		case i >= 0:
			f0 = strings.TrimSpace(s[:i])
			if f0 == "" {
				return nil, nil, nil,
					fmt.Errorf("line %d: blank name not allowed", line)
			}
			s = s[i+len(frDelim):]
		case frSpace:
			s = "" // allowed is from-NI only, no to-list
		default:
			return nil, nil, nil,
				fmt.Errorf("line %d: from-delimiter required", line)
		}
		fr := getNI(f0)
		var to []graph.NI
		if toSpace {
			for _, s := range strings.Fields(s) {
				to = append(to, getNI(s))
			}
		} else {
			f := strings.Split(s, toDelim)
			if last := len(f) - 1; strings.TrimSpace(f[last]) == "" {
				f = f[:last] // allow trailing delimiter
			}
			for _, s := range f {
				if s = strings.TrimSpace(s); s == "" {
					return nil, nil, nil,
						fmt.Errorf("line %d: blank name not allowed", line)
				}
				to = append(to, getNI(s))
			}
		}
		g[fr] = to
	}
}

/* order functions commented out for now.  the idea is that these would
be useful as subroutines for more complex formats.  uncomment as needed...

// ReadAdjacencyListOrder reads a graph from the simple text format written by
// WriteAdjacencyList.
//
// It constructs a graph with "order" nodes, and reads order lines from reader
// "r".  An error will occur if r cannot supply order lines.  A final newline
// is not required.
//
// Node IDs read from the file are interpreted directly as graph.NIs.
func ReadAdjacencyListOrder(r io.Reader, order int) (
	g graph.AdjacencyList, err error) {
	return ReadAdjacencyListOrderBase(r, order, 10)
}

// ReadAdjacencyListOrderBase reads a graph from the simple text format written
// by WriteAdjacencyListBase.
//
// It constructs a graph with "order" nodes, and reads order lines from reader
// "r".  An error will occur if r cannot supply order lines.  A final newline
// is not required.
//
// Node IDs read from the file are interpreted directly as graph.NIs.
//
// Argument base is passed to strconv.ParseInt.
func ReadAdjacencyListOrderBase(r io.Reader, order, base int) (
	g graph.AdjacencyList, err error) {
	g = make(graph.AdjacencyList, order)
	last := order - 1
	b := bufio.NewReader(r)
	for n := range g {
		if g[n], err = readToBase(b, base); err != nil {
			switch {
			case err != io.EOF:
				g = nil
			case n == last:
				err = nil
			}
			return
		}
	}
	return
}
*/

// ReadArcNIs reads a graph from the simple text format written by WriteArcNIs.
//
// It constructs a graph with an arc for each line of input, where a line
// has a from-NI and to-NI.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
func ReadArcNIs(r io.Reader, comment string) (graph.AdjacencyList, error) {
	return ReadArcNIsBase(r, comment, 10)
}

// ReadArcNIsBase reads a graph from the simple text format written by
// WriteArcNIsBase.
//
// It constructs a graph with an arc for each line of input, where a line
// has a from-NI and to-NI.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
//
// Argument base is passed to strconv.ParseInt.
func ReadArcNIsBase(r io.Reader, comment string, base int) (graph.AdjacencyList, error) {
	sep := bx(base)
	var max int64
	e := map[int][]graph.NI{}
	for b := bufio.NewReader(r); ; {
		s, err := b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if s == "" { // normal return
				g := make(graph.AdjacencyList, max+1)
				for fr := range g {
					g[fr] = e[fr]
				}
				return g, nil
			}
			// allow last line without \n
		}
		if comment > "" {
			if i := strings.Index(s, comment); i >= 0 {
				s = s[:i]
			}
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue // blank line
		}
		f := sep.Split(s, 2)
		if len(f) != 2 {
			return nil, fmt.Errorf("invalid: %q", s)
		}
		fr, err := strconv.ParseInt(strings.TrimSpace(f[0]), base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		if fr < 0 {
			return nil, fmt.Errorf("invalid: %q", s)
		}
		if fr > max {
			max = fr
		}
		to, err := strconv.ParseInt(strings.TrimSpace(f[1]), base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		if to > max {
			max = to
		}
		e[int(fr)] = append(e[int(fr)], graph.NI(to))
	}
}

// ReadArcNames reads a graph from the simple text format written by
// WriteArcNames.
//
// The format is similar to that of ReadArcNIs but reads an input graph of
// named nodes rather than NIs.  Unique node names are accumulated
// in a list parallel to the returned graph, naming each node.  A map, with
// the reverse mapping from name to NI is also accumulated.
//
// An empty string is not allowed as a node name. A final newline is not
// required.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
func ReadArcNames(r io.Reader, delim, comment string) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	ni = map[string]graph.NI{}
	getNI := func(s string) graph.NI {
		n, ok := ni[s]
		if !ok {
			n = graph.NI(len(g))
			g = append(g, nil)
			name = append(name, s)
			ni[s] = n
		}
		return n
	}
	delIndex := func(s string) int { return strings.Index(s, delim) }
	if delim == "" {
		delIndex = func(s string) int {
			return strings.IndexFunc(s, unicode.IsSpace)
		}
	}
	b := bufio.NewReader(r)
	s := ""
	for line := 1; ; line++ {
		s, err = b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, fmt.Errorf("line %d: %v", line, err)
			}
			if s == "" {
				err = nil
				return
			}
			// allow last line without \n
		}
		if comment > "" {
			if i := strings.Index(s, comment); i >= 0 {
				s = s[:i]
			}
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue // allow and ignore blank line
		}
		i := delIndex(s)
		if i < 0 {
			return nil, nil, nil,
				fmt.Errorf("line %d: delimiter required", line)
		}
		fn := strings.TrimSpace(s[:i])
		if fn == "" {
			return nil, nil, nil,
				fmt.Errorf("line %d: blank name not allowed", line)
		}
		tn := strings.TrimSpace(s[i+len(delim):])
		if tn == "" {
			return nil, nil, nil,
				fmt.Errorf("line %d: blank name not allowed", line)
		}
		fr := getNI(fn)
		g[fr] = append(g[fr], getNI(tn))
	}
}

// ReadLabeledAdjacencyList reads a graph from a simple text format.
//
// The format has a line for each node of the graph, each line consisting of
// a list of pairs of numbers, each pair being the NI and LI of a graph.Half
// of the node's to-list.  The from-node is implied by the (0-based) line
// number of input text.  Graph.Half elements read from the file are
// interpreted directly as graph.NIs and graph.LIs.
//
// ReadLabeledAdjacencyList reads to EOF.
//
// ReadLabeledAdjacencyList will read text written by WriteLabeledAdjacencyList.
func ReadLabeledAdjacencyList(r io.Reader) (graph.LabeledAdjacencyList, error) {
	return ReadLabeledAdjacencyListBase(r, 10)
}

// ReadLabeledAdjacencyListBase reads a graph from a simple text format with NIs
// in an arbitrary base.
//
// Like ReadLabeledAdjacencyList but parsing numbers in the given base.
//
// Argument base is passed to strconv.ParseInt.
func ReadLabeledAdjacencyListBase(r io.Reader, base int) (
	g graph.LabeledAdjacencyList, err error) {
	sep := bx(base)
	for b := bufio.NewReader(r); ; {
		to, err := readHalfBase(b, base, sep)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return g, nil
		}
		g = append(g, to)
	}
}

// read and parse a single line of a to-list.  if the last line is missing
// the newline, it returns the data and err == nil.  a subsequent call will
// return io.EOF.
func readHalfBase(r *bufio.Reader, base int, sep *regexp.Regexp) (to []graph.Half, err error) {
	s, err := r.ReadString('\n') // but allow last line without \n
	if err != nil {
		if err != io.EOF || s == "" {
			return
		}
	}
	if s == "\n" { // fast path for blank line
		return
	}
	// s is non-empty at this point
	f := sep.Split(s, -1)
	// some minimal possibilities:
	// single sep: "\n" -> ["", ""], the empty strings before and after sep
	// single non-sep: "0" -> ["0"], a non-empty string
	if f[0] == "" {
		f = f[1:] // toss "" before leading any separator
	}
	last := len(f) - 1
	if f[last] == "" {
		f = f[:last] // toss "" after trailing separator, usually the \n
	}
	if len(f)%2 != 0 {
		return nil, fmt.Errorf("odd data")
	}
	to = make([]graph.Half, len(f)/2)
	y := 0
	for x := range to {
		ni, err := strconv.ParseInt(f[y], base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		y++
		li, err := strconv.ParseInt(f[y], base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		y++
		to[x] = graph.Half{graph.NI(ni), graph.LI(li)}
	}
	return to, nil
}

// WriteAdjacencyList writes an adjacency list in a simple text format.
//
// A line is written for each node, consisting of the to-list, the NIs
// formatted as space separated base 10 numbers.  The NI of the from-node
// is not written but is implied by the line number.
//
// Returned is number of bytes written and error.
func WriteAdjacencyList(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	return WriteAdjacencyListBase(g, w, 10)
}

// WriteAdjacencyListBase writes an adjacency list in a simple text format.
//
// A line is written for each node, consisting of the to-list, the NIs
// formatted as space separated numbers formatted in the specified base.
// The NI of the from-node is not written but is implied by the line number.
//
// Argument base is passed to strconv.FormatInt.
//
// Returned is number of bytes written and error.
func WriteAdjacencyListBase(g graph.AdjacencyList, w io.Writer, base int) (
	n int, err error) {
	b := bufio.NewWriter(w)
	var c int
	for _, to := range g {
		if len(to) > 0 {
			c, err = b.WriteString(strconv.FormatInt(int64(to[0]), base))
			n += c
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to), base))
				n += c + 1
				if err != nil {
					return
				}
			}
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}

// WriteAdjacencyListNIs writes an adjacency list in a simple text format.
//
// The format is similar to Go keyed composite literals.  Each line has a
// from-NI as a "key" followed by ": " and a list of to-NIs.
// Nodes with no to-nodes generate no output, except a maximum NI with no
// to-nodes will be written with an empty to-list.
//
// Returned is number of bytes written and error.
func WriteAdjacencyListNIs(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	return WriteAdjacencyListNIsBase(g, w, ": ", " ", 10)
}

// WriteAdjacencyListNIsBase writes an adjacency list in a simple text format.
//
// The format is similar to Go keyed composite literals.  Each line has a
// from-NI as a "key" followed by a delimiter and a list of to-NIs.
// Nodes with no to-nodes generate no output, except a maximum NI with no
// to-nodes will be written with an empty to-list.
//
// Argument base is passed to strconv.FormatInt.  Argument frDelim is
// written after the from-NI.  Argument toDelim is written separating to-NIs.
//
// Returned is number of bytes written and error.
func WriteAdjacencyListNIsBase(g graph.AdjacencyList, w io.Writer,
	frDelim, toDelim string, base int) (n int, err error) {
	return writeAdjacencyListFrom(g, w, frDelim, toDelim,
		func(n graph.NI) string { return strconv.FormatInt(int64(n), base) },
		true)
}

func writeAdjacencyListFrom(g graph.AdjacencyList, w io.Writer,
	frDelim, toDelim string, format func(graph.NI) string,
	writeLast bool) (n int, err error) {
	if frDelim == "" {
		frDelim = ": "
	}
	if toDelim == "" {
		toDelim = " "
	}
	b := bufio.NewWriter(w)
	var c int
	last := len(g) - 1
	for i, to := range g {
		switch {
		case len(to) > 0:
			fr := graph.NI(i)
			c, err = b.WriteString(format(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(frDelim)
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(format(to[0]))
			n += c
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				c, err = b.WriteString(toDelim)
				n += c
				if err != nil {
					return
				}
				c, err = b.WriteString(format(to))
				n += c
				if err != nil {
					return
				}
			}
		case i == last && writeLast:
			fr := graph.NI(i)
			c, err = b.WriteString(format(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(
				strings.TrimRightFunc(frDelim, unicode.IsSpace))
			n += c
			if err != nil {
				return
			}
		default:
			continue // without writing \n
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}

// WriteAdjacencyListNames writes an adjacency list with named nodes.
//
// The format is similar to that of WriteAdjacencyListNIs but allows nodes
// to be written as names rather than NIs.
//
// Argument frDelim specifies a delimiter following the from-node.
// If blank, ": " is used.
//
// Argument toDelim specifies a delimiter to separate to-nodes.
// If blank " " is used.
//
// Argument name is a function to translate NIs to names.
//
// Returned is number of bytes written and error.
func WriteAdjacencyListNames(g graph.AdjacencyList, w io.Writer,
	frDelim, toDelim string, name func(graph.NI) string) (n int, err error) {
	return writeAdjacencyListFrom(g, w, frDelim, toDelim, name, false)
}

// WriteArcNIs writes arcs of an adjacency list in a simple text format.
//
// Each line has an arc, a from-NI and to-NI, separated with a space.
//
// Returned is number of bytes written and error.
func WriteArcNIs(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	return WriteArcNIsBase(g, w, " ", 10)
}

// WriteArcNIsBase writes arcs of an adjacency list in a simple text format.
//
// Each line has a full arc, a from-NI and to-NI, separated with a delimiter.
// If argument delim is the empty string, NIs will be separated with a space.
//
// Argument base is passed to strconv.FormatInt.
//
// Returned is number of bytes written and error.
func WriteArcNIsBase(g graph.AdjacencyList, w io.Writer,
	delim string, base int) (n int, err error) {
	if delim == "" {
		delim = " "
	}
	b := bufio.NewWriter(w)
	var c int
	for fr, to := range g {
		for _, to := range to {
			c, err = b.WriteString(strconv.FormatInt(int64(fr), base))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(delim)
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(strconv.FormatInt(int64(to), base))
			n += c
			if err != nil {
				return
			}
			if err = b.WriteByte('\n'); err != nil {
				return
			}
			n++
		}
	}
	b.Flush()
	return
}

// WriteArcNames writes arcs of an adjacency list with named nodes.
//
// The format is similar to that of WriteArcNIs but allows nodes
// to be written as names rather than NIs.
//
// If argument delim is the empty string, NIs will be separated with a space.
//
// Argument name is a function to translate NIs to names.
//
// Returned is number of bytes written and error.
func WriteArcNames(g graph.AdjacencyList, w io.Writer,
	delim string, name func(graph.NI) string) (n int, err error) {
	if delim == "" {
		delim = " "
	}
	b := bufio.NewWriter(w)
	var c int
	for fr, to := range g {
		for _, to := range to {
			c, err = b.WriteString(name(graph.NI(fr)))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(delim)
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(name(to))
			n += c
			if err != nil {
				return
			}
			if err = b.WriteByte('\n'); err != nil {
				return
			}
			n++
		}
	}
	b.Flush()
	return
}

// WriteUpper writes the "upper triangle" of an adjacency list.
//
// For an adjacency list representing an undirected graph, this writes
// an arc for each undirected edge but omits reciprocal arcs.
//
// A line is written for each node, consisting of to-NIs greater than or
// equal to the from-NI.  The from-NI is not written but is implied by the
// line number.
//
// Returned is number of bytes written and error.
func WriteUpper(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	return WriteUpperBase(g, w, 10)
}

// WriteUpperBase writes the "upper triangle" of an adjacency list.
//
// Like WriteUpper but formatting NIs in the given base.
// Argument base is passed to strconv.FormatInt.
//
// Returned is number of bytes written and error.
func WriteUpperBase(g graph.AdjacencyList, w io.Writer, base int) (
	n int, err error) {
	b := bufio.NewWriter(w)
	var c int
	for i, to := range g {
		fr := graph.NI(i)
		one := false
		for _, to := range to {
			if to >= fr {
				if one {
					if err = b.WriteByte(' '); err != nil {
						return
					}
					n++
				} else {
					one = true
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to), base))
				n += c
				if err != nil {
					return
				}
			}
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}

// WriteUpperNames writes the "upper triangle" of an adjacency list.
//
// Returned is number of bytes written and error.
func WriteUpperNames(g graph.AdjacencyList, w io.Writer,
	frDelim, toDelim string, format func(graph.NI) string) (n int, err error) {
	return writeUpper(g, w, frDelim, toDelim, format, false)
}

// WriteUpperNIs writes the "upper triangle" of an adjacency list.
//
// For an adjacency list representing an undirected graph, this writes
// an arc for each undirected edge but omits reciprocal arcs.
//
// A line is written for each from-node that has to-nodes >= the from-node.
// The format is similar to Go keyed composite literals.  Each line has the
// from-NI as a "key" followed by ": " and a list of to-NIs.  Nodes with no
// to-nodes >= the from-node generate no output, except a maximum NI with no
// to-nodes will be written with an empty to-list.
//
// Returned is number of bytes written and error.
func WriteUpperNIs(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	return WriteUpperNIsBase(g, w, ": ", " ", 10)
}

// WriteUpperNIsBase writes the "upper triangle" of an adjacency list.
//
// Like WriteUpperNis but formatting NIs in the given base.
// Argument base is passed to strconv.FormatInt.
//
// Returned is number of bytes written and error.
func WriteUpperNIsBase(g graph.AdjacencyList, w io.Writer,
	frDelim, toDelim string, base int) (n int, err error) {
	return writeUpper(g, w, frDelim, toDelim, func(n graph.NI) string {
		return strconv.FormatInt(int64(n), base)
	}, true)
}

func writeUpper(g graph.AdjacencyList, w io.Writer, frDelim, toDelim string,
	format func(graph.NI) string, writeLast bool) (n int, err error) {
	b := bufio.NewWriter(w)
	if frDelim == "" {
		frDelim = ": "
	}
	if toDelim == "" {
		toDelim = " "
	}
	var c int
	last := len(g) - 1
	for i, to := range g {
		fr := graph.NI(i)
		one := false
		for _, to := range to {
			if to >= fr {
				if !one {
					one = true
					c, err = b.WriteString(format(fr))
					n += c
					if err != nil {
						return
					}
					c, err = b.WriteString(frDelim)
					n += c
					if err != nil {
						return
					}
				} else {
					c, err = b.WriteString(toDelim)
					n += c
					if err != nil {
						return
					}
				}
				c, err = b.WriteString(format(to))
				n += c
				if err != nil {
					return
				}
			}
		}
		if writeLast && i == last && !one {
			one = true
			c, err = b.WriteString(format(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(
				strings.TrimRightFunc(frDelim, unicode.IsSpace))
			n += c
			if err != nil {
				return
			}
		}
		if one {
			if err = b.WriteByte('\n'); err != nil {
				return
			}
			n++
		}
	}
	b.Flush()
	return
}

// WriteLabeledAdjacencyList writes an adjacency list in a simple text format.
//
// A line is written for each node, consisting of the to-list, the NIs
// formatted as space separated base 10 numbers.  The NI of the from-node
// is not written but is implied by the line number.
//
// Returned is number of bytes written and error.
func WriteLabeledAdjacencyList(g graph.LabeledAdjacencyList, w io.Writer) (
	n int, err error) {
	return WriteLabeledAdjacencyListBase(g, w, 10)
}

// WriteLabeledAdjacencyListBase writes an adjacency list in a simple text
// format.
//
// A line is written for each node, consisting of the to-list, the NIs
// formatted as space separated numbers formatted in the specified base.
// The NI of the from-node is not written but is implied by the line number.
//
// Argument base is passed to strconv.FormatInt.
//
// Returned is number of bytes written and error.
func WriteLabeledAdjacencyListBase(g graph.LabeledAdjacencyList, w io.Writer,
	base int) (n int, err error) {
	b := bufio.NewWriter(w)
	var c int
	for _, to := range g {
		if len(to) > 0 {
			c, err = b.WriteString(strconv.FormatInt(int64(to[0].To), base))
			n += c
			if err != nil {
				return
			}
			if err = b.WriteByte(' '); err != nil {
				return
			}
			c, err = b.WriteString(strconv.FormatInt(int64(to[0].Label), base))
			n += c + 1
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to.To), base))
				n += c + 1
				if err != nil {
					return
				}
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to.Label), base))
				n += c + 1
				if err != nil {
					return
				}
			}
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}
