package classifier

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type Compiler interface {
	Compile(source Source) (Runner, error)
}

type Runner interface {
	Run(ctx context.Context, workflow string, flags Flags, t model.Torrent) (classification.Result, error)
}

type compiler struct {
	options      []compilerOption
	dependencies dependencies
}

type compilerContext struct {
	features
	celEnv        *cel.Env
	source        any
	path          []string
	workflowNames map[string]struct{}
}

type compilerOption func(Source, *compilerContext) error

type executionContext struct {
	context.Context
	dependencies
	flags     map[string]ref.Val
	workflows map[string]action
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

func (c compiler) Compile(source Source) (Runner, error) {
	ctx := &compilerContext{
		source:        source,
		workflowNames: source.workflowNames(),
	}
	source, sourceErr := decode[Source](*ctx)
	if sourceErr != nil {
		return nil, ctx.fatal(sourceErr)
	}
	for _, opt := range c.options {
		if err := opt(source, ctx); err != nil {
			return nil, ctx.fatal(err)
		}
	}
	workflowsCtx := ctx.child("workflows", source.Workflows)
	workflows := make(map[string]action)
	for name, src := range source.Workflows {
		a, err := ctx.compileAction(workflowsCtx.child(name, src))
		if err != nil {
			return nil, ctx.fatal(err)
		}
		workflows[name] = a
	}
	cfs := make(compiledFlags, len(source.FlagDefinitions))
	for k, def := range source.FlagDefinitions {
		rawVal, ok := source.Flags[k]
		if !ok {
			return nil, ctx.fatal(fmt.Errorf("missing value for flag '%q'", k))
		}
		val, err := def.celVal(rawVal)
		if err != nil {
			return nil, ctx.fatal(fmt.Errorf("invalid value for flag '%s': %w", k, err))
		}
		cfs[k] = val
	}
	return runner{
		dependencies:    c.dependencies,
		flagDefinitions: source.FlagDefinitions,
		compiledFlags:   cfs,
		workflows:       workflows,
	}, nil
}

func decodeTo[T any](ctx compilerContext, target *T) error {
	decoder, decoderErr := newDecoder(target)
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
