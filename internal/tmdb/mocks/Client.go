// Code generated by mockery v2.46.2. DO NOT EDIT.

package tmdb_mocks

import (
	context "context"

	tmdb "github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// FindByID provides a mock function with given fields: _a0, _a1
func (_m *Client) FindByID(_a0 context.Context, _a1 tmdb.FindByIDRequest) (tmdb.FindByIDResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 tmdb.FindByIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.FindByIDRequest) (tmdb.FindByIDResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.FindByIDRequest) tmdb.FindByIDResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(tmdb.FindByIDResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, tmdb.FindByIDRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_FindByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByID'
type Client_FindByID_Call struct {
	*mock.Call
}

// FindByID is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 tmdb.FindByIDRequest
func (_e *Client_Expecter) FindByID(_a0 interface{}, _a1 interface{}) *Client_FindByID_Call {
	return &Client_FindByID_Call{Call: _e.mock.On("FindByID", _a0, _a1)}
}

func (_c *Client_FindByID_Call) Run(run func(_a0 context.Context, _a1 tmdb.FindByIDRequest)) *Client_FindByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(tmdb.FindByIDRequest))
	})
	return _c
}

func (_c *Client_FindByID_Call) Return(_a0 tmdb.FindByIDResponse, _a1 error) *Client_FindByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_FindByID_Call) RunAndReturn(run func(context.Context, tmdb.FindByIDRequest) (tmdb.FindByIDResponse, error)) *Client_FindByID_Call {
	_c.Call.Return(run)
	return _c
}

// MovieDetails provides a mock function with given fields: _a0, _a1
func (_m *Client) MovieDetails(_a0 context.Context, _a1 tmdb.MovieDetailsRequest) (tmdb.MovieDetailsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for MovieDetails")
	}

	var r0 tmdb.MovieDetailsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.MovieDetailsRequest) (tmdb.MovieDetailsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.MovieDetailsRequest) tmdb.MovieDetailsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(tmdb.MovieDetailsResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, tmdb.MovieDetailsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_MovieDetails_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MovieDetails'
type Client_MovieDetails_Call struct {
	*mock.Call
}

// MovieDetails is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 tmdb.MovieDetailsRequest
func (_e *Client_Expecter) MovieDetails(_a0 interface{}, _a1 interface{}) *Client_MovieDetails_Call {
	return &Client_MovieDetails_Call{Call: _e.mock.On("MovieDetails", _a0, _a1)}
}

func (_c *Client_MovieDetails_Call) Run(run func(_a0 context.Context, _a1 tmdb.MovieDetailsRequest)) *Client_MovieDetails_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(tmdb.MovieDetailsRequest))
	})
	return _c
}

func (_c *Client_MovieDetails_Call) Return(_a0 tmdb.MovieDetailsResponse, _a1 error) *Client_MovieDetails_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_MovieDetails_Call) RunAndReturn(run func(context.Context, tmdb.MovieDetailsRequest) (tmdb.MovieDetailsResponse, error)) *Client_MovieDetails_Call {
	_c.Call.Return(run)
	return _c
}

// SearchMovie provides a mock function with given fields: _a0, _a1
func (_m *Client) SearchMovie(_a0 context.Context, _a1 tmdb.SearchMovieRequest) (tmdb.SearchMovieResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SearchMovie")
	}

	var r0 tmdb.SearchMovieResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.SearchMovieRequest) (tmdb.SearchMovieResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.SearchMovieRequest) tmdb.SearchMovieResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(tmdb.SearchMovieResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, tmdb.SearchMovieRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_SearchMovie_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchMovie'
type Client_SearchMovie_Call struct {
	*mock.Call
}

// SearchMovie is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 tmdb.SearchMovieRequest
func (_e *Client_Expecter) SearchMovie(_a0 interface{}, _a1 interface{}) *Client_SearchMovie_Call {
	return &Client_SearchMovie_Call{Call: _e.mock.On("SearchMovie", _a0, _a1)}
}

func (_c *Client_SearchMovie_Call) Run(run func(_a0 context.Context, _a1 tmdb.SearchMovieRequest)) *Client_SearchMovie_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(tmdb.SearchMovieRequest))
	})
	return _c
}

func (_c *Client_SearchMovie_Call) Return(_a0 tmdb.SearchMovieResponse, _a1 error) *Client_SearchMovie_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_SearchMovie_Call) RunAndReturn(run func(context.Context, tmdb.SearchMovieRequest) (tmdb.SearchMovieResponse, error)) *Client_SearchMovie_Call {
	_c.Call.Return(run)
	return _c
}

// SearchTv provides a mock function with given fields: _a0, _a1
func (_m *Client) SearchTv(_a0 context.Context, _a1 tmdb.SearchTvRequest) (tmdb.SearchTvResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SearchTv")
	}

	var r0 tmdb.SearchTvResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.SearchTvRequest) (tmdb.SearchTvResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.SearchTvRequest) tmdb.SearchTvResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(tmdb.SearchTvResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, tmdb.SearchTvRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_SearchTv_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchTv'
type Client_SearchTv_Call struct {
	*mock.Call
}

// SearchTv is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 tmdb.SearchTvRequest
func (_e *Client_Expecter) SearchTv(_a0 interface{}, _a1 interface{}) *Client_SearchTv_Call {
	return &Client_SearchTv_Call{Call: _e.mock.On("SearchTv", _a0, _a1)}
}

func (_c *Client_SearchTv_Call) Run(run func(_a0 context.Context, _a1 tmdb.SearchTvRequest)) *Client_SearchTv_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(tmdb.SearchTvRequest))
	})
	return _c
}

func (_c *Client_SearchTv_Call) Return(_a0 tmdb.SearchTvResponse, _a1 error) *Client_SearchTv_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_SearchTv_Call) RunAndReturn(run func(context.Context, tmdb.SearchTvRequest) (tmdb.SearchTvResponse, error)) *Client_SearchTv_Call {
	_c.Call.Return(run)
	return _c
}

// TvDetails provides a mock function with given fields: _a0, _a1
func (_m *Client) TvDetails(_a0 context.Context, _a1 tmdb.TvDetailsRequest) (tmdb.TvDetailsResponse, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for TvDetails")
	}

	var r0 tmdb.TvDetailsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.TvDetailsRequest) (tmdb.TvDetailsResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, tmdb.TvDetailsRequest) tmdb.TvDetailsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(tmdb.TvDetailsResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, tmdb.TvDetailsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_TvDetails_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TvDetails'
type Client_TvDetails_Call struct {
	*mock.Call
}

// TvDetails is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 tmdb.TvDetailsRequest
func (_e *Client_Expecter) TvDetails(_a0 interface{}, _a1 interface{}) *Client_TvDetails_Call {
	return &Client_TvDetails_Call{Call: _e.mock.On("TvDetails", _a0, _a1)}
}

func (_c *Client_TvDetails_Call) Run(run func(_a0 context.Context, _a1 tmdb.TvDetailsRequest)) *Client_TvDetails_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(tmdb.TvDetailsRequest))
	})
	return _c
}

func (_c *Client_TvDetails_Call) Return(_a0 tmdb.TvDetailsResponse, _a1 error) *Client_TvDetails_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_TvDetails_Call) RunAndReturn(run func(context.Context, tmdb.TvDetailsRequest) (tmdb.TvDetailsResponse, error)) *Client_TvDetails_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateApiKey provides a mock function with given fields: _a0
func (_m *Client) ValidateApiKey(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ValidateApiKey")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_ValidateApiKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateApiKey'
type Client_ValidateApiKey_Call struct {
	*mock.Call
}

// ValidateApiKey is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *Client_Expecter) ValidateApiKey(_a0 interface{}) *Client_ValidateApiKey_Call {
	return &Client_ValidateApiKey_Call{Call: _e.mock.On("ValidateApiKey", _a0)}
}

func (_c *Client_ValidateApiKey_Call) Run(run func(_a0 context.Context)) *Client_ValidateApiKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Client_ValidateApiKey_Call) Return(_a0 error) *Client_ValidateApiKey_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_ValidateApiKey_Call) RunAndReturn(run func(context.Context) error) *Client_ValidateApiKey_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
