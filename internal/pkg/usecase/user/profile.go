package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/google/uuid"
)

var ErrBadMIMEType = errors.New("bad mime type")

func (u *userCase) UpdateUserAvatar(ctx context.Context, userID int, avatar io.Reader, mimeType string) error {
	filename := uuid.New().String()
	dir := "upload/pins/" + time.Now().UTC().Format("2006/02/01/")
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return fmt.Errorf("create dir for upload file: %w", err)
	}
	piecesMimeType := strings.Split(mimeType, "/")
	if len(piecesMimeType) != 2 || piecesMimeType[0] != "image" {
		return ErrBadMIMEType
	}

	filePath := dir + filename + "." + piecesMimeType[1]
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create %s to save avatar to it: %w", filePath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, avatar)
	if err != nil {
		return fmt.Errorf("copy avatar in file %s: %w", filePath, err)
	}
	u.log.Info("upload file", log.F{"file", filePath})

	err = u.repo.EditUserAvatar(ctx, userID, "https://pinspire.online/"+filePath)
	if err != nil {
		return fmt.Errorf("edit user avatar: %w", err)
	}

	return nil
}
