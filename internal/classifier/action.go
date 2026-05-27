package classifier

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
)

func actions(defs ...actionDefinition) feature {
	return func(c *features) {
		c.actions = append(c.actions, defs...)
	}
}

type actionCompiler interface {
	compileAction(ctx compilerContext) (action, error)
}

type actionDefinition interface {
	HasJSONSchema
	name() string
	actionCompiler
}

func (c compilerContext) compileAction(ctx compilerContext) (action, error) {
	var rawActions []any

	isArray := false

	if s, ok := ctx.source.([]any); ok {
		rawActions = s
		isArray = true
	} else {
		rawActions = []any{ctx.source}
	}

	var actions []action

	var errs []error
outer:
	for i, rawAction := range rawActions {
		actionCtx := ctx
		if isArray {
			actionCtx = ctx.child(numericPathPart(i), rawAction)
		}
		for _, def := range c.actions {
			a, err := def.compileAction(actionCtx.child(def.name(), rawAction))
			if err == nil {
				actions = append(actions, a)
				continue outer
			}
			if asFatalCompilerError(err) != nil {
				return action{}, err
			}
		}
		errs = append(errs, fmt.Errorf("no action matched: %v", ctx.source))
	}

	if len(errs) > 0 {
		return action{}, errors.Join(errs...)
	}

	path := ctx.path

	return action{func(ctx executionContext) (classification.Result, error) {
		for i, a := range actions {
			ctx.logger = ctx._logger.Named(strings.Join(path, ".")+".["+fmt.Sprint(i)+"]")
			result, err := a.run(ctx)
			if err != nil {
				return classification.Result{}, err
			}
			ctx = ctx.withResult(result)
		}
		ctx.logger = ctx._logger
		return ctx.result, nil
	}}, errors.Join(errs...)
}

type action struct {
	run func(executionContext) (classification.Result, error)
}
