package pin

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
)

type PinUpdateData struct {
	Title       *string
	Description *string
	Public      *bool
	Tags        []string
}

func (p *pinCase) EditPinByID(ctx context.Context, pinID, userID int, updateData *PinUpdateData) error {
	data := pin.S{}
	if updateData.Title != nil {
		data["title"] = *updateData.Title
	}
	if updateData.Description != nil {
		data["description"] = *updateData.Description
	}
	if updateData.Public != nil {
		data["public"] = *updateData.Public
	}
	err := p.repo.EditPin(ctx, pinID, userID, data, updateData.Tags)
	if err != nil {
		return fmt.Errorf("edit pin by id: %w", err)
	}
	return nil
}
