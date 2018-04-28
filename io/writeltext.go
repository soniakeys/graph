// Copyright 2018 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Package graph/io provides graph readers and writers.
package io

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/soniakeys/graph"
)

// WriteLabeledAdjacencyList writes a lableed adjacency list as text.
//
// Fields of the receiver Text define how the text data is formatted.
// See documentation of the Text struct.
//
// Returned is number of bytes written and error.
func (t Text) WriteLabeledAdjacencyList(g graph.LabeledAdjacencyList,
	w io.Writer) (n int, err error) {
	switch t.Format {
	case Sparse:
		return t.writeLALSparse(g, w)
	case Dense:
		return t.writeLALDense(g, w)
	case Arcs:
		return t.writeLALArcs(g, w)
	}
	return 0, fmt.Errorf("format %d invalid", t.Format)
}

func (t Text) writeLALSparse(g graph.LabeledAdjacencyList,
	w io.Writer) (n int, err error) {
	return
}
func (t Text) writeLALArcs(g graph.LabeledAdjacencyList,
	w io.Writer) (n int, err error) {
	return
}
func (t Text) writeLALDense(g graph.LabeledAdjacencyList,
	w io.Writer) (n int, err error) {
	t.fixBase()
	b := bufio.NewWriter(w)
	var c int
	for _, to := range g {
		if len(to) > 0 {
			c, err = b.WriteString(strconv.FormatInt(int64(to[0].To), t.Base))
			n += c
			if err != nil {
				return
			}
			if err = b.WriteByte(' '); err != nil {
				return
			}
			c, err = b.WriteString(strconv.FormatInt(int64(to[0].Label), t.Base))
			n += c + 1
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to.To), t.Base))
				n += c + 1
				if err != nil {
					return
				}
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to.Label), t.Base))
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
