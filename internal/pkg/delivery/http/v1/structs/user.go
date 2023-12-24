package structs

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
)

//go:generate easyjson user.go

//easyjson:json
type UserInfo struct {
	ID           int    `json:"id" example:"123"`
	Username     string `json:"username" example:"Snapshot"`
	Avatar       string `json:"avatar" example:"/pic1"`
	Name         string `json:"name" example:"Bob"`
	Surname      string `json:"surname" example:"Dylan"`
	About        string `json:"about" example:"Cool guy"`
	IsSubscribed bool   `json:"is_subscribed" example:"true"`
	SubsCount    int    `json:"subscribers" example:"23"`
}

func (u *UserInfo) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if u != nil {
		u.Username = sanitizer.Sanitize(censor.Sanitize(u.Username))
		u.Name = sanitizer.Sanitize(censor.Sanitize(u.Name))
		u.Surname = sanitizer.Sanitize(censor.Sanitize(u.Surname))
		u.About = sanitizer.Sanitize(censor.Sanitize(u.About))
	}
}

//easyjson:json
type ProfileInfo struct {
	ID        int    `json:"id" example:"1"`
	Username  string `json:"username" example:"baobab"`
	Avatar    string `json:"avatar" example:"/pic1"`
	SubsCount int    `json:"subscribers" example:"12"`
}

func (p *ProfileInfo) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if p != nil {
		p.Username = sanitizer.Sanitize(censor.Sanitize(p.Username))
	}
}
