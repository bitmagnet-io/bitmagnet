package logging

type nop struct{}

func (nop) Debugf(string, ...any) {}
func (nop) Debugw(string, ...any) {}
func (nop) Infof(string, ...any)  {}
func (nop) Infow(string, ...any)  {}
func (nop) Warnf(string, ...any)  {}
func (nop) Warnw(string, ...any)  {}
func (nop) Errorf(string, ...any) {}
func (nop) Errorw(string, ...any) {}

var NopLogger = nop{}
