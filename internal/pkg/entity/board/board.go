package board

import (
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
)

//go:generate easyjson board.go

//easyjson:json
type Board struct {
	ID          int        `json:"id,omitempty" example:"15"`
	AuthorID    int        `json:"author_id,omitempty"`
	Title       string     `json:"title" example:"Sunny places"`
	Description string     `json:"description" example:"Sunny places desc"`
	Public      bool       `json:"public" example:"true"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

//easyjson:json
type BoardWithContent struct {
	BoardInfo  Board
	PinsNumber int
	Pins       []string
	TagTitles  []string
}

func (b *Board) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if b != nil {
		b.Title = sanitizer.Sanitize(censor.Sanitize(b.Title))
		b.Description = sanitizer.Sanitize(censor.Sanitize(b.Description))
	}
}
