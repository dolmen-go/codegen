// Package codegen provides utilities for writing short Go code generators from text templates
package codegen

import (
	"bytes"
	"go/format"
	"os"
	"sync"
	"text/template"
)

// Create runs the "text/template".Template with data, pass it through gofmt
// and saves it to filePath
func Create(filePath string, t *template.Template, data interface{}) (err error) {
	return (&CodeTemplate{}).Create(filePath, data)
}

type CodeTemplate struct {
	Template *template.Template // See "text/template"
	Buffer   bytes.Buffer       // Used for sharing allocated memory between multiple Create calls
	mu       sync.Mutex
}

// Parse creates a CodeTemplate from a "text/template" source
func Parse(codeTemplate string) (*CodeTemplate, error) {
	t, err := template.New("").Parse(codeTemplate)
	if err != nil {
		return nil, err
	}
	var tmpl CodeTemplate
	tmpl.Template = t
	return &tmpl, nil
}

// MustParse wraps Parse throwing errors as exception
func MustParse(codeTemplate string) *CodeTemplate {
	tmpl, err := Parse(codeTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// Create runs the template with data, pass it through gofmt
// and saves it to filePath
func (t *CodeTemplate) Create(filePath string, data interface{}) (err error) {
	// This anonymous function exists just to wrap the mutex protected block
	out, err := func() ([]byte, error) {
		// To protect t.Buffer
		t.mu.Lock()
		defer t.mu.Unlock()

		t.Buffer.Reset()

		if err := t.Template.Execute(&t.Buffer, data); err != nil {
			return nil, err
		}

		return format.Source(t.Buffer.Bytes())
	}()
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
