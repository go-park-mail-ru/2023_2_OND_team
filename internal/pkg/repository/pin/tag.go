package pin

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

var ErrNumberReceivedRows = errors.New("different number of received rows was expected")

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

func (p *pinRepoPG) getTagIdsByTitles(ctx context.Context, tx pgx.Tx, titles []string) ([]int, error) {
	sqlRow, args, err := p.sqlBuilder.Select("id").
		From("tag").
		Where(sq.Eq{"title": titles}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query row for select tags: %w", err)
	}

	rows, err := tx.Query(ctx, sqlRow, args...)
	if err != nil {
		return nil, fmt.Errorf("query to select tags id by title: %w", err)
	}
	defer rows.Close()

	idTag := 0
	tagIds := make([]int, 0, len(titles))
	for rows.Next() {
		err = rows.Scan(&idTag)
		if err != nil {
			return nil, fmt.Errorf("scan a tag id for add new pin with these tags: %w", err)
		}
		tagIds = append(tagIds, idTag)
	}

	if len(titles) != len(tagIds) {
		return nil, ErrNumberReceivedRows
	}
	return tagIds, nil
}

func (p *pinRepoPG) addTagsOnPin(ctx context.Context, tx pgx.Tx, tagIds []int, pinID int) error {
	insertBuilder := p.sqlBuilder.Insert("pin_tag").Columns("pin_id", "tag_id")
	for _, idTag := range tagIds {
		insertBuilder = insertBuilder.Values(pinID, idTag)
	}
	sqlRow, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build sql query row for insert link between pins and tags: %w", err)
	}

	_, err = tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("executing a query to insert link between pins and tags: %w", err)
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
