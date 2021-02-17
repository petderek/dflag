package dflag

import (
	"errors"
	"flag"
	"io"
	"reflect"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type encoder struct {
	ErrorHandling flag.ErrorHandling
	UsageText     string
	Output        io.Writer
}

func (e *encoder) Parse(i interface{}, args []string) error {
	t := reflect.TypeOf(i)
	if t == nil {
		return ErrInvalidInput
	}
	return nil
}

// ContinueOnError works the same as it does in the standard library
func ContinueOnError() func(*encoder) {
	return func(e *encoder) {
		e.ErrorHandling = flag.ContinueOnError
	}
}

// UsageText is the text to print before displaying the flagset options.
func UsageText(text string) func(*encoder) {
	return func(e *encoder) {
		e.UsageText = text
	}
}

func Output(w io.Writer) func(*encoder) {
	return func(e *encoder) {
		e.Output = w
	}
}

type option func(*encoder)

// Parse takes a pointer to a struct and modifies it
func Parse(i interface{}, opts ...option) {
	t := reflect.TypeOf(i)
	if t == nil {
		//todo: error
		return
	}
}
