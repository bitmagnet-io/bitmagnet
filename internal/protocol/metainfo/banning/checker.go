package banning

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func New() Checker {
	return combinedChecker{
		checkers: []Checker{
			nameLengthChecker{min: 8},
			sizeChecker{min: 1024},
			utf8Checker{},
		},
	}
}

type Checker interface {
	Check(metainfo.Info) error
}

type combinedChecker struct {
	checkers []Checker
}

func (c combinedChecker) Check(info metainfo.Info) error {
	return errors.Join(slice.Map(c.checkers, func(checker Checker) error {
		return checker.Check(info)
	})...)
}
