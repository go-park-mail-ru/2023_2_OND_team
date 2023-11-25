package v1

import (
	pb "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/api"
	roll "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/service"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type ServerGRPC struct {
	pb.UnimplementedRollServiceServer
	log         *logger.Logger
	rollService roll.Service
}

func NewServerGRPC(log *logger.Logger, rollService roll.Service) *ServerGRPC {
	return &ServerGRPC{
		log:         log,
		rollService: rollService,
	}
}

