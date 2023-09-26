package session

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
)

type Repository interface {
	CreateSessionForUser(ctx context.Context, userID int) (*session.Session, error)
	GetUserIDBySessionKey(ctx context.Context, key string) (int, error)
}
