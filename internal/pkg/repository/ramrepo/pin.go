package ramrepo

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type ramPinRepo struct {
	db *sql.DB
}

func NewRamPinRepo(db *sql.DB) *ramPinRepo {
	return &ramPinRepo{db}
}

func (r *ramPinRepo) GetNPinsAfterID(ctx context.Context, count int, afterPinID int) ([]pin.Pin, error) {
	return nil, nil
}

func (r *ramPinRepo) GetAuthorPin(ctx context.Context, pinID int) (*user.User, error) {
	return nil, nil
}
