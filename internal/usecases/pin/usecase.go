package pin

import (
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
)

type Usecase struct {
	repo repo.Repository
}

func New(repo repo.Repository) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) GetByID() (*entity.Pin, error) { return nil, nil }
