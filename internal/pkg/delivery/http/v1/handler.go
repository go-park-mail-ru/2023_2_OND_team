package v1

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/search"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type HandlerHTTP struct {
	log         *logger.Logger
	converter   converterHTTP
	authCase    auth.Usecase
	userCase    user.Usecase
	pinCase     pin.Usecase
	boardCase   board.Usecase
	subCase     subscription.Usecase
	searchCase  search.Usecase
	messageCase message.Usecase
	commentCase comment.Usecase
}

func New(log *logger.Logger, converter converterHTTP, hub UsecaseHub) *HandlerHTTP {
	return &HandlerHTTP{
		log:         log,
		converter:   converter,
		authCase:    hub.AuhtCase,
		userCase:    hub.UserCase,
		pinCase:     hub.PinCase,
		boardCase:   hub.BoardCase,
		subCase:     hub.SubscriptionCase,
		searchCase:  hub.SearchCase,
		messageCase: hub.MessageCase,
		commentCase: hub.CommentCase,
	}
}

type UsecaseHub struct {
	AuhtCase         auth.Usecase
	UserCase         user.Usecase
	PinCase          pin.Usecase
	BoardCase        board.Usecase
	SubscriptionCase subscription.Usecase
	SearchCase       search.Usecase
	MessageCase      message.Usecase
	CommentCase      comment.Usecase
}
