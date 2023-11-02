package user

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	ID       int         `json:"id" example:"123"`
	Username string      `json:"username" example:"Green"`
	Name     pgtype.Text `json:"name" example:"Peter"`
	Surname  pgtype.Text `json:"surname" example:"Green"`
	Email    string      `json:"email,omitempty" example:"digital@gmail.com"`
	Avatar   string      `json:"avatar" example:"pinspire.online/avatars/avatar.jpg"`
	AboutMe  pgtype.Text `json:"about_me"`
	Password string      `json:"password,omitempty" example:"pass123"`
} // @name User
