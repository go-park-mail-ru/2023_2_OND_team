package user

import (
	"context"
	"errors"
	"fmt"
	"io"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	repository "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/crypto"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
)

var ErrBadBody = errors.New("bad body avatar")

func (u *userCase) UpdateUserAvatar(ctx context.Context, userID int, mimeTypeAvatar string, sizeAvatar int64, avatar io.Reader) error {
	avatarProfile, err := u.UploadImage("avatars/", mimeTypeAvatar, sizeAvatar, avatar, check.BothSidesFallIntoRange(200, 1800))
	if err != nil {
		return fmt.Errorf("uploading an avatar when updating avatar profile: %w", err)
	}
	err = u.repo.EditUserAvatar(ctx, userID, avatarProfile)
	if err != nil {
		return fmt.Errorf("edit user avatar: %w", err)
	}

	return nil
}

func (u *userCase) GetAllProfileInfo(ctx context.Context, userID int) (*entity.User, error) {
	return u.repo.GetAllUserData(ctx, userID)
}

func (u *userCase) EditProfileInfo(ctx context.Context, userID int, updateData *ProfileUpdateData) error {
	updateFields := repository.S{}
	if updateData.Username != nil {
		updateFields["username"] = *updateData.Username
	}
	if updateData.Email != nil {
		updateFields["email"] = *updateData.Email
	}
	if updateData.Name != nil {
		updateFields["name"] = *updateData.Name
	}
	if updateData.Surname != nil {
		updateFields["surname"] = *updateData.Surname
	}
	if updateData.AboutMe != nil {
		updateFields["about_me"] = *updateData.AboutMe
	}
	if updateData.Password != nil {
		salt, err := crypto.NewRandomString(lenSalt)
		if err != nil {
			return fmt.Errorf("generating salt for registration: %w", err)
		}

		updateFields["password"] = salt + crypto.PasswordHash(*updateData.Password, salt, lenPasswordHash)
	}

	err := u.repo.EditUserInfo(ctx, userID, updateFields)
	if err != nil {
		return fmt.Errorf("edit profile info: %w", err)
	}
	return nil
}
