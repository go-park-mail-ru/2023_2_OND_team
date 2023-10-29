package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
	redis "github.com/redis/go-redis/v9"
)

var ErrMethodUnimplemented = errors.New("unimplemented")
var ErrExistsSession = errors.New("the session already exists")

type Repository interface {
	AddSession(ctx context.Context, session *session.Session) error
	GetSessionByKey(ctx context.Context, key string) (*session.Session, error)
	DeleteSessionByKey(ctx context.Context, key string) error
	DeleteAllSessionForUser(ctx context.Context, userID int) error
}

type sessionRepo struct {
	client *redis.Client
}

func NewSessionRepo(client *redis.Client) *sessionRepo {
	return &sessionRepo{client}
}

func (s *sessionRepo) AddSession(ctx context.Context, session *session.Session) error {
	res := s.client.SetNX(ctx, session.Key, session.UserID, time.Duration(session.Expire.Sub(time.Now().UTC())))
	if res.Err() != nil {
		return fmt.Errorf("add session in storage: %w", res.Err())
	}
	if !res.Val() {
		return ErrExistsSession
	}
	return nil
}

func (s *sessionRepo) GetSessionByKey(ctx context.Context, key string) (*session.Session, error) {
	res := s.client.Get(ctx, key)
	if res.Err() != nil {
		return nil, fmt.Errorf("get session by key from storage: %w", res.Err())
	}

	var err error
	sess := &session.Session{Key: key}
	sess.UserID, err = res.Int()
	if err != nil {
		return nil, fmt.Errorf("bad value for session in storage: %w", err)
	}
	return sess, nil
}

func (s *sessionRepo) DeleteSessionByKey(ctx context.Context, key string) error {
	res := s.client.Del(ctx, key)
	if res.Err() != nil {
		return fmt.Errorf("delete session by key from storage: %w", res.Err())
	}
	return nil
}

func (s *sessionRepo) DeleteAllSessionForUser(ctx context.Context, userID int) error {
	return ErrMethodUnimplemented
}
