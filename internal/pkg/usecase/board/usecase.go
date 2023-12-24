package board

import (
	"context"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
	userRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

//go:generate mockgen -destination=./mock/board_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	CreateNewBoard(ctx context.Context, newBoard entity.Board, tagTitles []string) (int, error)
	GetBoardsByUsername(ctx context.Context, username string) ([]entity.BoardWithContent, error)
	GetCertainBoard(ctx context.Context, boardID int) (entity.BoardWithContent, string, error)
	GetBoardInfoForUpdate(ctx context.Context, boardID int) (entity.Board, []string, error)
	UpdateBoardInfo(ctx context.Context, updatedBoard entity.Board, tagTitles []string) error
	DeleteCertainBoard(ctx context.Context, boardID int) error
	FixPinsOnBoard(ctx context.Context, boardID int, pinIds []int, userID int) error
	DeletePinFromBoard(ctx context.Context, boardID, pinID int) error
	CheckAvailabilityFeedPinCfgOnBoard(ctx context.Context, cfg pin.FeedPinConfig, userID int, isAuth bool) error
}

type boardUsecase struct {
	log       *logger.Logger
	boardRepo boardRepo.Repository
	userRepo  userRepo.Repository
}

func New(logger *logger.Logger, boardRepo boardRepo.Repository, userRepo userRepo.Repository) *boardUsecase {
	return &boardUsecase{log: logger, boardRepo: boardRepo, userRepo: userRepo}
}
