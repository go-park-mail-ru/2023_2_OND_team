package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Repository interface {
	AddNewUser(ctx context.Context, user *user.User) error
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUsernameAndAvatarByID(ctx context.Context, userID int) (username string, avatar string, err error)
}

type userRepoPG struct {
	db *pgxpool.Pool
}

func NewUserRepoPG(db *pgxpool.Pool) *userRepoPG {
	return &userRepoPG{db}
}

func (u *userRepoPG) AddNewUser(ctx context.Context, user *user.User) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction for add new user: %w", err)
	}

	row := tx.QueryRow(ctx, "INSERT INTO profile (email) VALUES ($1) RETURNING id;", user.Email)
	profileID := 0
	err = row.Scan(&profileID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("create a profile with the return of its id: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO auth (username, password, profile_id) VALUES ($1, $2, $3);", user.Username, user.Password, profileID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("linking credentials to a profile: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("confirmation of the transaction of adding a new user: %w", err)
	}
	return nil
}

func (u *userRepoPG) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	row := u.db.QueryRow(ctx, "SELECT profile_id, password, email FROM auth INNER JOIN profile ON profile.id = auth.profile_id WHERE username = $1;", username)
	user := &user.User{Username: username}
	err := row.Scan(&user.ID, &user.Password, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("getting a user from storage: %w", err)
	}
	return user, nil
}

func (r *userRepoPG) GetUsernameAndAvatarByID(ctx context.Context, userID int) (username string, avatar string, err error) {
	row := r.db.QueryRow(ctx, "SELECT username, avatar FROM auth INNER JOIN profile ON auth.profile_id = profile.id WHERE profile.id = $1;", userID)
	err = row.Scan(&username, &avatar)
	if err != nil {
		return "", "", fmt.Errorf("getting a username from storage by id: %w", err)
	}
	return
}
