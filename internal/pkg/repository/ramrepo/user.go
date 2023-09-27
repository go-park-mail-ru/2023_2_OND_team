package ramrepo

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type ramUserRepo struct {
	db *sql.DB
}

func NewRamUserRepo(db *sql.DB) *ramUserRepo {
	return &ramUserRepo{db}
}

func (r *ramUserRepo) AddNewUser(ctx context.Context, user *entity.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3);", user.Username, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("adding a new user to the ram repository: %w", err)
	}
	return nil
}

func (r *ramUserRepo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, password, email FROM users WHERE username = $1;", username)
	user := &entity.User{Username: username}
	err := row.Scan(&user.ID, &user.Password, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("getting a user from storage: %w", err)
	}
	return user, nil
}
