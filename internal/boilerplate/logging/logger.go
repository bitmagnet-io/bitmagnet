package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Config zap.Config
}

type Result struct {
	fx.Out
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
}

func New(params Params) (Result, error) {
	logger, err := params.Config.Build()
	if err != nil {
		return Result{}, err
	}
	sugar := logger.Sugar()

	return Result{
		Logger: logger,
		Sugar:  sugar,
	}, nil
}
