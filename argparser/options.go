package argparser

type options struct {
	defaultVal       interface{}
	shortDescription string
	usage            string
	required         bool
	validators       []ArgValidator
}

// Setter is function for setting options
type Setter func(*options)

// OptDefaultVal set defaultValue
func OptDefaultVal(defautVal interface{}) Setter {
	return func(args *options) {
		args.defaultVal = defautVal
	}
}

// OptShortDescription set short description, default is val type name
//
// ex. OptUsage: ./main -c <short description>
func OptShortDescription(desc string) Setter {
	return func(args *options) {
		args.shortDescription = desc
	}
}

// OptUsage set detailed description
//
// ex. OptUsage: ./main -c <short description>
//       -c string
//       will show detailed description here
func OptUsage(usage string) Setter {
	return func(args *options) {
		args.usage = usage
	}
}

// OptRequired set argument is required
func OptRequired() Setter {
	return func(args *options) {
		args.required = true
		OptValidator(requiredValidator)(args)
	}
}

// OptValidator add argument validator
func OptValidator(validator ArgValidator) Setter {
	return func(args *options) {
		args.validators = append(args.validators, validator)
	}
}
