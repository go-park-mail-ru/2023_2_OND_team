package chat

import (
	"context"
	"fmt"
	"strconv"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type EventMessageObjectID struct {
	Type      string
	MessageID int
	Err       error
}

func makeErrEventMessageObjectID(err error) EventMessageObjectID {
	return EventMessageObjectID{Err: err}
}

type Usecase interface {
	PublishNewMessage(ctx context.Context, userToWhom, msgID int) error
	PublishUpdateMessage(ctx context.Context, userToWhom, msgID int) error
	PublishDeleteMessage(ctx context.Context, userToWhom, msgID int) error
	SubscribeUserToAllChats(ctx context.Context, userID int) (<-chan EventMessageObjectID, error)
}

type realtimeCase struct {
	client realtime.RealTimeClient
	log    *logger.Logger
}

func New(client realtime.RealTimeClient, log *logger.Logger) *realtimeCase {
	return &realtimeCase{client, log}
}

func (r *realtimeCase) PublishNewMessage(ctx context.Context, userToWhom, msgID int) error {
	err := r.publishMessage(ctx, userToWhom, msgID, rt.EventType_EV_CREATE)
	if err != nil {
		r.log.Error(err.Error())
		return fmt.Errorf("publish new message: %w", err)
	}
	return nil
}

func (r *realtimeCase) PublishUpdateMessage(ctx context.Context, userToWhom, msgID int) error {
	err := r.publishMessage(ctx, userToWhom, msgID, rt.EventType_EV_UPDATE)
	if err != nil {
		r.log.Error(err.Error())
		return fmt.Errorf("publish update message: %w", err)
	}
	return nil
}

func (r *realtimeCase) PublishDeleteMessage(ctx context.Context, userToWhom, msgID int) error {
	err := r.publishMessage(ctx, userToWhom, msgID, rt.EventType_EV_DELETE)
	if err != nil {
		r.log.Error(err.Error())
		return fmt.Errorf("publish delete message: %w", err)
	}
	return nil
}

func (r *realtimeCase) SubscribeUserToAllChats(ctx context.Context, userToWhom int) (<-chan EventMessageObjectID, error) {
	chPack, err := r.client.Subscribe(ctx, []string{strconv.Itoa(userToWhom)})
	if err != nil {
		return nil, fmt.Errorf("subscribe user to all chats: %w", err)
	}

	chanEvMsg := make(chan EventMessageObjectID)
	go r.receiveFromSubClient(ctx, chPack, chanEvMsg)

	return chanEvMsg, nil
}

func (r *realtimeCase) receiveFromSubClient(ctx context.Context, subClient <-chan realtime.Pack, chanEvMsg chan<- EventMessageObjectID) {
	defer close(chanEvMsg)

	for pack := range subClient {
		if pack.Err != nil {
			chanEvMsg <- makeErrEventMessageObjectID(pack.Err)
			return
		}

		msg, ok := pack.Body.(*rt.Message_Object)
		if !ok {
			chanEvMsg <- makeErrEventMessageObjectID(realtime.ErrUnknownTypeObject)
			return
		}

		evMsgID := EventMessageObjectID{MessageID: int(msg.Object.GetId())}
		switch msg.Object.GetType() {
		case rt.EventType_EV_CREATE:
			evMsgID.Type = "create"
		case rt.EventType_EV_UPDATE:
			evMsgID.Type = "update"
		case rt.EventType_EV_DELETE:
			evMsgID.Type = "delete"
		}

		chanEvMsg <- evMsgID
	}
}

func (r *realtimeCase) publishMessage(ctx context.Context, userID, msgID int, t rt.EventType) error {
	return r.client.Publish(ctx, strconv.Itoa(userID), &rt.Message_Object{
		Object: &rt.EventObject{
			Type: t,
			Id:   int64(msgID),
		},
	})
}
