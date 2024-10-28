package classifier

import (
	_ "embed"
)

//go:embed classifier.core.yml
var classifierCoreYaml []byte
