package message

import (
	"context"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
)

type Repository interface {
	AddNewMessage(ctx context.Context, mes *entity.Message) error
	GetMessages(ctx context.Context, chat entity.Chat) ([]entity.Message, error)
	UpdateContentMessage(ctx context.Context, messageID int, newContent string) error
	DelMessage(ctx context.Context, messageID int) error
}
