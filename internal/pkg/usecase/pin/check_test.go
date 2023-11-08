package pin

import (
	"context"
	"testing"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
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

func TestIsAvailablePinForFixOnBoard(t *testing.T) {
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
	pin := &entity.Pin{Author: &user.User{ID: userID}, DeletedAt: pgtype.Timestamptz{Valid: true}}

	repo.EXPECT().
		GetPinByID(ctx, pinID, false).
		Return(pin, nil).
		Times(1)

	err = pinCase.IsAvailablePinForFixOnBoard(ctx, pinID, userID)
	require.Equal(t, ErrPinDeleted, err)
}
