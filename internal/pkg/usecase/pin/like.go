package pin

import (
	"context"
	"fmt"
)

func (p *pinCase) SetLikeFromUser(ctx context.Context, pinID, userID int) (int, error) {
	if err := p.isAvailablePinForSetLike(ctx, pinID, userID); err != nil {
		return 0, fmt.Errorf("set like from user: %w", err)
	}
	return p.repo.SetLike(ctx, pinID, userID)
}

func (p *pinCase) DeleteLikeFromUser(ctx context.Context, pinID, userID int) error {
	return p.repo.DelLike(ctx, pinID, userID)
}
