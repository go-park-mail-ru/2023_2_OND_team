package user

import "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"

type Usecase struct {
	repo user.Repository
}

func New(repo user.Repository) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) Register() error { return nil }
