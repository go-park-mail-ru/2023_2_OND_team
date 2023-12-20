package image

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/image"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	valid "github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
)

const PrefixURLImage = "https://pinspire.online:8081/"

var (
	ErrInvalidImage = errors.New("invalid images")
	ErrUploadFile   = errors.New("file upload failed")
)

//go:generate mockgen -destination=./mock/image_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	UploadImage(ctx context.Context, path string, mimeType string, size int64, image io.Reader, check check.CheckSize) (string, error)
}

type imageCase struct {
	log    *log.Logger
	repo   repo.Repository
	filter ImageFilter
}

func New(log *log.Logger, repo repo.Repository, filter ImageFilter) *imageCase {
	return &imageCase{log, repo, filter}
}

func (img *imageCase) UploadImage(ctx context.Context, path string, mimeType string, size int64, image io.Reader, check check.CheckSize) (string, error) {
	buf := bytes.NewBuffer(nil)

	extension, ok := valid.IsValidImage(io.TeeReader(image, buf), mimeType, check)
	if !ok {
		return "", ErrInvalidImage
	}
	io.Copy(buf, image)

	err := img.filter.Filter(ctx, buf.Bytes(), explicitLabels)
	if err != nil {
		if err == ErrExplicitImage {
			return "", err
		}
		return "", fmt.Errorf("upload image: %w", err)
	}

	filename, written, err := img.repo.SaveImage(path, extension, buf)
	if err != nil {
		return "", fmt.Errorf("upload image: %w", err)
	}
	if written != size {
		return "", ErrUploadFile
	}
	return PrefixURLImage + filename, nil
}
