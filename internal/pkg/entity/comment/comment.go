package comment

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
)

//go:generate easyjson comment.go
//easyjson:json
type Comment struct {
	ID      int         `json:"id"`
	Author  *user.User  `json:"author"`
	PinID   int         `json:"pinID"`
	Content pgtype.Text `json:"content"`
}

func (c *Comment) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if c != nil {
		if c.Author != nil {
			c.Author.Sanitize(sanitizer, censor)
		}
		c.Content = pgtype.Text{
			String: sanitizer.Sanitize(censor.Sanitize(c.Content.String)),
			Valid:  c.Content.Valid,
		}
	}
}
