package session

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
)

type Repository interface {
	AddSession(ctx context.Context, session *session.Session) error
	GetSessionByKey(ctx context.Context, key string) (*session.Session, error)
	DeleteSessionByKey(ctx context.Context, key string) error
	DeleteAllSessionForUser(ctx context.Context, userID int) error
}
