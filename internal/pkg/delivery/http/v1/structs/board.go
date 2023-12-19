package structs

import (
	"fmt"

	errHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/errors"
)

//go:generate easyjson board.go

// data for board creation/update
//
//easyjson:json
type BoardData struct {
	Title       *string  `json:"title" example:"new board"`
	Description *string  `json:"description" example:"long desc"`
	Public      *bool    `json:"public" example:"true"`
	Tags        []string `json:"tags" example:"['blue', 'car']"`
}

// board view for delivery layer
//
//easyjson:json
type CertainBoard struct {
	ID          int      `json:"board_id" example:"22"`
	AuthorID    int      `json:"author_id" example:"22"`
	Title       string   `json:"title" example:"new board"`
	Description string   `json:"description" example:"long desc"`
	CreatedAt   string   `json:"created_at" example:"07-11-2023"`
	PinsNumber  int      `json:"pins_number" example:"12"`
	Pins        []string `json:"pins" example:"['/pic1', '/pic2']"`
	Tags        []string `json:"tags" example:"['love', 'green']"`
}

//easyjson:json
type CertainBoardWithUsername struct {
	ID             int      `json:"board_id" example:"22"`
	AuthorID       int      `json:"author_id" example:"22"`
	AuthorUsername string   `json:"author_username" example:"Bob"`
	Title          string   `json:"title" example:"new board"`
	Description    string   `json:"description" example:"long desc"`
	CreatedAt      string   `json:"created_at" example:"07-11-2023"`
	PinsNumber     int      `json:"pins_number" example:"12"`
	Pins           []string `json:"pins" example:"['/pic1', '/pic2']"`
	Tags           []string `json:"tags" example:"['love', 'green']"`
}

//easyjson:json
type DeletePinFromBoard struct {
	PinID int `json:"pin_id" example:"22"`
}

func (data *BoardData) Validate() error {
	if data.Title == nil || *data.Title == "" {
		return errHTTP.ErrInvalidBoardTitle
	}
	if data.Description == nil {
		data.Description = new(string)
		*data.Description = ""
	}
	if data.Public == nil {
		return errHTTP.ErrEmptyPubOpt
	}
	if !isValidBoardTitle(*data.Title) {
		return errHTTP.ErrInvalidBoardTitle
	}
	if err := checkIsValidTagTitles(data.Tags); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), errHTTP.ErrInvalidTagTitles)
	}
	return nil
}
