package message

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/metadata"

	mess "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	messMS "github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/messenger/delivery/grpc"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
)

var ErrNoAccess = errors.New("there is no access to perform this action")

//go:generate mockgen -destination=./mock/message_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	GetUserChatsWithOtherUsers(ctx context.Context, userID, count, lastID int) (entity.FeedUserChats, int, error)
	SendMessage(ctx context.Context, userID int, mes *entity.Message) (int, error)
	GetMessagesFromChat(ctx context.Context, userID int, chat entity.Chat, count, lastID int) (feed []entity.Message, newLastID int, err error)
	UpdateContentMessage(ctx context.Context, userID int, mes *entity.Message) error
	DeleteMessage(ctx context.Context, userID, mesID int) error
	GetMessage(ctx context.Context, userID int, messageID int) (*entity.Message, error)
}

type messageCase struct {
	client mess.MessengerClient
}

func New(repo mess.MessengerClient) *messageCase {
	return &messageCase{repo}
}

func (m *messageCase) SendMessage(ctx context.Context, userID int, mes *entity.Message) (int, error) {
	msgID, err := m.client.SendMessage(setAuthenticatedMetadataCtx(ctx, userID), &mess.Message{
		UserFrom: int64(userID),
		UserTo:   int64(mes.To),
		Content:  mes.Content.String,
	})
	if err != nil {
		return 0, fmt.Errorf("send message by grpc client")
	}
	return int(msgID.GetId()), nil
}

func (m *messageCase) GetMessagesFromChat(ctx context.Context, userID int, chat entity.Chat, count, lastID int) (feed []entity.Message, newLastID int, err error) {
	feedMsg, err := m.client.MessageFromChat(setAuthenticatedMetadataCtx(ctx, userID), &mess.FeedMessageRequest{
		Chat: &mess.Chat{
			UserID1: int64(chat[0]),
			UserID2: int64(chat[1]),
		},
		Count:  int64(count),
		LastID: int64(count),
	})
	if err != nil {
		err = fmt.Errorf("get message by : %w", err)
	}
	if feedMsg == nil {
		return nil, 0, err
	}

	return convertFeedMessage(feedMsg), int(feedMsg.LastID), nil
}

func (m *messageCase) UpdateContentMessage(ctx context.Context, userID int, mes *entity.Message) error {
	if _, err := m.client.UpdateMessage(setAuthenticatedMetadataCtx(ctx, userID), &mess.Message{
		Id: &mess.MsgID{
			Id: int64(mes.ID),
		},
		Content: mes.Content.String,
	}); err != nil {
		return fmt.Errorf("update messege by grpc client")
	}
	return nil
}

func (m *messageCase) DeleteMessage(ctx context.Context, userID, mesID int) error {
	if _, err := m.client.DeleteMessage(setAuthenticatedMetadataCtx(ctx, userID), &mess.MsgID{Id: int64(mesID)}); err != nil {
		return fmt.Errorf("delete messege by grpc client")
	}
	return nil
}

func (m *messageCase) GetMessage(ctx context.Context, userID int, messageID int) (*entity.Message, error) {
	mes, err := m.client.GetMessage(setAuthenticatedMetadataCtx(ctx, userID), &mess.MsgID{Id: int64(messageID)})
	if err != nil {
		return nil, fmt.Errorf("get message by grpc client")
	}
	return &entity.Message{
		ID:   int(mes.GetId().Id),
		From: int(mes.GetUserFrom()),
		To:   int(mes.GetUserTo()),
		Content: pgtype.Text{
			String: mes.Content,
			Valid:  true,
		},
	}, nil
}

func (m *messageCase) GetUserChatsWithOtherUsers(ctx context.Context, userID, count, lastID int) (entity.FeedUserChats, int, error) {
	feed, err := m.client.UserChatsWithOtherUsers(setAuthenticatedMetadataCtx(ctx, userID), &mess.FeedChatRequest{
		Count:  int64(count),
		LastID: int64(lastID),
	})
	var errRes error
	if err != nil {
		errRes = fmt.Errorf("get user chats by grpc client: %w", err)
	}
	if feed == nil {
		return nil, 0, errRes
	}
	return convertFeedChat(feed), int(feed.GetLastID()), errRes
}

func setAuthenticatedMetadataCtx(ctx context.Context, userID int) context.Context {
	return metadata.AppendToOutgoingContext(ctx, messMS.AuthenticatedMetadataKey, strconv.FormatInt(int64(userID), 10))
}
