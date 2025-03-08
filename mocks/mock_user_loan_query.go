package mocks

import (
	reflect "reflect"

	models "github.com/candrap89/loanApi/models"
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

// Add these methods to the MockUserLoanQuery struct in the mock file

// GetUserLoanByCIF mocks base method
func (m *MockUserLoanQuery) GetUserLoanByCIF(cif string) ([]models.UserLoan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserLoanByCIF", cif)
	ret0, _ := ret[0].([]models.UserLoan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserLoanByCIF indicates an expected call of GetUserLoanByCIF
func (mr *MockUserLoanQueryMockRecorder) GetUserLoanByCIF(cif interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserLoanByCIF", reflect.TypeOf((*MockUserLoanQuery)(nil).GetUserLoanByCIF), cif)
}

// GetDelinquentUsers mocks base method
func (m *MockUserLoanQuery) GetDelinquentUsers() ([]models.UserLoan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelinquentUsers")
	ret0, _ := ret[0].([]models.UserLoan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelinquentUsers indicates an expected call of GetDelinquentUsers
func (mr *MockUserLoanQueryMockRecorder) GetDelinquentUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelinquentUsers", reflect.TypeOf((*MockUserLoanQuery)(nil).GetDelinquentUsers))
}
