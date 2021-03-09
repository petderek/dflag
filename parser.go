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

	logger  logger
	flagset flag.FlagSet
}

const (
	keyName     = "name"
	keyValue    = "value"
	keyUsage    = "usage"
	keyRequired = "required"
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

type parseNode struct {
	TypeField  reflect.StructField
	ValueField reflect.Value
	Pointer    interface{}
	Flag       *flag.Flag
	SetViaFlag bool
}

func (p parseNode) Name() string {
	name := p.TypeField.Tag.Get(keyName)
	if name == "" {
		name = strings.ToLower(p.TypeField.Name)
	}
	return name
}

func (p parseNode) Value() string {
	return p.TypeField.Tag.Get(keyValue)
}

func (p parseNode) Usage() string {
	return p.TypeField.Tag.Get(keyUsage)
}

func (p parseNode) Required() bool {
	if r := p.TypeField.Tag.Get(keyRequired); r != "" {
		required, err := strconv.ParseBool(r)
		if err != nil {
			panic("fix your struct: " + err.Error())
		}
		return required
	}
	return false
}

func (p *parser) parse(i interface{}) error {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	switch {
	case t == nil:
		p.log("input cannot be nil")
		return errInvalidInput
	case t.Kind() != reflect.Ptr || v.Kind() != reflect.Ptr:
		p.log("expected a pointer, but input is: type(%T) value(%#v)")
		return errInvalidInput
	case t.Elem().Kind() != reflect.Struct || v.Elem().Kind() != reflect.Struct:
		p.log("expected a struct, but input was: %s %s", t.Elem().Kind(), v.Elem().Kind())
		return errInvalidInput
	}

	nodes := make(map[string]*parseNode, t.Elem().NumField())
	for i := 0; i < t.Elem().NumField(); i++ {
		node := &parseNode{
			TypeField:  t.Elem().Field(i),
			ValueField: v.Elem().Field(i),
		}

		if !node.ValueField.CanSet() {
			p.log("Node isn't settaable: ", node.ValueField)
			continue
		}

		switch node.TypeField.Type.Kind() {
		case reflect.Int:
			num, e := strconv.Atoi(node.Value())
			if node.Value() != "" && e != nil {
				p.logger.Printf("node (%s) has bad value (%s). change the struct tags: %s", node.Name(), node.Value(), e)
				return errBadStruct
			}
			node.Pointer = p.flagset.Int(node.Name(), num, node.Usage())
		case reflect.String:
			node.Pointer = p.flagset.String(node.Name(), node.Value(), node.Usage())
		case reflect.Bool:
			boo, err := strconv.ParseBool(node.Value())
			if node.Value() != "" && err != nil {
				p.logger.Printf("node (%s) has bad value (%s). change the struct tags: %s", node.Name(), node.Value(), err)
				return errBadStruct
			}
			node.Pointer = p.flagset.Bool(node.Name(), boo, node.Usage())
		}
		nodes[node.Name()] = node
	}

	if err := p.flagset.Parse(p.Args[1:]); err != nil {
		p.logger.Printf("error parsing: %s", err)
		return ErrParsing
	}

	p.flagset.Visit(func(f *flag.Flag) {
		nodes[f.Name].SetViaFlag = true
	})

	p.flagset.VisitAll(func(f *flag.Flag) {
		nodes[f.Name].Flag = f
	})

	var err error
	for _, node := range nodes {
		if node.Required() && !node.SetViaFlag {
			_, _ = fmt.Fprintln(p.Output, "missing required arg: ", node.Name())
			err = ErrMissingArgument
		}
	}
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if !node.ValueField.CanSet() || node.Pointer == nil {
			continue
		}
		switch node.ValueField.Kind() {
		case reflect.Int:
			if i, ok := node.Pointer.(*int); ok {
				node.ValueField.SetInt(int64(*i))
			} else {
				p.logger.Printf("integer type assertion failed on node (%s) for value (%s)", node.Name(), node.Value())
				return errTypeAssertion
			}
		case reflect.String:
			if s, ok := node.Pointer.(*string); ok {
				node.ValueField.SetString(*s)
			} else {
				p.logger.Printf("string type assertion failed on node (%s) for value (%s)", node.Name(), node.Value())
				return errTypeAssertion
			}
		case reflect.Bool:
			if b, ok := node.Pointer.(*bool); ok {
				node.ValueField.SetBool(*b)
			} else {
				p.logger.Printf("bool type assertion failed on node (%s) for value (%s)", node.Name(), node.Value())
				return errTypeAssertion
			}
		default:
			p.logger.Printf("not recognized: %s %s %s", node.Name(), node.Value(), node.TypeField)
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

func (p *parser) log(s string, v ...interface{}) {
	if p.logger != nil {
		p.logger.Printf("dflag: "+s+"\n", v)
	}
}
