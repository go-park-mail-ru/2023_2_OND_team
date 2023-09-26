package pin

import "time"

type Pin struct {
	ID              int
	AuthorID        int
	Picture         string
	Title           string
	Description     string
	PublicationTime time.Time
}
