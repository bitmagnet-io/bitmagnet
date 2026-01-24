//go:build !wasip1

package fs

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

type fsProviderXDG struct {
	subPath string
}

func NewFSProviderXDG(subPath string) Provider {
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
	wd, err := os.Getwd()
	if err != nil {
		return afero.Afero{Fs: nopFs}
	}

	return afero.Afero{Fs: afero.NewReadOnlyFs(afero.NewBasePathFs(osFs, wd))}
}

func (fsProviderXDG) FSRoot() FS {
	return afero.Afero{Fs: afero.NewBasePathFs(osFs, "/")}
}
