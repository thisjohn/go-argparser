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
	return &ArgParser{
		flagSet: flag.NewFlagSet(name, flag.PanicOnError),
		metaMap: map[string]*meta{},
	}
}

// AddArgument defines how arg be parsed
//
// valPtr support bool, int, and string
func (p *ArgParser) AddArgument(valPtr interface{}, name string, setters ...Setter) error {
	// Default options
	ops := &options{
		validator: nullValidator,
	}
	for _, setter := range setters {
		setter(ops)
	}

	switch ptr := valPtr.(type) {
	case *bool:
		dv, ok := ops.defaultVal.(bool)
		if !ok {
			dv = false
		}
		p.flagSet.BoolVar(ptr, name, dv, ops.usage)

	case *int:
		dv, ok := ops.defaultVal.(int)
		if !ok {
			dv = 0
		}
		p.flagSet.IntVar(ptr, name, dv, ops.usage)

	case *string:
		dv, ok := ops.defaultVal.(string)
		if !ok {
			dv = ""
		}
		p.flagSet.StringVar(ptr, name, dv, ops.usage)

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
	// Handle and return panic error
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println(r)

			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	if err = p.flagSet.Parse(args); err != nil {
		return
	}

	err = p.validate()
	return
}

func (p *ArgParser) validate() error {
	for _, v := range p.metaMap {
		if pass := v.options.validator.Validate(v.valPtr); !pass {
			msg := fmt.Sprintf("flag validate failed: %s", v.name)
			return errors.New(msg)
		}
	}

	return nil
}
