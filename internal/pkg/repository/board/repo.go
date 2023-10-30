package board

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateBoard(ctx context.Context, board entity.Board, pinIDs []int, tagTitles []string) error
	// CreateBoardMembershipByPinIDs(ctx context.Context, pinIDs []int) error
	// SelectOwnBoardsByUserID(ctx context.Context, userID int) ([]entity.Board, error)
	// SelectUserBoardsByUserID(ctx context.Context, userID int) ([]entity.Board, error)
	// SelectBoardsByTitle(ctx context.Context, title string) ([]entity.Board, error)
	// SelctBoardsByTag(ctx context.Context, tagTitle string) ([]entity.Board, error)
	// SelectBoardTags(ctx context.Context, boardID int) (tagTitles []string, err error)
}

type BoardRepoPG struct {
	db         *pgxpool.Pool
	sqlBuilder squirrel.StatementBuilderType
}

func NewBoardRepoPG(db *pgxpool.Pool) *BoardRepoPG {
	return &BoardRepoPG{db: db, sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func (repo *BoardRepoPG) CreateBoard(ctx context.Context, board entity.Board, pinIDs []int, tagTitles []string) error {

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

func (repo *BoardRepoPG) insertBoard(ctx context.Context, tx pgx.Tx, board entity.Board) (int, error) {
	row := tx.QueryRow(ctx, CreateBoardQuery, board.AuthorID, board.Title, board.Description, board.Public)

	var newBoardID int
	err := row.Scan(&newBoardID)
	if err != nil {
		return 0, fmt.Errorf("scan result of insterting new board: %w", err)
	}
	return newBoardID, nil
}

func (repo *BoardRepoPG) addTagsToBoard(ctx context.Context, tx pgx.Tx, tagTitles []string, boardID int, isNewBoard bool) error {
	addTagsToBoardQuery := repo.sqlBuilder.
		Insert("board_tag").
		Columns("board_id", "tag_id").
		Select(
			squirrel.Select(strconv.FormatInt(int64(boardID), 10), "id").
				From("tag").
				Where(squirrel.Eq{"title": tagTitles}),
		)

	if !isNewBoard {
		addTagsToBoardQuery.Suffix("ON CONFLICT DO NOTHING")
	}

	sqlRow, args, err := addTagsToBoardQuery.ToSql()
	if err != nil {
		return fmt.Errorf("building sql query row for adding tags to board: %w", err)
	}

	cmdTag, err := tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("execute sql query to add tags to board: %w", err)
	}

	if cmdTag.RowsAffected() != int64(len(tagTitles)) {
		return fmt.Errorf("not all tags were inserted correctly")
	}

	return nil
}
