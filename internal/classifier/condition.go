package classifier

import (
	"errors"
)

func conditions(defs ...conditionDefinition) option {
	return func(_ workflowSource, c *compilerContext) error {
		c.conditions = append(c.conditions, defs...)
		return nil
	}
}

type conditionCompiler interface {
	compileCondition(ctx compilerContext) (condition, error)
}

type conditionDefinition interface {
	name() string
	conditionCompiler
}

func (c compilerContext) compileCondition(ctx compilerContext) (condition, error) {
	var errs []error
	for _, def := range c.conditions {
		c, err := def.compileCondition(ctx.child(def.name(), ctx.source))
		if err == nil {
			return c, nil
		}
		if asFatalCompilerError(err) != nil {
			return condition{}, err
		}
		errs = append(errs, err)
	}
	errs = append(errs, errors.New("no condition matched"))
	return condition{}, errors.Join(errs...)
}

type condition struct {
	check func(executionContext) (bool, error)
}
