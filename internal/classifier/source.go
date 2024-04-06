package classifier

type WorkflowSource struct {
	FlagTypes  flagTypes
	Flags      flags
	Keywords   keywordGroups
	Extensions extensionGroups
	Workflows  workflowSources
}

func (s WorkflowSource) merge(other WorkflowSource) (WorkflowSource, error) {
	flagDefs, err := s.FlagTypes.merge(other.FlagTypes)
	if err != nil {
		return WorkflowSource{}, err
	}
	return WorkflowSource{
		FlagTypes:  flagDefs,
		Flags:      s.Flags.merge(other.Flags),
		Keywords:   s.Keywords.merge(other.Keywords),
		Extensions: s.Extensions.merge(other.Extensions),
		Workflows:  s.Workflows.merge(other.Workflows),
	}, nil
}

type keywordGroups map[string][]string

func (g keywordGroups) merge(other keywordGroups) keywordGroups {
	result := make(keywordGroups)
	for k, v := range g {
		if _, ok := other[k]; ok {
			result[k] = append(v, other[k]...)
		} else {
			result[k] = v
		}
	}
	for k, v := range other {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}
	return result
}

type extensionGroups map[string][]string

func (g extensionGroups) merge(other extensionGroups) extensionGroups {
	result := make(extensionGroups)
	for k, v := range g {
		if _, ok := other[k]; ok {
			result[k] = append(v, other[k]...)
		} else {
			result[k] = v
		}
	}
	for k, v := range other {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}
	return result
}

type workflowSources map[string]any

func (s workflowSources) merge(other workflowSources) workflowSources {
	result := make(workflowSources)
	for k, v := range s {
		result[k] = v
	}
	for k, v := range other {
		result[k] = v
	}
	return result
}