package mocks

import (
	reflect "reflect"

	"github.com/golang/mock/gomock"
)

// MockUserLoanQuery is a mock of UserLoanQuery interface
type MockUserLoanQuery struct {
	ctrl     *gomock.Controller
	recorder *MockUserLoanQueryMockRecorder
}

// MockUserLoanQueryMockRecorder is the mock recorder for MockUserLoanQuery
type MockUserLoanQueryMockRecorder struct {
	mock *MockUserLoanQuery
}

// NewMockUserLoanQuery creates a new mock instance
func NewMockUserLoanQuery(ctrl *gomock.Controller) *MockUserLoanQuery {
	mock := &MockUserLoanQuery{ctrl: ctrl}
	mock.recorder = &MockUserLoanQueryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserLoanQuery) EXPECT() *MockUserLoanQueryMockRecorder {
	return m.recorder
}

// UpdateUserLoanOutstanding mocks base method
func (m *MockUserLoanQuery) UpdateUserLoanOutstanding(userId int, amount float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserLoanOutstanding", userId, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserLoanOutstanding indicates an expected call of UpdateUserLoanOutstanding
func (mr *MockUserLoanQueryMockRecorder) UpdateUserLoanOutstanding(userId, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserLoanOutstanding", reflect.TypeOf((*MockUserLoanQuery)(nil).UpdateUserLoanOutstanding), userId, amount)
}

// UpdateUserTodeliquent mocks base method
func (m *MockUserLoanQuery) UpdateUserTodeliquent(status bool, userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserTodeliquent", status, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserTodeliquent indicates an expected call of UpdateUserTodeliquent
func (mr *MockUserLoanQueryMockRecorder) UpdateUserTodeliquent(status, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserTodeliquent", reflect.TypeOf((*MockUserLoanQuery)(nil).UpdateUserTodeliquent), status, userId)
}
