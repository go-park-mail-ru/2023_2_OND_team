package pin

import "github.com/jackc/pgx/v5/pgtype"

type Pin struct {
	ID          int         `json:"id" example:"55"`
	AuthorID    int         `json:"author" example:"23"`
	Picture     string      `json:"picture" example:"pinspire/imgs/image.png"`
	Title       pgtype.Text `json:"-" example:"Nature's beauty"`
	Description pgtype.Text `json:"-" example:"about face"`
	Public      bool        `json:"public"`

	Tags      []Tag `json:"tags"`
	CountLike int   `json:"count_likes"`

	DeletedAt pgtype.Timestamptz `json:"-"`
} //@name Pin

func (p *Pin) SetTitle(title string) {
	p.Title = pgtype.Text{
		String: title,
		Valid:  true,
	}
}

func (p *Pin) SetDescription(description string) {
	p.Description = pgtype.Text{
		String: description,
		Valid:  true,
	}
}
