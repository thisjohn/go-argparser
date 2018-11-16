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
		"positive bool flag": {
			cmdLine:     "-b",
			expectedVal: true,
		},
		"default bool flag": {
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
		parser := newArgParserWithName("somename")
		parser.AddArgument(&val, "b", testcase.options...)

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
		"int flag": {
			cmdLine:     "-i 999",
			expectedVal: 999,
		},
		"default int flag": {
			cmdLine: "",
			options: []Setter{
				DefaultVal(777),
			},
			expectedVal: 777,
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
		parser := newArgParserWithName("somename")
		parser.AddArgument(&val, "i", testcase.options...)

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
		"float flag": {
			cmdLine:     "-f 3.7",
			expectedVal: 3.7,
		},
		"default float flag": {
			cmdLine: "",
			options: []Setter{
				DefaultVal(6.1),
			},
			expectedVal: 6.1,
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
		parser := newArgParserWithName("somename")
		parser.AddArgument(&val, "f", testcase.options...)

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
		"string flag": {
			cmdLine:     "-s foo",
			expectedVal: "foo",
		},
		"default string flag": {
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
		parser := newArgParserWithName("somename")
		parser.AddArgument(&val, "s", testcase.options...)

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
