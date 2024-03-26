package codegen_test

import (
	"log"
	"os"
	"testing"

	"github.com/dolmen-go/codegen"
)

func TestExample(t *testing.T) {
	os.Args = make([]string, 2)
	os.Args[1] = "foo"
	const filename = "main_foo.go"

	defer func() {
		f, err := os.Open(filename)
		if err != nil {
			t.Errorf("file " + filename + " not created!")
			t.Fail()
		}
		_ = f.Close()
		if err = os.Remove("main_foo.go"); err == nil {
			log.Printf("File %s removed.\n", filename)
		}
	}()

	Example()
}

// TestParseFailures to reach full coverage.
func TestParseFailures(t *testing.T) {
	_, err := codegen.Parse("{{")
	if err == nil {
		t.Fatal("Parse should fail.")
	}

	err = func() (r error) {
		defer func() {
			r = recover().(error)
		}()
		codegen.MustParse("{{")
		return
	}()
	if err == nil {
		t.Fatal("MustParse should panic.")
	}
}
