package ramrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type ramPinRepo struct {
	db *sql.DB
}

func NewRamPinRepo(db *sql.DB) *ramPinRepo {
	return &ramPinRepo{db}
}

func (r *ramPinRepo) GetSortedNPinsAfterID(ctx context.Context, count int, afterPinID int) ([]pin.Pin, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, picture FROM pin WHERE id > $1 ORDER BY id LIMIT $2;", afterPinID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins after %d: %w", count, afterPinID, err)
	}

	pins := []pin.Pin{}
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

func (r *ramPinRepo) GetAuthorPin(ctx context.Context, pinID int) (*user.User, error) {
	return nil, errors.New("unimplemented")
}
