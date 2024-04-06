package classifier

import (
	"errors"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
)

func actions(defs ...actionDefinition) option {
	return func(_ WorkflowSource, c *compilerContext) error {
		c.actions = append(c.actions, defs...)
		return nil
	}
}

type actionCompiler interface {
	compileAction(ctx compilerContext) (action, error)
}

type actionDefinition interface {
	Name() string
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
			a, err := def.compileAction(actionCtx.child(def.Name(), rawAction))
			if err == nil {
				actions = append(actions, a)
				continue outer
			} else {
				if asFatalCompilerError(err) != nil {
					return action{}, err
				}
			}
		}
		errs = append(errs, fmt.Errorf("no action matched: %v", ctx.source))
	}
	if len(errs) > 0 {
		return action{}, errors.Join(errs...)
	}
	return action{func(ctx executionContext) (classification.Result, error) {
		for _, a := range actions {
			result, err := a.run(ctx)
			if err != nil {
				return classification.Result{}, err
			}
			ctx = ctx.withResult(result)
		}
		return ctx.result, nil
	}}, errors.Join(errs...)
}

type action struct {
	run func(executionContext) (classification.Result, error)
}
