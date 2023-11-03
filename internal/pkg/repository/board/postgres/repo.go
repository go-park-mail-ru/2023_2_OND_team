package board

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	uEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardRepoPG struct {
	db         *pgxpool.Pool
	sqlBuilder squirrel.StatementBuilderType
}

func NewBoardRepoPG(db *pgxpool.Pool) *BoardRepoPG {
	return &BoardRepoPG{db: db, sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func (repo *BoardRepoPG) CreateBoard(ctx context.Context, board entity.Board, tagTitles []string) error {

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction for creating new board: %w", err)
	}

	newBoardId, err := repo.insertBoard(ctx, tx, board)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("inserting board within transaction: %w", err)
	}

	err = repo.insertTags(ctx, tx, tagTitles)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("inserting new tags within transaction: %w", err)
	}

	err = repo.addTagsToBoard(ctx, tx, tagTitles, newBoardId, true)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("adding new tags on board within transaction: %w", err)
	}

	tx.Commit(ctx)

	return nil
}

func (boardRepo *BoardRepoPG) GetBoardsByUserID(ctx context.Context, userID int, isAuthor bool) ([]dto.GetUserBoard, error) {
	getBoardsQuery := boardRepo.sqlBuilder.
		Select(
			"board.id",
			"board.title",
			"TO_CHAR(board.created_at, 'DD:MM:YYYY')",
			"COUNT(pin.id) AS pins_number",
			"ARRAY_REMOVE((ARRAY_AGG(pin.picture))[:3], NULL) AS pins").
		From("membership").
		JoinClause("FULL JOIN pin ON membership.pin_id = pin.id").
		JoinClause("FULL JOIN board ON membership.board_id = board.id").
		Where(squirrel.Eq{"board.author": userID})

	if !isAuthor {
		getBoardsQuery = getBoardsQuery.Where(squirrel.Eq{"board.public": true})
	}
	getBoardsQuery = getBoardsQuery.
		GroupBy(
			"board.id",
			"board.title",
			"board.created_at",
		).
		OrderBy("board.id ASC")

	sqlRow, args, err := getBoardsQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query for get boards by user id: %w", err)
	}

	rows, err := boardRepo.db.Query(ctx, sqlRow, args...)
	if err != nil {
		return nil, fmt.Errorf("making query for get boards by user id: %w", err)
	}
	defer rows.Close()

	boards := make([]dto.GetUserBoard, 0)
	for rows.Next() {
		board := dto.GetUserBoard{}
		err = rows.Scan(&board.BoardID, &board.Title, &board.CreatedAt, &board.PinsNumber, &board.Pins)
		if err != nil {
			return nil, fmt.Errorf("scanning the result of get boards by user id query: %w", err)
		}
		boards = append(boards, board)
	}

	return boards, nil
}

func (repo *BoardRepoPG) GetBoardByID(ctx context.Context, boardID int, hasAccess bool) (board dto.GetUserBoard, err error) {
	getBoardByIdQuery := repo.sqlBuilder.
		Select(
			"board.id",
			"board.title",
			"COALESCE(board.description, '')",
			"TO_CHAR(board.created_at, 'DD:MM:YYYY')",
			"COUNT(DISTINCT pin.id) AS pins_number",
			"ARRAY_REMOVE(ARRAY_AGG(DISTINCT pin.picture), NULL) AS pins",
			"ARRAY_REMOVE(ARRAY_AGG(DISTINCT tag.title), NULL) AS tag_titles").
		From("membership").
		JoinClause("FULL JOIN pin ON membership.pin_id = pin.id").
		JoinClause("FULL JOIN board ON membership.board_id = board.id").
		JoinClause("FULL JOIN board_tag ON board_tag.board_id = board.id").
		JoinClause("FULL JOIN tag ON board_tag.tag_id = tag.id").
		Where(squirrel.Eq{"board.id": boardID})

	if !hasAccess {
		getBoardByIdQuery = getBoardByIdQuery.Where(squirrel.Eq{"board.public": true})

	}
	getBoardByIdQuery = getBoardByIdQuery.GroupBy(
		"board.id",
		"board.title",
		"board.description",
		"board.created_at").
		OrderBy("board.id ASC")

	sqlRow, args, err := getBoardByIdQuery.ToSql()
	if err != nil {
		return dto.GetUserBoard{}, fmt.Errorf("building get board by id query: %w", err)
	}

	row := repo.db.QueryRow(ctx, sqlRow, args...)
	board = dto.GetUserBoard{}
	err = row.Scan(&board.BoardID, &board.Title, &board.Description, &board.CreatedAt, &board.PinsNumber, &board.Pins, &board.TagTitles)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return dto.GetUserBoard{}, repository.ErrNoData
		default:
			return dto.GetUserBoard{}, fmt.Errorf("scan result of get board by id query: %w", err)
		}
	}

	return board, nil
}

func (repo *BoardRepoPG) GetBoardAuthorByBoardID(ctx context.Context, boardID int) (int, error) {
	row := repo.db.QueryRow(ctx, SelectBoardAuthorByBoardIdQuery, boardID)
	var authorID int
	err := row.Scan(&authorID)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return 0, repository.ErrNoData
		default:
			return 0, fmt.Errorf("get board author by board id query: %w", err)
		}
	}
	return authorID, nil
}

func (repo *BoardRepoPG) GetContributorsByBoardID(ctx context.Context, boardID int) ([]uEntity.User, error) {
	rows, err := repo.db.Query(ctx, SelectBoardContributorsByBoardIdQuery, boardID)
	if err != nil {
		return nil, fmt.Errorf("select contributors by board id query: %w", err)
	}
	defer rows.Close()

	contributors := make([]uEntity.User, 0)
	for rows.Next() {
		var contributorID int
		err = rows.Scan(&contributorID)
		if err != nil {
			return nil, fmt.Errorf("scan result of get contributors by board id query: %w", err)
		}
		contributors = append(contributors, uEntity.User{ID: contributorID})
	}

	return contributors, nil
}

func (repo *BoardRepoPG) insertBoard(ctx context.Context, tx pgx.Tx, board entity.Board) (int, error) {
	row := tx.QueryRow(ctx, InsertBoardQuery, board.AuthorID, board.Title, board.Description, board.Public)

	var newBoardID int
	err := row.Scan(&newBoardID)
	if err != nil {
		return 0, fmt.Errorf("scan result of insterting new board: %w", err)
	}
	return newBoardID, nil
}
