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
	if err != nil && p.ErrorHandling == flag.ExitOnError {
		// ErrParsing indicates a client issue. Printing usage should
		// provide hints at how to fix
		if errors.Is(err, ErrParsing) {
			p.flagset.Usage()
		} else {
			// Other errors here indicate issues with the struct that
			// should be fixed by the developer.
			_, _ = fmt.Fprintf(p.Output, "Error in flag structure: %s", err)
		}
		os.Exit(2)
	}

	return err
}

func (p *parser) parse(i interface{}) error {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	switch {
	case t == nil:
		return ErrInvalidInput
	case t.Kind() != reflect.Ptr || v.Kind() != reflect.Ptr:
		return ErrInvalidInput
	case t.Elem().Kind() != reflect.Struct || v.Elem().Kind() != reflect.Struct:
		return ErrInvalidInput
	}

	pointers := make([]interface{}, t.Elem().NumField())
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

		switch field.Type.Kind() {
		case reflect.Int:
			num, e := strconv.Atoi(value)
			if value != "" && e != nil {
				return ErrBadStruct
			}
			pointers[i] = p.flagset.Int(name, num, usage)
		case reflect.String:
			pointers[i] = p.flagset.String(name, value, usage)
		case reflect.Bool:
			boo, err := strconv.ParseBool(value)
			if value != "" && err != nil {
				return ErrBadStruct
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
			continue
		}
		switch val.Kind() {
		case reflect.Int:
			if i, ok := ptr.(*int); ok {
				val.SetInt(int64(*i))
			} else {
				return ErrTypeAssertion
			}
		case reflect.String:
			if s, ok := ptr.(*string); ok {
				val.SetString(*s)
			} else {
				return ErrTypeAssertion
			}
		case reflect.Bool:
			if b, ok := ptr.(*bool); ok {
				val.SetBool(*b)
			} else {
				return ErrTypeAssertion
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
