package pin

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
)

type Pin struct {
	ID          int         `json:"id" example:"55"`
	Author      *user.User  `json:"author,omitempty" example:"23"`
	Picture     string      `json:"picture" example:"pinspire/imgs/image.png"`
	Title       pgtype.Text `json:"title" example:"Nature's beauty"`
	Description pgtype.Text `json:"description" example:"about face"`
	Public      bool        `json:"public"`

	Tags      []Tag `json:"tags,omitempty"`
	CountLike int   `json:"count_likes"`

	DeletedAt pgtype.Timestamptz `json:"-"`
} //@name Pin

func (p *Pin) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if p != nil {
		if p.Author != nil {
			p.Author.Sanitize(sanitizer, censor)
		}
		p.Title = pgtype.Text{
			String: sanitizer.Sanitize(censor.Sanitize(p.Title.String)),
			Valid:  p.Title.Valid,
		}
		p.Description = pgtype.Text{
			String: sanitizer.Sanitize(censor.Sanitize(p.Description.String)),
			Valid:  p.Description.Valid,
		}
	}
}

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
