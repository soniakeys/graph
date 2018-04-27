// Copyright 2018 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Package graph/io provides graph readers and writers.
package io

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/soniakeys/graph"
)

// WriteAdjacencyList writes an adjacency list as text.
//
// Fields of the receiver Text define how the text data is formatted.
// See documentation of the Text struct.
//
// Returned is number of bytes written and error.
func (t Text) WriteAdjacencyList(g graph.AdjacencyList, w io.Writer) (
	int, error) {
	switch t.Format {
	case Sparse:
		return t.writeALSparse(g, w)
	case Dense:
		return t.writeALDense(g, w)
	case Arcs:
		return t.writeALArcs(g, w)
	}
	return 0, fmt.Errorf("format %d invalid", t.Format)
}

func (t Text) writeALDense(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	if t.WriteArcs != All {
		return t.writeALDenseTriangle(g, w)
	}
	t.fixBase()
	b := bufio.NewWriter(w)
	var c int
	for _, to := range g {
		if len(to) > 0 {
			c, err = b.WriteString(strconv.FormatInt(int64(to[0]), t.Base))
			n += c
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to), t.Base))
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

func (t Text) writeALSparse(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	if t.FrDelim == "" {
		t.FrDelim = ": "
	}
	if t.ToDelim == "" {
		t.ToDelim = " "
	}
	writeLast := t.NodeName == nil
	if writeLast {
		t.fixBase()
		t.NodeName = func(n graph.NI) string {
			return strconv.FormatInt(int64(n), t.Base)
		}
	}
	if t.WriteArcs != All {
		return t.writeALSparseTriangle(g, w, writeLast)
	}
	b := bufio.NewWriter(w)
	var c int
	last := len(g) - 1
	for i, to := range g {
		switch {
		case len(to) > 0:
			fr := graph.NI(i)
			c, err = b.WriteString(t.NodeName(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.FrDelim)
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.NodeName(to[0]))
			n += c
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				c, err = b.WriteString(t.ToDelim)
				n += c
				if err != nil {
					return
				}
				c, err = b.WriteString(t.NodeName(to))
				n += c
				if err != nil {
					return
				}
			}
		case i == last && writeLast:
			fr := graph.NI(i)
			c, err = b.WriteString(t.NodeName(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(
				strings.TrimRightFunc(t.FrDelim, unicode.IsSpace))
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

func (t Text) writeALArcs(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	t.fixBase()
	if t.FrDelim == "" {
		t.FrDelim = " "
	}
	if t.NodeName == nil {
		t.NodeName = func(n graph.NI) string {
			return strconv.FormatInt(int64(n), t.Base)
		}
	}
	b := bufio.NewWriter(w)
	var c int
	for fr, to := range g {
		for _, to := range to {
			c, err = b.WriteString(t.NodeName(graph.NI(fr)))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.FrDelim)
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.NodeName(to))
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

func (t Text) writeALDenseTriangle(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	t.fixBase()
	b := bufio.NewWriter(w)
	var c int
	p := func(fr, to graph.NI) bool { return to >= fr }
	if t.WriteArcs == Lower {
		p = func(fr, to graph.NI) bool { return to <= fr }
	}
	for i, to := range g {
		fr := graph.NI(i)
		one := false
		for _, to := range to {
			if p(to, fr) {
				if one {
					if err = b.WriteByte(' '); err != nil {
						return
					}
					n++
				} else {
					one = true
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to), t.Base))
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

func (t Text) writeALSparseTriangle(g graph.AdjacencyList, w io.Writer, writeLast bool) (
	n int, err error) {
	b := bufio.NewWriter(w)
	var c int
	p := func(fr, to graph.NI) bool { return to >= fr }
	if t.WriteArcs == Lower {
		p = func(fr, to graph.NI) bool { return to <= fr }
	}
	last := len(g) - 1
	for i, to := range g {
		fr := graph.NI(i)
		one := false
		for _, to := range to {
			if p(fr, to) {
				if !one {
					one = true
					c, err = b.WriteString(t.NodeName(fr))
					n += c
					if err != nil {
						return
					}
					c, err = b.WriteString(t.FrDelim)
					n += c
					if err != nil {
						return
					}
				} else {
					c, err = b.WriteString(t.ToDelim)
					n += c
					if err != nil {
						return
					}
				}
				c, err = b.WriteString(t.NodeName(to))
				n += c
				if err != nil {
					return
				}
			}
		}
		if writeLast && i == last && !one {
			one = true
			c, err = b.WriteString(t.NodeName(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(
				strings.TrimRightFunc(t.FrDelim, unicode.IsSpace))
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
