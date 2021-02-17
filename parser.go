package dflag

import (
	"flag"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type parser struct {
	ErrorHandling flag.ErrorHandling
	UsageText     string
	Output        io.Writer
}

const (
	keyName  = "name"
	keyValue = "value"
	keyUsage = "usage"
)

func (e *parser) Parse(i interface{}, args []string) error {
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

	fs := flag.NewFlagSet("", e.ErrorHandling)
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
			if e != nil {
				return ErrBadStruct
			}
			pointers[i] = fs.Int(name, num, usage)
		case reflect.String:
			pointers[i] = fs.String(name, value, usage)
		case reflect.Bool:
			boo, err := strconv.ParseBool(value)
			if err != nil {
				return ErrBadStruct
			}
			pointers[i] = fs.Bool(name, boo, usage)
		}
	}

	if err := fs.Parse(args); err != nil {
		return err
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
