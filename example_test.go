package codegen_test

import (
	"github.com/dolmen-go/codegen"
	"log"
	"os"
)

func ExampleMustParse() {
	const template = `
{{/**/}}//+build {{.}}

package main

import "os"

func main() {
	os.StdOut.WriteString("Hello, {{.}}!\n")
}
`

	tmpl := codegen.MustParse(template)
	for _, tag := range os.Args[1:] {
		if err := tmpl.CreateFile("main_"+tag+".go", tag); err != nil {
			log.Fatal(err)
		}
	}
}
