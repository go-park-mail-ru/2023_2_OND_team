package message

import (
	"context"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
)

//go:generate mockgen -destination=./mock/message_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	SendMessage(ctx context.Context, mes *entity.Message) error
	GetMessagesFromChat(ctx context.Context, chat entity.Chat, lastID int) ([]entity.Message, error)
	UpdateContentMessage(ctx context.Context, mes *entity.Message) error
	DeleteMessage(ctx context.Context, mesID int) error
}
