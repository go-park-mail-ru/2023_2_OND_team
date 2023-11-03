package board

import (
	"context"

	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/microcosm-cc/bluemonday"
)

type Usecase interface {
	CreateNewBoard(ctx context.Context, createBoardObj dto.CreateBoard) error
	GetBoardsByUsername(ctx context.Context, username string) ([]dto.GetUserBoard, error)
	GetCertainBoardByID(ctx context.Context, boardID int) (dto.GetUserBoard, error)
	// CheckBoardContributor(ctx context.Context, userID int) (bool, error)
}

type BoardUsecase struct {
	log       *logger.Logger
	boardRepo boardRepo.Repository
	userRepo  userRepo.Repository
	sanitizer *bluemonday.Policy
}

func New(logger *logger.Logger, boardRepo boardRepo.Repository, userRepo userRepo.Repository, sanitizer *bluemonday.Policy) *BoardUsecase {
	return &BoardUsecase{log: logger, boardRepo: boardRepo, userRepo: userRepo, sanitizer: sanitizer}
}
