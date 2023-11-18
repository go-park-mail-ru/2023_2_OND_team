package pin

import (
	"context"
	"fmt"
)

func (p *pinRepoPG) IsAvailableToUserAsContributorBoard(ctx context.Context, pinID, userID int) (bool, error) {
	row := p.db.QueryRow(ctx, SelectCheckAvailability, pinID, userID)
	var check bool
	err := row.Scan(&check)
	if err != nil {
		return false, fmt.Errorf("check available pin for user at the storage layer: %w", err)
	}
	return check, nil
}
