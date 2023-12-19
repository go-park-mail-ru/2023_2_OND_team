package notification

import notify "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/notification"

type Option interface {
	apply(*notificationClient)
}

type funcOption func(*notificationClient)

func (f funcOption) apply(cl *notificationClient) {
	f(cl)
}

func Register(notifier notify.Notifier) Option {
	return funcOption(func(cl *notificationClient) {
		cl.notifiers[notifier.Type()] = notifier
	})
}
