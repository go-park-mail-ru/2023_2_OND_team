package board

import (
	"context"

	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/microcosm-cc/bluemonday"
)

//go:generate mockgen -destination=./mock/board_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	CreateNewBoard(ctx context.Context, newBoard dto.BoardData) (int, error)
	GetBoardsByUsername(ctx context.Context, username string) ([]dto.UserBoard, error)
	GetCertainBoard(ctx context.Context, boardID int) (dto.UserBoard, error)
	UpdateBoardInfo(ctx context.Context, updatedData dto.BoardData) error
	DeleteCertainBoard(ctx context.Context, boardID int) error
	FixPinsOnBoard(ctx context.Context, boardID int, pinIds []int, userID int) error
}

type boardUsecase struct {
	log       *logger.Logger
	boardRepo boardRepo.Repository
	userRepo  userRepo.Repository
	sanitizer *bluemonday.Policy
}

func New(logger *logger.Logger, boardRepo boardRepo.Repository, userRepo userRepo.Repository, sanitizer *bluemonday.Policy) *boardUsecase {
	return &boardUsecase{log: logger, boardRepo: boardRepo, userRepo: userRepo, sanitizer: sanitizer}
}
