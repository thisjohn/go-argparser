package argparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullValidator(t *testing.T) {
	v := nullValidator

	tests := map[string]struct {
		val interface{}
	}{
		"bool value": {
			val: true,
		},
		"empty bool value": {
			val: false,
		},
		"int value": {
			val: 1,
		},
		"empty int value": {
			val: 0,
		},
		"string value": {
			val: "somestring",
		},
		"empty string value": {
			val: "",
		},
	}

	for k, tt := range tests {
		pass := v.Validate(tt.val)
		assert.True(t, pass, k)
	}
}

func TestRequiredValidator(t *testing.T) {
	v := requiredValidator

	tests := map[string]struct {
		val          interface{}
		expectedPass bool
	}{
		"bool value": {
			val:          true,
			expectedPass: true,
		},
		"empty bool value": {
			val:          false,
			expectedPass: true,
		},
		"int value": {
			val:          1,
			expectedPass: true,
		},
		"empty int value": {
			val:          0,
			expectedPass: false,
		},
		"string value": {
			val:          "somestring",
			expectedPass: true,
		},
		"empty string value": {
			val:          "",
			expectedPass: false,
		},
	}

	for k, tt := range tests {
		pass := v.Validate(tt.val)
		assert.Equal(t, tt.expectedPass, pass, k)
	}
}
