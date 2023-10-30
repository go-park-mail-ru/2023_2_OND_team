package board

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
)

func (bCase *BoardUsecase) CreateNewBoard(ctx context.Context, createBoardObj CreateBoard) error {
	if !isValidBoardTitle(createBoardObj.Title) {
		return fmt.Errorf("invalid board title")
	}
	if !isValidTagTitles(createBoardObj.TagTitles) {
		return fmt.Errorf("invalid board tag titles")
	}

	err := bCase.BoardRepo.CreateBoard(ctx, board.Board{
		AuthorID:    createBoardObj.AuthorID,
		Title:       createBoardObj.Title,
		Description: createBoardObj.Description,
		Public:      createBoardObj.Public,
	}, createBoardObj.PinIDs, createBoardObj.TagTitles)

	if err != nil {
		return fmt.Errorf("create new board usecase: %w", err)
	}
	return nil
}
