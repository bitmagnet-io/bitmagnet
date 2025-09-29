package classifier

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protobuf"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type Compiler interface {
	Compile(source Source) (Runner, error)
}

type Runner interface {
	Run(ctx context.Context, workflow Workflow, flags Flags, t model.Torrent) (classification.Result, error)
}

type compiler struct {
	options      []compilerOption
	dependencies dependencies
}

type jsonSpec = json_spec.ParseContext

type compilerContext struct {
	jsonSpec
	features
	celEnv        *cel.Env
	workflowNames map[Workflow]struct{}
}

type compilerOption func(Source, *compilerContext) error

type executionContext struct {
	context.Context
	dependencies
	flags     map[string]ref.Val
	workflows map[Workflow]action
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
	c.jsonSpec = c.jsonSpec.Child(pathPart, source)

	return c
}

func (c compiler) Compile(source Source) (Runner, error) {
	ctx := &compilerContext{
		jsonSpec: json_spec.ParseContext{
			Source:     source,
			KeyMatcher: json_spec.KeyMatcherSnake,
		},
		workflowNames: source.workflowNames(),
	}
	source, sourceErr := json_spec.Decode[Source](ctx.jsonSpec)

	if sourceErr != nil {
		return nil, ctx.Fatal(sourceErr)
	}

	for _, opt := range c.options {
		if err := opt(source, ctx); err != nil {
			return nil, ctx.Fatal(err)
		}
	}

	workflowsCtx := ctx.child("workflows", source.Workflows)
	workflows := make(map[Workflow]action)

	for name, src := range source.Workflows {
		a, err := compileAction(workflowsCtx.child(string(name), src))
		if err != nil {
			return nil, ctx.Fatal(err)
		}

		workflows[name] = a
	}

	cfs := make(compiledFlags, len(source.FlagDefinitions))

	for k, def := range source.FlagDefinitions {
		rawVal, ok := source.Flags[k]
		if !ok {
			return nil, ctx.Fatal(fmt.Errorf("missing value for flag '%q'", k))
		}

		val, err := def.celVal(rawVal)
		if err != nil {
			return nil, ctx.Fatal(fmt.Errorf("invalid value for flag '%s': %w", k, err))
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
