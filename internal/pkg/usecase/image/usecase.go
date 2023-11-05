package image

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	valid "github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
)

const PrefixURLImage = "httsp://pinspire.online:8081/"

var ErrInvalidImage = errors.New("invalid images")
var ErrUploadFile = errors.New("file upload failed")

type Usecase interface {
	UploadImage(path string, mimeType string, size int64, image io.Reader, check check.CheckSize) (string, error)
}

type imageCase struct {
	log  *log.Logger
	repo repo.Repository
}

func New(log *log.Logger, repo repo.Repository) *imageCase {
	return &imageCase{log, repo}
}

func (img *imageCase) UploadImage(path string, mimeType string, size int64, image io.Reader, check check.CheckSize) (string, error) {
	buf := bytes.NewBuffer(nil)

	extension, ok := valid.IsValidImage(io.TeeReader(image, buf), mimeType, check)
	if !ok {
		return "", ErrInvalidImage
	}

	io.Copy(buf, image)

	filename, written, err := img.repo.SaveImage(path, extension, buf)
	if err != nil {
		return "", fmt.Errorf("upload image: %w", err)
	}
	if written != size {
		return "", ErrUploadFile
	}
	return "https://pinspire.online:8081/" + filename, nil
}
