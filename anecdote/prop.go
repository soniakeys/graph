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
