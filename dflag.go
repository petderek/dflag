package dflag

import (
	"errors"
	"flag"
	"io"
	"os"
)

var (
	ErrParsing         = errors.New("error parsing values")
	ErrMissingArgument = errors.New("missing a required argument")
	errInvalidInput    = errors.New("invalid input")
	errBadStruct       = errors.New("bad struct tags")
	errTypeAssertion   = errors.New("error asserting types")
)

// ContinueOnError works the same as it does in the standard library.
// Normally, calls to Parse() will os.Exit() on error. Passing this option
// in allows clients to handle the error themselves.
func ContinueOnError() func(*parser) {
	return func(e *parser) {
		e.ErrorHandling = flag.ContinueOnError
	}
}

// Output allows clients to change the stream that Usage() writes to. By
// default, it is stderr.
func Output(w io.Writer) func(*parser) {
	return func(e *parser) {
		e.Output = w
	}
}

// UsageText provides the text that will be printed before printing the
// command defaults. By default, it will just be "usage: ".
func UsageText(text string) func(*parser) {
	return func(p *parser) {
		p.UsageText = text
	}
}

type option func(*parser)

var state = &parser{
	ErrorHandling: flag.ExitOnError,
	Output:        os.Stderr,
	Args:          os.Args,
	logger:        &noopLogger{},
}

// Parse parses the command-line flags from os.Args[1:]. Must be called after
// all flags are defined and before flags are accessed by the program.
// i must be a pointer to a struct.
func Parse(i interface{}, opts ...option) error {
	for _, opt := range opts {
		opt(state)
	}
	return state.Parse(i)
}

// Args returns the non-flag command-line arguments.
func Args() []string {
	return state.flagset.Args()
}

// Arg returns the i'th command-line argument. Arg(0) is the first remaining
// argument after flags have been processed. Arg returns an empty string if
// the requested element does not exist.
func Arg(i int) string {
	return state.flagset.Arg(i)
}
