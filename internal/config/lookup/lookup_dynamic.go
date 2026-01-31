package lookup

import "strings"

func Dynamic(source string, values map[string]any) Lookup {
	return lookupDynamic{
		source: source,
		values: values,
	}
}

type lookupDynamic struct {
	source string
	values map[string]any
}

func (i lookupDynamic) Lookup(path []string) (Result, bool, error) {
	rawValue, found := i.values[strings.Join(path, ".")]

	return Result{source: i.source, value: rawValue}, found, nil
}
