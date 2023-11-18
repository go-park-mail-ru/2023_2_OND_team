package board

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
)

func (bCase *boardUsecase) DeleteCertainBoard(ctx context.Context, boardID int) error {
	boardAuthorID, err := bCase.boardRepo.GetBoardAuthorByBoardID(ctx, boardID)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return ErrNoSuchBoard
		default:
			return fmt.Errorf("delete certain board: %w", err)
		}
	}

	currUserID, loggedIn := ctx.Value(auth.KeyCurrentUserID).(int)
	if !(loggedIn && currUserID == boardAuthorID) {
		return ErrNoAccess
	}

	err = bCase.boardRepo.DeleteBoardByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("delete certain board: %w", err)
	}

	return nil
}
