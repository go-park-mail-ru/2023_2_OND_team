package user

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/crypto"
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
}

type userCase struct {
	log  *logger.Logger
	repo repo.Repository
}

func New(log *logger.Logger, repo repo.Repository) *userCase {
	return &userCase{log, repo}
}

func (u *userCase) Register(ctx context.Context, user *entity.User) error {
	salt, err := crypto.NewRandomString(lenSalt)
	if err != nil {
		return fmt.Errorf("generating salt for registration: %w", err)
	}

	user.Password = salt + crypto.PasswordHash(user.Password, salt, lenPasswordHash)
	err = u.repo.AddNewUser(ctx, user)
	if err != nil {
		return fmt.Errorf("user registration: %w", err)
	}
	return nil
}

func (u *userCase) Authentication(ctx context.Context, credentials userCredentials) (*entity.User, error) {
	user, err := u.repo.GetUserByUsername(ctx, credentials.Username)
	if err != nil {
		return nil, fmt.Errorf("user authentication: %w", err)
	}
	salt := user.Password[:lenSalt]
	if crypto.PasswordHash(credentials.Password, salt, lenPasswordHash) != user.Password[lenSalt:] {
		return nil, ErrUserAuthentication
	}
	user.Password = ""
	return user, nil
}

func (u *userCase) FindOutUsernameAndAvatar(ctx context.Context, userID int) (username string, avatar string, err error) {
	return u.repo.GetUsernameAndAvatarByID(ctx, userID)
}