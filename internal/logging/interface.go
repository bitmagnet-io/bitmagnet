package logging

import "go.uber.org/zap"

type Logger interface {
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
}

var _ Logger = (*zap.SugaredLogger)(nil)
