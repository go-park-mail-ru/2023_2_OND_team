package pin

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/validator/image/check"
)

var ErrBadMIMEType = errors.New("bad mime type")

//go:generate mockgen -destination=./mock/pin_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	SelectNewPins(ctx context.Context, count, minID, maxID int) ([]entity.Pin, int, int)
	SelectUserPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, int, int)
	CreateNewPin(ctx context.Context, pin *entity.Pin, mimeTypePicture string, sizePicture int64, picture io.Reader) error
	DeletePinFromUser(ctx context.Context, pinID, userID int) error
	SetLikeFromUser(ctx context.Context, pinID, userID int) (int, error)
	CheckUserHasSetLike(ctx context.Context, pinID, userID int) (bool, error)
	DeleteLikeFromUser(ctx context.Context, pinID, userID int) (int, error)
	EditPinByID(ctx context.Context, pinID, userID int, updateData *PinUpdateData) error
	ViewAnPin(ctx context.Context, pinID, userID int) (*entity.Pin, error)
	IsAvailablePinForFixOnBoard(ctx context.Context, pinID, userID int) error
	IsAvailableBatchPinForFixOnBoard(ctx context.Context, pinID []int, userID int) error

	SelectUserLikedPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, int, int)
	ViewFeedPin(ctx context.Context, userID int, cfg pin.FeedPinConfig) (pin.FeedPin, error)
}

type pinCase struct {
	image.Usecase
	log  *log.Logger
	repo repo.Repository
}

func New(log *log.Logger, imgCase image.Usecase, repo repo.Repository) *pinCase {
	return &pinCase{
		Usecase: imgCase,
		log:     log,
		repo:    repo,
	}
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

func (p *pinCase) CreateNewPin(ctx context.Context, pin *entity.Pin, mimeTypePicture string, sizePicture int64, picture io.Reader) error {
	picturePin, err := p.UploadImage("pins/", mimeTypePicture, sizePicture, picture, check.BothSidesFallIntoRange(200, 1800))
	if err != nil {
		return fmt.Errorf("uploading an avatar when creating pin: %w", err)
	}
	pin.Picture = picturePin

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

func (p *pinCase) ViewFeedPin(ctx context.Context, userID int, cfg pin.FeedPinConfig) (pin.FeedPin, error) {
	_, hasBoard := cfg.Board()
	user, hasUser := cfg.User()

	if !hasBoard && (userID == UserUnknown || !hasUser || userID != user) && cfg.Protection != pin.FeedProtectionPublic {
		return pin.FeedPin{}, ErrForbiddenAction
	}

	return p.repo.GetFeedPins(ctx, cfg)
}
