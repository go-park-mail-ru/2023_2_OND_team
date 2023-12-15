package notification

import (
	"context"
	"fmt"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/notification"
)

func (n *notificationClient) NotifyCommentLeftOnPin(ctx context.Context, commentID int) error {
	notifier, ok := n.notifiers[entity.NotifyComment]
	if !ok {
		n.log.Error(ErrNotifierNotRegistered.Error())
		return ErrNotifierNotRegistered
	}

	chanName, data, err := notifier.ChannelNameForPublishWithData(ctx, commentID)
	if err != nil {
		n.log.Error(err.Error())
		return fmt.Errorf("notify comment left on pin: %w", err)
	}

	err = n.client.Publish(ctx, chanName, &rt.Message_Content{
		Content: &rt.EventMap{
			Type: int64(entity.NotifyComment),
			M:    data,
		},
	})
	if err != nil {
		n.log.Error(err.Error())
		return fmt.Errorf("publish to client: %w", err)
	}

	return nil
}
