package board

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *boardRepoPG) insertTags(ctx context.Context, tx pgx.Tx, titles []string) error {
	if len(titles) == 0 {
		return nil
	}

	insertTagsQuery := repo.sqlBuilder.
		Insert("tag").
		Columns("title")
	for _, title := range titles {
		insertTagsQuery = insertTagsQuery.Values(title)
	}
	sqlRow, args, err := insertTagsQuery.
		Suffix("ON CONFLICT (title) DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql row query while inserting tags: %w", err)
	}

	_, err = tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("making insertTags query: %w", err)
	}
	return nil
}

func (repo *boardRepoPG) addTagsToBoard(ctx context.Context, tx pgx.Tx, tagTitles []string, boardID int, isNewBoard bool) error {
	if !isNewBoard {
		err := repo.deleteCurrentBoardTags(ctx, tx, boardID)
		if err != nil {
			return fmt.Errorf("add tags to board within transaction: %w", err)
		}
	}

	addTagsToBoardQuery := repo.sqlBuilder.
		Insert("board_tag").
		Columns("board_id", "tag_id").
		Select(
			squirrel.Select(strconv.FormatInt(int64(boardID), 10), "id").
				From("tag").
				Where(squirrel.Eq{"title": tagTitles}),
		)

	if !isNewBoard {
		addTagsToBoardQuery = addTagsToBoardQuery.Suffix("ON CONFLICT DO NOTHING")
	}

	sqlRow, args, err := addTagsToBoardQuery.ToSql()
	if err != nil {
		return fmt.Errorf("building sql query row for adding tags to board: %w", err)
	}

	status, err := tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("execute sql query to add tags to board: %w", err)
	}

	if isNewBoard && int(status.RowsAffected()) != len(tagTitles) {
		return ErrIncorrectNumberRowsAffcted
	}

	return nil
}

func (repo *boardRepoPG) deleteCurrentBoardTags(ctx context.Context, tx pgx.Tx, boardID int) error {
	_, err := tx.Exec(ctx, DeleteCurrentBoardTags, boardID)
	if err != nil {
		return fmt.Errorf("deleting current board tags: %w", err)
	}

	return nil
}
