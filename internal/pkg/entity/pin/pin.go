package pin

import "time"

type Pin struct {
	ID              int
	Picture         string
	Title           string
	Description     string
	PublicationTime time.Time
}
