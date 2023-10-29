package pin

import "context"

func (p *pinCase) SetLikeFromUser(ctx context.Context, pinID, userID int) error {
	return p.repo.SetLike(ctx, pinID, userID)
}

func (p *pinCase) DeleteLikeFromUser(ctx context.Context, pinID, userID int) error {
	return p.repo.DelLike(ctx, pinID, userID)
}
