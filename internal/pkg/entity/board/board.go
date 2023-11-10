package board

import (
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type Board struct {
	ID          int       `json:"id,omitempty" example:"15"`
	AuthorID    int       `json:"-"`
	Title       string    `json:"title" example:"Sunny places"`
	Description string    `json:"description" example:"Sunny places desc"`
	Public      bool      `json:"public" example:"true"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
}

func (board *Board) Sanitize(sanitizer *bluemonday.Policy) {
	sanitizer.Sanitize(board.Title)
	sanitizer.Sanitize(board.Description)
}
