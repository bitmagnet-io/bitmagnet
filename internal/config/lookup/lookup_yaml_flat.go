package lookup

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type lookupYAMLFlat struct {
	source  string
	nodeMap map[string]yaml.Node
}

func (l lookupYAMLFlat) Lookup(path []string) (Result, bool, error) {
	node, ok := l.nodeMap[strings.Join(path, ".")]

	return Result{
		source: l.source,
		value:  node,
	}, ok, nil
}
