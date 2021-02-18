package dflag

import (
	"flag"
	"reflect"
	"testing"
)

type happyCase struct {
	Happy string
	sad   string
}

func TestHappyCase(t *testing.T) {
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args:          args("-happy", "yes"),
	}
	var h happyCase
	err := e.Parse(&h)
	if err != nil {
		t.Error("Error calling parse: ", err)
	}
	if h.Happy != "yes" {
		t.Errorf("Expected value to be \"yes\", was actually: \"%s\"", h.Happy)
	}
}

func TestInvalidValues(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panicd at the disco")
		}
	}()
	var (
		str     = "string"
		num     = 12
		bool    = true
		iface   interface{}
		nilface interface{} = nil
	)
	cases := []interface{}{
		str,
		num,
		bool,
		iface,
		nilface,
		&str,
		&num,
		&bool,
		&iface,
		&nilface,
		(*error)(nil),
		struct{}{},
	}
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args:          args("-happy", "yes"),
	}
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			err := e.Parse(tc)
			if err != ErrInvalidInput {
				t.Error("Expected invalid input, got: ", err)
			}
		})
	}
}

type testAnnotated struct {
	Annotated string `name:"foo" value:"bar" usage:"baz"`
}

func TestAnnotated(t *testing.T) {
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args:          args("na"),
	}
	var annotated testAnnotated
	err := e.Parse(&annotated)
	if err != nil {
		t.Error("Error parsing: ", err)
	}
	if annotated.Annotated != "bar" {
		t.Errorf("Expected \"bar\", but actually \"%s\"", annotated.Annotated)
	}
}

type testStrings struct {
	Zero  string `name:"" value:"" usage:""`
	One   string `name:"" value:"" usage:"one"`
	Two   string `name:"" value:"one" usage:""`
	Three string `name:"" value:"one" usage:"one"`
	Four  string `name:"onefour" value:"" usage:""`
	Five  string `name:"onefive" value:"" usage:"1"`
	Six   string `name:"onesix" value:"one" usage:""`
	Seven string `name:"oneseven" value:"one" usage:"one"`
}

func TestFlagParsed(t *testing.T) {
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args: args(
			"-zero", "0",
			"-one", "1",
			"-two", "2",
			"-three", "3",
			"-onefour", "4",
			"-onefive", "5",
			"-onesix", "6",
			"-oneseven", "7",
		),
	}
	var data testStrings
	err := e.Parse(&data)
	if err != nil {
		t.Error("error: ", err)
	}

	if !reflect.DeepEqual(data, testStrings{
		Zero:  "0",
		One:   "1",
		Two:   "2",
		Three: "3",
		Four:  "4",
		Five:  "5",
		Six:   "6",
		Seven: "7",
	}) {
		t.Error("not deeply equal: ", data)
	}
}

func args(a ...string) []string {
	return append([]string{"cmd"}, a...)
}
