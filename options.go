package argparser // import "github.com/thisjohn/go-argparser"

type options struct {
	defaultVal interface{}
	usage      string
	validator  ArgValidator
}

// Setter represents a setter function for setting options
type Setter func(*options)

// DefaultVal set defaultValue
func DefaultVal(defautVal interface{}) Setter {
	return func(args *options) {
		args.defaultVal = defautVal
	}
}

// Usage set usage
func Usage(usage string) Setter {
	return func(args *options) {
		args.usage = usage
	}
}

// Required set requiredValidator
func Required() Setter {
	return func(args *options) {
		args.validator = requiredValidator
	}
}

// Validator set validator
func Validator(validator ArgValidator) Setter {
	return func(args *options) {
		args.validator = validator
	}
}
