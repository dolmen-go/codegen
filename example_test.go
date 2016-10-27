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
	os.Stdout.WriteString("Hello, {{.}}!\n")
}
`

	tmpl := codegen.MustParse(template)
	for _, tag := range os.Args[1:] {
		f := "main_" + tag + ".go"
		if err := tmpl.CreateFile(f, tag); err != nil {
			log.Fatal(err)
		}
		log.Printf("File %s created.\n", f)
	}
}
