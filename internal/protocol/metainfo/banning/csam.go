package banning

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
)

var csamWords = []string{
	"pedo",
	"pedofile",
	"pedofilia",
	"preteen",
	"pthc",
	"ptsc",
	"lsbar",
	"lsm",
	"underage",
	"hebefilia",
	"opva",
	"child porn",
}

var csamRegex = regex.NewRegexFromNames(csamWords...)

type csamChecker struct{}

func (c csamChecker) Check(info metainfo.Info) error {
	if csamRegex.MatchString(info.BestName()) {
		return errors.New("torrent appears to contain CSAM")
	}
	return nil
}
