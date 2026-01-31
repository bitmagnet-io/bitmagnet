package json_spec

import (
	"errors"
	"fmt"
	"strings"
)

type ParseContext struct {
	KeyMatcher KeyMatcher
	Source     any
	Path       []string
}

func (c ParseContext) Child(pathPart string, source any) ParseContext {
	c.Source = source
	newPath := make([]string, len(c.Path), len(c.Path)+1)
	copy(newPath, c.Path)
	newPath = append(newPath, pathPart)
	c.Path = newPath

	return c
}

func (c ParseContext) Error(cause error) error {
	if AsCompilerError(cause) != nil {
		return cause
	}

	return CompilerError{c.Path, cause}
}

func (c ParseContext) Fatal(cause error) error {
	if AsFatalCompilerError(cause) != nil {
		return cause
	}

	cErr := AsCompilerError(cause)
	if cErr != nil {
		return FatalCompilerError{CompilerError: *cErr}
	}

	return FatalCompilerError{CompilerError{c.Path, cause}}
}

type CompilerError struct {
	path  []string
	cause error
}

func (e CompilerError) Error() string {
	return fmt.Sprintf("compiler error at path '%s': %s", strings.Join(e.path, "."), e.cause)
}

func (e CompilerError) Unwrap() error {
	return e.cause
}

func AsCompilerError(err error) *CompilerError {
	ue := &CompilerError{}
	if ok := errors.As(err, ue); ok {
		return ue
	}

	return nil
}

type FatalCompilerError struct {
	CompilerError
}

func (e FatalCompilerError) Unwrap() error {
	return e.CompilerError
}

func AsFatalCompilerError(err error) *FatalCompilerError {
	ue := &FatalCompilerError{}
	if ok := errors.As(err, ue); ok {
		return ue
	}

	return nil
}

func NumericPathPart(num int) string {
	return fmt.Sprintf("[%d]", num)
}
