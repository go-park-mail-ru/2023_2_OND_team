package board

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	uEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	repoBoard "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/internal/pgtype"
)

type boardRepoPG struct {
	db         pgtype.PgxPoolIface
	sqlBuilder squirrel.StatementBuilderType
}

func NewBoardRepoPG(db pgtype.PgxPoolIface) *boardRepoPG {
	return &boardRepoPG{db: db, sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func (repo *boardRepoPG) CreateBoard(ctx context.Context, board entity.Board, tagTitles []string) (int, error) {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("starting transaction for creating new board: %w", err)
	}

	newBoardId, err := repo.insertBoard(ctx, tx, board)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("inserting board within transaction: %w", err)
	}

	err = repo.insertTags(ctx, tx, tagTitles)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("inserting new tags within transaction: %w", err)
	}

	err = repo.addTagsToBoard(ctx, tx, tagTitles, newBoardId, true)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("adding new tags on board within transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit transaction for create new board: %w", err)
	}

	return newBoardId, nil
}

func (boardRepo *boardRepoPG) GetBoardsByUserID(ctx context.Context, userID int, isAuthor bool, accessableBoardsIDs []int) ([]entity.BoardWithContent, error) {
	getBoardsQuery := boardRepo.sqlBuilder.
		Select(
			"board.id",
			"board.title",
			"COALESCE(board.description, '')",
			"board.created_at",
			"COUNT(DISTINCT pin.id) FILTER (WHERE pin.deleted_at IS NULL) AS pins_number",
			"COALESCE(ARRAY_REMOVE((ARRAY_AGG(DISTINCT pin.picture) FILTER (WHERE pin.deleted_at IS NULL))[:3], NULL), array[]::text[]) AS pins",
			"ARRAY_REMOVE(ARRAY_AGG(DISTINCT tag.title), NULL) AS tag_titles").
		From("board").
		LeftJoin("membership ON board.id = membership.board_id").
		LeftJoin("pin ON membership.pin_id = pin.id").
		LeftJoin("board_tag ON board.id = board_tag.board_id").
		LeftJoin("tag ON board_tag.tag_id = tag.id").
		Where(squirrel.Eq{"board.deleted_at": nil}).
		Where(squirrel.Eq{"board.author": userID})

	if !isAuthor {
		getBoardsQuery = getBoardsQuery.Where(
			squirrel.Or{
				squirrel.Eq{"board.public": true},
				squirrel.Eq{"board.id": accessableBoardsIDs},
			},
		)
	}
	getBoardsQuery = getBoardsQuery.
		GroupBy(
			"board.id",
			"board.title",
			"board.description",
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

	boards := make([]entity.BoardWithContent, 0)
	for rows.Next() {
		board := entity.BoardWithContent{}
		err = rows.Scan(&board.BoardInfo.ID, &board.BoardInfo.Title, &board.BoardInfo.Description, &board.BoardInfo.CreatedAt, &board.PinsNumber, &board.Pins, &board.TagTitles)
		if err != nil {
			return nil, fmt.Errorf("scanning the result of get boards by user id query: %w", err)
		}
		boards = append(boards, board)
	}

	return boards, nil
}

func (repo *boardRepoPG) GetBoardByID(ctx context.Context, boardID int, hasAccess bool) (board entity.BoardWithContent, err error) {
	getBoardByIdQuery := repo.sqlBuilder.
		Select(
			"board.id",
			"board.author",
			"board.title",
			"COALESCE(board.description, '')",
			"board.created_at",
			"COUNT(DISTINCT pin.id) FILTER (WHERE pin.deleted_at IS NULL) AS pins_number",
			"COALESCE(ARRAY_REMOVE(ARRAY_AGG(DISTINCT pin.picture) FILTER (WHERE pin.deleted_at IS NULL), NULL), array[]::text[]) AS pins",
			"ARRAY_REMOVE(ARRAY_AGG(DISTINCT tag.title), NULL) AS tag_titles").
		From("board").
		LeftJoin("board_tag ON board.id = board_tag.board_id").
		LeftJoin("tag ON board_tag.tag_id = tag.id").
		LeftJoin("membership ON board.id = membership.board_id").
		LeftJoin("pin ON membership.pin_id = pin.id").
		Where(squirrel.Eq{"board.deleted_at": nil}).
		Where(squirrel.Eq{"board.id": boardID})

	if !hasAccess {
		getBoardByIdQuery = getBoardByIdQuery.Where(squirrel.Eq{"board.public": true})
	}
	getBoardByIdQuery = getBoardByIdQuery.GroupBy(
		"board.id",
		"board.author",
		"board.title",
		"board.description",
		"board.created_at").
		OrderBy("board.id ASC")

	sqlRow, args, err := getBoardByIdQuery.ToSql()
	if err != nil {
		return entity.BoardWithContent{}, fmt.Errorf("building get board by id query: %w", err)
	}

	row := repo.db.QueryRow(ctx, sqlRow, args...)
	board = entity.BoardWithContent{}
	err = row.Scan(&board.BoardInfo.ID, &board.BoardInfo.AuthorID, &board.BoardInfo.Title, &board.BoardInfo.Description, &board.BoardInfo.CreatedAt, &board.PinsNumber, &board.Pins, &board.TagTitles)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return entity.BoardWithContent{}, repository.ErrNoData
		default:
			return entity.BoardWithContent{}, fmt.Errorf("scan result of get board by id query: %w", err)
		}
	}

	return board, nil
}

func (repo *boardRepoPG) GetBoardInfoForUpdate(ctx context.Context, boardID int, hasAccess bool) (entity.Board, []string, error) {
	getBoardByIdQuery := repo.sqlBuilder.
		Select(
			"board.title",
			"COALESCE(board.description, '')",
			"board.public",
			"ARRAY_REMOVE(ARRAY_AGG(DISTINCT tag.title), NULL) AS tag_titles").
		From("board").
		LeftJoin("board_tag ON board.id = board_tag.board_id").
		LeftJoin("tag ON board_tag.tag_id = tag.id").
		Where(squirrel.Eq{"board.deleted_at": nil}).
		Where(squirrel.Eq{"board.id": boardID})

	if !hasAccess {
		getBoardByIdQuery = getBoardByIdQuery.Where(squirrel.Eq{"board.public": true})
	}
	getBoardByIdQuery = getBoardByIdQuery.GroupBy(
		"board.id",
		"board.title",
		"board.description").
		OrderBy("board.id ASC")

	sqlRow, args, err := getBoardByIdQuery.ToSql()
	if err != nil {
		return entity.Board{}, nil, fmt.Errorf("building get board info for update query: %w", err)
	}

	row := repo.db.QueryRow(ctx, sqlRow, args...)
	board := entity.Board{}
	tagTitles := make([]string, 0)
	err = row.Scan(&board.Title, &board.Description, &board.Public, &tagTitles)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return entity.Board{}, nil, repository.ErrNoData
		default:
			return entity.Board{}, nil, fmt.Errorf("scan result of get board by id query: %w", err)
		}
	}

	return board, tagTitles, nil
}

func (repo *boardRepoPG) GetBoardAuthorByBoardID(ctx context.Context, boardID int) (int, error) {
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

func (repo *boardRepoPG) GetContributorsByBoardID(ctx context.Context, boardID int) ([]uEntity.User, error) {
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

func (repo *boardRepoPG) GetContributorBoardsIDs(ctx context.Context, contributorID int) ([]int, error) {
	rows, err := repo.db.Query(ctx, GetContributorBoardsIDs, contributorID)
	if err != nil {
		return nil, fmt.Errorf("get contributor boardsIDs query: %w", err)
	}
	defer rows.Close()

	boardsIDs := make([]int, 0)
	for rows.Next() {
		var boardID int
		err = rows.Scan(&boardID)
		if err != nil {
			return nil, fmt.Errorf("get contributor boardsIDs query: %w", err)
		}
		boardsIDs = append(boardsIDs, boardID)
	}

	return boardsIDs, nil
}

func (repo *boardRepoPG) UpdateBoard(ctx context.Context, newBoardData entity.Board, tagTitles []string) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("start update board transaction: %w", err)
	}

	err = repo.insertTags(ctx, tx, tagTitles)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("update board: insert tags within transaction - %w", err)
	}

	err = repo.addTagsToBoard(ctx, tx, tagTitles, newBoardData.ID, false)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("update board: add tags to board within transaction - %w", err)
	}

	status, err := repo.db.Exec(ctx, UpdateBoardByIdQuery, newBoardData.Title, newBoardData.Description, newBoardData.Public, newBoardData.ID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("update board: edit board data within transaction - %w", err)
	}

	if status.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return repository.ErrNoData
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction for update board: %w", err)
	}
	return nil
}

func (repo *boardRepoPG) DeleteBoardByID(ctx context.Context, boardID int) error {
	status, err := repo.db.Exec(ctx, DeleteBoardByIdQuery, time.Now(), boardID)
	if err != nil {
		return fmt.Errorf("delete board by id: %w", err)
	}

	if status.RowsAffected() == 0 {
		return repository.ErrNoDataAffected
	}

	return nil
}

func (repo *boardRepoPG) insertBoard(ctx context.Context, tx pgx.Tx, board entity.Board) (int, error) {
	row := tx.QueryRow(ctx, InsertBoardQuery, board.AuthorID, board.Title, board.Description, board.Public)

	var newBoardID int
	err := row.Scan(&newBoardID)
	if err != nil {
		return 0, fmt.Errorf("scan result of insterting new board: %w", err)
	}
	return newBoardID, nil
}

func (repo *boardRepoPG) AddPinsOnBoard(ctx context.Context, boardID int, pinIds []int) error {
	insertBuilder := repo.sqlBuilder.Insert("membership").Columns("pin_id", "board_id")
	for _, pinID := range pinIds {
		insertBuilder = insertBuilder.Values(pinID, boardID)
	}
	sqlRow, args, err := insertBuilder.Suffix("ON CONFLICT (pin_id, board_id) DO NOTHING").ToSql()
	if err != nil {
		return fmt.Errorf("build sql query for add pins on board: %w", err)
	}

	_, err = repo.db.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("insert membership for add pins on board: %w", err)
	}
	return nil
}

func (b *boardRepoPG) GetProtectionStatusBoard(ctx context.Context, boardID int) (repoBoard.ProtectionBoard, error) {
	var isPublic bool
	err := b.db.QueryRow(ctx, SelectProtectionStatusBoard, boardID).Scan(&isPublic)
	if err != nil {
		return 0, fmt.Errorf("get status board in storage: %w", err)
	}
	if isPublic {
		return repoBoard.ProtectionPublic, nil
	}
	return repoBoard.ProtectionPrivate, nil
}
