package argparser // import "github.com/thisjohn/go-argparser"

type options struct {
	defaultVal       interface{}
	shortDescription string
	usage            string
	required         bool
	validator        ArgValidator
}

// Setter is function for setting options
type Setter func(*options)

// DefaultVal set defaultValue
func DefaultVal(defautVal interface{}) Setter {
	return func(args *options) {
		args.defaultVal = defautVal
	}
}

// ShortDescription set short description, default is val type name
//
// ex. Usage: ./main -c <short description>
func ShortDescription(desc string) Setter {
	return func(args *options) {
		args.shortDescription = desc
	}
}

// Usage set detailed description
//
// ex. Usage: ./main -c <short description>
//       -c string
//       will show detailed description here
func Usage(usage string) Setter {
	return func(args *options) {
		args.usage = usage
	}
}

// Required set argument is required
func Required() Setter {
	return func(args *options) {
		args.required = true
		Validator(requiredValidator)(args)
	}
}

// Validator set argument validator
func Validator(validator ArgValidator) Setter {
	return func(args *options) {
		args.validator = validator
	}
}
