package board

import "time"

type Board struct {
	ID          int
	AuthorID    int
	Title       string
	Description string
	Public      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
