package argparser // import "github.com/thisjohn/go-argparser"

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type meta struct {
	valPtr  interface{}
	name    string
	options *options
}

// ArgParser wraps `flag.FlagSet` and able to pass arg validator
type ArgParser struct {
	flagSet *flag.FlagSet
	metaMap map[string]*meta
}

// NewArgParser creates `ArgParser`
func NewArgParser() *ArgParser {
	return newArgParserWithName(os.Args[0])
}

func newArgParserWithName(name string) *ArgParser {
	flagset := flag.NewFlagSet(name, flag.ContinueOnError)
	flagset.SetOutput(&nullWriter{}) // Output nothing

	return &ArgParser{
		flagSet: flagset,
		metaMap: map[string]*meta{},
	}
}

// Usage prints err if any and help message
func (p *ArgParser) Usage(anyErr error) {
	if anyErr != nil {
		fmt.Fprintln(os.Stderr, anyErr.Error())
	}

	output := p.flagSet.Output()
	p.flagSet.SetOutput(os.Stderr)
	p.flagSet.Usage()
	p.flagSet.SetOutput(output) // Restore output
}

// AddArgument defines how arg be parsed
//
// valPtr support bool, int, float, and string
func (p *ArgParser) AddArgument(valPtr interface{}, name string, setters ...Setter) error {
	// Default options
	ops := &options{
		validator: nullValidator,
	}
	for _, setter := range setters {
		setter(ops)
	}

	usage := ops.usage
	if ops.required {
		usage = "(required) " + usage
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

	p.metaMap[name] = &meta{valPtr: valPtr, name: name, options: ops}

	return nil
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
	for _, v := range p.metaMap {
		if err := v.options.validator.Validate(v.valPtr); err != nil {
			msg := fmt.Sprintf("%s: %s", err, v.name)
			return errors.New(msg)
		}
	}

	return nil
}

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
