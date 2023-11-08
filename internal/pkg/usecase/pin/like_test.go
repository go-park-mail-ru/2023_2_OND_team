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
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func TestIsAvailableBatchPinForFixOnBoard(t *testing.T) {
	testCases := []struct {
		Name    string
		Pins    []entity.Pin
		UserID  int
		WantErr error
	}{
		{
			Name:    "one deleted pin",
			Pins:    []entity.Pin{{DeletedAt: pgtype.Timestamptz{Valid: true}}},
			UserID:  12,
			WantErr: ErrPinDeleted,
		},
		{
			Name: "all available",
			Pins: []entity.Pin{
				{Author: &user.User{ID: 34}, DeletedAt: pgtype.Timestamptz{Valid: false}, Public: true},
				{Author: &user.User{ID: 12}, Public: false},
				{Author: &user.User{ID: 34}, Public: true},
			},
			UserID:  12,
			WantErr: nil,
		},
		{
			Name: "one not available",
			Pins: []entity.Pin{
				{Author: &user.User{ID: 34}, DeletedAt: pgtype.Timestamptz{Valid: false}, Public: true},
				{Author: &user.User{ID: 12}, Public: false},
				{Author: &user.User{ID: 34}, Public: false},
			},
			UserID:  12,
			WantErr: ErrForbiddenAction,
		},
	}
	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			test := test
			t.Parallel()
			err := isAvailableBatchPinForFixOnBoard(test.UserID, test.Pins...)
			require.Equal(t, test.WantErr, err)
		})
	}
}
func TestSetLikeOnAvailablePin(t *testing.T) {
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
	wantCountLike := 12
	pinID, userID := 123, 1
	pin := &entity.Pin{
		ID:     pinID,
		Author: &user.User{ID: 0},
		Public: false,
	}

	repo.EXPECT().
		GetPinByID(ctx, pinID, false).
		Return(pin, nil).
		Times(1)
	repo.EXPECT().
		IsAvailableToUserAsContributorBoard(ctx, pinID, userID).
		Return(true, nil).
		Times(1)
	repo.EXPECT().
		SetLike(ctx, pinID, userID).
		Return(wantCountLike, nil).
		Times(1)

	actualCoutnLike, actualErr := pinCase.SetLikeFromUser(ctx, pinID, userID)
	require.NoError(t, actualErr)
	require.Equal(t, wantCountLike, actualCoutnLike)

	pin.Author.ID = userID
	repo.EXPECT().
		GetPinByID(ctx, pinID, false).
		Return(pin, nil).
		Times(1)
	repo.EXPECT().
		SetLike(ctx, pinID, userID).
		Return(wantCountLike, nil).
		Times(1)

	actualCoutnLike, actualErr = pinCase.SetLikeFromUser(ctx, pinID, userID)
	require.NoError(t, actualErr)
	require.Equal(t, wantCountLike, actualCoutnLike)
}

func TestSetLikeOnNotAvailablePin(t *testing.T) {
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
	wantCountLike := 0
	pinID, userID := 123, 1
	pin := &entity.Pin{
		ID:     pinID,
		Author: &user.User{ID: 0},
		Public: false,
	}

	repo.EXPECT().
		GetPinByID(ctx, pinID, false).
		Return(pin, nil).
		Times(1)
	repo.EXPECT().
		IsAvailableToUserAsContributorBoard(ctx, pinID, userID).
		Return(false, nil).
		Times(1)

	actualCoutnLike, actualErr := pinCase.SetLikeFromUser(ctx, pinID, userID)
	require.ErrorIs(t, actualErr, ErrPinNotAccess)
	require.Equal(t, wantCountLike, actualCoutnLike)

	wantErr := errors.New("returned IsAvailableToUserAsContributorBoard")
	repo.EXPECT().
		GetPinByID(ctx, pinID, false).
		Return(pin, nil).
		Times(1)
	repo.EXPECT().
		IsAvailableToUserAsContributorBoard(ctx, pinID, userID).
		Return(false, wantErr).
		Times(1)

	actualCoutnLike, actualErr = pinCase.SetLikeFromUser(ctx, pinID, userID)
	require.ErrorIs(t, actualErr, wantErr)
	require.EqualError(t, actualErr, "set like from user: fail check available pin: returned IsAvailableToUserAsContributorBoard")
	require.Equal(t, wantCountLike, actualCoutnLike)

	pin.DeletedAt.Valid = true
	repo.EXPECT().
		GetPinByID(ctx, pinID, false).
		Return(pin, nil).
		Times(1)

	actualCoutnLike, actualErr = pinCase.SetLikeFromUser(ctx, pinID, userID)
	require.ErrorIs(t, actualErr, ErrPinDeleted)
	require.EqualError(t, actualErr, "set like from user: pin has been deleted")
	require.Equal(t, wantCountLike, actualCoutnLike)
}

func TestDeleteLikeFromUser(t *testing.T) {
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
	pinID, userID := 123, 1
	wantCountLike := 999
	repo.EXPECT().
		DelLike(ctx, pinID, userID).
		Return(nil).
		Times(1)

	actualCountLike, err := pinCase.DeleteLikeFromUser(ctx, pinID, userID)
	require.NoError(t, err)
	require.Equal(t, wantCountLike, actualCountLike)
}

func TestCheckUserHasSetLike(t *testing.T) {
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
	pinID, userID := 123, 1

	repo.EXPECT().
		IsSetLike(ctx, pinID, userID).
		Return(true, nil).
		Times(1)

	has, err := pinCase.CheckUserHasSetLike(ctx, pinID, userID)
	require.NoError(t, err)
	require.True(t, has)
}
