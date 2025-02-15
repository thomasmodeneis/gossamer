// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ChainSafe/gossamer/dot/network (interfaces: TransactionHandler)

// Package network is a generated GoMock package.
package network

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

// MockTransactionHandler is a mock of TransactionHandler interface.
type MockTransactionHandler struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionHandlerMockRecorder
}

// MockTransactionHandlerMockRecorder is the mock recorder for MockTransactionHandler.
type MockTransactionHandlerMockRecorder struct {
	mock *MockTransactionHandler
}

// NewMockTransactionHandler creates a new mock instance.
func NewMockTransactionHandler(ctrl *gomock.Controller) *MockTransactionHandler {
	mock := &MockTransactionHandler{ctrl: ctrl}
	mock.recorder = &MockTransactionHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionHandler) EXPECT() *MockTransactionHandlerMockRecorder {
	return m.recorder
}

// HandleTransactionMessage mocks base method.
func (m *MockTransactionHandler) HandleTransactionMessage(arg0 peer.ID, arg1 *TransactionMessage) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleTransactionMessage", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleTransactionMessage indicates an expected call of HandleTransactionMessage.
func (mr *MockTransactionHandlerMockRecorder) HandleTransactionMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleTransactionMessage", reflect.TypeOf((*MockTransactionHandler)(nil).HandleTransactionMessage), arg0, arg1)
}

// TransactionsCount mocks base method.
func (m *MockTransactionHandler) TransactionsCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionsCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// TransactionsCount indicates an expected call of TransactionsCount.
func (mr *MockTransactionHandlerMockRecorder) TransactionsCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionsCount", reflect.TypeOf((*MockTransactionHandler)(nil).TransactionsCount))
}
