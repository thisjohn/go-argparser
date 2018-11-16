package argparser // import "github.com/thisjohn/go-argparser"

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type meta struct {
	valPtr interface{}
	name   string
	ops    *options
}

func (m *meta) flag() string {
	return "-" + m.name
}

// ArgParser wraps `flag.FlagSet`
type ArgParser struct {
	name    string
	flagSet *flag.FlagSet
	metas   []*meta
}

// NewArgParser creates `ArgParser`
func NewArgParser() *ArgParser {
	return newArgParserWithName(os.Args[0])
}

func newArgParserWithName(name string) *ArgParser {
	flagset := flag.NewFlagSet(name, flag.ContinueOnError)
	flagset.SetOutput(&nullWriter{}) // Output nothing

	return &ArgParser{
		name:    name,
		flagSet: flagset,
		metas:   []*meta{},
	}
}

// Usage prints help message and err if any
func (p *ArgParser) Usage(anyErr error) {
	// Set output to stderr
	oriOutput := p.flagSet.Output()
	output := os.Stderr
	p.flagSet.SetOutput(output)

	// Print error
	if anyErr != nil {
		fmt.Fprintln(output, anyErr.Error())
	}

	// Print usage
	var requiredFlags []string
	for _, v := range p.metas {
		if v.ops.required {
			flagText := v.flag()
			if _, ok := v.valPtr.(*bool); !ok {
				flagText = fmt.Sprintf("%s <%s>", flagText, v.ops.usage)
			}
			requiredFlags = append(requiredFlags, flagText)
		}
	}
	fmt.Fprintf(output, "Usage: %s %s%s\n",
		p.name,
		strings.Join(requiredFlags, " "),
		func() string { // has more options?
			if len(p.metas) > len(requiredFlags) {
				return " [...]"
			}
			return ""
		}(),
	)

	p.flagSet.PrintDefaults()

	// Restore output
	p.flagSet.SetOutput(oriOutput)
}

// AddArgument defines how arg be parsed
//
// valPtr support bool, int, float, and string
func (p *ArgParser) AddArgument(valPtr interface{}, name string, setters ...Setter) error {
	// Default options
	ops := &options{
		usage:     p.defaultArgUsage(valPtr),
		validator: nullValidator,
	}
	for _, setter := range setters {
		setter(ops)
	}

	usage := ops.usage
	if ops.required {
		usage = "(Required) " + usage
	}

	switch ptr := valPtr.(type) {
	case *bool:
		dv, ok := ops.defaultVal.(bool)
		if !ok {
			dv = false
		}
		p.flagSet.BoolVar(ptr, name, dv, usage)

	case *int:
		dv, ok := ops.defaultVal.(int)
		if !ok {
			dv = 0
		}
		p.flagSet.IntVar(ptr, name, dv, usage)

	case *float64:
		dv, ok := ops.defaultVal.(float64)
		if !ok {
			dv = 0.0
		}
		p.flagSet.Float64Var(ptr, name, dv, usage)

	case *string:
		dv, ok := ops.defaultVal.(string)
		if !ok {
			dv = ""
		}
		p.flagSet.StringVar(ptr, name, dv, usage)

	default:
		return errors.New("Unknown type of valPtr")
	}

	p.metas = append(p.metas, &meta{valPtr: valPtr, name: name, ops: ops})

	return nil
}

func (p *ArgParser) defaultArgUsage(valPtr interface{}) string {
	switch valPtr.(type) {
	case *bool:
		return "somebool"
	case *int:
		return "someint"
	case *float64:
		return "somefloat"
	case *string:
		return "somestring"
	}
	return "something"
}

// Parse command line args
func (p *ArgParser) Parse() error {
	return p.parseWithArgs(os.Args[1:]...)
}

func (p *ArgParser) parseWithArgs(args ...string) (err error) {
	if err = p.flagSet.Parse(args); err != nil {
		return
	}

	err = p.validate()
	return
}

func (p *ArgParser) validate() error {
	for _, v := range p.metas {
		if err := v.ops.validator.Validate(v.valPtr); err != nil {
			msg := fmt.Sprintf("%s: %s", err, v.flag())
			return errors.New(msg)
		}
	}

	return nil
}

// ----------------
//  Null io.Writer
// ----------------

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
