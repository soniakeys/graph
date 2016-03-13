// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package graph

import (
	"errors"
	"math"
	"math/rand"
)

// Euclidean generates a random simple graph on the Euclidean plane.
//
// Nodes are associated with coordinates uniformly distributed from (0,0)
// to (1,1).  Arcs are added between random nodes with a bias toward connecting
// nearer nodes.
//
// Unfortunately the function has a few "knobs".
// The returned graph will have order nNodes and arc size nArcs.  The affinity
// argument controls the bias toward connecting nearer nodes.  The function
// selects random pairs of nodes as a candidate arc then rejects the candidate
// if the nodes fail an affinity test.  Also parallel arcs are rejected.
// As more affine or denser graphs are requested, rejections increase,
// increasing run time.  The patience argument controls the number of arc
// rejections allowed before the function gives up and returns an error.
//
// Returned is a directed simple graph and associated positions indexed by
// node number.
//
// See also LabeledEuclidean.
func Euclidean(nNodes, nArcs int, affinity float64, patience int, r *rand.Rand) (g Directed, pos []struct{ X, Y float64 }, err error) {
	a := make(AdjacencyList, nNodes) // graph
	// generate random positions
	pos = make([]struct{ X, Y float64 }, nNodes)
	for i := range pos {
		pos[i].X = r.Float64()
		pos[i].Y = r.Float64()
	}
	// arcs
	var tooFar, dup int
arc:
	for i := 0; i < nArcs; {
		if tooFar == nArcs*patience {
			err = errors.New("affinity not found")
			return
		}
		if dup == nArcs*patience {
			err = errors.New("overcrowding")
			return
		}
		n1 := NI(r.Intn(nNodes))
		var n2 NI
		for {
			n2 = NI(r.Intn(nNodes))
			if n2 != n1 { // no graph loops
				break
			}
		}
		c1 := &pos[n1]
		c2 := &pos[n2]
		dist := math.Hypot(c2.X-c1.X, c2.Y-c1.Y)
		if dist*affinity > r.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		for _, nb := range a[n1] {
			if nb == n2 { // no parallel arcs
				dup++
				continue arc
			}
		}
		a[n1] = append(a[n1], n2)
		i++
	}
	g = Directed{a}
	return
}

// LabeledEuclidean generates a random simple graph on the Euclidean plane.
//
// Arc label values in the returned graph g are indexes into the return value
// wt.  Wt is the Euclidean distance between the from and to nodes of the arc.
//
// Otherwise the function arguments and return values are the same as for
// function Euclidean.  See Euclidean.
func LabeledEuclidean(nNodes, nArcs int, affinity float64, patience int, r *rand.Rand) (g LabeledDirected, pos []struct{ X, Y float64 }, wt []float64, err error) {
	a := make(LabeledAdjacencyList, nNodes) // graph
	wt = make([]float64, nArcs)             // arc weights
	// generate random positions
	pos = make([]struct{ X, Y float64 }, nNodes)
	for i := range pos {
		pos[i].X = r.Float64()
		pos[i].Y = r.Float64()
	}
	// arcs
	var tooFar, dup int
arc:
	for i := 0; i < nArcs; {
		if tooFar == nArcs*patience {
			err = errors.New("affinity not found")
			return
		}
		if dup == nArcs*patience {
			err = errors.New("overcrowding")
			return
		}
		n1 := NI(r.Intn(nNodes))
		var n2 NI
		for {
			n2 = NI(r.Intn(nNodes))
			if n2 != n1 { // no graph loops
				break
			}
		}
		c1 := &pos[n1]
		c2 := &pos[n2]
		dist := math.Hypot(c2.X-c1.X, c2.Y-c1.Y)
		if dist*affinity > r.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		for _, nb := range a[n1] {
			if nb.To == n2 { // no parallel arcs
				dup++
				continue arc
			}
		}
		wt[i] = dist
		a[n1] = append(a[n1], Half{n2, LI(i)})
		i++
	}
	g = LabeledDirected{a}
	return
}
