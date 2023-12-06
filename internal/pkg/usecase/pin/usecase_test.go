package pin

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin/mock"
	mockImage "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func TestCreateNewPin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := mock.NewMockRepository(ctrl)
	imgCase := mockImage.NewMockUsecase(ctrl)
	pinCase := New(log, imgCase, repo)
	mimeType, size := "image/webp", int64(45)
	filename := "filename.webp"
	pin := &entity.Pin{
		ID: 34,
	}

	imgCase.EXPECT().
		UploadImage("pins/", mimeType, size, nil, gomock.Any()).
		Return(filename, nil).
		Times(1)

	pin.Picture = filename
	repo.EXPECT().
		AddNewPin(ctx, pin).
		Return(nil).
		Times(1)

	err = pinCase.CreateNewPin(ctx, pin, "image/webp", size, nil)
	require.NoError(t, err)
}

func TestDeletePinFromUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := mock.NewMockRepository(ctrl)
	pinCase := New(log, nil, repo)
	pinID, userID := 8, 16

	wantErr := errors.New("returned err")
	repo.EXPECT().
		DeletePin(ctx, pinID, userID).
		Return(wantErr).
		Times(1)

	actualErr := pinCase.DeletePinFromUser(ctx, pinID, userID)
	require.Equal(t, wantErr, actualErr)
}

func TestViewAnPin(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log, err := logger.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := mock.NewMockRepository(ctrl)
	pinCase := New(log, nil, repo)
	pinID, userID := 44, 90
	countLike := 22
	tags := []entity.Tag{{Title: "good"}, {Title: "home"}}
	wantPin := entity.Pin{
		Title:  pgtype.Text{String: "someone else 's public pin", Valid: true},
		Author: &user.User{ID: 100},
		Public: true,
	}
	actualPin := new(entity.Pin)
	*actualPin = wantPin

	repo.EXPECT().GetPinByID(ctx, pinID, true).Return(actualPin, nil).Times(1)
	repo.EXPECT().GetCountLikeByPinID(ctx, pinID).Return(countLike, nil).Times(1)
	repo.EXPECT().GetTagsByPinID(ctx, pinID).Return(tags, nil).Times(1)

	wantPin.CountLike = countLike
	wantPin.Tags = tags

	actualPin, actualErr := pinCase.ViewAnPin(ctx, pinID, userID)
	require.NoError(t, actualErr)
	require.Equal(t, wantPin, *actualPin)
}
