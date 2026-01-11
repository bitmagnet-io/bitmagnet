package lookup

import (
	"strings"

	"github.com/bitmagnet-io/bitmagnet/pkg/env"
)

func LookupEnv(source string, varsLookup env.VarsLookup) Lookup {
	return lookupEnv{
		source:     source,
		VarsLookup: varsLookup,
	}
}

type lookupEnv struct {
	source string
	env.VarsLookup
}

func (l lookupEnv) Lookup(path []string) (Result, bool, error) {
	envKey := strings.ToUpper(strings.ReplaceAll(strings.Join(path, "_"), ".", "_"))
	envValue, found := l.LookupVar(envKey)

	return Result{source: l.source, value: envValue}, found, nil
}
