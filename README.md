# GO Command Line Arguments Parser

`ArgParser` wraps `flag.FlagSet`

## Example

```go
package main

import (
	"errors"
	"fmt"

	"github.com/thisjohn/go-argparser/argparser"
)

type arguments struct {
	h      bool
	text   string
	num    float64
	evenno int
	extra  string
}

func main() {
	args := &arguments{}

	// Parse arguments
	parser := argparser.NewArgParser().
		// alias -h
		EnableHelpArgument(&args.h).

		// -t
		AddArgument(&args.text, "t",
			argparser.OptRequired(),
			argparser.OptShortDescription("text"),
			argparser.OptUsage("Description for text"),
		).

		// -num
		AddArgument(&args.num, "num",
			argparser.OptDefaultVal(9.9),
			argparser.OptUsage("Description for num"),
		).

		// -evenno
		AddArgument(&args.evenno, "evenno",
			argparser.OptUsage("Description for evenno"),
			argparser.OptValidator(&evenNumberValidator{}),
		).

		// trailing args
		AddNonFlagArgument("extra", "Description for extra", true)
	err := parser.Parse()

	if args.h { // Print help first
		parser.PrintUsage(nil)
		return
	}
	if err != nil { // Then, check error
		parser.PrintUsage(err)
		return
	}

	fmt.Println(args, parser.Args())
}

type evenNumberValidator struct{}

func (v *evenNumberValidator) Validate(val interface{}) error {
	var no int
	switch tv := val.(type) {
	case int:
		no = tv
	case *int:
		no = *tv
	default:
		return errors.New("val must be integer")
	}

	if no%2 != 0 {
		msg := fmt.Sprintf("val is not an even number (%d)", no)
		return errors.New(msg)
	}
	return nil
}
```

Command line:
```console
$ ./main -h
Usage: ./main -t <text> [...] extra
  -evenno int
        Description for evenno
  -h    Help
  -num float
        Description for num (default 9.9)
  -t string
        (Required) Description for text
  extra (Required) Description for extra

$ ./main
val is required: -t

$ ./main -t
flag needs an argument: -t

$ ./main -t sometext -num notnumber
invalid value "notnumber" for flag -num: strconv.ParseFloat: parsing "notnumber": invalid syntax

$ ./main -t sometext -evenno 7
val is not an even number (7): -evenno

$ ./main -t sometext lollol
&{false sometext 9.9 0 } [lollol]
```
