package pin

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
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

// unused
func (p *pinRepoPG) GetSortedUserLikedPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, error) {
	rows, err := p.db.Query(ctx, SelectUserLikedPinsLimit, userID, minID, maxID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins: %w", count, err)
	}

	pins := make([]entity.Pin, 0, count)
	pin := entity.Pin{}
	for rows.Next() {
		err := rows.Scan(&pin.ID, &pin.Picture, &pin.Public)
		if err != nil {
			return pins, fmt.Errorf("scan to receive %d pins: %w", count, err)
		}
		pins = append(pins, pin)
	}

	return pins, nil
}
