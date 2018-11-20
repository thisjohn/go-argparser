package argparser

import (
	"errors"
)

var nullValidator = &nullArgValidator{} // Default validator
var requiredValidator = &requiredArgValidator{}

// ArgValidator is an interface to validate argument value
type ArgValidator interface {
	Validate(val interface{}) error
}

type nullArgValidator struct {
}

// Validate implements the interface of `ArgValidator`
func (v *nullArgValidator) Validate(val interface{}) error {
	// Always pass
	return nil
}

type requiredArgValidator struct {
}

// Validate implements the interface of `ArgValidator`
func (v *requiredArgValidator) Validate(val interface{}) error {
	switch tv := val.(type) {
	// don't care
	case bool:
		return nil
	case *bool:
		return nil

	// number cannot be zero
	case int, float64:
		if tv != 0 {
			return nil
		}
	case *int:
		if *tv != 0 {
			return nil
		}
	case *float64:
		if *tv != 0 {
			return nil
		}

	// string cannot be empty
	case string:
		if tv != "" {
			return nil
		}
	case *string:
		if *tv != "" {
			return nil
		}
	}

	return errors.New("val is required")
}
