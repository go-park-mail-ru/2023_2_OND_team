package board

import (
	"context"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	uEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
)

//go:generate mockgen -destination=./mock/board_mock.go -package=mock -source=repo.go Repository
type Repository interface {
	CreateBoard(ctx context.Context, board entity.Board, tagTitles []string) (int, error)
	GetBoardsByUserID(ctx context.Context, userID int, isAuthor bool, accessableBoardsIDs []int) ([]dto.UserBoard, error)
	GetBoardByID(ctx context.Context, boardID int, hasAccess bool) (board dto.UserBoard, err error)
	GetBoardAuthorByBoardID(ctx context.Context, boardID int) (int, error)
	GetContributorsByBoardID(ctx context.Context, boardID int) ([]uEntity.User, error)
	GetContributorBoardsIDs(ctx context.Context, contributorID int) ([]int, error)
	GetBoardInfoForUpdate(ctx context.Context, boardID int, hasAccess bool) (entity.Board, []string, error)
	UpdateBoard(ctx context.Context, newBoardData entity.Board, tagTitles []string) error
	DeleteBoardByID(ctx context.Context, boardID int) error
	RoleUserHaveOnThisBoard(ctx context.Context, boardID int, userID int) (UserRole, error)
	AddPinsOnBoard(ctx context.Context, boardID int, pinIds []int) error
	GerProtectionStatusBoard(ctx context.Context, boardID int) (ProtectionBoard, error)
}

type UserRole uint8

const (
	RegularUser UserRole = 1 << iota
	Subscriber
	ContributorForReading
	ContributorForAdding
	Author
)

type ProtectionBoard uint8

const (
	_ ProtectionBoard = iota
	ProtectionPublic
	ProtectionPrivate
)
