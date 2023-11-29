// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	message "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	gomock "github.com/golang/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// DeleteMessage mocks base method.
func (m *MockUsecase) DeleteMessage(ctx context.Context, userID, mesID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", ctx, userID, mesID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessage indicates an expected call of DeleteMessage.
func (mr *MockUsecaseMockRecorder) DeleteMessage(ctx, userID, mesID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockUsecase)(nil).DeleteMessage), ctx, userID, mesID)
}

// GetMessagesFromChat mocks base method.
func (m *MockUsecase) GetMessagesFromChat(ctx context.Context, chat message.Chat, count, lastID int) ([]message.Message, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessagesFromChat", ctx, chat, count, lastID)
	ret0, _ := ret[0].([]message.Message)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMessagesFromChat indicates an expected call of GetMessagesFromChat.
func (mr *MockUsecaseMockRecorder) GetMessagesFromChat(ctx, chat, count, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessagesFromChat", reflect.TypeOf((*MockUsecase)(nil).GetMessagesFromChat), ctx, chat, count, lastID)
}

// SendMessage mocks base method.
func (m *MockUsecase) SendMessage(ctx context.Context, mes *message.Message) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, mes)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockUsecaseMockRecorder) SendMessage(ctx, mes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockUsecase)(nil).SendMessage), ctx, mes)
}

// UpdateContentMessage mocks base method.
func (m *MockUsecase) UpdateContentMessage(ctx context.Context, userID int, mes *message.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContentMessage", ctx, userID, mes)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContentMessage indicates an expected call of UpdateContentMessage.
func (mr *MockUsecaseMockRecorder) UpdateContentMessage(ctx, userID, mes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContentMessage", reflect.TypeOf((*MockUsecase)(nil).UpdateContentMessage), ctx, userID, mes)
}
