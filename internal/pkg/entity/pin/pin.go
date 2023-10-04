package pin

import "time"

type Pin struct {
	ID              int       `json:"id" example:"55"`
	AuthorID        int       `json:"-" example:"23"`
	Picture         string    `json:"picture" example:"pinspire/imgs/image.png"`
	Title           string    `json:"-" example:"Nature's beauty"`
	Description     string    `json:"-" example:"about face"`
	PublicationTime time.Time `json:"-"`
} //@name Pin
