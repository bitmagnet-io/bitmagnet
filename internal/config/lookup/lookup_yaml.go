package lookup

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type lookupYAML struct {
	source string
	node   yaml.Node
}

func (l lookupYAML) Lookup(path []string) (Result, bool, error) {
	currentNode := l.node

	currentIndex := 0

outer:
	for {
		if currentIndex == len(path) {
			var value any
			err := currentNode.Decode(&value)
			return Result{source: l.source, value: value}, true, err
		}

		if currentNode.Kind != yaml.MappingNode {
			return Result{source: l.source}, true, fmt.Errorf(
				"%w: %w: %w: %s: expected map at path %s",
				Err,
				ErrLookup,
				ErrInvalidStructure,
				l.source,
				strings.Join(path[:currentIndex], "."),
			)
		}

		for j := 0; j+1 < len(currentNode.Content); j += 2 {
			keyNode := currentNode.Content[j]
			if keyNode.Tag == "!!str" && keyNode.Value == path[currentIndex] {
				currentIndex++
				currentNode = *currentNode.Content[j+1]
				continue outer
			}
		}

		return Result{source: l.source}, false, nil
	}
}
