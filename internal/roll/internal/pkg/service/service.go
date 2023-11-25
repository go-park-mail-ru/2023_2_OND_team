package roll

import (
	"context"

	roll "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/entity"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/repository"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/microcosm-cc/bluemonday"
)

type Service interface {
	FillRoll(ctx context.Context, ans []roll.RollAnswer) error
	HasUserFilledRoll(ctx context.Context, userID, rollID int) (bool, error)
	GetHistStat(ctx context.Context, rollID, questionID int) ([]roll.HistStatObj, error)
}

type rollService struct {
	log       *logger.Logger
	rollRepo  repo.Repository
	sanitizer *bluemonday.Policy
}

func New(logger *logger.Logger, repo repo.Repository, sanitizer *bluemonday.Policy) *rollService {
	return &rollService{log: logger, rollRepo: repo, sanitizer: sanitizer}
}

func (s *rollService) FillRoll(ctx context.Context, ans []roll.RollAnswer) error {
	return s.rollRepo.InsertRollAnswer(ctx, ans)
}

func (s *rollService) HasUserFilledRoll(ctx context.Context, userID, rollID int) (bool, error) {
	return s.rollRepo.CheckUserFilledRoll(ctx, userID, rollID)
}

func (s *rollService) GetHistStat(ctx context.Context, rollID, questionID int) ([]roll.HistStatObj, error) {
	return s.rollRepo.GetHistStat(ctx, rollID, questionID)
}
