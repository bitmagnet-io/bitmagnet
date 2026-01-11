package cmd

import (
	"errors"
	"fmt"
)

var (
	Err = errors.New("arg")

	// Compilation errors:

	ErrCompilation = errors.New("compilation failed")

	ErrZeroValue            = errors.New("zero value provided for introspection")
	ErrNonStructValue       = errors.New("non-struct value provided for introspection")
	ErrCmdNotEmbedded       = fmt.Errorf("should be embedded: %T", Cmd{})
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidAbbr          = errors.New("invalid abbreviation")
	ErrUnknownCmdTag        = fmt.Errorf("unknown struct tag for %T", Cmd{})
	ErrUnknownParamTag      = errors.New("unknown struct tag for parameter")
	ErrUnsupportedParamType = errors.New("unsupported parameter type")

	// Execution errors:

	ErrExecution     = errors.New("execution failed")
	ErrHelp          = errors.New("help failed")
	ErrUninitialized = errors.New("command uninitialized")

	ErrSetup    = errors.New("setup")
	ErrTeardown = errors.New("teardown")
	ErrOnError  = errors.New("on error")

	// Input errors:

	ErrInvalidArgs = errors.New("invalid arguments")

	ErrUnknownParam       = errors.New("unknown parameter")
	ErrRepeatedParam      = errors.New("repeated parameter")
	ErrUnexpectedArgument = errors.New("unexpected argument")
	ErrMissingParamValue  = errors.New("missing parameter value")
	ErrParseParamValue    = errors.New("failed to parse parameter value")
)
