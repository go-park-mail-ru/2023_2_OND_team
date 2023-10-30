package board

import (
	"context"

	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type CreateBoard struct {
	Title       string   `json:"title" example:"Sunny places"`
	Description string   `json:"description" example:"long description"`
	AuthorID    int      `json:"author_id" example:"45"`
	Public      bool     `json:"public" example:"true"`
	PinIDs      []int    `json:"pin_ids" example:"[1, 2, 3]"`
	TagTitles   []string `json:"tags" example:"['flowers', 'sunrise']"`
} //@name Board

type Usecase interface {
	CreateNewBoard(ctx context.Context, createBoardObj CreateBoard) error
	// GetOwnBoards()
	// GetUserBoards()
}

type BoardUsecase struct {
	log       *logger.Logger
	BoardRepo boardRepo.Repository
}

func New(logger *logger.Logger, boardRepo boardRepo.Repository) *BoardUsecase {
	return &BoardUsecase{log: logger, BoardRepo: boardRepo}
}
