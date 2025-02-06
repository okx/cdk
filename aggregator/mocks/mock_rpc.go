// Code generated by mockery v2.39.0. DO NOT EDIT.

package mocks

import (
	types "github.com/0xPolygon/cdk/rpc/types"
	mock "github.com/stretchr/testify/mock"
	
)

// RPCInterfaceMock is an autogenerated mock type for the RPCInterface type
type RPCInterfaceMock struct {
	mock.Mock
}

// GetBatch provides a mock function with given fields: batchNumber
func (_m *RPCInterfaceMock) GetBatch(batchNumber uint64) (*types.RPCBatch, error) {
	
	ret := _m.Called(batchNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetBatch")
	}

	var r0 *types.RPCBatch
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*types.RPCBatch, error)); ok {
		return rf(batchNumber)
	}
	if rf, ok := ret.Get(0).(func(uint64) *types.RPCBatch); ok {
		r0 = rf(batchNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.RPCBatch)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(batchNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWitness provides a mock function with given fields: batchNumber, fullWitness
func (_m *RPCInterfaceMock) GetWitness(batchNumber uint64, fullWitness bool) ([]byte, error) {
	ret := _m.Called(batchNumber, fullWitness)

	if len(ret) == 0 {
		panic("no return value specified for GetWitness")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64, bool) ([]byte, error)); ok {
		return rf(batchNumber, fullWitness)
	}
	if rf, ok := ret.Get(0).(func(uint64, bool) []byte); ok {
		r0 = rf(batchNumber, fullWitness)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, bool) error); ok {
		r1 = rf(batchNumber, fullWitness)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBatch provides a mock function
func (_m *RPCInterfaceMock) GetLatestBatch() (*types.RPCBatch, error) {
	ret := _m.Called()
	var r0 *types.RPCBatch
	var r1 error
	if rf, ok := ret.Get(0).(func() (*types.RPCBatch, error)); ok {
		return rf()
	}
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*types.RPCBatch)
	}
	r1 = ret.Error(1)
	return r0, r1
}

// NewRPCInterfaceMock creates a new instance of RPCInterfaceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRPCInterfaceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *RPCInterfaceMock {
	mock := &RPCInterfaceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
