package user

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/crypto"
)

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
