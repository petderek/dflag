package dflag

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type parser struct {
	ErrorHandling flag.ErrorHandling
	UsageText     string
	Output        io.Writer
	Args          []string

	flagset flag.FlagSet
}

const (
	keyName  = "name"
	keyValue = "value"
	keyUsage = "usage"
)

func (p *parser) Parse(i interface{}) error {
	p.flagset = *flag.NewFlagSet(p.Args[0], p.ErrorHandling)
	p.flagset.SetOutput(p.Output)
	p.flagset.Usage = p.getUsage()

	err := p.parse(i)
	if err == nil {
		return nil
	}

	// These errors indicate that the struct is malformed and should be fixed
	// by the developer.
	if errors.Is(err, errInvalidInput) || errors.Is(err, errBadStruct) || errors.Is(err, errTypeAssertion) {
		panic(err)
	}

	if p.ErrorHandling == flag.ExitOnError {
		// ErrParsing indicates a client issue. Printing usage should
		// provide hints at how to fix
		p.flagset.Usage()
		os.Exit(2)
	}

	return err
}

func (p *parser) parse(i interface{}) error {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	switch {
	case t == nil:
		return errInvalidInput
	case t.Kind() != reflect.Ptr || v.Kind() != reflect.Ptr:
		return errInvalidInput
	case t.Elem().Kind() != reflect.Struct || v.Elem().Kind() != reflect.Struct:
		return errInvalidInput
	}

	pointers := make([]interface{}, t.Elem().NumField())
	required := make([]bool, t.Elem().NumField())
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		if !v.Elem().Field(i).CanSet() {
			continue
		}
		name := field.Tag.Get(keyName)
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		value := field.Tag.Get(keyValue)
		usage := field.Tag.Get(keyUsage)

		tkns := getTokens(field.Tag)
		required[i] = tkns.hasValue("required")

		switch field.Type.Kind() {
		case reflect.Int:
			num, e := strconv.Atoi(value)
			if value != "" && e != nil {
				return errBadStruct
			}
			pointers[i] = p.flagset.Int(name, num, usage)
		case reflect.String:
			pointers[i] = p.flagset.String(name, value, usage)
		case reflect.Bool:
			boo, err := strconv.ParseBool(value)
			if value != "" && err != nil {
				return errBadStruct
			}
			pointers[i] = p.flagset.Bool(name, boo, usage)
		}
	}

	if err := p.flagset.Parse(p.Args[1:]); err != nil {
		return ErrParsing
	}

	for i, ptr := range pointers {
		val := v.Elem().Field(i)
		if !val.CanSet() || ptr == nil {
			if required[i] {
				return ErrMissingArgument
			}
			continue
		}
		switch val.Kind() {
		case reflect.Int:
			if i, ok := ptr.(*int); ok {
				val.SetInt(int64(*i))
			} else {
				return errTypeAssertion
			}
		case reflect.String:
			if s, ok := ptr.(*string); ok {
				val.SetString(*s)
			} else {
				return errTypeAssertion
			}
		case reflect.Bool:
			if b, ok := ptr.(*bool); ok {
				val.SetBool(*b)
			} else {
				return errTypeAssertion
			}
		}
	}
	return nil
}

func (p *parser) getUsage() func() {
	var once sync.Once
	usage := "usage: "
	if p.UsageText != "" {
		usage = p.UsageText
	}
	return func() {
		once.Do(func() {
			_, _ = fmt.Fprint(p.Output, usage, "\n")
			p.flagset.PrintDefaults()
		})
	}
}

type tokens struct {
	data []string
}

// dflag tokens are of the form:
// `dflag:"token1,token2"`
func getTokens(tag reflect.StructTag) tokens {
	t, ok := tag.Lookup("dflag")
	if !ok {
		return tokens{}
	}
	data := strings.Split(t, ",")
	return tokens{data: data}
}

func (t tokens) hasValue(value string) bool {
	for _, v := range t.data {
		if v == value {
			return true
		}
	}
	return false
}
