package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

type (
	FS interface {
		afero.Fs
		Chtimes(name string, atime time.Time, mtime time.Time) error
		DirExists(path string) (bool, error)
		Exists(path string) (bool, error)
		IsDir(path string) (bool, error)
		IsEmpty(path string) (bool, error)
		ReadDir(dirname string) ([]os.FileInfo, error)
		ReadFile(filename string) ([]byte, error)
		Walk(root string, walkFn filepath.WalkFunc) error
		WriteFile(filename string, data []byte, perm os.FileMode) error
	}

	File = afero.File

	FSConfigProvider interface {
		FSConfig() FS
	}

	FSDataProvider interface {
		FSData() FS
	}

	FSCurrentProvider interface {
		FSCurrent() FS
	}

	FSRootProvider interface {
		FSRoot() FS
	}

	FSProvider interface {
		FSConfigProvider
		FSDataProvider
		FSCurrentProvider
		FSRootProvider
	}
)

var (
	ErrNotExist = fs.ErrNotExist

	osFs  = afero.NewOsFs()
	nopFs = afero.NewReadOnlyFs(afero.NewMemMapFs())
)

type fsProviderXDG struct {
	subPath string
}

func NewFSProviderXDG(subPath string) FSProvider {
	return fsProviderXDG{
		subPath: subPath,
	}
}

func (p fsProviderXDG) FSConfig() FS {
	return afero.Afero{Fs: afero.NewReadOnlyFs(afero.NewBasePathFs(osFs, filepath.Join(xdg.ConfigHome, p.subPath)))}
}

func (p fsProviderXDG) FSData() FS {
	return afero.Afero{Fs: afero.NewBasePathFs(osFs, filepath.Join(xdg.DataHome, p.subPath))}
}

func (fsProviderXDG) FSCurrent() FS {
	return afero.Afero{Fs: afero.NewReadOnlyFs(afero.NewBasePathFs(osFs, "."))}
}

func (fsProviderXDG) FSRoot() FS {
	return afero.Afero{Fs: afero.NewBasePathFs(osFs, "/")}
}

type FSProviderNop struct{}

func (FSProviderNop) FSConfig() FS {
	return afero.Afero{Fs: nopFs}
}

func (FSProviderNop) FSData() FS {
	return afero.Afero{Fs: nopFs}
}

func (FSProviderNop) FSCurrent() FS {
	return afero.Afero{Fs: nopFs}
}

func (FSProviderNop) FSRoot() FS {
	return afero.Afero{Fs: nopFs}
}
