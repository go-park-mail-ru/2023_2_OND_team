package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func TestGetProfileInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	usecase := mock.NewMockUsecase(ctrl)
	hander := New(log, nil, usecase, nil, nil)
	rec := httptest.NewRecorder()
	request := &http.Request{}
	ctxExp := context.WithValue(ctx, auth.KeyCurrentUserID, 122)
	wantUser := user.User{
		ID: 122,
	}

	usecase.EXPECT().GetAllProfileInfo(ctxExp, 122).
		Return(&user.User{ID: 122}, nil).
		Times(1)

	hander.GetProfileInfo(rec, request.WithContext(ctxExp))
	res := rec.Result()
	defer res.Body.Close()

	actualBody := JsonResponse{Body: &user.User{}}
	err = json.NewDecoder(res.Body).Decode(&actualBody)
	require.NoError(t, err)
	fmt.Println(actualBody.Body)
	wantBody := JsonResponse{
		Status:  "ok",
		Message: "user data has been successfully received",
		Body:    &wantUser,
	}
	require.Equal(t, wantBody, actualBody)
}
