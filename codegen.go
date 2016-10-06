// Package codegen provides utilities for writing short Go code generators from text templates
package codegen

import (
	"bytes"
	"go/format"
	"os"
	"text/template"
)

func Create(filePath string, t *template.Template, data interface{}) (err error) {
	return (&CodeTemplate{}).Create(filePath, data)
}

type CodeTemplate struct {
	Template *template.Template
	Buffer   bytes.Buffer
}

func Parse(codeTemplate string) (*CodeTemplate, error) {
	t, err := template.New("").Parse(codeTemplate)
	if err != nil {
		return nil, err
	}
	var tmpl CodeTemplate
	tmpl.Template = t
	return &tmpl, nil
}

func MustParse(codeTemplate string) *CodeTemplate {
	tmpl, err := Parse(codeTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func (t *CodeTemplate) Create(filePath string, data interface{}) (err error) {

	t.Buffer.Reset()

	err = t.Template.Execute(&t.Buffer, data)
	if err != nil {
		return
	}

	out, err := format.Source(t.Buffer.Bytes())
	if err != nil {
		return
	}

	f, err := os.Create(filePath)
	if err == nil {
		defer f.Close()
		_, err = f.Write(out)
	}
	return
}
