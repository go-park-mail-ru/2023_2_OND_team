package pin

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

type Repository interface {
	GetNPinsFromDate(ctx context.Context, count int, date time.Time) ([]pin.Pin, error)
}
