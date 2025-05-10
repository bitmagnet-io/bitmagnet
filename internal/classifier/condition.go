package classifier

import (
	"errors"
)

func conditions(defs ...conditionDefinition) feature {
	return func(c *features) {
		c.conditions = append(c.conditions, defs...)
	}
}

type conditionCompiler interface {
	compileCondition(ctx compilerContext) (condition, error)
}

type conditionDefinition interface {
	HasJSONSchema
	name() string
	conditionCompiler
}

func (c compilerContext) compileCondition(ctx compilerContext) (condition, error) {
	//nolint:prealloc
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
