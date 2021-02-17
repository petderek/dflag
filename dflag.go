package dflag

import (
	"errors"
	"flag"
	"io"
	"os"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrBadStruct     = errors.New("bad struct")
	ErrParsing       = errors.New("error parsing flags")
	ErrTypeAssertion = errors.New("error asserting types")
)

// ContinueOnError works the same as it does in the standard library
func ContinueOnError() func(*parser) {
	return func(e *parser) {
		e.ErrorHandling = flag.ContinueOnError
	}
}

func Output(w io.Writer) func(*parser) {
	return func(e *parser) {
		e.Output = w
	}
}

type option func(*parser)

// Parse takes a pointer to a struct and modifies it to include annotated
// flag values
func Parse(i interface{}, opts ...option) error {
	enc := &parser{
		ErrorHandling: flag.ExitOnError,
		Output:        os.Stderr,
	}
	for _, opt := range opts {
		opt(enc)
	}
	return enc.Parse(i, os.Args[1:])
}
