package mocks

import (
	reflect "reflect"

	"github.com/candrap89/loanApi/models"
	"github.com/golang/mock/gomock"
)

type MockTransactionQuery struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionQueryMockRecorder
}

type MockTransactionQueryMockRecorder struct {
	mock *MockTransactionQuery
}

func NewMockTransactionQuery(ctrl *gomock.Controller) *MockTransactionQuery {
	mock := &MockTransactionQuery{ctrl: ctrl}
	mock.recorder = &MockTransactionQueryMockRecorder{mock}
	return mock
}

// Add methods to MockTransactionQuery here...

func (m *MockTransactionQuery) EXPECT() *MockTransactionQueryMockRecorder {
	return m.recorder
}

func (m *MockTransactionQuery) InsertTransaction(transaction models.TransactionHistory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransaction", transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockTransactionQueryMockRecorder) InsertTransaction(transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransaction", reflect.TypeOf((*MockTransactionQuery)(nil).InsertTransaction), transaction)
}
