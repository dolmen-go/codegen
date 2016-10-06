package codegen

/*
import (
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
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
		_ = os.Remove("main_foo.go")
	}()

	// How could
	ExampleMustParse()
}
*/
