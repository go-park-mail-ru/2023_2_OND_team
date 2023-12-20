package search

import (
	"strings"
	"unicode"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/microcosm-cc/bluemonday"
)

//go:generate easyjson search.go

type Template string

func (t *Template) Validate() bool {
	if len(*t) == 0 || len(*t) > 40 {
		return false
	}
	for _, sym := range *t {
		if !(unicode.IsNumber(sym) || unicode.IsLetter(sym) || unicode.IsPunct(sym) || unicode.IsSpace(sym)) {
			return false
		}
	}
	return true
}

func (t *Template) GetSubStrings(sep string) []string {
	return strings.Split(string(*t), sep)
}

//easyjson:json
type BoardForSearch struct {
	BoardHeader board.Board
	PinsNumber  int      `json:"pins_number"`
	PreviewPins []string `json:"pins"`
}

//easyjson:json
type PinForSearch struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Picture string `json:"picture"`
	Likes   int    `json:"likes"`
}

//easyjson:json
type UserForSearch struct {
	ID                      int    `json:"id"`
	Username                string `json:"username"`
	Avatar                  string `json:"avatar"`
	SubsCount               int    `json:"subscribers"`
	HasSubscribeFromCurUser bool   `json:"is_subscribed"`
}

func (u *UserForSearch) Sanitize(sanitizer *bluemonday.Policy) {
	sanitizer.Sanitize(u.Username)
}

func (b *BoardForSearch) Sanitize(sanitizer *bluemonday.Policy) {
	sanitizer.Sanitize(b.BoardHeader.Title)
	sanitizer.Sanitize(b.BoardHeader.Description)
}

func (p *PinForSearch) Sanitize(sanitizer *bluemonday.Policy) {
	sanitizer.Sanitize(p.Title)
}

type SearchOpts struct {
	General GeneralOpts
	SortBy  string
}

type GeneralOpts struct {
	Template   Template
	SortOrder  string
	CurrUserID int
	Count      int
	Offset     int
}
