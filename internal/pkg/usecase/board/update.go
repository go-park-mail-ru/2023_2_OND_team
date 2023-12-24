package board

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	repoBoard "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
)

func (bCase *boardUsecase) GetBoardInfoForUpdate(ctx context.Context, boardID int) (entity.Board, []string, error) {
	boardAuthorID, err := bCase.boardRepo.GetBoardAuthorByBoardID(ctx, boardID)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return entity.Board{}, nil, ErrNoSuchBoard
		default:
			return entity.Board{}, nil, fmt.Errorf("get certain board info for update: %w", err)
		}
	}

	boardContributors, err := bCase.boardRepo.GetContributorsByBoardID(ctx, boardID)
	if err != nil {
		return entity.Board{}, nil, fmt.Errorf("get certain board info for update: %w", err)
	}

	boardContributorsIDs := make([]int, 0, len(boardContributors))

	for _, contributor := range boardContributors {
		boardContributorsIDs = append(boardContributorsIDs, contributor.ID)
	}

	var hasAccess bool
	currUserID, loggedIn := ctx.Value(auth.KeyCurrentUserID).(int)
	if loggedIn && (currUserID == boardAuthorID || isContributor(boardContributorsIDs, currUserID)) {
		hasAccess = true
	}

	board, tagTitles, err := bCase.boardRepo.GetBoardInfoForUpdate(ctx, boardID, hasAccess)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return entity.Board{}, nil, ErrNoSuchBoard
		default:
			return entity.Board{}, nil, fmt.Errorf("get certain board: %w", err)
		}
	}

	return board, tagTitles, nil
}

func (bCase *boardUsecase) UpdateBoardInfo(ctx context.Context, updatedBoard entity.Board, tagTitles []string) error {
	boardAuthorID, err := bCase.boardRepo.GetBoardAuthorByBoardID(ctx, updatedBoard.ID)
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

	err = bCase.boardRepo.UpdateBoard(ctx, board.Board{
		ID:          updatedBoard.ID,
		Title:       updatedBoard.Title,
		Description: updatedBoard.Description,
		Public:      updatedBoard.Public,
	}, tagTitles)
	if err != nil {
		return fmt.Errorf("update certain board: %w", err)
	}
	return nil
}

func (b *boardUsecase) FixPinsOnBoard(ctx context.Context, boardID int, pinIds []int, userID int) error {
	role, err := b.boardRepo.RoleUserHaveOnThisBoard(ctx, boardID, userID)
	if err != nil {
		return fmt.Errorf("get role for fix pins: %w", err)
	}
	if role&(repoBoard.Author|repoBoard.ContributorForAdding) == 0 {
		return ErrNoAccess
	}

	err = b.boardRepo.AddPinsOnBoard(ctx, boardID, pinIds)
	if err != nil {
		return fmt.Errorf("fix pins on board: %w", err)
	}
	return nil
}
