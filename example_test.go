package codegen_test

import (
	"github.com/dolmen-go/codegen"
	"os"
)

func ExampleMustParse() {
	const template = `
// +build {{tag}}
package main

import "os"

func main() {
	os.StdOut.WriteString("Hello, {{tag}}!\n")
}
`

	tmpl := codegen.MustParse(template)
	for _, tag := range os.Args[1:] {
		tmpl.Create("main_"+tag+".go", tag)
	}
}
