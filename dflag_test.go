package dflag

import (
	"flag"
	"testing"
)

type happyCase struct {
	happy string
}

func TestHappyCase(t *testing.T) {
	e := &encoder{
		ErrorHandling: flag.ContinueOnError,
	}
	var h happyCase
	err := e.Parse(&h, args("-happy", "yes"))
	if err != nil {
		t.Error("Error calling parse: ", err)
	}
	if h.happy != "yes" {
		t.Errorf("Expected value to be \"yes\", was actually: \"%s\"", h.happy)
	}
}

func TestInvalidValues(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panicd at the disco")
		}
	}()
	var (
		str  = "string"
		num  = 12
		bool = true
		iface interface{}
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
	e := &encoder{
		ErrorHandling: flag.ContinueOnError,
	}
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			err := e.Parse(tc, args("-happy", "yes"))
			if err != ErrInvalidInput {
				t.Error("Expected invalid input, got: ", err)
			}
		})
	}
}

type testAnnotated struct {
	annotated string `name:"foo" value:"bar" usage:"baz"`
}

func TestAnnotated(t *testing.T) {

}

type testStrings struct {
	one   string `name:"" value:"" usage:""`
	two   string `name:"" value:"" usage:""`
	three string `name:"" value:"" usage:""`
	four  string `name:"" value:"" usage:""`
	five  string `name:"" value:"" usage:""`
	six   string `name:"" value:"" usage:""`
}

func TestFlagParsed(t *testing.T) {
	t.Log("ok")
}

var (
	testcmd = []string{"testcmd"}
)

func args(a ...string) []string {
	return append(testcmd, a...)
}
