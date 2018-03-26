// Copyright 2018 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Package graph/io provides graph readers and writers.
package io

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/soniakeys/graph"
)

// ReadAdjacencyListBase reads a graph from the simple text format written
// by WriteAdjacencyListBase.
//
// It constructs a graph with a nodes for each line of input, reading to EOF.
// A final newline is not required.
//
// Node IDs read from the file are interpreted directly as graph.NIs.
func ReadAdjacencyList(r io.Reader) (g graph.AdjacencyList, err error) {
	return ReadAdjacencyListBase(r, 10)
}

// ReadAdjacencyListBase reads a graph from the simple text format written
// by WriteAdjacencyListBase.
//
// It constructs a graph with a nodes for each line of input, reading to EOF.
// A final newline is not required.
//
// Node IDs read from the file are interpreted directly as graph.NIs.
//
// Argument base is passed to strconv.ParseInt.
func ReadAdjacencyListBase(r io.Reader, base int) (
	g graph.AdjacencyList, err error) {
	for b := bufio.NewReader(r); ; {
		to, err := readToBase(b, base)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return g, nil
		}
		g = append(g, to)
	}
}

// read a parse a single line.  if the last line is missing the newline, it
// returns the data and err == nil.  a subsequent call will return err=io.EOF.
func readToBase(r *bufio.Reader, base int) (to []graph.NI, err error) {
	s, err := r.ReadString('\n') // but allow last line without \n
	if err != nil {
		if err != io.EOF || s == "" {
			return
		}
	}
	if f := strings.Fields(s); len(f) > 0 {
		to = make([]graph.NI, len(f))
		for x, s := range f {
			i, err := strconv.ParseInt(s, base, graph.NIBits)
			if err != nil {
				return nil, err
			}
			to[x] = graph.NI(i)
		}
	}
	return to, nil
}

// ReadAdjacencyListNIs reads a graph from the simple text format written by
// WriteAdjacencyListNIs.
//
// It constructs a graph with at least "order" nodes, but reads from "r" until
// EOF.  A graph of the specified order is allocated initially, but the graph
// will be expaneded as needed.  Zero can be passed for order.
//
// The format is similar to Go keyed composite literals.  Each line must have
// a from-NI as a "key" followed by a list of to-NIs.  Any string of characters
// not in "-0123456789" functions as a delimiter between NIs.
// A final newline is not required.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
func ReadAdjacencyListNIs(r io.Reader, order int, comment string) (
	g graph.AdjacencyList, err error) {
	return ReadAdjacencyListNIsBase(r, order, comment, 10)
}

var sep = regexp.MustCompile("[^-0123456789]+")

// ReadAdjacencyListNIsBase reads a graph from the simple text format
// written by WriteAdjacencyListNIsBase.
//
// It constructs a graph with at least "order" nodes, but reads from "r" until
// EOF.  A graph of the specified order is allocated initially, but the graph
// will be expaneded as needed.  Zero can be passed for order.
//
// The format is similar to Go keyed composite literals.  Each line must have
// a from-NI as a "key" followed by a list of to-NIs.  Any string of characters
// not in "-0123456789" functions as a delimiter between NIs.
// A final newline is not required.
//
// Argument base is passed to strconv.ParseInt.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
func ReadAdjacencyListNIsBase(r io.Reader,
	order int, comment string, base int) (g graph.AdjacencyList, err error) {
	g = make(graph.AdjacencyList, order)
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

// ReadAdjacencyListNames reads a graph from the simple text format
// written by WriteAdjacencyListNames.
//
// The format is similar to that of ReadAdjacencyListNIs but reads an input
// graph of named nodes rather than NIs.  Unique node names are accumulated
// in a list parallel to the returned graph, naming each node.  A map, with
// the reverse mapping from name to NI is also accumulated.
//
// Argument frDelim specifies a delimiter following the from-node.  If an
// empty string is passed, a space character will delimit the from-node.
//
// Argument toDelim specifies a delimiter to separate to-nodes.  If an empty
// string is passed, then any string of whitespace will separate to-nodes.
//
// Node names are also stripped of leading and trailing space.  An empty
// string is not allowed as a node name. A final newline is not required.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
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
	if frDelim == "" {
		frDelim = " "
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
		switch i := strings.Index(s, frDelim); {
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
// formatted as space separated numbers.  The NI of the from-node
// is not written but is implied by the line number.
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
	return WriteAdjacencyListNIsBase(g, w, 10, ": ", " ")
}

// WriteAdjacencyListNIsBase writes an adjacency list in a simple text format.
//
// The format is similar to Go keyed composite literals.  Each line has a
// from-NI as a "key" followed by ": " and a list of to-NIs.
// Nodes with no to-nodes generate no output, except a maximum NI with no
// to-nodes will be written with an empty to-list.
//
// Argument base is passed to strconv.FormatInt.  Argument keyDelim is
// written after the from-NI.  Argument toDelim is written separating to-NIs.
//
// Returned is number of bytes written and error.
func WriteAdjacencyListNIsBase(g graph.AdjacencyList, w io.Writer,
	base int, keyDelim, toDelim string) (n int, err error) {
	b := bufio.NewWriter(w)
	var c int
	last := len(g) - 1
	for fr, to := range g {
		if len(to) == 0 && fr < last {
			continue
		}
		c, err = b.WriteString(strconv.FormatInt(int64(fr), base))
		n += c
		if err != nil {
			return
		}
		c, err = b.WriteString(keyDelim)
		n += c
		if err != nil {
			return
		}
		if fr < last || len(to) > 0 {
			c, err = b.WriteString(strconv.FormatInt(int64(to[0]), base))
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

// WriteAdjacencyListNames writes an adjacency list with named nodes.
//
// The format is similar to that of WriteAdjacencyListNIs but allows nodes
// to be written as names rather than NIs.
//
// Argument frDelim specifies a delimiter following the from-node.
// It cannot be an empty string.
//
// Argument toDelim specifies a delimiter to separate to-nodes.
// It cannot be an empty string.
//
// Argument name is a function to translate NIs to names.
//
// Returned is number of bytes written and error.
func WriteAdjacencyListNames(g graph.AdjacencyList, w io.Writer,
	frDelim, toDelim string, name func(graph.NI) string) (n int, err error) {
	if frDelim == "" {
		return 0, fmt.Errorf("empty frDelim")
	}
	if toDelim == "" {
		return 0, fmt.Errorf("empty toDelim")
	}
	b := bufio.NewWriter(w)
	var c int
	for fr, to := range g {
		if len(to) == 0 {
			continue
		}
		c, err = b.WriteString(name(graph.NI(fr)))
		n += c
		if err != nil {
			return
		}
		c, err = b.WriteString(frDelim)
		n += c
		if err != nil {
			return
		}
		c, err = b.WriteString(name(to[0]))
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
			c, err = b.WriteString(name(to))
			n += c
			if err != nil {
				return
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
