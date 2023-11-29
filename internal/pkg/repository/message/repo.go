package message

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/internal/pgtype"
)

//go:generate mockgen -destination=./mock/message_mock.go -package=mock -source=repo.go Repository
type Repository interface {
	GetMessageByID(ctx context.Context, mesID int) (*entity.Message, error)
	AddNewMessage(ctx context.Context, mes *entity.Message) (int, error)
	GetMessages(ctx context.Context, chat entity.Chat, count, lastID int) ([]entity.Message, error)
	UpdateContentMessage(ctx context.Context, messageID int, newContent string) error
	DelMessage(ctx context.Context, messageID int) error
	GetUserChats(ctx context.Context, userID, count, lastID int) (entity.FeedUserChats, error)
}

type messageRepo struct {
	db pgtype.PgxPoolIface
}

func NewMessageRepo(db pgtype.PgxPoolIface) *messageRepo {
	return &messageRepo{db}
}

func (m *messageRepo) GetMessageByID(ctx context.Context, mesID int) (*entity.Message, error) {
	message := &entity.Message{ID: mesID}
	err := m.db.QueryRow(ctx, SelectMessageByID, mesID).Scan(&message.From, &message.To, &message.Content)
	if err != nil {
		return nil, fmt.Errorf("get message by id from storage: %w", err)
	}
	return message, nil
}

func (m *messageRepo) AddNewMessage(ctx context.Context, mes *entity.Message) (int, error) {
	err := m.db.QueryRow(ctx, InsertMessage, mes.From, mes.To, mes.Content).Scan(&mes.ID)
	if err != nil {
		return 0, fmt.Errorf("add new message in storage: %w", err)
	}
	return mes.ID, nil
}

func (m *messageRepo) GetMessages(ctx context.Context, chat entity.Chat, count, lastID int) ([]entity.Message, error) {
	rows, err := m.db.Query(ctx, SelectMessageFromChat, lastID, chat[0], chat[1], count)
	if err != nil {
		return nil, fmt.Errorf("get message for chat from storage: %w", err)
	}

	message := entity.Message{}
	messages := make([]entity.Message, 0, count)
	for rows.Next() {
		err = rows.Scan(&message.ID, &message.From, &message.To, &message.Content)
		if err != nil {
			return messages, fmt.Errorf("scan selected message: %w", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (m *messageRepo) UpdateContentMessage(ctx context.Context, messageID int, newContent string) error {
	_, err := m.db.Exec(ctx, UpdateMessageContent, newContent, messageID)
	if err != nil {
		return fmt.Errorf("update content message in storage: %w", err)
	}
	return nil
}

func (m *messageRepo) DelMessage(ctx context.Context, messageID int) error {
	_, err := m.db.Exec(ctx, UpdateMessageStatusToDeleted, messageID)
	if err != nil {
		return fmt.Errorf("delete message from storage: %w", err)
	}
	return nil
}

func (m *messageRepo) GetUserChats(ctx context.Context, userID, count, lastID int) (entity.FeedUserChats, error) {
	rows, err := m.db.Query(ctx, SelectUserChats, userID, lastID, count)
	if err != nil {
		return nil, fmt.Errorf("get user chats in storage: %w", err)
	}
	defer rows.Close()

	feed := make(entity.FeedUserChats, 0, count)
	chat := entity.ChatWithUser{}
	for rows.Next() {
		if err = rows.Scan(&chat.MessageLastID, &chat.WichWhomChat.ID,
			&chat.WichWhomChat.Username, &chat.WichWhomChat.Avatar); err != nil {

			return feed, fmt.Errorf("scan chat with user for feed: %w", err)
		}
		feed = append(feed, chat)
	}
	return feed, nil
}
