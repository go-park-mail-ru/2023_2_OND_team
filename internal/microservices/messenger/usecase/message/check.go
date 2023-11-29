package message

import (
	"context"
	"fmt"
)

func (m *messageCase) isAvailableForChanges(ctx context.Context, userID, mesID int) (bool, error) {
	mes, err := m.repo.GetMessageByID(ctx, mesID)
	if err != nil {
		return false, fmt.Errorf("get message for check available: %w", err)
	}
	if mes.From == userID {
		return true, nil
	}
	return false, nil
}
