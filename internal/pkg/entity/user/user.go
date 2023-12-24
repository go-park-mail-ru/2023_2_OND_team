package user

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
	"github.com/jackc/pgx/v5/pgtype"
)

//go:generate easyjson user.go

const UserUnknown = -1

//easyjson:json
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

func (u *User) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if u != nil {
		u.Username = sanitizer.Sanitize(censor.Sanitize(u.Username))
		u.Email = sanitizer.Sanitize(censor.Sanitize(u.Email))
		u.Name = pgtype.Text{
			String: sanitizer.Sanitize(censor.Sanitize(u.Name.String)),
			Valid:  u.Name.Valid,
		}
		u.Surname = pgtype.Text{
			String: sanitizer.Sanitize(censor.Sanitize(u.Surname.String)),
			Valid:  u.Surname.Valid,
		}
		u.AboutMe = pgtype.Text{
			String: sanitizer.Sanitize(censor.Sanitize(u.AboutMe.String)),
			Valid:  u.AboutMe.Valid,
		}
	}
}

//easyjson:json
type SubscriptionUser struct {
	ID                      int    `json:"id"`
	Username                string `json:"username"`
	Avatar                  string `json:"avatar"`
	HasSubscribeFromCurUser bool   `json:"is_subscribed"`
}

func (u *SubscriptionUser) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if u != nil {
		u.Username = sanitizer.Sanitize(censor.Sanitize(u.Username))
	}
}

type SubscriptionOpts struct {
	UserID int
	Count  int
	LastID int
	Filter string
}
