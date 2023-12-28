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
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/realtime/chat"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrNoAccess = errors.New("there is no access to perform this action")
var ErrRealTimeDisable = errors.New("realtime disable")
var ErrUnknowObj = errors.New("unknow object")

//go:generate mockgen -destination=./mock/message_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	GetUserChatsWithOtherUsers(ctx context.Context, userID, count, lastID int) (entity.FeedUserChats, int, error)
	SendMessage(ctx context.Context, userID int, mes *entity.Message) (int, error)
	GetMessagesFromChat(ctx context.Context, userID int, chat entity.Chat, count, lastID int) (feed []entity.Message, newLastID int, err error)
	UpdateContentMessage(ctx context.Context, userID int, mes *entity.Message) error
	DeleteMessage(ctx context.Context, userID int, mes *entity.Message) error
	GetMessage(ctx context.Context, userID int, messageID int) (*entity.Message, error)
	SubscribeUserToAllChats(ctx context.Context, userID int) (<-chan EventMessage, error)
}

type EventMessage struct {
	Type    string
	Message *entity.Message
	Err     error
}

func makeErrEventMessage(err error) EventMessage {
	return EventMessage{Err: err}
}

type messageCase struct {
	client           mess.MessengerClient
	realtimeChatCase chat.Usecase
	log              *logger.Logger
	realtimeIsEnable bool
}

func New(log *logger.Logger, cl mess.MessengerClient, rtChatCase chat.Usecase) *messageCase {
	m := &messageCase{
		client: cl,
		log:    log,
	}

	if rtChatCase != nil {
		m.realtimeChatCase = rtChatCase
		m.realtimeIsEnable = true
	}

	return m
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

	if m.realtimeIsEnable {
		go m.realtimeChatCase.PublishNewMessage(ctx, mes.To, int(msgID.GetId()))
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
		LastID: int64(lastID),
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

	if m.realtimeIsEnable {
		go m.realtimeChatCase.PublishUpdateMessage(ctx, mes.To, mes.ID)
	}

	return nil
}

func (m *messageCase) DeleteMessage(ctx context.Context, userID int, mes *entity.Message) error {
	if _, err := m.client.DeleteMessage(setAuthenticatedMetadataCtx(ctx, userID), &mess.MsgID{Id: int64(mes.ID)}); err != nil {
		return fmt.Errorf("delete messege by grpc client")
	}

	if m.realtimeIsEnable {
		go m.realtimeChatCase.PublishDeleteMessage(ctx, mes.To, mes.ID)
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

func (m *messageCase) SubscribeUserToAllChats(ctx context.Context, userID int) (<-chan EventMessage, error) {
	if !m.realtimeIsEnable {
		return nil, ErrRealTimeDisable
	}

	subClient, err := m.realtimeChatCase.SubscribeUserToAllChats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("subscribe: %w", err)
	}

	chanEvMsg := make(chan EventMessage)
	go m.receiveFromSubClient(ctx, userID, subClient, chanEvMsg)
	return chanEvMsg, nil
}

func (m *messageCase) receiveFromSubClient(ctx context.Context, userID int, subClient <-chan chat.EventMessageObjectID, chanEvMsg chan<- EventMessage) {
	defer close(chanEvMsg)

	var (
		evMsg EventMessage
		err   error
	)
	for msgObjID := range subClient {
		if msgObjID.Err != nil {
			chanEvMsg <- makeErrEventMessage(fmt.Errorf("receive from subcribtion client: %w", msgObjID.Err))
			return
		}

		evMsg = EventMessage{
			Type: msgObjID.Type,
		}

		evMsg.Message, err = m.GetMessage(ctx, userID, msgObjID.MessageID)
		if err != nil {
			m.log.Error(err.Error())
			evMsg.Err = err
		} else if evMsg.Type == "delete" {
			evMsg.Message.Content.String = ""
		}

		chanEvMsg <- evMsg
	}
}

func setAuthenticatedMetadataCtx(ctx context.Context, userID int) context.Context {
	return metadata.AppendToOutgoingContext(ctx, messMS.AuthenticatedMetadataKey, strconv.FormatInt(int64(userID), 10))
}
