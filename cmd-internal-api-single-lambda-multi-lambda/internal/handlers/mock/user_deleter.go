// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/jha-captech/user-microservice/internal/models"
)

// MockUserDeleter is an autogenerated mock type for the userDeleter type
type MockUserDeleter struct {
	mock.Mock
}

type MockUserDeleter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserDeleter) EXPECT() *MockUserDeleter_Expecter {
	return &MockUserDeleter_Expecter{mock: &_m.Mock}
}

// DeleteUser provides a mock function with given fields: ctx, ID
func (_m *MockUserDeleter) DeleteUser(ctx context.Context, ID int) error {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUserDeleter_DeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUser'
type MockUserDeleter_DeleteUser_Call struct {
	*mock.Call
}

// DeleteUser is a helper method to define mock.On call
//   - ctx context.Context
//   - ID int
func (_e *MockUserDeleter_Expecter) DeleteUser(ctx interface{}, ID interface{}) *MockUserDeleter_DeleteUser_Call {
	return &MockUserDeleter_DeleteUser_Call{Call: _e.mock.On("DeleteUser", ctx, ID)}
}

func (_c *MockUserDeleter_DeleteUser_Call) Run(run func(ctx context.Context, ID int)) *MockUserDeleter_DeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockUserDeleter_DeleteUser_Call) Return(_a0 error) *MockUserDeleter_DeleteUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUserDeleter_DeleteUser_Call) RunAndReturn(run func(context.Context, int) error) *MockUserDeleter_DeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// FetchUser provides a mock function with given fields: ctx, ID
func (_m *MockUserDeleter) FetchUser(ctx context.Context, ID int) (models.User, error) {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for FetchUser")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (models.User, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) models.User); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserDeleter_FetchUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchUser'
type MockUserDeleter_FetchUser_Call struct {
	*mock.Call
}

// FetchUser is a helper method to define mock.On call
//   - ctx context.Context
//   - ID int
func (_e *MockUserDeleter_Expecter) FetchUser(ctx interface{}, ID interface{}) *MockUserDeleter_FetchUser_Call {
	return &MockUserDeleter_FetchUser_Call{Call: _e.mock.On("FetchUser", ctx, ID)}
}

func (_c *MockUserDeleter_FetchUser_Call) Run(run func(ctx context.Context, ID int)) *MockUserDeleter_FetchUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockUserDeleter_FetchUser_Call) Return(_a0 models.User, _a1 error) *MockUserDeleter_FetchUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserDeleter_FetchUser_Call) RunAndReturn(run func(context.Context, int) (models.User, error)) *MockUserDeleter_FetchUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserDeleter creates a new instance of MockUserDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserDeleter {
	mock := &MockUserDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
