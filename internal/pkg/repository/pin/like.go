package pin

import (
	"context"
	"fmt"
)

func (p *pinRepoPG) SetLike(ctx context.Context, pinID, userID int) error {
	_, err := p.db.Exec(ctx, InsertLikePinFromUser, pinID, userID)
	if err != nil {
		return fmt.Errorf("insert like to pin from user in storage: %w", err)
	}
	return nil
}

func (p *pinRepoPG) DelLike(ctx context.Context, pinID, userID int) error {
	_, err := p.db.Exec(ctx, DeleteLikePinFromUser, pinID, userID)
	if err != nil {
		return fmt.Errorf("delete like to pin from user in storage: %w", err)
	}
	return nil
}

func (p *pinRepoPG) GetCountLikeByPinID(ctx context.Context, pinID int) (int, error) {
	row := p.db.QueryRow(ctx, SelectCountLikePin, pinID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get count like by pin id: %w", err)
	}
	return count, nil
}
