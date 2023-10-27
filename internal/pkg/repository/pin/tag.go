package pin

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

var ErrNumberAffectedRows = errors.New("different number of affected rows was expected")

func (p *pinRepoPG) addTags(ctx context.Context, tx pgx.Tx, titles []string) error {
	insertBuilder := p.sqlBuilder.Insert("tag").Columns("title")
	for _, title := range titles {
		insertBuilder = insertBuilder.Values(title)
	}
	sqlRow, args, err := insertBuilder.
		Suffix("ON CONFLICT (title) DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql query row for insert tags: %w", err)
	}

	_, err = tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("executing a query to insert tags: %w", err)
	}
	return nil
}

func (p *pinRepoPG) addTagsByTitleOnPin(ctx context.Context, tx pgx.Tx, titles []string, pinID int) error {
	sqlRow, args, err := p.sqlBuilder.Insert("pin_tag").
		Columns("pin_id", "tag_id").
		Select(
			sq.Select(strconv.FormatInt(int64(pinID), 10), "id").
				From("tag").
				Where(sq.Eq{"title": titles}),
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql query row for insert link between pins and tags: %w", err)
	}

	commTag, err := tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("executing a query to insert link between pins and tags: %w", err)
	}

	if commTag.RowsAffected() != int64(len(titles)) {
		return ErrNumberAffectedRows
	}
	return nil
}

func fetchTitles(tags []pin.Tag) []string {
	titles := make([]string, 0, len(tags))
	for _, tag := range tags {
		titles = append(titles, tag.Title)
	}
	return titles
}
