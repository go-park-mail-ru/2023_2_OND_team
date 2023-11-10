package ramrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	repository "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
)

type ramPinRepo struct {
	db *sql.DB
}

func NewRamPinRepo(db *sql.DB) *ramPinRepo {
	return &ramPinRepo{db}
}

func (r *ramPinRepo) GetSortedNewNPins(ctx context.Context, count, minID, maxID int) ([]pin.Pin, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, picture FROM pin WHERE id > $1 ORDER BY id LIMIT $2;", maxID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins after %d: %w", count, maxID, err)
	}

	pins := make([]pin.Pin, 0, count)
	pin := pin.Pin{}
	for rows.Next() {
		err := rows.Scan(&pin.ID, &pin.Picture)
		if err != nil {
			return pins, fmt.Errorf("scan to receive %d pins after %d: %w", count, minID, err)
		}
		pins = append(pins, pin)
	}

	return pins, nil
}

func (r *ramPinRepo) GetAuthorPin(ctx context.Context, pinID int) (*user.User, error) {
	return nil, ErrMethodUnimplemented
}

func (r *ramPinRepo) AddNewPin(ctx context.Context, pin *pin.Pin) error {
	return ErrMethodUnimplemented
}

func (r *ramPinRepo) DeletePin(ctx context.Context, pinID, userID int) error {
	return ErrMethodUnimplemented
}

func (r *ramPinRepo) SetLike(ctx context.Context, pinID, userID int) (int, error) {
	return 0, ErrMethodUnimplemented
}

func (r *ramPinRepo) DelLike(ctx context.Context, pinID, userID int) (int, error) {
	return 0, ErrMethodUnimplemented
}

func (r *ramPinRepo) EditPinTags(ctx context.Context, pinID, userID int, titlePins []string) error {
	return ErrMethodUnimplemented
}

func (r *ramPinRepo) EditPin(ctx context.Context, pinID int, updateData repository.S, titleTags []string) error {
	return ErrMethodUnimplemented
}

func (r *ramPinRepo) GetPinByID(ctx context.Context, pinID int, revealAuthor bool) (*pin.Pin, error) {
	return nil, ErrMethodUnimplemented
}

func (r *ramPinRepo) IsAvailableToUserAsContributorBoard(ctx context.Context, pinID, userID int) (bool, error) {
	return false, ErrMethodUnimplemented
}

func (r *ramPinRepo) GetCountLikeByPinID(ctx context.Context, pinID int) (int, error) {
	return 0, ErrMethodUnimplemented
}

func (r *ramPinRepo) GetTagsByPinID(ctx context.Context, pinID int) ([]pin.Tag, error) {
	return nil, ErrMethodUnimplemented
}

func (r *ramPinRepo) GetSortedUserPins(ctx context.Context, userID, count, minID, maxID int) ([]pin.Pin, error) {
	return nil, ErrMethodUnimplemented
}

func (r *ramPinRepo) IsSetLike(ctx context.Context, pinID, userID int) (bool, error) {
	return false, ErrMethodUnimplemented
}

func (r *ramPinRepo) GetBatchPinByID(ctx context.Context, pinID []int) ([]pin.Pin, error) {
	return nil, ErrMethodUnimplemented
}

func (r *ramPinRepo) GetSortedUserLikedPins(ctx context.Context, userID, count, minID, maxID int) ([]pin.Pin, error) {
	return nil, ErrMethodUnimplemented
}
