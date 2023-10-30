package board

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (repo *BoardRepoPG) insertTags(ctx context.Context, tx pgx.Tx, titles []string) error {
	insertTagsQuery := repo.sqlBuilder.
		Insert("tag").
		Columns("title")
	for _, title := range titles {
		insertTagsQuery = insertTagsQuery.Values(title)
	}
	sqlRow, args, err := insertTagsQuery.
		Suffix("ON CONFLICT (title) DO NOTHING").
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql row query while inserting tags: %w", err)
	}

	cmdTag, err := tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("making insertTags query: %w", err)
	}
	if cmdTag.RowsAffected() != int64(len(titles)) {
		return fmt.Errorf("checking rows affected after insertTags: %w", err)
	}

	return nil
}
