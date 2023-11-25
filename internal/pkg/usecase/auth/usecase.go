package auth

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	authRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/auth"
)

type Usecase interface {
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, username, password string) (*session.Session, error)
	GetUserIDBySession(ctx context.Context, sess *session.Session) (int, error)
	Logout(ctx context.Context, sess *session.Session) error
}

type authCase struct {
	repo authRepo.Repository
}

func New(repo authRepo.Repository) *authCase {
	return &authCase{repo}
}

func (ac *authCase) Register(ctx context.Context, user *entity.User) error {
	return ac.repo.Register(ctx, user)
}

func (ac *authCase) Logout(ctx context.Context, sess *session.Session) error {
	return ac.repo.Logout(ctx, sess)
}

func (ac *authCase) Login(ctx context.Context, username, password string) (*session.Session, error) {
	return ac.repo.Login(ctx, username, password)
}

func (ac *authCase) GetUserIDBySession(ctx context.Context, sess *session.Session) (int, error) {
	return ac.repo.GetUserID(ctx, sess)
}
