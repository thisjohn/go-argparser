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
		err := v.Validate(tt.val)
		assert.NoError(t, err, k)
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
		err := v.Validate(tt.val)
		if tt.expectedPass {
			assert.NoError(t, err, k)
		} else {
			assert.Error(t, err, k)
		}
	}
}
