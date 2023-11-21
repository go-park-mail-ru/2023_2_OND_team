package v1

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/subscription"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type HandlerHTTP struct {
	log       *logger.Logger
	userCase  user.Usecase
	pinCase   pin.Usecase
	boardCase board.Usecase
	subCase   subscription.Usecase
	sm        session.SessionManager
}

func New(log *logger.Logger, sm session.SessionManager, user user.Usecase, pin pin.Usecase, board board.Usecase, sub subscription.Usecase) *HandlerHTTP {
	return &HandlerHTTP{
		log:       log,
		userCase:  user,
		pinCase:   pin,
		boardCase: board,
		subCase:   sub,
		sm:        sm,
	}
}
