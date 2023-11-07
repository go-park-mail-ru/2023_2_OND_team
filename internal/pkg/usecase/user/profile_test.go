package user

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	repoUser "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user/mock"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func TestUpdateUserAvatar(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := repo.NewMockRepository(ctrl)
	imageUsecase := usecase.NewMockUsecase(ctrl)
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}
	usecase := New(log, imageUsecase, userRepo)
	userID := 12
	var image io.Reader = bytes.NewBuffer(make([]byte, 2))

	imageUsecase.EXPECT().
		UploadImage("avatars/", "image/png", int64(128), image, gomock.Any()).
		Return("https://pinspire.online/upload/avatars/2023/avatar.png", nil).
		Times(1)

	userRepo.EXPECT().
		EditUserAvatar(ctx, userID, "https://pinspire.online/upload/avatars/2023/avatar.png").
		Return(nil).
		Times(1)

	err = usecase.UpdateUserAvatar(ctx, userID, "image/png", 128, image)
	require.NoError(t, err)

	expErr := errors.New("upload avatar error")
	imageUsecase.EXPECT().
		UploadImage("avatars/", "image/jpeg", int64(1), image, gomock.Any()).
		Return("", expErr).
		Times(1)

	err = usecase.UpdateUserAvatar(ctx, userID, "image/jpeg", 1, image)
	require.ErrorIs(t, err, expErr)
	require.EqualError(t, err, "uploading an avatar when updating avatar profile: upload avatar error")

	imageUsecase.EXPECT().
		UploadImage("avatars/", "", int64(-1), image, gomock.Any()).
		Return("", nil).
		Times(1)

	userRepo.EXPECT().
		EditUserAvatar(ctx, 144, "").
		Return(expErr).
		Times(1)
	err = usecase.UpdateUserAvatar(ctx, 144, "", -1, image)
	require.ErrorIs(t, err, expErr)
	require.EqualError(t, err, "edit user avatar: upload avatar error")
}

func TestGetAllProfileInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repo.NewMockRepository(ctrl)
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}
	usecase := New(log, nil, repoMock)
	wantUser := &user.User{}
	wantErr := errors.New("err info")

	repoMock.EXPECT().
		GetAllUserData(ctx, 11).
		Return(wantUser, wantErr).
		Times(1)
	actualUser, err := usecase.GetAllProfileInfo(ctx, 11)
	require.Equal(t, wantUser, actualUser)
	require.Equal(t, wantErr, err)
}

func TestEditProfileInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repo.NewMockRepository(ctrl)
	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}
	usecase := New(log, nil, repoMock)
	updateData := &ProfileUpdateData{
		Username: new(string),
		Email:    new(string),
		Name:     new(string),
		Surname:  new(string),
		AboutMe:  new(string),
	}

	repoMock.EXPECT().
		EditUserInfo(ctx, 5, repoUser.S{
			"username": "",
			"email":    "",
			"name":     "",
			"surname":  "",
			"about_me": "",
		}).
		Return(nil).
		Times(1)
	err = usecase.EditProfileInfo(ctx, 5, updateData)
	require.NoError(t, err)

	wantErr := errors.New("edit profile error")
	repoMock.EXPECT().
		EditUserInfo(ctx, 5, repoUser.S{}).
		Return(wantErr).
		Times(1)
	err = usecase.EditProfileInfo(ctx, 5, &ProfileUpdateData{})
	require.ErrorIs(t, err, wantErr)
	require.EqualError(t, err, "edit profile info: edit profile error")
}
