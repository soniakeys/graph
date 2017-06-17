// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// Alternative search implementation, using code generation.
package gen

import (
	"io"
	"text/template"
)

type Data struct {
	Pkg   string // required
	Funcs []FuncData
}

type FuncData struct {
	Func          string // required
	NodeVisitor   bool
	OkNodeVisitor bool
}

func DepthFirst(src io.Writer, data *Data) error {
	tmpl := template.Must(template.New("parms").Parse(tsParms))
	tmpl = template.Must(tmpl.New("funcDef").Parse(tsFuncDef))
	tmpl = template.Must(tmpl.New("nodeVis").Parse(tsNodeVis))
	tmpl = template.Must(tmpl.New("recurse").Parse(tsRecurse))
	tmpl = template.Must(tmpl.New("funcEnd").Parse(tsFuncEnd))
	tmpl = template.Must(tmpl.New("").Parse(ts))
	return tmpl.Execute(src, data)
}

var (
	ts = `
package {{.Pkg}}

import "github.com/soniakeys/bits"
import "github.com/soniakeys/graph"

{{range .Funcs}}
func {{.Func}}(g graph.AdjacencyList, start graph.NI {{template "parms" .}}) {
	{{template "funcDef" .}}
		{{template "nodeVis" .}}
		for _, to := range g[n] {
			{{template "recurse" .}}
		}
		{{template "funcEnd" .}}
}
{{end}}`

	tsParms = `
{{- if .NodeVisitor}},
	nv func(graph.NI){{end}}
{{- if .OkNodeVisitor}},
	oknv func(graph.NI) bool{{end -}}
`
	tsFuncDef = `
	b := bits.New(len(g))
	var df func(graph.NI){{if .OkNodeVisitor}} bool{{end}}
	df = func(n graph.NI){{if .OkNodeVisitor}} bool{{end}} {
		if b.Bit(int(n)) != 0 {
			return{{if .OkNodeVisitor}} true{{end}}
		}
		b.SetBit(int(n), 1)
`
	tsNodeVis = `
{{- if .NodeVisitor}}
		nv(n){{end}}
{{- if .OkNodeVisitor}}
		if !oknv(n) {
			return false
		}{{end}}
`
	tsRecurse = `
			{{if .OkNodeVisitor}}if !{{end}}df(to)
{{- if .OkNodeVisitor}} {
				return false
			}{{end}}
`
	tsFuncEnd = `
{{- if .OkNodeVisitor}}
		return true{{end}}
	}
	df(start)
`
)
