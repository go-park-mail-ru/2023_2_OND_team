package pin

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

func (p *pinCase) SetLikeFromUser(ctx context.Context, pinID, userID int) (int, error) {
	if err := p.isAvailablePinForSetLike(ctx, pinID, userID); err != nil {
		return 0, fmt.Errorf("set like from user: %w", err)
	}
	return p.repo.SetLike(ctx, pinID, userID)
}

func (p *pinCase) DeleteLikeFromUser(ctx context.Context, pinID, userID int) (int, error) {
	return p.repo.DelLike(ctx, pinID, userID)
}

func (p *pinCase) CheckUserHasSetLike(ctx context.Context, pinID, userID int) (bool, error) {
	return p.repo.IsSetLike(ctx, pinID, userID)
}

// unused
func (p *pinCase) SelectUserLikedPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, int, int) {
	pins, err := p.repo.GetSortedUserPins(ctx, userID, count, minID, maxID)
	if err != nil {
		p.log.Error(err.Error())
	}
	if len(pins) == 0 {
		return []entity.Pin{}, minID, maxID
	}
	return pins, pins[len(pins)-1].ID, pins[0].ID
}
