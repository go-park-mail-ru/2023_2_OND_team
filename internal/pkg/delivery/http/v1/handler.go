package v1

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type HandlerHTTP struct {
	log      *logger.Logger
	userCase user.Usecase
	pinCase  pin.Usecase
	imgCase  image.Usecase
	sm       session.SessionManager
}

func New(log *logger.Logger, sm session.SessionManager, user user.Usecase, pin pin.Usecase, img image.Usecase) *HandlerHTTP {
	return &HandlerHTTP{
		log:      log,
		userCase: user,
		pinCase:  pin,
		imgCase:  img,
		sm:       sm,
	}
}
