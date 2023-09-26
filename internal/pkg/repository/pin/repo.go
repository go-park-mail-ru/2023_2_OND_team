package pin

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Repository interface {
	GetNPinsAfterID(ctx context.Context, count int, afterPinID int) ([]pin.Pin, error)
	GetAuthorPin(ctx context.Context, pinID int) (*user.User, error)
}
