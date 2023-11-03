package board

import (
	"context"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	uEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
)

type Repository interface {
	CreateBoard(ctx context.Context, board entity.Board, tagTitles []string) error
	GetBoardsByUserID(ctx context.Context, userID int, isAuthor bool, accessableBoardsIDs []int) ([]dto.UserBoard, error)
	GetBoardByID(ctx context.Context, boardID int, hasAccess bool) (board dto.UserBoard, err error)
	GetBoardAuthorByBoardID(ctx context.Context, boardID int) (int, error)
	GetContributorsByBoardID(ctx context.Context, boardID int) ([]uEntity.User, error)
	GetContributorBoardsIDs(ctx context.Context, contributorID int) ([]int, error)
	UpdateBoard(ctx context.Context, newBoardData entity.Board, tagTitles []string) error
	// AddContributor
	// GetBoardsByTitle(ctx context.Context, title string) ([]entity.Board, error)
	// GetBoardsByTag(ctx context.Context, tagTitle string) ([]entity.Board, error)
	// GetBoardTags(ctx context.Context, boardID int) (tagTitles []string, err error)
}
