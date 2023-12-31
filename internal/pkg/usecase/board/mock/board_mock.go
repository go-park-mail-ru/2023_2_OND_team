// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	board "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	pin "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
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

// CheckAvailabilityFeedPinCfgOnBoard mocks base method.
func (m *MockUsecase) CheckAvailabilityFeedPinCfgOnBoard(ctx context.Context, cfg pin.FeedPinConfig, userID int, isAuth bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAvailabilityFeedPinCfgOnBoard", ctx, cfg, userID, isAuth)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckAvailabilityFeedPinCfgOnBoard indicates an expected call of CheckAvailabilityFeedPinCfgOnBoard.
func (mr *MockUsecaseMockRecorder) CheckAvailabilityFeedPinCfgOnBoard(ctx, cfg, userID, isAuth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAvailabilityFeedPinCfgOnBoard", reflect.TypeOf((*MockUsecase)(nil).CheckAvailabilityFeedPinCfgOnBoard), ctx, cfg, userID, isAuth)
}

// CreateNewBoard mocks base method.
func (m *MockUsecase) CreateNewBoard(ctx context.Context, newBoard board.Board, tagTitles []string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewBoard", ctx, newBoard, tagTitles)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewBoard indicates an expected call of CreateNewBoard.
func (mr *MockUsecaseMockRecorder) CreateNewBoard(ctx, newBoard, tagTitles interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewBoard", reflect.TypeOf((*MockUsecase)(nil).CreateNewBoard), ctx, newBoard, tagTitles)
}

// DeleteCertainBoard mocks base method.
func (m *MockUsecase) DeleteCertainBoard(ctx context.Context, boardID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCertainBoard", ctx, boardID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCertainBoard indicates an expected call of DeleteCertainBoard.
func (mr *MockUsecaseMockRecorder) DeleteCertainBoard(ctx, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCertainBoard", reflect.TypeOf((*MockUsecase)(nil).DeleteCertainBoard), ctx, boardID)
}

// DeletePinFromBoard mocks base method.
func (m *MockUsecase) DeletePinFromBoard(ctx context.Context, boardID, pinID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePinFromBoard", ctx, boardID, pinID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePinFromBoard indicates an expected call of DeletePinFromBoard.
func (mr *MockUsecaseMockRecorder) DeletePinFromBoard(ctx, boardID, pinID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePinFromBoard", reflect.TypeOf((*MockUsecase)(nil).DeletePinFromBoard), ctx, boardID, pinID)
}

// FixPinsOnBoard mocks base method.
func (m *MockUsecase) FixPinsOnBoard(ctx context.Context, boardID int, pinIds []int, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FixPinsOnBoard", ctx, boardID, pinIds, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// FixPinsOnBoard indicates an expected call of FixPinsOnBoard.
func (mr *MockUsecaseMockRecorder) FixPinsOnBoard(ctx, boardID, pinIds, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FixPinsOnBoard", reflect.TypeOf((*MockUsecase)(nil).FixPinsOnBoard), ctx, boardID, pinIds, userID)
}

// GetBoardInfoForUpdate mocks base method.
func (m *MockUsecase) GetBoardInfoForUpdate(ctx context.Context, boardID int) (board.Board, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardInfoForUpdate", ctx, boardID)
	ret0, _ := ret[0].(board.Board)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBoardInfoForUpdate indicates an expected call of GetBoardInfoForUpdate.
func (mr *MockUsecaseMockRecorder) GetBoardInfoForUpdate(ctx, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardInfoForUpdate", reflect.TypeOf((*MockUsecase)(nil).GetBoardInfoForUpdate), ctx, boardID)
}

// GetBoardsByUsername mocks base method.
func (m *MockUsecase) GetBoardsByUsername(ctx context.Context, username string) ([]board.BoardWithContent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardsByUsername", ctx, username)
	ret0, _ := ret[0].([]board.BoardWithContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoardsByUsername indicates an expected call of GetBoardsByUsername.
func (mr *MockUsecaseMockRecorder) GetBoardsByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardsByUsername", reflect.TypeOf((*MockUsecase)(nil).GetBoardsByUsername), ctx, username)
}

// GetCertainBoard mocks base method.
func (m *MockUsecase) GetCertainBoard(ctx context.Context, boardID int) (board.BoardWithContent, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertainBoard", ctx, boardID)
	ret0, _ := ret[0].(board.BoardWithContent)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCertainBoard indicates an expected call of GetCertainBoard.
func (mr *MockUsecaseMockRecorder) GetCertainBoard(ctx, boardID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertainBoard", reflect.TypeOf((*MockUsecase)(nil).GetCertainBoard), ctx, boardID)
}

// UpdateBoardInfo mocks base method.
func (m *MockUsecase) UpdateBoardInfo(ctx context.Context, updatedBoard board.Board, tagTitles []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBoardInfo", ctx, updatedBoard, tagTitles)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBoardInfo indicates an expected call of UpdateBoardInfo.
func (mr *MockUsecaseMockRecorder) UpdateBoardInfo(ctx, updatedBoard, tagTitles interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBoardInfo", reflect.TypeOf((*MockUsecase)(nil).UpdateBoardInfo), ctx, updatedBoard, tagTitles)
}
