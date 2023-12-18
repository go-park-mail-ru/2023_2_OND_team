package grpc

import (
	"context"

	mess "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/messenger/usecase/message"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const AuthenticatedMetadataKey = "user_id"

type MessengerServer struct {
	mess.UnimplementedMessengerServer

	log         *logger.Logger
	messageCase message.Usecase
}

func New(log *logger.Logger, msgCase message.Usecase) MessengerServer {
	return MessengerServer{
		log:         log,
		messageCase: msgCase,
	}
}

func (m MessengerServer) UserChatsWithOtherUsers(ctx context.Context, r *mess.FeedChatRequest) (*mess.FeedChat, error) {
	userID := ctx.Value(auth.KeyCurrentUserID).(int)

	feed, lastID, err := m.messageCase.GetUserChatsWithOtherUsers(ctx, userID, int(r.GetCount()), int(r.GetLastID()))
	if err != nil {
		m.log.Error(err.Error())
	}

	return &mess.FeedChat{Chats: convertFeedChat(feed), LastID: int64(lastID)}, nil
}

func (m MessengerServer) SendMessage(ctx context.Context, msg *mess.Message) (*mess.MsgID, error) {
	userID := ctx.Value(auth.KeyCurrentUserID).(int)

	msgID, err := m.messageCase.SendMessage(ctx, &entity.Message{
		From:    userID,
		To:      int(msg.UserTo),
		Content: pgtype.Text{String: msg.GetContent(), Valid: true},
	})

	if err != nil {
		m.log.Error(err.Error())
		return nil, status.Error(codes.Internal, "send message error")
	}

	return &mess.MsgID{Id: int64(msgID)}, nil
}

func (m MessengerServer) MessageFromChat(ctx context.Context, r *mess.FeedMessageRequest) (*mess.FeedMessage, error) {
	feed, lastID, err := m.messageCase.GetMessagesFromChat(ctx, entity.Chat{int(r.Chat.GetUserID1()), int(r.Chat.GetUserID2())},
		int(r.GetCount()), int(r.GetLastID()))
	if err != nil {
		m.log.Error(err.Error())
	}

	return &mess.FeedMessage{
		Messages: convertFeedMessage(feed),
		LastID:   int64(lastID),
	}, nil
}

func (m MessengerServer) UpdateMessage(ctx context.Context, msg *mess.Message) (*empty.Empty, error) {
	userID := ctx.Value(auth.KeyCurrentUserID).(int)

	err := m.messageCase.UpdateContentMessage(ctx, userID, &entity.Message{
		ID:      int(msg.Id.GetId()),
		Content: pgtype.Text{String: msg.Content, Valid: true},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "update message error")
	}

	return &empty.Empty{}, nil
}

func (m MessengerServer) DeleteMessage(ctx context.Context, msgID *mess.MsgID) (*empty.Empty, error) {
	userID := ctx.Value(auth.KeyCurrentUserID).(int)

	err := m.messageCase.DeleteMessage(ctx, userID, int(msgID.GetId()))
	if err != nil {
		m.log.Error(err.Error())
		return nil, status.Error(codes.Internal, "delete message")
	}

	return &empty.Empty{}, nil
}

func (m MessengerServer) GetMessage(ctx context.Context, msgID *mess.MsgID) (*mess.Message, error) {
	msg, err := m.messageCase.GetMessage(ctx, int(msgID.GetId()))
	if err != nil {
		m.log.Error(err.Error())
		return nil, status.Error(codes.Internal, "get message")
	}

	resMsg := &mess.Message{
		Id:       msgID,
		UserFrom: int64(msg.From),
		UserTo:   int64(msg.To),
		Content:  msg.Content.String,
	}

	if msg.DeletedAt.Valid {
		resMsg.DeletedAt = timestamppb.New(msg.DeletedAt.Time)
	}

	return resMsg, nil
}
