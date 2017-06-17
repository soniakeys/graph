// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
	"github.com/soniakeys/graph/traverse"
)

func DFSmall() (string, string) {
	b := bits.New(chungLuSmall.Order())
	err := traverse.DepthFirst(chungLuSmall, chungLuSmallCCRep,
		traverse.Visited(&b))
	if err != nil {
		return "DepthFirst", err.Error()
	}
	return "DepthFirst", chungLuSmallCCTag
}

func DFLarge() (string, string) {
	b := bits.New(chungLuLarge.Order())
	err := traverse.DepthFirst(chungLuLarge, chungLuLargeCCRep,
		traverse.Visited(&b))
	if err != nil {
		return "DepthFirst", err.Error()
	}
	return "DepthFirst", chungLuLargeCCTag
}

func BFSmall() (string, string) {
	traverse.BreadthFirst(chungLuSmall, chungLuSmallCCRep)
	return "BreadthFirst", chungLuSmallCCTag
}

func BFLarge() (string, string) {
	traverse.BreadthFirst(chungLuLarge, chungLuLargeCCRep)
	return "BreadthFirst", chungLuLargeCCTag
}

func AltBFSmall() (string, string) {
	g := chungLuSmall.AdjacencyList
	alt.BreadthFirst(g, g, chungLuSmallCCma, chungLuSmallCCRep, nil,
		func(graph.NI) bool { return true })
	return "DO BreadthFirst", chungLuSmallCCTag
}

func AltBFLarge() (string, string) {
	g := chungLuLarge.AdjacencyList
	alt.BreadthFirst(g, g, chungLuLargeCCma, chungLuLargeCCRep, nil,
		func(graph.NI) bool { return true })
	return "DO BreadthFirst", chungLuLargeCCTag
}
