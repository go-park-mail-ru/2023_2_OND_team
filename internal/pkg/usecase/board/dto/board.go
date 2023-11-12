package board

import "github.com/microcosm-cc/bluemonday"

type BoardData struct {
	ID          int
	Title       string
	Description string
	AuthorID    int
	Public      bool
	TagTitles   []string
}

type UserBoard struct {
	BoardID     int      `json:"board_id" example:"15"`
	AuthorID    int      `json:"author_id,omitempty" example:"15"`
	Title       string   `json:"title" example:"Sunny places"`
	Description string   `json:"description" example:"Sunny places"`
	CreatedAt   string   `json:"created_at" example:"08.10.2020"`
	PinsNumber  int      `json:"pins_number" example:"10"`
	Pins        []string `json:"pins" example:"['/upload/pin/pic1', '/upload/pin/pic2']"`
	TagTitles   []string `json:"tags" example:"['flowers', 'sunrise']"`
}

func (uBoard *UserBoard) Sanitize(sanitizer *bluemonday.Policy) {
	sanitizer.Sanitize(uBoard.Title)
	sanitizer.Sanitize(uBoard.Description)
	for id, title := range uBoard.TagTitles {
		uBoard.TagTitles[id] = sanitizer.Sanitize(title)
	}
}
