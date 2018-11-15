package argparser // import "github.com/thisjohn/go-argparser"

var nullValidator = &nullArgValidator{} // Default validator
var requiredValidator = &requiredArgValidator{}

// ArgValidator is an interface to validate arg value
type ArgValidator interface {
	Validate(val interface{}) bool
}

type nullArgValidator struct {
}

// Validate implements the interface of `ArgValidator`
func (v *nullArgValidator) Validate(val interface{}) bool {
	// Always pass
	return true
}

type requiredArgValidator struct {
}

// Validate implements the interface of `ArgValidator`
func (v *requiredArgValidator) Validate(val interface{}) bool {
	switch ptr := val.(type) {
	// don't care
	case bool:
		return true
	case *bool:
		return true

		// int cannot be zero
	case int:
		return ptr != 0
	case *int:
		return *ptr != 0

		// string cannot be empty
	case string:
		return ptr != ""
	case *string:
		return *ptr != ""
	}

	// Unknown types
	return false
}
