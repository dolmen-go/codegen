/*
Copyright 2024 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package codegen provides utilities for writing short Go code generators from a text/template.
package codegen

import (
	"bytes"
	"errors"
	"go/format"
	"os"
	"regexp"
	"sync"
	"text/template"
)

// GeneratedCodeRegexp checks the code generation standard
// defined at https://golang.org/s/generatedcode.
var GeneratedCodeRegexp = regexp.MustCompile(`(?m:^// Code generated .* DO NOT EDIT\.$)`)

// CreateFile runs the "text/template".Template with data, pass it through gofmt
// and saves it to filePath.
func CreateFile(filePath string, t *template.Template, data interface{}) (err error) {
	return (&CodeTemplate{Template: t}).CreateFile(filePath, data)
}

// CodeTemplate is the precompiled template for generating one or multiple Go source files.
type CodeTemplate struct {
	Template *template.Template // See "text/template"
	Buffer   bytes.Buffer       // Used for sharing allocated memory between multiple CreateFile calls
	mu       sync.Mutex
}

// Parse creates a CodeTemplate from a "text/template" source.
//
// The expansion of the template is expected to be valid a Go source file
// containing the code generation standard tag. See GeneratedCodeRegexp.
func Parse(codeTemplate string) (*CodeTemplate, error) {
	t, err := template.New("").Parse(codeTemplate)
	if err != nil {
		return nil, err
	}
	return &CodeTemplate{Template: t}, nil
}

// MustParse wraps Parse throwing errors as exception.
func MustParse(codeTemplate string) *CodeTemplate {
	tmpl, err := Parse(codeTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// CreateFile runs the template with data, pass it through gofmt
// and saves it to filePath.
//
// The code generation standard at https://golang.org/s/generatedcode is enforced.
func (t *CodeTemplate) CreateFile(filePath string, data interface{}) error {
	// This anonymous function exists just to wrap the mutex protected block
	out, err := func() ([]byte, error) {
		// To protect t.Buffer
		t.mu.Lock()
		defer t.mu.Unlock()

		t.Buffer.Reset()

		if err := t.Template.Execute(&t.Buffer, data); err != nil {
			return nil, err
		}

		code := t.Buffer.Bytes()

		// Enforce code generation standard https://golang.org/s/generatedcode
		if !GeneratedCodeRegexp.Match(code) {
			return nil, errors.New("output does not follow standard defined at https://golang.org/s/generatedcode")
		}

		return format.Source(code)
	}()
	if err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err == nil {
		defer f.Close()
		_, err = f.Write(out)
	}
	return err
}
