package user

import (
	"context"
	"errors"
	"io"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrUserAuthentication = errors.New("user authentication")

const (
	lenSalt         = 16
	lenPasswordHash = 64
)

type Usecase interface {
	Register(ctx context.Context, user *entity.User) error
	Authentication(ctx context.Context, credentials userCredentials) (*entity.User, error)
	FindOutUsernameAndAvatar(ctx context.Context, userID int) (username string, avatar string, err error)
	UpdateUserAvatar(ctx context.Context, userID int, avatar io.Reader, mimeType string) error
}

type userCase struct {
	log  *logger.Logger
	repo repo.Repository
}

func New(log *logger.Logger, repo repo.Repository) *userCase {
	return &userCase{log, repo}
}
