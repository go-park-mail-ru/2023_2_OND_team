package user

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	ID       int         `json:"id,omitempty" example:"123"`
	Username string      `json:"username" example:"Green"`
	Name     pgtype.Text `json:"name,omitempty" example:"Peter"`
	Surname  pgtype.Text `json:"surname,omitempty" example:"Green"`
	Email    string      `json:"email,omitempty" example:"digital@gmail.com"`
	Avatar   string      `json:"avatar" example:"pinspire.online/avatars/avatar.jpg"`
	AboutMe  pgtype.Text `json:"about_me,omitempty"`
	Password string      `json:"password,omitempty" example:"pass123"`
} // @name User

type SubscriptionUser struct {
	ID                      int    `json:"id"`
	Username                string `json:"username"`
	Avatar                  string `json:"avatar"`
	HasSubscribeFromCurUser bool   `json:"is_subscribed"`
}

type SubscriptionOpts struct {
	UserID int
	Count  int
	LastID int
	Filter string
}
