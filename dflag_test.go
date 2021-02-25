package dflag

import (
	"flag"
	"reflect"
	"strings"
	"testing"
)

type happyCase struct {
	Happy string
	Hippy int
	Hoppy bool
	sad   string
}

func TestHappyCase(t *testing.T) {
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args:          args("-happy", "yes", "-hippy", "42", "-hoppy", "true"),
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

func TestDynamicDefaults(t *testing.T) {
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args:          args("-hippy", "42", "-hoppy", "true"),
	}
	var h happyCase
	h.Happy = "potato"
	err := e.Parse(&h)
	if err != nil {
		t.Error("Error calling parse: ", err)
	}
	if h.Happy != "potato" {
		t.Errorf("Expected value to be \"potato\", was actually: \"%s\"", h.Happy)
	}
}

type dupes struct {
	Dupe string
	DUPE string
}

func TestDupes(t *testing.T) {
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args:          args("-happy", "yes"),
	}
	var d dupes
	shouldPanic(t, func() {
		_ = e.Parse(&d)
	})
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
		b       = true
		iface   interface{}
		nilface interface{} = nil
	)
	cases := []interface{}{
		str,
		num,
		b,
		iface,
		nilface,
		&str,
		&num,
		&b,
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
			shouldPanic(t, func() {
				t.Log(e.Parse(tc))
			})
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

type requiredArgs struct {
	Zero int    `required:"true"`
	One  bool   `required:"T"` // strconv.ParseBool accepts T as true
	Two  string `required:"1"` // strconv.ParseBool accepts 1 as true
}

func TestRequiredArgsHappy(t *testing.T) {
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args: args(
			"-zero=0",
			"-one",
			"-two=ok",
		),
	}

	if err := e.Parse(&requiredArgs{}); err != nil {
		t.Error("Expected error to be nil, was: ", err)
	}
}

func TestRequiredArgsSad(t *testing.T) {
	builder := &strings.Builder{}
	e := &parser{
		ErrorHandling: flag.ContinueOnError,
		Args:          args(),
		Output:        builder,
	}

	if err := e.Parse(&requiredArgs{}); err == nil {
		t.Error("Expected error to be to be thrown, but was nil.")
	}
	t.Log(builder.String())
}

func args(a ...string) []string {
	return append([]string{"cmd"}, a...)
}

func shouldPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Log("Panicked! But it was expected.")
		} else {
			t.Error("expected a panic, but was none")
		}
	}()
	f()
}
