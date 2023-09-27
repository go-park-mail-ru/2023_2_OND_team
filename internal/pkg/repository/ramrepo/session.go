package ramrepo

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
)

type ramSessionRepo struct {
	db *sql.DB
}

func NewRamSessionRepo(db *sql.DB) *ramSessionRepo {
	return &ramSessionRepo{db}
}

func (r *ramSessionRepo) AddSession(ctx context.Context, session *entity.Session) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO session (session_key, user_id, expire) VALUES ($1, $2, $3);",
		session.Key, session.UserID, session.Expire)
	if err != nil {
		return fmt.Errorf("save session in ram repository: %w", err)
	}
	return nil
}

func (r *ramSessionRepo) GetSessionByKey(ctx context.Context, key string) (*entity.Session, error) {
	row := r.db.QueryRowContext(ctx, "SELECT user_id, expire FROM session WHERE session_key = $1;", key)
	session := &entity.Session{Key: key}
	err := row.Scan(&session.UserID, &session.Expire)
	if err != nil {
		return nil, fmt.Errorf("get session from ram repository: %w", err)
	}
	return session, nil
}

func (r *ramSessionRepo) DeleteSessionByKey(ctx context.Context, key string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM session WHERE session_key = $1;", key)
	if err != nil {
		return fmt.Errorf("delete session by key from ram repository: %w", err)
	}
	return nil
}

func (r *ramSessionRepo) DeleteAllSessionForUser(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM session WHERE user_id = $1;", userID)
	if err != nil {
		return fmt.Errorf("delete session by user id from ram repository: %w", err)
	}
	return nil
}
