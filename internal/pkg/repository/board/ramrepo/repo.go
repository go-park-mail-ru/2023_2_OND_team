package board

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	uEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
)

type BoardRepoRam struct {
	db         *sql.DB
	sqlBuilder squirrel.StatementBuilderType
}

func NewBoardRepoRam(db *sql.DB) *BoardRepoRam {
	return &BoardRepoRam{db: db, sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func (repo *BoardRepoRam) CreateBoard(ctx context.Context, board entity.Board, tagTitles []string) error {
	return repository.ErrMethodUnimplemented
}

func (repo *BoardRepoRam) GetBoardsByUserID(ctx context.Context, userID int, isAuthor bool) ([]dto.GetUserBoard, error) {
	return nil, repository.ErrMethodUnimplemented
}

func (repo *BoardRepoRam) GetBoardByID(ctx context.Context, boardID int, hasAccess bool) (board dto.GetUserBoard, err error) {
	return dto.GetUserBoard{}, repository.ErrMethodUnimplemented
}

func (repo *BoardRepoRam) GetBoardAuthorByBoardID(ctx context.Context, boardID int) (int, error) {
	return 0, repository.ErrMethodUnimplemented
}

func (repo *BoardRepoRam) GetContributorsByBoardID(ctx context.Context, boardID int) ([]uEntity.User, error) {
	return nil, repository.ErrMethodUnimplemented
}
