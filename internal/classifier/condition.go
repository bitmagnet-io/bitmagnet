package classifier

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
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
	json_spec.HasJSONSchema
	name() string
	conditionCompiler
}

func compileCondition(ctx compilerContext) (condition, error) {
	var errs []error

	for _, def := range ctx.conditions {
		c, err := def.compileCondition(ctx.child(def.name(), ctx.Source))
		if err == nil {
			return c, nil
		}

		if json_spec.AsFatalCompilerError(err) != nil {
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
