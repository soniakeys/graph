// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package gen_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/soniakeys/graph/gen"
)

func Example() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := `
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}`
	// data for generating source
	data := &gen.Data{
		Pkg: "main",
		Funcs: []gen.FuncData{{
			Func:        "dfZ",
			NodeVisitor: true,
		}},
	}
	call := `
	dfZ(g, 0, func(n graph.NI) {
		fmt.Println("visit", n)
	})`
	run(g, data, call)
	// Output:
	// visit 0
	// visit 1
	// visit 2
	// visit 3
}

func run(g string, data *gen.Data, call string) {
	// a temp work dir
	dir, err := ioutil.TempDir("", "gen")
	if err != nil {
		log.Fatal(err)
	}
	// a file for the generated source
	f, err := os.Create(path.Join(dir, "dfz.go"))
	if err != nil {
		log.Fatal(err)
	}
	// generate the source
	if err := gen.DepthFirst(f, data); err != nil {
		log.Fatal(err)
	}
	// and a test program
	ioutil.WriteFile(path.Join(dir, "main.go"), []byte(`
package main

import (
	"fmt"

	"github.com/soniakeys/graph"
)

func main() {`+g+call+`}`), 0444)
	// run the test program
	cmd := exec.Command("go", "run",
		path.Join(dir, "dfz.go"),
		path.Join(dir, "main.go"))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	runErr, _ := ioutil.ReadAll(stderr)
	runOut, _ := ioutil.ReadAll(stdout)
	if len(runErr) > 0 {
		log.Println(string(runErr))
	}
	fmt.Println(string(runOut))
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	os.RemoveAll(dir)
}
