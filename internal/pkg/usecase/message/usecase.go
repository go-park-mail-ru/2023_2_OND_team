package message

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	mesRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/message"
)

var ErrNoAccess = errors.New("there is no access to perform this action")

//go:generate mockgen -destination=./mock/message_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	SendMessage(ctx context.Context, mes *entity.Message) (int, error)
	GetMessagesFromChat(ctx context.Context, chat entity.Chat, count, lastID int) (feed []entity.Message, newLastID int, err error)
	UpdateContentMessage(ctx context.Context, userID int, mes *entity.Message) error
	DeleteMessage(ctx context.Context, userID, mesID int) error
	GetMessage(ctx context.Context, messageID int) (*entity.Message, error)
}

type messageCase struct {
	repo mesRepo.Repository
}

func New(repo mesRepo.Repository) *messageCase {
	return &messageCase{repo}
}

func (m *messageCase) SendMessage(ctx context.Context, mes *entity.Message) (int, error) {
	return m.repo.AddNewMessage(ctx, mes)
}

func (m *messageCase) GetMessagesFromChat(ctx context.Context, chat entity.Chat, count, lastID int) (feed []entity.Message, newLastID int, err error) {
	feed, err = m.repo.GetMessages(ctx, chat, count, lastID)
	if err != nil {
		err = fmt.Errorf("get message: %w", err)
	}
	if len(feed) != 0 {
		newLastID = feed[len(feed)-1].ID
	}
	return
}

func (m *messageCase) UpdateContentMessage(ctx context.Context, userID int, mes *entity.Message) error {
	if ok, err := m.isAvailableForChanges(ctx, userID, mes.ID); err != nil {
		return fmt.Errorf("update message: %w", err)
	} else if !ok {
		return ErrNoAccess
	}
	return m.repo.UpdateContentMessage(ctx, mes.ID, mes.Content.String)
}

func (m *messageCase) DeleteMessage(ctx context.Context, userID, mesID int) error {
	if ok, err := m.isAvailableForChanges(ctx, userID, mesID); err != nil {
		return fmt.Errorf("delete message: %w", err)
	} else if !ok {
		return ErrNoAccess
	}
	return m.repo.DelMessage(ctx, mesID)
}

func (m *messageCase) GetMessage(ctx context.Context, messageID int) (*entity.Message, error) {
	return m.repo.GetMessageByID(ctx, messageID)
}
