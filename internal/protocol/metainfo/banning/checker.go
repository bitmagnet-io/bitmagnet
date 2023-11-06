package banning

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Checkers []Checker `group:"metainfo_banning_checkers"`
}

type Result struct {
	fx.Out
	Checker Checker `name:"metainfo_banning_checker"`
}

func New(p Params) Result {
	checkers := p.Checkers
	checkers = append(
		checkers,
		csamChecker{},
		nameLengthChecker{min: 8},
		sizeChecker{min: 1024},
		utf8Checker{},
	)
	return Result{
		Checker: combinedChecker{
			checkers: checkers,
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
	var errs []error
	for _, checker := range c.checkers {
		errs = append(errs, checker.Check(info))
	}
	return errors.Join(errs...)
}
