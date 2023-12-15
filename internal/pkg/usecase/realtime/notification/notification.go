package notification

import (
	"context"
	"errors"
	"fmt"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/notification"
	notify "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/notification"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrNotifierNotRegistered = errors.New("notifier with this type not registered")

type Usecase interface {
	NotifyCommentLeftOnPin(ctx context.Context, commentID int) error
}

type notificationClient struct {
	client    realtime.RealTimeClient
	log       *logger.Logger
	notifiers map[entity.NotifyType]notify.Notifier
}

func New(cl realtime.RealTimeClient, log *logger.Logger, opts ...Option) *notificationClient {
	client := &notificationClient{
		client:    cl,
		log:       log,
		notifiers: make(map[entity.NotifyType]notify.Notifier),
	}

	for _, opt := range opts {
		opt.apply(client)
	}

	return client
}

func (n *notificationClient) SubscribeOnAllNotifications(ctx context.Context, userID int) (<-chan *entity.NotifyMessage, error) {
	setChans := make(map[string]struct{})
	for t, notifier := range n.notifiers {
		nameChans, err := notifier.ChannelsNameForSubscribe(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("receiving name channels for subscribe on %s notifier: %w", entity.TypeString(t), err)
		}

		for _, name := range nameChans {
			setChans[name] = struct{}{}
		}
	}

	uniqChans := make([]string, 0, len(setChans))

	for nameChan := range setChans {
		uniqChans = append(uniqChans, nameChan)
	}

	chanPack, err := n.client.Subscribe(ctx, uniqChans)
	if err != nil {
		return nil, fmt.Errorf("subscribe on all notifications: %w", err)
	}

	chanNotifyMsg := make(chan *entity.NotifyMessage)

	go n.pipelineNotify(chanPack, chanNotifyMsg)

	return chanNotifyMsg, nil
}

func (n *notificationClient) pipelineNotify(chRecv <-chan realtime.Pack, chSend chan<- *entity.NotifyMessage) {
	defer close(chSend)

	for pack := range chRecv {
		if pack.Err != nil {
			chSend <- entity.NewNotifyMessageWithError(pack.Err)
			return
		}

		notifyData, ok := pack.Body.(*rt.Message_Content)
		if !ok {
			chSend <- entity.NewNotifyMessageWithError(realtime.ErrUnknownTypeObject)
			return
		}

		notifier, ok := n.notifiers[entity.NotifyType(notifyData.Content.GetType())]
		if !ok {
			chSend <- entity.NewNotifyMessageWithError(ErrNotifierNotRegistered)
			return
		}

		msg, err := notifier.MessageNotify(notifyData.Content.GetM())
		if err != nil {
			chSend <- entity.NewNotifyMessageWithError(err)
			return
		}

		chSend <- msg
	}
}
