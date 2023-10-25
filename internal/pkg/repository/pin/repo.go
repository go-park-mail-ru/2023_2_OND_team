package pin

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Repository interface {
	GetSortedNPinsAfterID(ctx context.Context, count int, afterPinID int) ([]pin.Pin, error)
	GetAuthorPin(ctx context.Context, pinID int) (*user.User, error)
}

type pinRepoPG struct {
	db *pgxpool.Pool
}

func NewPinRepoPG(db *pgxpool.Pool) *pinRepoPG {
	return &pinRepoPG{db}
}

func (p *pinRepoPG) GetSortedNPinsAfterID(ctx context.Context, count int, afterPinID int) ([]pin.Pin, error) {
	rows, err := p.db.Query(ctx, SelectAfterIdWithLimit, afterPinID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins after %d: %w", count, afterPinID, err)
	}

	pins := make([]pin.Pin, 0, count)
	pin := pin.Pin{}
	for rows.Next() {
		err := rows.Scan(&pin.ID, &pin.Picture)
		if err != nil {
			return pins, fmt.Errorf("scan to receive %d pins after %d: %w", count, afterPinID, err)
		}
		pins = append(pins, pin)
	}

	return pins, nil
}

func (r *pinRepoPG) GetAuthorPin(ctx context.Context, pinID int) (*user.User, error) {
	return nil, errors.New("unimplemented")
}
