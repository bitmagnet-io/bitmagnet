package registry

import "github.com/bitmagnet-io/bitmagnet/internal/workers/worker"

type WorkerState struct {
	worker.StateInfo
	RequiredBy worker.DependencyMap
	DependsOn  worker.DependencyMap
	Autostart  bool
}
