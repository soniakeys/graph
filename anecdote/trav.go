// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func DFSmall() (string, string) {
	chungLuSmall.DepthFirst(chungLuSmallCCRep, &bits.Bits{}, nil)
	return "DepthFirst", chungLuSmallCCTag
}

func DFLarge() (string, string) {
	chungLuLarge.DepthFirst(chungLuLargeCCRep, &bits.Bits{}, nil)
	return "DepthFirst", chungLuLargeCCTag
}

func BFSmall() (string, string) {
	chungLuSmall.BreadthFirst(chungLuSmallCCRep, nil, nil,
		func(graph.NI) bool { return true })
	return "BreadthFirst", chungLuSmallCCTag
}

func BFLarge() (string, string) {
	chungLuLarge.BreadthFirst(chungLuLargeCCRep, nil, nil,
		func(graph.NI) bool { return true })
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
