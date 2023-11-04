package board

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
)

func (bCase *BoardUsecase) UpdateBoardInfo(ctx context.Context, updatedData dto.BoardData) error {
	boardAuthorID, err := bCase.boardRepo.GetBoardAuthorByBoardID(ctx, updatedData.ID)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return ErrNoSuchBoard
		default:
			return fmt.Errorf("update certain board: %w", err)
		}
	}

	currUserID, loggedIn := ctx.Value(auth.KeyCurrentUserID).(int)
	if !(loggedIn && currUserID == boardAuthorID) {
		return ErrNoAccess
	}

	if !bCase.isValidBoardTitle(updatedData.Title) {
		return ErrInvalidBoardTitle
	}
	if err := bCase.checkIsValidTagTitles(updatedData.TagTitles); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidTagTitles)
	}
	bCase.sanitizer.Sanitize(updatedData.Description)

	err = bCase.boardRepo.UpdateBoard(ctx, board.Board{
		ID:          updatedData.ID,
		Title:       updatedData.Title,
		Description: updatedData.Description,
		Public:      updatedData.Public,
	}, updatedData.TagTitles)
	if err != nil {
		return fmt.Errorf("update certain board: %w", err)
	}
	return nil
}
