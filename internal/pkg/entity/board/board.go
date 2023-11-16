package board

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type Board struct {
	ID          int        `json:"id,omitempty" example:"15"`
	AuthorID    int        `json:"-"`
	Title       string     `json:"title" example:"Sunny places"`
	Description string     `json:"description" example:"Sunny places desc"`
	Public      bool       `json:"public" example:"true"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

type BoardWithContent struct {
	BoardInfo  Board
	PinsNumber int
	Pins       []string
	TagTitles  []string
}

func (b *Board) Sanitize(sanitizer *bluemonday.Policy) {
	sanitizer.Sanitize(b.Title)
	sanitizer.Sanitize(b.Description)
}

func (b *BoardWithContent) Sanitize(sanitizer *bluemonday.Policy) {
	b.BoardInfo.Sanitize(sanitizer)
	for id, title := range b.TagTitles {
		b.TagTitles[id] = sanitizer.Sanitize(title)
	}
}
