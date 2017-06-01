// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"github.com/soniakeys/graph"
)

func DFSmall() (string, string) {
	chungLuSmall.DepthFirst(chungLuSmallCCRep, &graph.Bits{}, nil)
	return "DepthFirst", chungLuSmallCCTag
}

func DFLarge() (string, string) {
	chungLuLarge.DepthFirst(chungLuLargeCCRep, &graph.Bits{}, nil)
	return "DepthFirst", chungLuLargeCCTag
}

func BFSmall() (string, string) {
	chungLuSmall.BreadthFirst(chungLuSmallCCRep, nil, nil,
		func(graph.NI) bool { return true } )
	return "BreadthFirst", chungLuSmallCCTag
}

func BFLarge() (string, string) {
	chungLuLarge.BreadthFirst(chungLuLargeCCRep, nil, nil,
		func(graph.NI) bool { return true } )
	return "BreadthFirst", chungLuLargeCCTag
}
