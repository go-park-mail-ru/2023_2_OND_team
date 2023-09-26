package user

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Repository interface {
	AddNewUser(ctx context.Context, user *user.User) error
	GetUserByPasswordLogin(ctx context.Context, password, login string) (*user.User, error)
}
