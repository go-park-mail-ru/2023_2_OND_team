package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/crypto"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

const SessionLifeTime = 30 * 24 * time.Hour

const lenSessionKey = 16

var ErrExpiredSession = errors.New("session lifetime expired")

type SessionManager struct {
	log  *logger.Logger
	repo repo.Repository
}

func New(log *logger.Logger, repo repo.Repository) *SessionManager {
	return &SessionManager{log, repo}
}

func (sm *SessionManager) CreateNewSessionForUser(ctx context.Context, userID int) (*session.Session, error) {
	sessionKey, err := crypto.NewRandomStr(lenSessionKey)
	if err != nil {
		return nil, fmt.Errorf("session key generation for new session: %w", err)
	}

	session := &session.Session{
		Key:    sessionKey,
		UserID: userID,
		Expire: time.Now().UTC().Add(SessionLifeTime),
	}
	err = sm.repo.AddSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("creating a new session by the session manager: %w", err)
	}
	return session, nil
}

func (sm *SessionManager) GetUserIDBySessionKey(ctx context.Context, sessionKey string) (int, error) {
	session, err := sm.repo.GetSessionByKey(ctx, sessionKey)
	if err != nil {
		return 0, fmt.Errorf("getting a session by a manager: %w", err)
	}
	if time.Now().UTC().After(session.Expire) {
		return 0, ErrExpiredSession
	}
	return session.UserID, nil
}

func (sm *SessionManager) DeleteUserSession(ctx context.Context, key string) error {
	return sm.repo.DeleteSessionByKey(ctx, key)
}
