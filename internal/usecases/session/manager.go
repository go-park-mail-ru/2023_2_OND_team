package session

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/session"
)

type SessionManager struct {
	repo session.Repository
}

func New(repo session.Repository) *SessionManager {
	return &SessionManager{repo}
}

func (sm *SessionManager) AddSession(userID int) error {
	return nil
}

func (sm *SessionManager) GetUserBySessionKey(sessionKey string) (*user.User, error) {
	return nil, nil
}
