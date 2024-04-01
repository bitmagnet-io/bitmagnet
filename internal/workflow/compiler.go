package workflow

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/google/cel-go/cel"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type Compiler interface {
	Compile(payload any) (Workflow, error)
}

type Workflow interface {
	Run(context.Context, model.Torrent) (classification.Result, error)
}

type compiler []option

type compilerContext struct {
	celEnv     *cel.Env
	conditions []conditionDefinition
	actions    []actionDefinition
	source     any
	path       []string
	vars       map[string]any
}

type option func(workflowSource, *compilerContext) error

type executionContext struct {
	context.Context
	torrent   model.Torrent
	torrentPb *protobuf.Torrent
	result    classification.Result
	resultPb  *protobuf.Classification
}

func (c executionContext) withResult(result classification.Result) executionContext {
	c.result = result
	c.resultPb = protobuf.NewClassification(result)
	return c
}

func (c compilerContext) child(pathPart string, source any) compilerContext {
	c.source = source
	newPath := make([]string, len(c.path), len(c.path)+1)
	copy(newPath, c.path)
	newPath = append(newPath, pathPart)
	c.path = newPath
	return c
}

func (c compilerContext) error(cause error) error {
	if asCompilerError(cause) != nil {
		return cause
	}
	return compilerError{c.path, cause}
}

func (c compilerContext) fatal(cause error) error {
	if asFatalCompilerError(cause) != nil {
		return cause
	}
	cErr := asCompilerError(cause)
	if cErr != nil {
		return fatalCompilerError{compilerError: *cErr}
	}
	return fatalCompilerError{compilerError{c.path, cause}}
}

func (r compiler) Compile(payload any) (Workflow, error) {
	ctx := &compilerContext{source: payload}
	source, sourceErr := decode[workflowSource](*ctx)
	if sourceErr != nil {
		return nil, ctx.fatal(sourceErr)
	}
	for _, c := range r {
		if err := c(source, ctx); err != nil {
			return nil, ctx.fatal(err)
		}
	}
	a, err := ctx.compileAction(ctx.child("actions", source.Actions))
	if err != nil {
		return nil, ctx.fatal(err)
	}
	return workflow{action: a}, nil
}

type keywordGroups map[string][]string

type extensionGroups map[string][]string

type workflowSource struct {
	Name            string
	FlagDefinitions flagDefinitions
	Flags           flags
	Keywords        keywordGroups
	Extensions      extensionGroups
	Actions         any
}

type workflow struct {
	action action
}

func decodeTo[T any](ctx compilerContext, target *T) error {
	decoder, decoderErr := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: target,
		MatchName: func(mapKey, fieldName string) bool {
			return mapKey == strcase.ToSnake(fieldName)
		},
		ErrorUnused: true,
	})
	if decoderErr != nil {
		return ctx.error(decoderErr)
	}
	return decoder.Decode(ctx.source)
}

func decode[T any](ctx compilerContext) (T, error) {
	var target T
	err := decodeTo(ctx, &target)
	return target, err
}

type compilerError struct {
	path  []string
	cause error
}

func (e compilerError) Error() string {
	return fmt.Sprintf("compiler error at path '%s': %s", strings.Join(e.path, "."), e.cause)
}

func (e compilerError) Unwrap() error {
	return e.cause
}

func asCompilerError(err error) *compilerError {
	ue := &compilerError{}
	if ok := errors.As(err, ue); ok {
		return ue
	}
	return nil
}

type fatalCompilerError struct {
	compilerError
}

func (e fatalCompilerError) Unwrap() error {
	return e.compilerError
}

func asFatalCompilerError(err error) *fatalCompilerError {
	ue := &fatalCompilerError{}
	if ok := errors.As(err, ue); ok {
		return ue
	}
	return nil
}

func numericPathPart(num int) string {
	return fmt.Sprintf("[%d]", num)
}
