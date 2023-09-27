package pin

import (
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

func (u *Usecase) GetByID() (*entity.Pin, error) { return nil, nil }
