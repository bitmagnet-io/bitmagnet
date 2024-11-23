package classifier

type Source struct {
	Schema          string          `json:"$schema,omitempty" yaml:"$schema,omitempty"`
	Workflows       workflowSources `json:"workflows"`
	FlagDefinitions flagDefinitions `json:"flag_definitions"`
	Flags           Flags           `json:"flags"`
	Keywords        keywordGroups   `json:"keywords"`
	Extensions      extensionGroups `json:"extensions"`
	Plugins         pluginGroups    `json:"plugins,omitempty"`
}

func (s Source) merge(other Source) (Source, error) {
	flagDefs, err := s.FlagDefinitions.merge(other.FlagDefinitions)
	if err != nil {
		return Source{}, err
	}
	return Source{
		FlagDefinitions: flagDefs,
		Flags:           s.Flags.merge(other.Flags),
		Keywords:        s.Keywords.merge(other.Keywords),
		Extensions:      s.Extensions.merge(other.Extensions),
		Workflows:       s.Workflows.merge(other.Workflows),
		Plugins:         s.Plugins.merge(other.Plugins),
	}, nil
}

func (s Source) workflowNames() map[string]struct{} {
	result := make(map[string]struct{})
	for k := range s.Workflows {
		result[k] = struct{}{}
	}
	return result
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

type pluginSource struct {
	Url    string  `json:"url"`
	ApiKey *string `json:"apikey,omitempty"`
	Start  *string `json:"start,omitempy"`
	End    *string `json:"end,omitempy"`
	Days   *int    `json:"days,omitempty"`
}

type pluginGroups []struct {
	Source pluginSource `json:"source"`
	Flag   string       `json:"flag"`
}

func (s pluginGroups) merge(other pluginGroups) pluginGroups {
	var result pluginGroups
	for _, v := range s {
		result = append(result, v)
	}
	for _, v := range other {
		result = append(result, v)
	}
	return result
}
