#Graph

A graph library with goals of speed and simplicity, Graph implements
graph algorithms on graphs of zero-based integer node IDs.

[![GoDoc](https://godoc.org/github.com/soniakeys/graph?status.svg)](https://godoc.org/github.com/soniakeys/graph) [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/soniakeys/graph) [![GoSearch](http://go-search.org/badge?id=github.com%2Fsoniakeys%2Fgraph)](http://go-search.org/view?id=github.com%2Fsoniakeys%2Fgraph)[![Build Status](https://travis-ci.org/soniakeys/graph.svg?branch=master)](https://travis-ci.org/soniakeys/graph)

Status, 7 Jun 2017:  Go 1.9 is bringing features I want.  Until it comes out,
building latest commit requires Go tip.

###Non-source files of interest

The directory [tutorials](tutorials) is a work in progress - there are only
a few tutorials there yet - but the concept is to provide some topical
walk-throughs to supplement godoc.  The source-based godoc documentation
remains the primary documentation.

* [Dijkstra's algorithm](tutorials/dijkstra.md)
* [AdjacencyList types](tutorials/adjacencylist.md)
* [Missing methods](tutorials/missingmethods.md)

The directory [bench](bench) is another work in progress.  The concept is
to present some plots showing benchmark performance approaching some
theoretical asymptote.

[hacking.md](hacking.md) has some information about how the library is
developed, built, and tested.  It might be of interest if for example you
plan to fork or contribute to the the repository.

###Test coverage
5 Nov 2016
```
graph          95.5%
graph/df       20.7%
graph/dot      77.5%
graph/treevis  79.4%
```
