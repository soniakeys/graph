// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"github.com/soniakeys/graph"
)

var chungLuSmallCCRep graph.NI
var chungLuSmallCCma int
var chungLuSmallCCTag string

func CCSmall() (string, string) {
	reps, orders, ma := chungLuSmall.ConnectedComponentReps()
	max := 0
	for i, o := range orders {
		if o > max {
			max = o
			chungLuSmallCCRep = reps[i]
			chungLuSmallCCma = ma[i]
		}
	}
	chungLuSmallCCTag = "ChungLu giant component " + h(max) + " nds"
	return "Connected Components", chungLuSmallTag
}

var chungLuLargeCCRep graph.NI
var chungLuLargeCCma int
var chungLuLargeCCTag string

func CCLarge() (string, string) {
	reps, orders, ma := chungLuLarge.ConnectedComponentReps()
	max := 0
	for i, o := range orders {
		if o > max {
			max = o
			chungLuLargeCCRep = reps[i]
			chungLuLargeCCma = ma[i]
		}
	}
	chungLuLargeCCTag = "ChungLu giant component " + h(max) + " nds"
	return "Connected Components", chungLuLargeTag
}

var eucSmallSCCRep graph.NI
var eucSmallSCCTag string

func SCCPathSmall() (string, string) {
	max := 0
	eucSmall.SCCPathBased(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucSmallSCCRep = c[0]
		}
		return true
	})
	eucSmallSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC path based", eucSmallTag
}

var eucLargeSCCRep graph.NI
var eucLargeSCCTag string

func SCCPathLarge() (string, string) {
	max := 0
	eucLarge.SCCPathBased(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucLargeSCCRep = c[0]
		}
		return true
	})
	eucLargeSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC path based", eucLargeTag
}

func SCCPearceSmall() (string, string) {
	max := 0
	eucSmall.SCCPearce(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucSmallSCCRep = c[0]
		}
		return true
	})
	eucSmallSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC Pearce", eucSmallTag
}

func SCCPearceLarge() (string, string) {
	max := 0
	eucLarge.SCCPearce(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucLargeSCCRep = c[0]
		}
		return true
	})
	eucLargeSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC Pearce", eucLargeTag
}

func SCCTarjanSmall() (string, string) {
	max := 0
	eucSmall.SCCTarjan(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucSmallSCCRep = c[0]
		}
		return true
	})
	eucSmallSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC Tarjan", eucSmallTag
}

func SCCTarjanLarge() (string, string) {
	max := 0
	eucLarge.SCCTarjan(func(c []graph.NI) bool {
		if len(c) > max {
			max = len(c)
			eucLargeSCCRep = c[0]
		}
		return true
	})
	eucLargeSCCTag = "Euclidean giant component " + h(max) + " nds"
	return "SCC Tarjan", eucLargeTag
}
