package registry

import (
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
)

type WorkerState struct {
	worker.StateInfo
	RequiredBy ref.Set
	DependsOn  ref.Set
	Autostart  bool
}
