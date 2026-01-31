package classifier

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/classifier/classification"
	"github.com/bitmagnet-io/bitmagnet/internal/json_spec"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

const findMatchName = "find_match"

type findMatchAction struct{}

func (findMatchAction) name() string {
	return findMatchName
}

var findMatchSpec = json_spec.SingleKeyValue[[]any]{
	Key: findMatchName,
	ValueSpec: json_spec.MustSucceed[[]any]{Typed: json_spec.List[any]{
		ItemSpec: json_spec.Generic[any]{
			Schema: json_schema.MustNew(
				json_schema.RefDefinition("action_single"),
			),
		},
	}},
	Description: "Iterate through a series of actions to find the first that does not return an unmatched error",
}

func (findMatchAction) compile(ctx compilerContext) (action, error) {
	payload, err := findMatchSpec.Parse(ctx.jsonSpec)
	if err != nil {
		return action{}, ctx.Error(err)
	}

	actions := make([]action, len(payload))

	for i, actionPayload := range payload {
		a, err := compileAction(ctx.child(json_spec.NumericPathPart(i), actionPayload))
		if err != nil {
			return action{}, err
		}

		actions[i] = a
	}

	path := ctx.Path

	return action{
		func(ctx executionContext) (classification.Result, error) {
			for _, action := range actions {
				result, err := action.run(ctx)
				if err != nil {
					if errors.Is(err, classification.ErrUnmatched) {
						continue
					}

					return classification.Result{}, classification.RuntimeError{
						Cause: err,
						Path:  path,
					}
				}

				return result, nil
			}

			return ctx.result, nil
		},
	}, nil
}

func (findMatchAction) JSONSchema() json_schema.JSONSchema {
	return findMatchSpec.JSONSchema()
}
