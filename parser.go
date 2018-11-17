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
	name         string
	flagSet      *flag.FlagSet
	metas        []*meta // flag arguments
	nonFlagMetas []*meta // non-flag arguments

	errors []error
}

// NewArgParser creates `ArgParser`
func NewArgParser() *ArgParser {
	return newArgParserWithName(os.Args[0])
}

func newArgParserWithName(name string) *ArgParser {
	flagset := flag.NewFlagSet(name, flag.ContinueOnError)
	flagset.SetOutput(&nullWriter{}) // Output nothing

	return &ArgParser{
		name:         name,
		flagSet:      flagset,
		metas:        []*meta{},
		nonFlagMetas: []*meta{},
		errors:       []error{},
	}
}

// EnableHelpArgument alias to `AddArgument(..., "h", Usage("Help"))`
func (p *ArgParser) EnableHelpArgument(valPtr *bool) *ArgParser {
	return p.AddArgument(valPtr, "h", Usage("Help"))
}

// AddArgument defines how flag argument be parsed
//
// valPtr support types of bool, int, float, and string
func (p *ArgParser) AddArgument(valPtr interface{}, name string, setters ...Setter) *ArgParser {
	// Default options
	ops := &options{
		shortDescription: p.defaultShortDescription(valPtr),
		validator:        nullValidator,
	}
	for _, setter := range setters {
		setter(ops)
	}

	usage := ops.usage
	if ops.required {
		usage = "(Required) " + usage
	}

	var errMsg string
	switch ptr := valPtr.(type) {
	case *bool:
		dv, ok := ops.defaultVal.(bool)
		if !ok && ops.defaultVal != nil {
			errMsg = "Type mismatch between valPtr and defaultVal"
			break
		}
		p.flagSet.BoolVar(ptr, name, dv, usage)

	case *int:
		dv, ok := ops.defaultVal.(int)
		if !ok && ops.defaultVal != nil {
			errMsg = "Type mismatch between valPtr and defaultVal"
			break
		}
		p.flagSet.IntVar(ptr, name, dv, usage)

	case *float64:
		dv, ok := ops.defaultVal.(float64)
		if !ok && ops.defaultVal != nil {
			errMsg = "Type mismatch between valPtr and defaultVal"
			break
		}
		p.flagSet.Float64Var(ptr, name, dv, usage)

	case *string:
		dv, ok := ops.defaultVal.(string)
		if !ok && ops.defaultVal != nil {
			errMsg = "Type mismatch between valPtr and defaultVal"
			break
		}
		p.flagSet.StringVar(ptr, name, dv, usage)

	default:
		errMsg = "Unknown type of valPtr"
	}

	if errMsg != "" {
		errMsg += ": " + name
		p.errors = append(p.errors, errors.New(errMsg))
		return p
	}

	p.metas = append(p.metas, &meta{valPtr: valPtr, name: name, ops: ops})
	return p
}

// AddNonFlagArgument defines non-flag argument.
// NOTE: Used for building usage
func (p *ArgParser) AddNonFlagArgument(name string, usage string, required bool) *ArgParser {
	// Default options
	ops := &options{
		usage:    usage,
		required: required,
	}

	if ops.required {
		ops.usage = "(Required) " + ops.usage
	}

	p.nonFlagMetas = append(p.nonFlagMetas, &meta{valPtr: nil, name: name, ops: ops})
	return p
}

func (p *ArgParser) defaultShortDescription(valPtr interface{}) string {
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
	if len(p.errors) > 0 {
		return p.errors[0] // Only return the first error
	}

	if err = p.flagSet.Parse(args); err != nil {
		return
	}
	if err = p.validate(); err != nil {
		return
	}

	nRequiredNonFlags := func() int {
		n := 0
		for _, v := range p.nonFlagMetas {
			if v.ops.required {
				n++
			}
		}
		return n
	}()
	if len(p.Args()) < nRequiredNonFlags {
		err = errors.New("Insufficient number of non-flag arguments")
		return
	}

	return
}

func (p *ArgParser) validate() error {
	for _, v := range p.metas { // Validate flag arguments
		if err := v.ops.validator.Validate(v.valPtr); err != nil {
			msg := fmt.Sprintf("%s: %s", err, v.flag())
			return errors.New(msg)
		}
	}

	return nil
}

// Args returns the non-flag arguments
func (p *ArgParser) Args() []string {
	return p.flagSet.Args()
}

// PrintUsage and err if any
func (p *ArgParser) PrintUsage(anyErr error) {
	// Set output to stderr
	oriOutput := p.flagSet.Output()
	output := os.Stderr
	p.flagSet.SetOutput(output)

	// Print error
	if anyErr != nil {
		fmt.Fprintln(output, anyErr.Error())
	}

	// Print usage

	// - Build texts
	var flagTexts []string
	var hasOptionalFlag bool
	for _, v := range p.metas { // flag arguments
		if v.ops.required {
			text := v.flag()
			if _, ok := v.valPtr.(*bool); !ok {
				text = fmt.Sprintf("%s <%s>", text, v.ops.shortDescription)
			}
			flagTexts = append(flagTexts, text)
		} else {
			hasOptionalFlag = true
		}
	}

	var nonFlagTexts []string
	for _, v := range p.nonFlagMetas { // non-flag arguments
		if v.ops.required {
			nonFlagTexts = append(nonFlagTexts, v.name)
		}
	}

	// - Prints
	fmt.Fprintf(output, "Usage: %s%s%s%s\n",
		p.name,
		func() string {
			leadSpace := ""
			if len(flagTexts) > 0 {
				leadSpace = " "
			}
			return leadSpace + strings.Join(flagTexts, " ")
		}(),
		func() string {
			if hasOptionalFlag {
				return " [...]"
			}
			return ""
		}(),
		func() string {
			leadSpace := ""
			if len(nonFlagTexts) > 0 {
				leadSpace = " "
			}
			return leadSpace + strings.Join(nonFlagTexts, " ")
		}(),
	)

	p.flagSet.PrintDefaults()

	for _, v := range p.nonFlagMetas {
		fmt.Fprintf(output, "  %s\t%s\n", v.name, v.ops.usage)
	}

	// Restore output
	p.flagSet.SetOutput(oriOutput)
}

// ----------------
//  Null io.Writer
// ----------------

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
