package comment

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

//go:generate easyjson comment.go
//easyjson:json
type Comment struct {
	ID      int         `json:"id"`
	Author  *user.User  `json:"author"`
	PinID   int         `json:"pinID"`
	Content pgtype.Text `json:"content"`
}
