package notification

import (
	"context"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/notification"
)

type M map[string]string

type TypeNotifier interface {
	Type() entity.NotifyType
}

type Notifier interface {
	TypeNotifier

	ChannelNameForPublishWithData(ctx context.Context, entityID int) (string, M, error)
	ChannelsNameForSubscribe(ctx context.Context, userID int) ([]string, error)
	MessageNotify(data M) (*entity.NotifyMessage, error)
}

type NotifyBuilder interface {
	TypeNotifier

	BuildNotifyMessage(data any) (*entity.NotifyMessage, error)
}
