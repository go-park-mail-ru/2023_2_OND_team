package service

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type Service struct {
	log      *logger.Logger
	userCase *user.Usecase
	pinCase  *pin.Usecase
	sm       *session.SessionManager
}

func New(log *logger.Logger, sm *session.SessionManager, user *user.Usecase, pin *pin.Usecase) *Service {
	return &Service{
		log:      log,
		userCase: user,
		pinCase:  pin,
		sm:       sm,
	}
}
