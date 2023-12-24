package board

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
)

func (bCase *boardUsecase) GetBoardsByUsername(ctx context.Context, username string) ([]entity.BoardWithContent, error) {
	userID, err := bCase.userRepo.GetUserIdByUsername(ctx, username)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return nil, ErrInvalidUsername
		default:
			return nil, fmt.Errorf("get user id by username in get boards usecase: %w", err)
		}
	}

	var isAuthor bool
	currUserID, loggedIn := ctx.Value(auth.KeyCurrentUserID).(int)
	if loggedIn && currUserID == userID {
		isAuthor = true
	}

	contributorBoardsIDs, err := bCase.boardRepo.GetContributorBoardsIDs(ctx, currUserID)
	if err != nil {
		return nil, fmt.Errorf("get contributor boards in get boards by username usecase: %w", err)
	}

	boards, err := bCase.boardRepo.GetBoardsByUserID(ctx, userID, isAuthor, contributorBoardsIDs)
	if err != nil {
		return nil, fmt.Errorf("get boards by user id usecase: %w", err)
	}

	return boards, nil
}

func (bCase *boardUsecase) GetCertainBoard(ctx context.Context, boardID int) (entity.BoardWithContent, string, error) {
	boardAuthorID, err := bCase.boardRepo.GetBoardAuthorByBoardID(ctx, boardID)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return entity.BoardWithContent{}, "", ErrNoSuchBoard
		default:
			return entity.BoardWithContent{}, "", fmt.Errorf("get certain board: %w", err)
		}
	}

	boardContributors, err := bCase.boardRepo.GetContributorsByBoardID(ctx, boardID)
	if err != nil {
		return entity.BoardWithContent{}, "", fmt.Errorf("get certain board: %w", err)
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

	board, username, err := bCase.boardRepo.GetBoardByID(ctx, boardID, hasAccess)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return entity.BoardWithContent{}, "", ErrNoSuchBoard
		default:
			return entity.BoardWithContent{}, "", fmt.Errorf("get certain board: %w", err)
		}
	}
	return board, username, nil
}

func isContributor(contributorsIDs []int, userID int) bool {
	for _, contributorID := range contributorsIDs {
		if contributorID == userID {
			return true
		}
	}
	return false
}
