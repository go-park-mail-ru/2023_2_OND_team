package pin

import (
	"context"
	"testing"

	repository "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin/mock"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestEditPinByID(t *testing.T) {
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
	tags := []string{"new", "tag"}

	updateData := &PinUpdateData{
		Title:       new(string),
		Description: new(string),
		Tags:        tags,
		Public:      new(bool),
	}

	repo.EXPECT().
		EditPin(ctx, pinID, repository.S{
			"title":       "",
			"description": "",
			"public":      false,
		}, tags).
		Return(nil).
		Times(1)

	err = pinCase.EditPinByID(ctx, pinID, userID, updateData)
	require.NoError(t, err)
}
