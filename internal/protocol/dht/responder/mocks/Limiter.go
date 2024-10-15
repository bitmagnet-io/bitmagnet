// Code generated by mockery v2.46.3. DO NOT EDIT.

package responder_mocks

import (
	netip "net/netip"

	mock "github.com/stretchr/testify/mock"
)

// Limiter is an autogenerated mock type for the Limiter type
type Limiter struct {
	mock.Mock
}

type Limiter_Expecter struct {
	mock *mock.Mock
}

func (_m *Limiter) EXPECT() *Limiter_Expecter {
	return &Limiter_Expecter{mock: &_m.Mock}
}

// Allow provides a mock function with given fields: addr
func (_m *Limiter) Allow(addr netip.Addr) bool {
	ret := _m.Called(addr)

	if len(ret) == 0 {
		panic("no return value specified for Allow")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(netip.Addr) bool); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Limiter_Allow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Allow'
type Limiter_Allow_Call struct {
	*mock.Call
}

// Allow is a helper method to define mock.On call
//   - addr netip.Addr
func (_e *Limiter_Expecter) Allow(addr interface{}) *Limiter_Allow_Call {
	return &Limiter_Allow_Call{Call: _e.mock.On("Allow", addr)}
}

func (_c *Limiter_Allow_Call) Run(run func(addr netip.Addr)) *Limiter_Allow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(netip.Addr))
	})
	return _c
}

func (_c *Limiter_Allow_Call) Return(_a0 bool) *Limiter_Allow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Limiter_Allow_Call) RunAndReturn(run func(netip.Addr) bool) *Limiter_Allow_Call {
	_c.Call.Return(run)
	return _c
}

// NewLimiter creates a new instance of Limiter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLimiter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Limiter {
	mock := &Limiter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
