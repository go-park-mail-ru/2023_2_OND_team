package realtime

import (
	"context"
	"errors"
	"fmt"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
)

var ErrUnknownTypeObject = errors.New("unknown type")

const (
	_topicChat         = "chat"
	_topicNotification = "notification"
)

type RealTimeClient interface {
	Subscribe(ctx context.Context, nameChans []string) (<-chan Pack, error)
	Publish(ctx context.Context, chanName string, object any) error
}

type Pack struct {
	Body any
	Err  error
}

type realtimeClient struct {
	client rt.RealTimeClient
	topic  string
}

func NewRealTimeChatClient(client rt.RealTimeClient) realtimeClient {
	return realtimeClient{
		client: client,
		topic:  _topicChat,
	}
}

func NewRealTimeNotificationClient(client rt.RealTimeClient) realtimeClient {
	return realtimeClient{
		client: client,
		topic:  _topicNotification,
	}
}

func (r realtimeClient) Publish(ctx context.Context, chanName string, object any) error {
	pubMsg := &rt.PublishMessage{
		Channel: &rt.Channel{
			Topic: r.topic,
			Name:  chanName,
		},
		Message: &rt.Message{},
	}

	switch body := object.(type) {
	case *rt.Message_Object:
		pubMsg.Message.Body = body
	case *rt.Message_Content:
		pubMsg.Message.Body = body
	default:
		return ErrUnknownTypeObject
	}

	_, err := r.client.Publish(ctx, pubMsg)
	if err != nil {
		return fmt.Errorf("publish as a realtime client: %w", err)
	}
	return nil
}

func (r realtimeClient) Subscribe(ctx context.Context, nameChans []string) (<-chan Pack, error) {
	chans := &rt.Channels{
		Chans: make([]*rt.Channel, len(nameChans)),
	}

	for _, name := range nameChans {
		chans.Chans = append(chans.Chans, &rt.Channel{Topic: r.topic, Name: name})
	}

	subClient, err := r.client.Subscribe(ctx, chans)
	if err != nil {
		return nil, fmt.Errorf("subscribe as a realtime client: %w", err)
	}

	ch := make(chan Pack)
	go runServeSubscribeClient(subClient, ch)

	return ch, nil
}

func runServeSubscribeClient(client rt.RealTime_SubscribeClient, ch chan<- Pack) {
	defer close(ch)

	var (
		mes *rt.Message
		err error
	)

	for {
		mes, err = client.Recv()
		if err != nil {
			ch <- Pack{Err: err}
			return
		}

		ch <- Pack{Body: mes.GetBody()}
	}
}
