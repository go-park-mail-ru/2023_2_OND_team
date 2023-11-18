package board

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
)

func (bCase *boardUsecase) CreateNewBoard(ctx context.Context, newBoard entity.Board, tagTitles []string) (int, error) {
	newBoardID, err := bCase.boardRepo.CreateBoard(ctx, entity.Board{
		AuthorID:    newBoard.AuthorID,
		Title:       newBoard.Title,
		Description: newBoard.Description,
		Public:      newBoard.Public,
	}, tagTitles)

	if err != nil {
		return 0, fmt.Errorf("create new board: %w", err)
	}
	return newBoardID, nil
}
