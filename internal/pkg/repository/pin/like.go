package pin

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (p *pinRepoPG) SetLike(ctx context.Context, pinID, userID int) (int, error) {
	row := p.db.QueryRow(ctx, InsertLikePinFromUser, pinID, userID)
	var currCountLike int
	err := row.Scan(&currCountLike)
	if err != nil {
		return 0, fmt.Errorf("insert like to pin from user in storage: %w", err)
	}
	return currCountLike + 1, nil
}

func (p *pinRepoPG) DelLike(ctx context.Context, pinID, userID int) (int, error) {
	row := p.db.QueryRow(ctx, DeleteLikePinFromUser, pinID, userID)
	var currCountLike int
	err := row.Scan(&currCountLike)
	if err != nil {
		return 0, fmt.Errorf("delete like to pin from user in storage: %w", err)
	}
	return currCountLike - 1, nil
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

func (p *pinRepoPG) IsSetLike(ctx context.Context, pinID, userID int) (bool, error) {
	row := p.db.QueryRow(ctx, SelectCheckSetLike, pinID, userID)
	var check int
	err := row.Scan(&check)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("is set like to pin in storage: %w", err)
	}
	return true, nil
}
