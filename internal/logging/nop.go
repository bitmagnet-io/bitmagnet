package logging

type nop struct{}

func (nop) Debugf(string, ...any) {}
func (nop) Infof(string, ...any)  {}
func (nop) Warnf(string, ...any)  {}
func (nop) Errorf(string, ...any) {}

var NopLogger = nop{}
