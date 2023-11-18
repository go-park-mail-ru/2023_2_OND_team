package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/crypto"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func TestRegisterMock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockRepository(ctrl)
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}
	usecase := New(log, nil, repoMock)
	userRegister := &user.User{}

	repoMock.EXPECT().
		AddNewUser(ctx, userRegister).
		Return(nil).
		Times(1)

	err = usecase.Register(ctx, userRegister)
	assert.NoError(t, err)

	expErr := errors.New("repo error")
	repoMock.EXPECT().
		AddNewUser(ctx, userRegister).
		Return(expErr).
		Times(1)

	err = usecase.Register(ctx, userRegister)
	require.ErrorIs(t, err, expErr)
	require.EqualError(t, err, "user registration: repo error")
}

func TestAuthentication(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockRepository(ctrl)
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}
	usecase := New(log, nil, repoMock)
	cred := UserCredentials{
		Password: "1234",
		Username: "test",
	}
	userAuth := &user.User{
		ID:       4,
		Username: "test username",
		Password: "$$$$$$$$$$$$$$$$$$$$$$$$$$$$",
	}

	repoMock.EXPECT().
		GetUserByUsername(ctx, cred.Username).
		Return(userAuth, nil).
		Times(1)

	actualUser, err := usecase.Authentication(ctx, cred)
	require.Nil(t, actualUser)
	require.Equal(t, err, ErrUserAuthentication)

	expErr := errors.New("get user error")
	repoMock.EXPECT().
		GetUserByUsername(ctx, cred.Username).
		Return(userAuth, expErr).
		Times(1)

	actualUser, err = usecase.Authentication(ctx, cred)
	require.Nil(t, actualUser)
	require.ErrorIs(t, err, expErr)
	require.EqualError(t, err, "user authentication: get user error")

	salt, err := crypto.NewRandomString(lenSalt)
	if err != nil {
		t.Fatal(err)
	}

	userAuth.Password = salt + crypto.PasswordHash(cred.Password, salt, lenPasswordHash)
	repoMock.EXPECT().
		GetUserByUsername(ctx, cred.Username).
		Return(userAuth, nil).
		Times(1)

	actualUser, err = usecase.Authentication(ctx, cred)
	require.NoError(t, err)
	require.Equal(t, actualUser.Password, "")
	require.Equal(t, userAuth, actualUser)
}

func TestFindOutUsernameAndAvatar(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockRepository(ctrl)
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}
	usecase := New(log, nil, repoMock)
	userID := 3

	repoMock.EXPECT().
		GetUsernameAndAvatarByID(ctx, userID).
		Return("TestUsername", "avatar.png", nil).
		Times(1)
	username, avatar, err := usecase.FindOutUsernameAndAvatar(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, "TestUsername", username)
	require.Equal(t, "avatar.png", avatar)

	expErr := errors.New("error find out")
	repoMock.EXPECT().
		GetUsernameAndAvatarByID(ctx, userID).
		Return("", "", expErr).
		Times(1)
	username, avatar, err = usecase.FindOutUsernameAndAvatar(ctx, userID)
	require.Equal(t, expErr, err)
	require.Equal(t, "", username)
	require.Equal(t, "", avatar)
}
