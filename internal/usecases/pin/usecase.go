package pin

import (
	"context"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type Usecase struct {
	log  *logger.Logger
	repo repo.Repository
}

func New(log *logger.Logger, repo repo.Repository) *Usecase {
	return &Usecase{log, repo}
}

func (u *Usecase) SelectNewPins(ctx context.Context, count, lastID int) ([]entity.Pin, int) {
	pins, err := u.repo.GetSortedNPinsAfterID(ctx, count, lastID)
	if err != nil {
		u.log.Error(err.Error())
	}
	if len(pins) == 0 {
		return []entity.Pin{}, lastID
	}
	return pins, pins[len(pins)-1].ID
}
