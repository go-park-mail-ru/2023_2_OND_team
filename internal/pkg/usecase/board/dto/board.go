package board

type BoardData struct {
	Title       string   `json:"title" example:"Sunny places"`
	Description string   `json:"description" example:"long description"`
	AuthorID    int      `json:"author_id,omitempty" example:"45"`
	Public      bool     `json:"public" example:"true"`
	TagTitles   []string `json:"tags" example:"['flowers', 'sunrise']"`
} //@name NewBoardData

type UserBoard struct {
	BoardID     int      `json:"board_id" example:"15"`
	Title       string   `json:"title" example:"Sunny places"`
	Description string   `json:"description,omitempty" example:"Sunny places"`
	CreatedAt   string   `json:"created_at" example:"08.10.2020"`
	PinsNumber  int      `json:"pins_number" example:"10"`
	Pins        []string `json:"pins" example:"['/upload/pin/pic1', '/upload/pin/pic2']"`
	TagTitles   []string `json:"tags,omitempty" example:"['flowers', 'sunrise']"`
} //@name UserBoard
