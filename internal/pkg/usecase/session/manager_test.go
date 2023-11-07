package session

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateNewSessionForUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	sessRepo := mock.NewMockRepository(ctrl)

	sm := New(log, sessRepo)
	sessRepo.EXPECT().
		AddSession(ctx, gomock.Any()).
		Return(nil).
		Times(1)

	expUserID := 32
	s, err := sm.CreateNewSessionForUser(ctx, expUserID)
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, expUserID, s.UserID)

	expErr := errors.New("err")
	sessRepo.EXPECT().
		AddSession(ctx, gomock.Any()).
		Return(expErr).
		Times(1)

	s, err = sm.CreateNewSessionForUser(ctx, 0)
	require.ErrorIs(t, err, expErr)
	require.Nil(t, s)
}

func TestDeleteUserSession(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	sessRepo := mock.NewMockRepository(ctrl)

	sm := New(log, sessRepo)
	expKey := "session-key"
	expErr := errors.New("err")

	sessRepo.EXPECT().
		DeleteSessionByKey(ctx, expKey).
		Return(expErr).
		Times(1)

	err = sm.DeleteUserSession(ctx, expKey)
	require.ErrorIs(t, err, expErr)
}
func TestGetUserIDBySessionKey(t *testing.T) {
	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	db, err := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer db.Close()

	sm := New(log, ramrepo.NewRamSessionRepo(db))

	testCases := []struct {
		name        string
		session_key string
		expUserId   int
		expErr      error
	}{
		{
			"providing valid session key",
			"461afabf38b3147c",
			1,
			nil,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			id, err := sm.GetUserIDBySessionKey(context.Background(), tCase.session_key)
			require.Equal(t, tCase.expErr, err)
			require.Equal(t, tCase.expUserId, id)
		})
	}
}
