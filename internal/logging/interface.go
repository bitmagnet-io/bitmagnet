package logging

import "go.uber.org/zap"

type Logger interface {
	Debugf(string, ...any)
	Debugw(string, ...any)
	Infof(string, ...any)
	Infow(string, ...any)
	Warnf(string, ...any)
	Warnw(string, ...any)
	Errorf(string, ...any)
	Errorw(string, ...any)
}

var _ Logger = (*zap.SugaredLogger)(nil)
