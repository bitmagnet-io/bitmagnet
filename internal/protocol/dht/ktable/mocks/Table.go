// Code generated by mockery v2.35.2. DO NOT EDIT.

package ktable_mocks

import (
	ktable "github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	btree "github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable/btree"

	mock "github.com/stretchr/testify/mock"

	netip "net/netip"

	protocol "github.com/bitmagnet-io/bitmagnet/internal/protocol"

	time "time"
)

// Table is an autogenerated mock type for the Table type
type Table struct {
	mock.Mock
}

type Table_Expecter struct {
	mock *mock.Mock
}

func (_m *Table) EXPECT() *Table_Expecter {
	return &Table_Expecter{mock: &_m.Mock}
}

// BatchCommand provides a mock function with given fields: commands
func (_m *Table) BatchCommand(commands ...ktable.Command) {
	_va := make([]interface{}, len(commands))
	for _i := range commands {
		_va[_i] = commands[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Table_BatchCommand_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BatchCommand'
type Table_BatchCommand_Call struct {
	*mock.Call
}

// BatchCommand is a helper method to define mock.On call
//   - commands ...ktable.Command
func (_e *Table_Expecter) BatchCommand(commands ...interface{}) *Table_BatchCommand_Call {
	return &Table_BatchCommand_Call{Call: _e.mock.On("BatchCommand",
		append([]interface{}{}, commands...)...)}
}

func (_c *Table_BatchCommand_Call) Run(run func(commands ...ktable.Command)) *Table_BatchCommand_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]ktable.Command, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(ktable.Command)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Table_BatchCommand_Call) Return() *Table_BatchCommand_Call {
	_c.Call.Return()
	return _c
}

func (_c *Table_BatchCommand_Call) RunAndReturn(run func(...ktable.Command)) *Table_BatchCommand_Call {
	_c.Call.Return(run)
	return _c
}

// DropNode provides a mock function with given fields: id, reason
func (_m *Table) DropNode(id protocol.ID, reason error) bool {
	ret := _m.Called(id, reason)

	var r0 bool
	if rf, ok := ret.Get(0).(func(protocol.ID, error) bool); ok {
		r0 = rf(id, reason)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Table_DropNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DropNode'
type Table_DropNode_Call struct {
	*mock.Call
}

// DropNode is a helper method to define mock.On call
//   - id protocol.ID
//   - reason error
func (_e *Table_Expecter) DropNode(id interface{}, reason interface{}) *Table_DropNode_Call {
	return &Table_DropNode_Call{Call: _e.mock.On("DropNode", id, reason)}
}

func (_c *Table_DropNode_Call) Run(run func(id protocol.ID, reason error)) *Table_DropNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(protocol.ID), args[1].(error))
	})
	return _c
}

func (_c *Table_DropNode_Call) Return(_a0 bool) *Table_DropNode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_DropNode_Call) RunAndReturn(run func(protocol.ID, error) bool) *Table_DropNode_Call {
	_c.Call.Return(run)
	return _c
}

// FilterKnownAddrs provides a mock function with given fields: addrs
func (_m *Table) FilterKnownAddrs(addrs []netip.Addr) []netip.Addr {
	ret := _m.Called(addrs)

	var r0 []netip.Addr
	if rf, ok := ret.Get(0).(func([]netip.Addr) []netip.Addr); ok {
		r0 = rf(addrs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]netip.Addr)
		}
	}

	return r0
}

// Table_FilterKnownAddrs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FilterKnownAddrs'
type Table_FilterKnownAddrs_Call struct {
	*mock.Call
}

// FilterKnownAddrs is a helper method to define mock.On call
//   - addrs []netip.Addr
func (_e *Table_Expecter) FilterKnownAddrs(addrs interface{}) *Table_FilterKnownAddrs_Call {
	return &Table_FilterKnownAddrs_Call{Call: _e.mock.On("FilterKnownAddrs", addrs)}
}

func (_c *Table_FilterKnownAddrs_Call) Run(run func(addrs []netip.Addr)) *Table_FilterKnownAddrs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]netip.Addr))
	})
	return _c
}

func (_c *Table_FilterKnownAddrs_Call) Return(_a0 []netip.Addr) *Table_FilterKnownAddrs_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_FilterKnownAddrs_Call) RunAndReturn(run func([]netip.Addr) []netip.Addr) *Table_FilterKnownAddrs_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateNodeID provides a mock function with given fields:
func (_m *Table) GenerateNodeID() protocol.ID {
	ret := _m.Called()

	var r0 protocol.ID
	if rf, ok := ret.Get(0).(func() protocol.ID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(protocol.ID)
		}
	}

	return r0
}

// Table_GenerateNodeID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateNodeID'
type Table_GenerateNodeID_Call struct {
	*mock.Call
}

// GenerateNodeID is a helper method to define mock.On call
func (_e *Table_Expecter) GenerateNodeID() *Table_GenerateNodeID_Call {
	return &Table_GenerateNodeID_Call{Call: _e.mock.On("GenerateNodeID")}
}

func (_c *Table_GenerateNodeID_Call) Run(run func()) *Table_GenerateNodeID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Table_GenerateNodeID_Call) Return(_a0 protocol.ID) *Table_GenerateNodeID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_GenerateNodeID_Call) RunAndReturn(run func() protocol.ID) *Table_GenerateNodeID_Call {
	_c.Call.Return(run)
	return _c
}

// GetClosestNodes provides a mock function with given fields: id
func (_m *Table) GetClosestNodes(id protocol.ID) []ktable.Node {
	ret := _m.Called(id)

	var r0 []ktable.Node
	if rf, ok := ret.Get(0).(func(protocol.ID) []ktable.Node); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ktable.Node)
		}
	}

	return r0
}

// Table_GetClosestNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetClosestNodes'
type Table_GetClosestNodes_Call struct {
	*mock.Call
}

// GetClosestNodes is a helper method to define mock.On call
//   - id protocol.ID
func (_e *Table_Expecter) GetClosestNodes(id interface{}) *Table_GetClosestNodes_Call {
	return &Table_GetClosestNodes_Call{Call: _e.mock.On("GetClosestNodes", id)}
}

func (_c *Table_GetClosestNodes_Call) Run(run func(id protocol.ID)) *Table_GetClosestNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(protocol.ID))
	})
	return _c
}

func (_c *Table_GetClosestNodes_Call) Return(_a0 []ktable.Node) *Table_GetClosestNodes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_GetClosestNodes_Call) RunAndReturn(run func(protocol.ID) []ktable.Node) *Table_GetClosestNodes_Call {
	_c.Call.Return(run)
	return _c
}

// GetHashOrClosestNodes provides a mock function with given fields: id
func (_m *Table) GetHashOrClosestNodes(id protocol.ID) ktable.GetHashOrClosestNodesResult {
	ret := _m.Called(id)

	var r0 ktable.GetHashOrClosestNodesResult
	if rf, ok := ret.Get(0).(func(protocol.ID) ktable.GetHashOrClosestNodesResult); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(ktable.GetHashOrClosestNodesResult)
	}

	return r0
}

// Table_GetHashOrClosestNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHashOrClosestNodes'
type Table_GetHashOrClosestNodes_Call struct {
	*mock.Call
}

// GetHashOrClosestNodes is a helper method to define mock.On call
//   - id protocol.ID
func (_e *Table_Expecter) GetHashOrClosestNodes(id interface{}) *Table_GetHashOrClosestNodes_Call {
	return &Table_GetHashOrClosestNodes_Call{Call: _e.mock.On("GetHashOrClosestNodes", id)}
}

func (_c *Table_GetHashOrClosestNodes_Call) Run(run func(id protocol.ID)) *Table_GetHashOrClosestNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(protocol.ID))
	})
	return _c
}

func (_c *Table_GetHashOrClosestNodes_Call) Return(_a0 ktable.GetHashOrClosestNodesResult) *Table_GetHashOrClosestNodes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_GetHashOrClosestNodes_Call) RunAndReturn(run func(protocol.ID) ktable.GetHashOrClosestNodesResult) *Table_GetHashOrClosestNodes_Call {
	_c.Call.Return(run)
	return _c
}

// GetNodesForSampleInfoHashes provides a mock function with given fields: n
func (_m *Table) GetNodesForSampleInfoHashes(n int) []ktable.Node {
	ret := _m.Called(n)

	var r0 []ktable.Node
	if rf, ok := ret.Get(0).(func(int) []ktable.Node); ok {
		r0 = rf(n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ktable.Node)
		}
	}

	return r0
}

// Table_GetNodesForSampleInfoHashes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNodesForSampleInfoHashes'
type Table_GetNodesForSampleInfoHashes_Call struct {
	*mock.Call
}

// GetNodesForSampleInfoHashes is a helper method to define mock.On call
//   - n int
func (_e *Table_Expecter) GetNodesForSampleInfoHashes(n interface{}) *Table_GetNodesForSampleInfoHashes_Call {
	return &Table_GetNodesForSampleInfoHashes_Call{Call: _e.mock.On("GetNodesForSampleInfoHashes", n)}
}

func (_c *Table_GetNodesForSampleInfoHashes_Call) Run(run func(n int)) *Table_GetNodesForSampleInfoHashes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *Table_GetNodesForSampleInfoHashes_Call) Return(_a0 []ktable.Node) *Table_GetNodesForSampleInfoHashes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_GetNodesForSampleInfoHashes_Call) RunAndReturn(run func(int) []ktable.Node) *Table_GetNodesForSampleInfoHashes_Call {
	_c.Call.Return(run)
	return _c
}

// GetOldestNodes provides a mock function with given fields: cutoff, n
func (_m *Table) GetOldestNodes(cutoff time.Time, n int) []ktable.Node {
	ret := _m.Called(cutoff, n)

	var r0 []ktable.Node
	if rf, ok := ret.Get(0).(func(time.Time, int) []ktable.Node); ok {
		r0 = rf(cutoff, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ktable.Node)
		}
	}

	return r0
}

// Table_GetOldestNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOldestNodes'
type Table_GetOldestNodes_Call struct {
	*mock.Call
}

// GetOldestNodes is a helper method to define mock.On call
//   - cutoff time.Time
//   - n int
func (_e *Table_Expecter) GetOldestNodes(cutoff interface{}, n interface{}) *Table_GetOldestNodes_Call {
	return &Table_GetOldestNodes_Call{Call: _e.mock.On("GetOldestNodes", cutoff, n)}
}

func (_c *Table_GetOldestNodes_Call) Run(run func(cutoff time.Time, n int)) *Table_GetOldestNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Time), args[1].(int))
	})
	return _c
}

func (_c *Table_GetOldestNodes_Call) Return(_a0 []ktable.Node) *Table_GetOldestNodes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_GetOldestNodes_Call) RunAndReturn(run func(time.Time, int) []ktable.Node) *Table_GetOldestNodes_Call {
	_c.Call.Return(run)
	return _c
}

// Origin provides a mock function with given fields:
func (_m *Table) Origin() protocol.ID {
	ret := _m.Called()

	var r0 protocol.ID
	if rf, ok := ret.Get(0).(func() protocol.ID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(protocol.ID)
		}
	}

	return r0
}

// Table_Origin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Origin'
type Table_Origin_Call struct {
	*mock.Call
}

// Origin is a helper method to define mock.On call
func (_e *Table_Expecter) Origin() *Table_Origin_Call {
	return &Table_Origin_Call{Call: _e.mock.On("Origin")}
}

func (_c *Table_Origin_Call) Run(run func()) *Table_Origin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Table_Origin_Call) Return(_a0 protocol.ID) *Table_Origin_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_Origin_Call) RunAndReturn(run func() protocol.ID) *Table_Origin_Call {
	_c.Call.Return(run)
	return _c
}

// PutHash provides a mock function with given fields: id, peers, options
func (_m *Table) PutHash(id protocol.ID, peers []ktable.HashPeer, options ...ktable.HashOption) btree.PutResult {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, id, peers)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 btree.PutResult
	if rf, ok := ret.Get(0).(func(protocol.ID, []ktable.HashPeer, ...ktable.HashOption) btree.PutResult); ok {
		r0 = rf(id, peers, options...)
	} else {
		r0 = ret.Get(0).(btree.PutResult)
	}

	return r0
}

// Table_PutHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutHash'
type Table_PutHash_Call struct {
	*mock.Call
}

// PutHash is a helper method to define mock.On call
//   - id protocol.ID
//   - peers []ktable.HashPeer
//   - options ...ktable.HashOption
func (_e *Table_Expecter) PutHash(id interface{}, peers interface{}, options ...interface{}) *Table_PutHash_Call {
	return &Table_PutHash_Call{Call: _e.mock.On("PutHash",
		append([]interface{}{id, peers}, options...)...)}
}

func (_c *Table_PutHash_Call) Run(run func(id protocol.ID, peers []ktable.HashPeer, options ...ktable.HashOption)) *Table_PutHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]ktable.HashOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(ktable.HashOption)
			}
		}
		run(args[0].(protocol.ID), args[1].([]ktable.HashPeer), variadicArgs...)
	})
	return _c
}

func (_c *Table_PutHash_Call) Return(_a0 btree.PutResult) *Table_PutHash_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_PutHash_Call) RunAndReturn(run func(protocol.ID, []ktable.HashPeer, ...ktable.HashOption) btree.PutResult) *Table_PutHash_Call {
	_c.Call.Return(run)
	return _c
}

// PutNode provides a mock function with given fields: id, addr, options
func (_m *Table) PutNode(id protocol.ID, addr netip.AddrPort, options ...ktable.NodeOption) btree.PutResult {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, id, addr)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 btree.PutResult
	if rf, ok := ret.Get(0).(func(protocol.ID, netip.AddrPort, ...ktable.NodeOption) btree.PutResult); ok {
		r0 = rf(id, addr, options...)
	} else {
		r0 = ret.Get(0).(btree.PutResult)
	}

	return r0
}

// Table_PutNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutNode'
type Table_PutNode_Call struct {
	*mock.Call
}

// PutNode is a helper method to define mock.On call
//   - id protocol.ID
//   - addr netip.AddrPort
//   - options ...ktable.NodeOption
func (_e *Table_Expecter) PutNode(id interface{}, addr interface{}, options ...interface{}) *Table_PutNode_Call {
	return &Table_PutNode_Call{Call: _e.mock.On("PutNode",
		append([]interface{}{id, addr}, options...)...)}
}

func (_c *Table_PutNode_Call) Run(run func(id protocol.ID, addr netip.AddrPort, options ...ktable.NodeOption)) *Table_PutNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]ktable.NodeOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(ktable.NodeOption)
			}
		}
		run(args[0].(protocol.ID), args[1].(netip.AddrPort), variadicArgs...)
	})
	return _c
}

func (_c *Table_PutNode_Call) Return(_a0 btree.PutResult) *Table_PutNode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_PutNode_Call) RunAndReturn(run func(protocol.ID, netip.AddrPort, ...ktable.NodeOption) btree.PutResult) *Table_PutNode_Call {
	_c.Call.Return(run)
	return _c
}

// SampleHashesAndNodes provides a mock function with given fields:
func (_m *Table) SampleHashesAndNodes() ktable.SampleHashesAndNodesResult {
	ret := _m.Called()

	var r0 ktable.SampleHashesAndNodesResult
	if rf, ok := ret.Get(0).(func() ktable.SampleHashesAndNodesResult); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ktable.SampleHashesAndNodesResult)
	}

	return r0
}

// Table_SampleHashesAndNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SampleHashesAndNodes'
type Table_SampleHashesAndNodes_Call struct {
	*mock.Call
}

// SampleHashesAndNodes is a helper method to define mock.On call
func (_e *Table_Expecter) SampleHashesAndNodes() *Table_SampleHashesAndNodes_Call {
	return &Table_SampleHashesAndNodes_Call{Call: _e.mock.On("SampleHashesAndNodes")}
}

func (_c *Table_SampleHashesAndNodes_Call) Run(run func()) *Table_SampleHashesAndNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Table_SampleHashesAndNodes_Call) Return(_a0 ktable.SampleHashesAndNodesResult) *Table_SampleHashesAndNodes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_SampleHashesAndNodes_Call) RunAndReturn(run func() ktable.SampleHashesAndNodesResult) *Table_SampleHashesAndNodes_Call {
	_c.Call.Return(run)
	return _c
}

// Stats provides a mock function with given fields:
func (_m *Table) Stats() ktable.Stats {
	ret := _m.Called()

	var r0 ktable.Stats
	if rf, ok := ret.Get(0).(func() ktable.Stats); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ktable.Stats)
	}

	return r0
}

// Table_Stats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stats'
type Table_Stats_Call struct {
	*mock.Call
}

// Stats is a helper method to define mock.On call
func (_e *Table_Expecter) Stats() *Table_Stats_Call {
	return &Table_Stats_Call{Call: _e.mock.On("Stats")}
}

func (_c *Table_Stats_Call) Run(run func()) *Table_Stats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Table_Stats_Call) Return(_a0 ktable.Stats) *Table_Stats_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Table_Stats_Call) RunAndReturn(run func() ktable.Stats) *Table_Stats_Call {
	_c.Call.Return(run)
	return _c
}

// NewTable creates a new instance of Table. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTable(t interface {
	mock.TestingT
	Cleanup(func())
}) *Table {
	mock := &Table{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
