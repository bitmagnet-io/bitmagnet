package app

import "go.uber.org/fx"

type Builder interface {
	Build(options ...fx.Option) *fx.App
}
