package plugin

import (
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/spf13/afero"
)

type PermissionFS struct{}

func (PermissionFS) permission() {}

func (PermissionFS) jsonSchema() json_schema.JSONSchema {
	return json_schema.MustNew(
		json_schema.Typed(json_schema.TypeObject),
	)
}

func (p PermissionFS) build(b *instanceBuilder) {
	b.moduleConfig = b.moduleConfig.WithFS(&aferoFSAdapter{
		Fs: afero.NewBasePathFs(
			b.env.FSData(),
			filepath.Join("plugin", b.manifest.Name),
		),
	})
}

type aferoFSAdapter struct {
	mkdirOnce sync.Once
	afero.Fs
}

func (a *aferoFSAdapter) Open(name string) (fs.File, error) {
	var err error

	a.mkdirOnce.Do(func() {
		err = a.MkdirAll(".", 0755)
	})

	if err != nil {
		return nil, err
	}

	return a.Fs.Open(name)
}
