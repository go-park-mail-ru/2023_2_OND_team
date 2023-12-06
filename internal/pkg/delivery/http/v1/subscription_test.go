package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	auth "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	repo_sub "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/subscription"
	mock_sub "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/subscription/mock"
	repo_user "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	mock_user "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user/mock"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/require"
)

type (
	subscribe            func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int)
	checkUserExistence   func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int)
	getUserSubscriptions func(mockRepo *mock_sub.MockRepository, ctx context.Context, userID, count, lastID, currUserID int)
)

func getHttpRequest(method string, url string, body io.Reader, isAuth bool, currUserID int) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	if isAuth {
		req = req.WithContext(context.WithValue(req.Context(), auth.KeyCurrentUserID, currUserID))
	}
	return req
}

func TestGetOpts(t *testing.T) {

	baseUrl := "http://test.com/api/v1/subscription"

	cases := []struct {
		name    string
		inReq   *http.Request
		expOpts *user.SubscriptionOpts
		wantErr bool
		expErr  error
	}{
		{
			name:  "valid request from authenticated user",
			inReq: getHttpRequest(http.MethodGet, baseUrl+"?userID=2&count=20&lastID=23&view=subscribers", nil, true, 13),
			expOpts: &user.SubscriptionOpts{
				UserID: 2,
				Count:  20,
				LastID: 23,
				Filter: "subscribers",
			},
		},
		{
			name:  "valid request from authenticated user, no userID",
			inReq: getHttpRequest(http.MethodGet, baseUrl+"?userID=&count=20&lastID=23&view=subscribers", nil, true, 13),
			expOpts: &user.SubscriptionOpts{
				UserID: 13,
				Count:  20,
				LastID: 23,
				Filter: "subscribers",
			},
		},
		{
			name:  "valid request from unauthenticated user, no userID",
			inReq: getHttpRequest(http.MethodGet, baseUrl+"?userID=&count=20&lastID=23&view=subscribers", nil, false, -1),
			expOpts: &user.SubscriptionOpts{
				UserID: 0,
				Count:  20,
				LastID: 23,
				Filter: "subscribers",
			},
		},
		{
			name:  "valid request with default params",
			inReq: getHttpRequest(http.MethodGet, baseUrl+"?userID=1&count=&lastID=&view=subscribers", nil, false, 23),
			expOpts: &user.SubscriptionOpts{
				UserID: 1,
				Count:  20,
				LastID: 1 << 30,
				Filter: "subscribers",
			},
		},
		{
			name:    "incorrect query",
			inReq:   getHttpRequest(http.MethodGet, baseUrl+"?userID=123wa1&count=2~&lastID=33~&view=subscribers", nil, false, 23),
			wantErr: true,
			expErr:  &ErrInvalidQueryParam{map[string]string{"userID": "123wa1", "count": "2~", "lastID": "33~"}},
		},
		{
			name:    "incorrect filter",
			inReq:   getHttpRequest(http.MethodGet, baseUrl+"?userID=123&count=2&lastID=33&view=subs", nil, false, 23),
			wantErr: true,
			expErr:  &ErrInvalidQueryParam{map[string]string{"view": "subs"}},
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actualOpts, err := GetOpts(test.inReq)
			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expOpts, actualOpts)
			}
		})
	}
}

// func TestHandler_GetSubscriptionInfoForUser(t *testing.T) {
// 	log, err := logger.New()
// 	if err != nil {
// 		t.Fatalf("test: log init - %s", err.Error())
// 	}

// 	san := bluemonday.UGCPolicy()
// 	baseUrl := "http://test.com/api/v1/subscription"

// 	cases := []struct {
// 		name        string
// 		url         string
// 		userID      int
// 		count       int
// 		lastID      int
// 		view        string
// 		inCtx       context.Context
// 		setSubMock  getUserSubscriptions
// 		setUserMock checkUserExistence
// 		expResp     JsonResponse
// 		wantErr     bool
// 		expErr      JsonErrResponse
// 	}{
// 		{
// 			name: "valid request from authenticated user",
// 			url:  baseUrl + "?userID=2&count=20&lastID=1&view=subscribers",
// 			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
// 				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(nil).Times(1)
// 			},
// 			inCtx: context.WithValue(context.Background(), auth.KeyCurrentUserID, 13),
// 			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, userID, count, lastID, currUserID int) {
// 				mockRepo.EXPECT().GetUserSubscriptions(ctx, userID, count, lastID, currUserID).Return([]user.SubscriptionUser{{2, "baobab", "/pic1", true}}, nil).Times(1)
// 			},
// 			expResp: JsonResponse{
// 				Status:  "ok",
// 				Message: "got subscription info successfully",
// 				Body:    []user.SubscriptionUser{{ID: 2, Username: "baobab", Avatar: "/pic1", HasSubscribeFromCurUser: true}},
// 			},
// 		},
// 		{
// 			name:  "valid request from authenticated user, no userID",
// 			url:   baseUrl + "?userID=2&count=20&lastID=1&view=subscribers",
// 			inCtx: context.WithValue(context.Background(), auth.KeyCurrentUserID, 13),
// 			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
// 				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(nil).Times(1)
// 			},
// 			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, userID, count, lastID, currUserID int) {
// 				mockRepo.EXPECT().GetUserSubscriptions(ctx, userID, count, lastID, currUserID).Return([]user.SubscriptionUser{{2, "baobab", "/pic1", false}}, nil).Times(1)
// 			},
// 			expResp: JsonResponse{
// 				Status:  "ok",
// 				Message: "got subscription info successfully",
// 				Body:    []user.SubscriptionUser{{ID: 2, Username: "baobab", Avatar: "/pic1", HasSubscribeFromCurUser: false}},
// 			},
// 		},
// 	}

// 	for _, test := range cases {
// 		test := test
// 		t.Run(test.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctl := gomock.NewController(t)
// 			defer ctl.Finish()

// 			mockSubRepo := mock_sub.NewMockRepository(ctl)
// 			mockUserRepo := mock_user.NewMockRepository(ctl)
// 			currUserID, _ := test.inCtx.Value(auth.KeyCurrentUserID).(int)
// 			test.setSubMock(mockSubRepo, test.inCtx, test.userID, test.count, test.lastID, currUserID)
// 			test.setUserMock(mockUserRepo, test.inCtx, test.userID)

// 			r := httptest.NewRequest(http.MethodGet, test.url, nil)
// 			r = r.WithContext(test.inCtx)
// 			w := httptest.NewRecorder()
// 			handler := New(log, UsecaseHub{nil, nil, nil, nil, usecase.New(log, mockSubRepo, mockUserRepo, san), nil, nil})
// 			handler.GetSubscriptionInfoForUser(w, r)

// 			if test.wantErr {
// 				var resp JsonErrResponse
// 				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
// 				require.Equal(t, test.expErr, resp)
// 			} else {
// 				var resp JsonResponse
// 				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
// 				require.Equal(t, test.expResp, resp)
// 			}
// 		})
// 	}
// }

func TestHandler_Subscribe(t *testing.T) {

	log, err := logger.New()
	if err != nil {
		t.Fatalf("test: log init - %s", err.Error())
	}

	san := bluemonday.UGCPolicy()

	cases := []struct {
		name          string
		inContentType string
		inBody        string
		inTo          int
		inCtx         context.Context
		setSubMock    subscribe
		setUserMock   checkUserExistence
		expResp       JsonResponse
		wantErr       bool
		expErr        JsonErrResponse
	}{
		{
			name:          "valid request",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
				mockRepo.EXPECT().CreateSubscriptionUser(ctx, from, to).Return(nil).Times(1)
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(nil).Times(1)
			},
			expResp: JsonResponse{
				Status:  "ok",
				Message: "subscribed successfully",
				Body:    nil,
			},
		},
		{
			name:          "invalid content-type",
			inContentType: "application/jsn",
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
			},
			inCtx:   context.Background(),
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&ErrInvalidContentType{"application/json"}).Error(),
				Code:    "bad_input",
			},
		},
		{
			name:          "self subsrciption",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 2),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
			},
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&usecase.ErrSelfSubscription{}).Error(),
				Code:    "bad_input",
			},
		},
		{
			name:          "already subscribed",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
				mockRepo.EXPECT().CreateSubscriptionUser(ctx, from, to).Return(&repo_sub.ErrSubscriptionAlreadyExist{}).Times(1)
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(nil).Times(1)
			},
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&repo_sub.ErrSubscriptionAlreadyExist{}).Error(),
				Code:    "already_exists",
			},
		},
		{
			name:          "non-existing user",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(&repo_user.ErrNonExistingUser{}).Times(1)
			},
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&repo_user.ErrNonExistingUser{}).Error(),
				Code:    "not_found",
			},
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockSubRepo := mock_sub.NewMockRepository(ctl)
			mockUserRepo := mock_user.NewMockRepository(ctl)
			from, _ := test.inCtx.Value(auth.KeyCurrentUserID).(int)
			test.setSubMock(mockSubRepo, test.inCtx, from, test.inTo)
			test.setUserMock(mockUserRepo, test.inCtx, test.inTo)

			r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(test.inBody))
			r = r.WithContext(test.inCtx)
			r.Header.Set("Content-Type", test.inContentType)
			w := httptest.NewRecorder()
			handler := New(log, UsecaseHub{nil, nil, nil, nil, usecase.New(log, mockSubRepo, mockUserRepo, san), nil, nil})
			handler.Subscribe(w, r)

			if test.wantErr {
				var resp JsonErrResponse
				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
				require.Equal(t, test.expErr, resp)
			} else {
				var resp JsonResponse
				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
				require.Equal(t, test.expResp, resp)
			}
		})
	}
}

func TestHandler_Unsubsribe(t *testing.T) {
	log, err := logger.New()
	if err != nil {
		t.Fatalf("test: log init - %s", err.Error())
	}

	san := bluemonday.UGCPolicy()

	cases := []struct {
		name          string
		inContentType string
		inBody        string
		inTo          int
		inCtx         context.Context
		setSubMock    subscribe
		setUserMock   checkUserExistence
		expResp       JsonResponse
		wantErr       bool
		expErr        JsonErrResponse
	}{
		{
			name:          "valid request",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
				mockRepo.EXPECT().DeleteSubscriptionUser(ctx, from, to).Return(nil).Times(1)
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(nil).Times(1)
			},
			expResp: JsonResponse{
				Status:  "ok",
				Message: "unsubscribed successfully",
				Body:    nil,
			},
		},
		{
			name:          "invalid content-type",
			inContentType: "application/jsn",
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
			},
			inCtx:   context.Background(),
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&ErrInvalidContentType{"application/json"}).Error(),
				Code:    "bad_input",
			},
		},
		{
			name:          "self unsubsrciption",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 2),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
			},
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&usecase.ErrSelfUnsubscription{}).Error(),
				Code:    "bad_input",
			},
		},
		{
			name:          "non-existing subscription",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
				mockRepo.EXPECT().DeleteSubscriptionUser(ctx, from, to).Return(&repo_sub.ErrNonExistingSubscription{}).Times(1)
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(nil).Times(1)
			},
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&repo_sub.ErrNonExistingSubscription{}).Error(),
				Code:    "not_found",
			},
		},
		{
			name:          "non-existing user",
			inContentType: "application/json",
			inBody:        `{"to": 2}`,
			inTo:          2,
			inCtx:         context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			setSubMock: func(mockRepo *mock_sub.MockRepository, ctx context.Context, from, to int) {
			},
			setUserMock: func(mockRepo *mock_user.MockRepository, ctx context.Context, userID int) {
				mockRepo.EXPECT().CheckUserExistence(ctx, userID).Return(&repo_user.ErrNonExistingUser{}).Times(1)
			},
			wantErr: true,
			expErr: JsonErrResponse{
				Status:  "error",
				Message: (&repo_user.ErrNonExistingUser{}).Error(),
				Code:    "not_found",
			},
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockSubRepo := mock_sub.NewMockRepository(ctl)
			mockUserRepo := mock_user.NewMockRepository(ctl)
			from, _ := test.inCtx.Value(auth.KeyCurrentUserID).(int)
			test.setSubMock(mockSubRepo, test.inCtx, from, test.inTo)
			test.setUserMock(mockUserRepo, test.inCtx, test.inTo)

			r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(test.inBody))
			r = r.WithContext(test.inCtx)
			r.Header.Set("Content-Type", test.inContentType)
			w := httptest.NewRecorder()
			handler := New(log, UsecaseHub{nil, nil, nil, nil, usecase.New(log, mockSubRepo, mockUserRepo, san), nil, nil})
			handler.Unsubscribe(w, r)

			if test.wantErr {
				var resp JsonErrResponse
				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
				require.Equal(t, test.expErr, resp)
			} else {
				var resp JsonResponse
				require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
				require.Equal(t, test.expResp, resp)
			}
		})
	}
}
