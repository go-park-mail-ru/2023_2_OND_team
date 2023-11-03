package board

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
)

func (bCase *BoardUsecase) GetBoardsByUsername(ctx context.Context, username string) ([]dto.GetUserBoard, error) {
	if !bCase.isValidUsername(username) {
		return nil, ErrInvalidUsername
	}

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
	boards, err := bCase.boardRepo.GetBoardsByUserID(ctx, userID, isAuthor)
	if err != nil {
		return nil, fmt.Errorf("get boards by user id usecase: %w", err)
	}

	return boards, nil
}

func (bCase *BoardUsecase) GetCertainBoardByID(ctx context.Context, boardID int) (dto.GetUserBoard, error) {
	boardAuthorID, err := bCase.boardRepo.GetBoardAuthorByBoardID(ctx, boardID)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return dto.GetUserBoard{}, ErrNoSuchBoard
		default:
			return dto.GetUserBoard{}, fmt.Errorf("get certain board by id: %w", err)
		}
	}

	boardContributors, err := bCase.boardRepo.GetContributorsByBoardID(ctx, boardID)
	if err != nil {
		return dto.GetUserBoard{}, fmt.Errorf("get certain board by id usecase: %w", err)
	}

	boardContributorsIDs := make([]int, 0, len(boardContributors))
	func() {
		for _, contributor := range boardContributors {
			boardContributorsIDs = append(boardContributorsIDs, contributor.ID)
		}
	}()

	var hasAccess bool
	currUserID, loggedIn := ctx.Value(auth.KeyCurrentUserID).(int)
	if loggedIn && (currUserID == boardAuthorID || slices.Contains(boardContributorsIDs, currUserID)) {
		hasAccess = true
	}

	fmt.Println(loggedIn, currUserID, boardAuthorID, hasAccess)
	board, err := bCase.boardRepo.GetBoardByID(ctx, boardID, hasAccess)
	if err != nil {
		switch err {
		case repository.ErrNoData:
			return dto.GetUserBoard{}, ErrNoSuchBoard
		default:
			return dto.GetUserBoard{}, fmt.Errorf("get certain board by id usecase: %w", err)
		}
	}

	return board, nil
}
