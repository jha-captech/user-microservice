// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// MocksLogger is an autogenerated mock type for the sLogger type
type MocksLogger struct {
	mock.Mock
}

type MocksLogger_Expecter struct {
	mock *mock.Mock
}

func (_m *MocksLogger) EXPECT() *MocksLogger_Expecter {
	return &MocksLogger_Expecter{mock: &_m.Mock}
}

// Debug provides a mock function with given fields: msg, args
func (_m *MocksLogger) Debug(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MocksLogger_Debug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debug'
type MocksLogger_Debug_Call struct {
	*mock.Call
}

// Debug is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *MocksLogger_Expecter) Debug(msg interface{}, args ...interface{}) *MocksLogger_Debug_Call {
	return &MocksLogger_Debug_Call{Call: _e.mock.On("Debug",
		append([]interface{}{msg}, args...)...)}
}

func (_c *MocksLogger_Debug_Call) Run(run func(msg string, args ...interface{})) *MocksLogger_Debug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MocksLogger_Debug_Call) Return() *MocksLogger_Debug_Call {
	_c.Call.Return()
	return _c
}

func (_c *MocksLogger_Debug_Call) RunAndReturn(run func(string, ...interface{})) *MocksLogger_Debug_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields: msg, args
func (_m *MocksLogger) Error(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MocksLogger_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type MocksLogger_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *MocksLogger_Expecter) Error(msg interface{}, args ...interface{}) *MocksLogger_Error_Call {
	return &MocksLogger_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{msg}, args...)...)}
}

func (_c *MocksLogger_Error_Call) Run(run func(msg string, args ...interface{})) *MocksLogger_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MocksLogger_Error_Call) Return() *MocksLogger_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *MocksLogger_Error_Call) RunAndReturn(run func(string, ...interface{})) *MocksLogger_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Info provides a mock function with given fields: msg, args
func (_m *MocksLogger) Info(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MocksLogger_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type MocksLogger_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *MocksLogger_Expecter) Info(msg interface{}, args ...interface{}) *MocksLogger_Info_Call {
	return &MocksLogger_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{msg}, args...)...)}
}

func (_c *MocksLogger_Info_Call) Run(run func(msg string, args ...interface{})) *MocksLogger_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MocksLogger_Info_Call) Return() *MocksLogger_Info_Call {
	_c.Call.Return()
	return _c
}

func (_c *MocksLogger_Info_Call) RunAndReturn(run func(string, ...interface{})) *MocksLogger_Info_Call {
	_c.Call.Return(run)
	return _c
}

// Warn provides a mock function with given fields: msg, args
func (_m *MocksLogger) Warn(msg string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MocksLogger_Warn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warn'
type MocksLogger_Warn_Call struct {
	*mock.Call
}

// Warn is a helper method to define mock.On call
//   - msg string
//   - args ...interface{}
func (_e *MocksLogger_Expecter) Warn(msg interface{}, args ...interface{}) *MocksLogger_Warn_Call {
	return &MocksLogger_Warn_Call{Call: _e.mock.On("Warn",
		append([]interface{}{msg}, args...)...)}
}

func (_c *MocksLogger_Warn_Call) Run(run func(msg string, args ...interface{})) *MocksLogger_Warn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MocksLogger_Warn_Call) Return() *MocksLogger_Warn_Call {
	_c.Call.Return()
	return _c
}

func (_c *MocksLogger_Warn_Call) RunAndReturn(run func(string, ...interface{})) *MocksLogger_Warn_Call {
	_c.Call.Return(run)
	return _c
}

// NewMocksLogger creates a new instance of MocksLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMocksLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *MocksLogger {
	mock := &MocksLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
