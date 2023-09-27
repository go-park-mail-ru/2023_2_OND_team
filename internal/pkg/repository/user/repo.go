package user

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type UserCredentials struct {
	Username string
	Password string
}

type Repository interface {
	AddNewUser(ctx context.Context, user *user.User) error
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
}
