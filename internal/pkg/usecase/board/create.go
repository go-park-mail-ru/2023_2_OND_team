package board

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
)

func (bCase *BoardUsecase) CreateNewBoard(ctx context.Context, newBoard dto.BoardData) error {
	if !bCase.isValidBoardTitle(newBoard.Title) {
		return ErrInvalidBoardTitle
	}
	if err := bCase.checkIsValidTagTitles(newBoard.TagTitles); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidTagTitles)
	}

	err := bCase.boardRepo.CreateBoard(ctx, entity.Board{
		AuthorID:    newBoard.AuthorID,
		Title:       newBoard.Title,
		Description: newBoard.Description,
		Public:      newBoard.Public,
	}, newBoard.TagTitles)

	if err != nil {
		return fmt.Errorf("create new board usecase: %w", err)
	}
	return nil
}
