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

func (p *pinRepoPG) GetTagsByPinID(ctx context.Context, pinID int) ([]pin.Tag, error) {
	rows, err := p.db.Query(ctx, SelectTagsByPinID, pinID)
	if err != nil {
		return nil, fmt.Errorf("get pin tags by its id in storage: %w", err)
	}
	defer rows.Close()

	tags := []pin.Tag{}
	tag := pin.Tag{}
	for rows.Next() {
		err = rows.Scan(&tag.Title)
		if err != nil {
			return tags, fmt.Errorf("scan title tag for get pin tags in storage: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

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

func (p *pinRepoPG) addTagsByTitleOnPin(ctx context.Context, tx pgx.Tx, titles []string, pinID int, newPin bool) error {
	insertQuery := p.sqlBuilder.Insert("pin_tag").
		Columns("pin_id", "tag_id").
		Select(
			sq.Select(strconv.FormatInt(int64(pinID), 10), "id").
				From("tag").
				Where(sq.Eq{"title": titles}),
		)
	if !newPin {
		insertQuery = insertQuery.Suffix("ON CONFLICT (pin_id, tag_id) DO NOTHING")
	}
	sqlRow, args, err := insertQuery.ToSql()
	if err != nil {
		return fmt.Errorf("build sql query row for insert link between pins and tags: %w", err)
	}

	commTag, err := tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("executing a query to insert link between pins and tags: %w", err)
	}

	if commTag.RowsAffected() != int64(len(titles)) && newPin {
		return ErrNumberAffectedRows
	}
	return nil
}

func (p *pinRepoPG) updateSetOfTagsInPin(ctx context.Context, tx pgx.Tx, pinID int, titles []string) error {
	if len(titles) == 0 {
		return removeAllTagsFromPin(ctx, tx, pinID)
	}

	err := p.addTags(ctx, tx, titles)
	if err != nil {
		return fmt.Errorf("add tags for update set tags fo pin: %w", err)
	}

	err = p.deleteAllTagsExcept(ctx, tx, pinID, titles)
	if err != nil {
		return fmt.Errorf("delete tags for updates: %w", err)
	}

	err = p.addTagsByTitleOnPin(ctx, tx, titles, pinID, false)
	if err != nil {
		return fmt.Errorf("add tags for updates: %w", err)
	}

	return nil
}

func (p *pinRepoPG) deleteAllTagsExcept(ctx context.Context, tx pgx.Tx, pinID int, except []string) error {
	if len(except) == 0 {
		return removeAllTagsFromPin(ctx, tx, pinID)
	}

	sqlRow, args, err := p.sqlBuilder.Delete("pin_tag").
		Where(sq.Eq{"pin_id": pinID}).
		Where(sq.Expr("tag_id NOT IN (?)", sq.Select("id").
			From("tag").
			Where(sq.Eq{"title": except}))).
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql row to delete all tags except those specified: %w", err)
	}

	_, err = tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("delete all tags except those specified: %w", err)
	}
	return nil
}

func removeAllTagsFromPin(ctx context.Context, tx pgx.Tx, pinID int) error {
	_, err := tx.Exec(ctx, DeleteAllTagsFromPin, pinID)
	if err != nil {
		return fmt.Errorf("delete all tags from pin: %w", err)
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
