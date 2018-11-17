package argparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgParser_Bool(t *testing.T) {
	// Prepare
	testcases := map[string]struct {
		cmdLine     string
		options     []Setter
		expectedVal bool
	}{
		"empty bool flag": {
			cmdLine:     "",
			expectedVal: false,
		},
		"positive bool arg": {
			cmdLine:     "-b",
			expectedVal: true,
		},
		"default bool": {
			cmdLine: "",
			options: []Setter{
				DefaultVal(true),
			},
			expectedVal: true,
		},
	}

	// Act
	for k, testcase := range testcases {
		var val bool
		parser := newArgParserWithName("somename").
			AddArgument(&val, "b", testcase.options...)

		err := parser.parseWithArgs(strings.Split(testcase.cmdLine, " ")...)

		// Assert
		assrt := assert.New(t)
		assrt.NoError(err)
		assrt.Equal(testcase.expectedVal, val, k)
	}
}

func TestArgParser_Int(t *testing.T) {
	// Prepare
	testcases := map[string]struct {
		cmdLine          string
		options          []Setter
		expectedParseErr bool
		expectedVal      int
	}{
		"empty int flag": {
			cmdLine:     "",
			expectedVal: 0,
		},
		"int arg": {
			cmdLine:     "-i 999",
			expectedVal: 999,
		},
		"non-int arg": {
			cmdLine:          "-i notanumber",
			expectedParseErr: true,
		},
		"default int": {
			cmdLine: "",
			options: []Setter{
				DefaultVal(777),
			},
			expectedVal: 777,
		},
		"mismatch default int": {
			cmdLine: "",
			options: []Setter{
				DefaultVal("notanumber"),
			},
			expectedParseErr: true,
		},
		"required int flag": {
			cmdLine: "",
			options: []Setter{
				Required(),
			},
			expectedParseErr: true,
		},
	}

	// Act
	for k, testcase := range testcases {
		var val int
		parser := newArgParserWithName("somename").
			AddArgument(&val, "i", testcase.options...)

		err := parser.parseWithArgs(strings.Split(testcase.cmdLine, " ")...)

		// Assert
		assrt := assert.New(t)
		if testcase.expectedParseErr {
			assrt.Error(err, k)
		} else {
			assrt.NoError(err, k)
		}
		assrt.Equal(testcase.expectedVal, val, k)
	}
}

func TestArgParser_Float(t *testing.T) {
	// Prepare
	testcases := map[string]struct {
		cmdLine          string
		options          []Setter
		expectedParseErr bool
		expectedVal      float64
	}{
		"empty float flag": {
			cmdLine:     "",
			expectedVal: 0.0,
		},
		"float arg": {
			cmdLine:     "-f 3.7",
			expectedVal: 3.7,
		},
		"non-float arg": {
			cmdLine:          "-f notafloat",
			expectedParseErr: true,
		},
		"default float": {
			cmdLine: "",
			options: []Setter{
				DefaultVal(6.1),
			},
			expectedVal: 6.1,
		},
		"mismatch default float": {
			cmdLine: "",
			options: []Setter{
				DefaultVal("notafloat"),
			},
			expectedParseErr: true,
		},
		"required float flag": {
			cmdLine: "",
			options: []Setter{
				Required(),
			},
			expectedParseErr: true,
		},
	}

	// Act
	for k, testcase := range testcases {
		var val float64
		parser := newArgParserWithName("somename").
			AddArgument(&val, "f", testcase.options...)

		err := parser.parseWithArgs(strings.Split(testcase.cmdLine, " ")...)

		// Assert
		assrt := assert.New(t)
		if testcase.expectedParseErr {
			assrt.Error(err, k)
		} else {
			assrt.NoError(err, k)
		}
		assrt.Equal(testcase.expectedVal, val, k)
	}
}

func TestArgParser_String(t *testing.T) {
	// Prepare
	testcases := map[string]struct {
		cmdLine          string
		options          []Setter
		expectedParseErr bool
		expectedVal      string
	}{
		"empty string flag": {
			cmdLine:     "",
			expectedVal: "",
		},
		"string arg": {
			cmdLine:     "-s foo",
			expectedVal: "foo",
		},
		"default string": {
			cmdLine: "",
			options: []Setter{
				DefaultVal("bar"),
			},
			expectedVal: "bar",
		},
		"required string flag": {
			cmdLine: "",
			options: []Setter{
				Required(),
			},
			expectedParseErr: true,
		},
	}

	// Act
	for k, testcase := range testcases {
		var val string
		parser := newArgParserWithName("somename").
			AddArgument(&val, "s", testcase.options...)

		err := parser.parseWithArgs(strings.Split(testcase.cmdLine, " ")...)

		// Assert
		assrt := assert.New(t)
		if testcase.expectedParseErr {
			assrt.Error(err, k)
		} else {
			assrt.NoError(err, k)
		}
		assrt.Equal(testcase.expectedVal, val, k)
	}
}
