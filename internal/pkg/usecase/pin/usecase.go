package pin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/google/uuid"
)

var ErrBadMIMEType = errors.New("bad mime type")

type Usecase interface {
	SelectNewPins(ctx context.Context, count, minID, maxID int) ([]entity.Pin, int, int)
	CreateNewPin(ctx context.Context, pin *entity.Pin, picture io.Reader, mimeType string) error
	DeletePinFromUser(ctx context.Context, pinID, userID int) error
	SetLikeFromUser(ctx context.Context, pinID, userID int) error
	DeleteLikeFromUser(ctx context.Context, pinID, userID int) error
	EditPinByID(ctx context.Context, pinID, userID int, updateData *pinUpdateData) error
	ViewAnPin(ctx context.Context, pinID, userID int) (*entity.Pin, error)
}

type pinCase struct {
	log  *log.Logger
	repo repo.Repository
}

func New(log *log.Logger, repo repo.Repository) *pinCase {
	return &pinCase{log, repo}
}

func (p *pinCase) SelectNewPins(ctx context.Context, count, minID, maxID int) ([]entity.Pin, int, int) {
	pins, err := p.repo.GetSortedNewNPins(ctx, count, minID, maxID)
	if err != nil {
		p.log.Error(err.Error())
	}
	if len(pins) == 0 {
		return []entity.Pin{}, minID, maxID
	}
	return pins, pins[len(pins)-1].ID, pins[0].ID
}

func (p *pinCase) CreateNewPin(ctx context.Context, pin *entity.Pin, picture io.Reader, mimeType string) error {
	filename := uuid.New().String()
	dir := "upload/pins/" + time.Now().UTC().Format("2006/01/02/")
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

	_, err = io.Copy(file, picture)
	if err != nil {
		return fmt.Errorf("copy avatar in file %s: %w", filePath, err)
	}
	p.log.Info("upload file", log.F{"file", filePath})

	pin.Picture = "https://pinspire.online:8081/" + filePath

	err = p.repo.AddNewPin(ctx, pin)
	if err != nil {
		return fmt.Errorf("add new pin: %w", err)
	}

	return nil
}

func (p *pinCase) DeletePinFromUser(ctx context.Context, pinID, userID int) error {
	return p.repo.DeletePin(ctx, pinID, userID)
}

func (p *pinCase) ViewAnPin(ctx context.Context, pinID, userID int) (*entity.Pin, error) {
	pin, err := p.repo.GetPinByID(ctx, pinID)
	if err != nil {
		return nil, fmt.Errorf("get a pin to view: %w", err)
	}

	if err := p.isAvailablePinForViewingUser(ctx, pin, userID); err != nil {
		return nil, fmt.Errorf("view pin: %w", err)
	}

	pin.CountLike, err = p.repo.GetCountLikeByPinID(ctx, pinID)
	if err != nil {
		p.log.Error(err.Error())
	}
	pin.Tags, err = p.repo.GetTagsByPinID(ctx, pinID)
	if err != nil {
		p.log.Error(err.Error())
	}

	return pin, nil
}
