package banning

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
)

type nameLengthChecker struct {
	min int
}

func (c nameLengthChecker) Check(info metainfo.Info) error {
	if len(info.BestName()) < c.min {
		return errors.New("name too short")
	}

	return nil
}
