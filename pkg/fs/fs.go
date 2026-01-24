package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

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

	ConfigProvider interface {
		FSConfig() FS
	}

	DataProvider interface {
		FSData() FS
	}

	CurrentProvider interface {
		FSCurrent() FS
	}

	RootProvider interface {
		FSRoot() FS
	}

	Provider interface {
		ConfigProvider
		DataProvider
		CurrentProvider
		RootProvider
	}
)

var (
	ErrNotExist = fs.ErrNotExist

	osFs  = afero.NewOsFs()
	nopFs = afero.NewReadOnlyFs(afero.NewMemMapFs())
)

type ProviderNop struct{}

func (ProviderNop) FSConfig() FS {
	return afero.Afero{Fs: nopFs}
}

func (ProviderNop) FSData() FS {
	return afero.Afero{Fs: nopFs}
}

func (ProviderNop) FSCurrent() FS {
	return afero.Afero{Fs: nopFs}
}

func (ProviderNop) FSRoot() FS {
	return afero.Afero{Fs: nopFs}
}
