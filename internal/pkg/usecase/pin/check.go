package pin

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

var (
	ErrPinNotAccess = errors.New("pin is not available")
	ErrPinDeleted   = errors.New("")
)

const UserUnknown = -1

func (p *pinCase) isAvailablePinForViewingUser(ctx context.Context, pin *entity.Pin, userID int) error {
	if pin.DeletedAt.Valid {
		return ErrPinDeleted
	}

	if pin.Public || pin.AuthorID == userID {
		return nil
	}
	if userID == UserUnknown {
		return ErrPinNotAccess
	}

	ok, err := p.repo.IsAvailableToUserAsContributorBoard(ctx, pin.ID, userID)
	if err != nil {
		return fmt.Errorf("fail check available pin: %w", err)
	}

	if !ok {
		return ErrPinNotAccess
	}
	return nil
}

func (p *pinCase) isAvailablePinForSetLike(ctx context.Context, pinID, userID int) error {
	pin, err := p.repo.GetPinByID(ctx, pinID)
	if err != nil {
		return fmt.Errorf("get a pin to check for the availability of a like: %w", err)
	}

	return p.isAvailablePinForViewingUser(ctx, pin, userID)
}
