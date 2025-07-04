package socket

import (
	"fmt"
	"maps"
	"net/netip"
	"slices"
	"sync"
)

type Adapter func(netip.AddrPort) Runner

var adaptersMu sync.Mutex

var adapters = make(map[string]Adapter)

func addAdapter(name string, adapter Adapter) {
	adaptersMu.Lock()
	defer adaptersMu.Unlock()

	adapters[name] = adapter
}

func AdapterNames() []string {
	names := slices.Collect(maps.Keys(adapters))

	slices.Sort(names)

	return names
}

func GetAdapter(name string) (Adapter, error) {
	adaptersMu.Lock()
	defer adaptersMu.Unlock()

	adapter, ok := adapters[name]
	if !ok {
		return nil, fmt.Errorf("%w: %w: %s", Err, ErrUnknownAdapter, name)
	}

	return adapter, nil
}
