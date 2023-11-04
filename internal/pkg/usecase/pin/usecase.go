package pin

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrBadMIMEType = errors.New("bad mime type")

type Usecase interface {
	SelectNewPins(ctx context.Context, count, minID, maxID int) ([]entity.Pin, int, int)
	SelectUserPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, int, int)
	CreateNewPin(ctx context.Context, pin *entity.Pin) error
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

func (p *pinCase) SelectUserPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, int, int) {
	pins, err := p.repo.GetSortedUserPins(ctx, userID, count, minID, maxID)
	if err != nil {
		p.log.Error(err.Error())
	}
	if len(pins) == 0 {
		return []entity.Pin{}, minID, maxID
	}
	return pins, pins[len(pins)-1].ID, pins[0].ID
}

func (p *pinCase) CreateNewPin(ctx context.Context, pin *entity.Pin) error {
	pin.Picture = "https://pinspire.online:8081/" + pin.Picture

	err := p.repo.AddNewPin(ctx, pin)
	if err != nil {
		return fmt.Errorf("add new pin: %w", err)
	}

	return nil
}

func (p *pinCase) DeletePinFromUser(ctx context.Context, pinID, userID int) error {
	return p.repo.DeletePin(ctx, pinID, userID)
}

func (p *pinCase) ViewAnPin(ctx context.Context, pinID, userID int) (*entity.Pin, error) {
	pin, err := p.repo.GetPinByID(ctx, pinID, true)
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
