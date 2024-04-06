package classifier

type features struct {
	conditions []conditionDefinition
	actions    []actionDefinition
}

type feature func(*features)

func newFeatures(fs ...feature) features {
	result := features{}
	for _, f := range fs {
		f(&result)
	}
	return result
}

func compilerFeatures(features features) compilerOption {
	return func(_ Source, c *compilerContext) error {
		c.features = features
		return nil
	}
}

var defaultFeatures = newFeatures(
	conditions(
		andCondition{},
		orCondition{},
		expressionCondition{},
	),
	actions(
		addTagAction{},
		attachLocalContentByIdAction{},
		attachLocalContentBySearchAction{},
		attachTmdbContentByIdAction{},
		attachTmdbContentBySearchAction{},
		deleteAction{},
		findMatchAction{},
		ifElseAction{},
		noMatchAction{},
		parseDateAction{},
		parseVideoContentAction{},
		runWorkflowAction{},
		setContentTypeAction{},
	),
)
