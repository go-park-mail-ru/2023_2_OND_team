package v1

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type HandlerHTTP struct {
	log         *logger.Logger
	authCase    auth.Usecase
	userCase    user.Usecase
	pinCase     pin.Usecase
	boardCase   board.Usecase
	subCase     subscription.Usecase
	messageCase message.Usecase
}

func New(log *logger.Logger, hub UsecaseHub) *HandlerHTTP {
	return &HandlerHTTP{
		log:         log,
		authCase:    hub.AuhtCase,
		userCase:    hub.UserCase,
		pinCase:     hub.PinCase,
		boardCase:   hub.BoardCase,
		subCase:     hub.SubscriptionCase,
		messageCase: hub.MessageCase,
	}
}

type UsecaseHub struct {
	AuhtCase         auth.Usecase
	UserCase         user.Usecase
	PinCase          pin.Usecase
	BoardCase        board.Usecase
	SubscriptionCase subscription.Usecase
	MessageCase      message.Usecase
}
