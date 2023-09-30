package pin

import "time"

type Pin struct {
	ID              int       `json:"id" example:"55"`
	AuthorID        int       `json:"authorId" example:"23"`
	Picture         string    `json:"picture" example:"pinspire/imgs/image.png"`
	Title           string    `json:"title" example:"Nature's beauty"`
	Description     string    `json:"description" example:"about face"`
	PublicationTime time.Time `json:"publicationTime"`
} //@name Pin
