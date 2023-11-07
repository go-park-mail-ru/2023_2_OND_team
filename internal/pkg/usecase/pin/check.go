package pin

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

var (
	ErrPinNotAccess    = errors.New("pin is not available")
	ErrPinDeleted      = errors.New("pin has been deleted")
	ErrForbiddenAction = errors.New("this action is not available to the user")
	ErrEmptyBatch      = errors.New("an empty batch was received")
	ErrSizeBatch       = errors.New("the batch size exceeds the maximum possible")
)

const MaxSizeBatchPin = 100

const UserUnknown = -1

func (p *pinCase) IsAvailablePinForFixOnBoard(ctx context.Context, pinID, userID int) error {
	pin, err := p.repo.GetPinByID(ctx, pinID, false)
	if err != nil {
		return err
	}

	return isAvailableBatchPinForFixOnBoard(userID, *pin)
}

func (p *pinCase) IsAvailableBatchPinForFixOnBoard(ctx context.Context, pinID []int, userID int) error {
	if len(pinID) == 0 {
		return ErrEmptyBatch
	}
	if len(pinID) > MaxSizeBatchPin {
		return ErrSizeBatch
	}

	pins, err := p.repo.GetBatchPinByID(ctx, pinID)
	if err != nil {
		return fmt.Errorf("get batch pin for chekc available: %w", err)
	}
	if err = isAvailableBatchPinForFixOnBoard(userID, pins...); err != nil {
		return fmt.Errorf("one of the pins turned out to be inaccessible for fixing on the board: %w", err)
	}
	return nil
}

func (p *pinCase) isAvailablePinForViewingUser(ctx context.Context, pin *entity.Pin, userID int) error {
	if pin.DeletedAt.Valid {
		return ErrPinDeleted
	}

	if pin.Public || pin.Author.ID == userID {
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
	pin, err := p.repo.GetPinByID(ctx, pinID, false)
	if err != nil {
		return fmt.Errorf("get a pin to check for the availability of a like: %w", err)
	}

	return p.isAvailablePinForViewingUser(ctx, pin, userID)
}

func isAvailableBatchPinForFixOnBoard(userID int, pins ...entity.Pin) error {
	for ind := range pins {
		if pins[ind].DeletedAt.Valid {
			return ErrPinDeleted
		}
		if !pins[ind].Public && pins[ind].Author.ID != userID {
			return ErrForbiddenAction
		}
	}
	return nil
}
