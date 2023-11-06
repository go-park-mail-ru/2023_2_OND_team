package board

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
)

func (bCase *boardUsecase) CreateNewBoard(ctx context.Context, newBoard dto.BoardData) (int, error) {
	if !bCase.isValidBoardTitle(newBoard.Title) {
		return 0, ErrInvalidBoardTitle
	}
	if err := bCase.checkIsValidTagTitles(newBoard.TagTitles); err != nil {
		return 0, fmt.Errorf("%s: %w", err.Error(), ErrInvalidTagTitles)
	}
	bCase.sanitizer.Sanitize(newBoard.Description)

	newBoardID, err := bCase.boardRepo.CreateBoard(ctx, entity.Board{
		AuthorID:    newBoard.AuthorID,
		Title:       newBoard.Title,
		Description: newBoard.Description,
		Public:      newBoard.Public,
	}, newBoard.TagTitles)

	if err != nil {
		return 0, fmt.Errorf("create new board: %w", err)
	}
	return newBoardID, nil
}
