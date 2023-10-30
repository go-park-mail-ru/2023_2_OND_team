package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Repository interface {
	AddNewUser(ctx context.Context, user *user.User) error
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUsernameAndAvatarByID(ctx context.Context, userID int) (username string, avatar string, err error)
	EditUserAvatar(ctx context.Context, userID int, avatar string) error
	GetAllUserData(ctx context.Context, userID int) (*user.User, error)
	EditUserInfo(ctx context.Context, userID int, updateFields S) error
}

type S map[string]any

type userRepoPG struct {
	db *pgxpool.Pool
}

func NewUserRepoPG(db *pgxpool.Pool) *userRepoPG {
	return &userRepoPG{db}
}

func (u *userRepoPG) AddNewUser(ctx context.Context, user *user.User) error {
	_, err := u.db.Exec(ctx, InsertNewUser, user.Username, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("add a new profile in storage: %w", err)
	}
	return nil
}

func (u *userRepoPG) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	row := u.db.QueryRow(ctx, SelectAuthByUsername, username)
	user := &user.User{Username: username}
	err := row.Scan(&user.ID, &user.Password, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("getting a user from storage: %w", err)
	}
	return user, nil
}

func (u *userRepoPG) GetUsernameAndAvatarByID(ctx context.Context, userID int) (username string, avatar string, err error) {
	row := u.db.QueryRow(ctx, SelectUsernameAndAvatar, userID)
	err = row.Scan(&username, &avatar)
	if err != nil {
		return "", "", fmt.Errorf("getting a username from storage by id: %w", err)
	}
	return
}

func (u *userRepoPG) EditUserAvatar(ctx context.Context, userID int, avatar string) error {
	_, err := u.db.Exec(ctx, UpdateAvatarProfile, avatar, userID)
	if err != nil {
		return fmt.Errorf("edit user avatar: %w", err)
	}
	return nil
}

func (u *userRepoPG) GetAllUserData(ctx context.Context, userID int) (*user.User, error) {
	row := u.db.QueryRow(ctx, SelectUserDataExceptPassword, userID)
	user := &user.User{ID: userID}
	err := row.Scan(&user.Username, &user.Email, &user.Avatar, &user.Name, &user.Surname, &user.AboutMe)
	if err != nil {
		return nil, fmt.Errorf("get user info by id in storage: %w", err)
	}
	return user, nil
}

func (u *userRepoPG) EditUserInfo(ctx context.Context, userID int, updateFields S) error {
	sqlRow, args, err := sq.Update("profile").
		SetMap(updateFields).
		Where("id = ?", userID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("build sql query row: %w", err)
	}

	_, err = u.db.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("update user info in the storage: %w", err)
	}
	return nil
}
