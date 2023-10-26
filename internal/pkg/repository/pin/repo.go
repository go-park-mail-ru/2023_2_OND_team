package pin

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Repository interface {
	GetSortedNPinsAfterID(ctx context.Context, count int, afterPinID int) ([]entity.Pin, error)
	GetAuthorPin(ctx context.Context, pinID int) (*user.User, error)
	AddNewPin(ctx context.Context, pin *entity.Pin) error
}

type pinRepoPG struct {
	db         *pgxpool.Pool
	sqlBuilder sq.StatementBuilderType
}

func NewPinRepoPG(db *pgxpool.Pool) *pinRepoPG {
	return &pinRepoPG{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (p *pinRepoPG) GetSortedNPinsAfterID(ctx context.Context, count int, afterPinID int) ([]entity.Pin, error) {
	rows, err := p.db.Query(ctx, SelectAfterIdWithLimit, afterPinID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins after %d: %w", count, afterPinID, err)
	}

	pins := make([]entity.Pin, 0, count)
	pin := entity.Pin{}
	for rows.Next() {
		err := rows.Scan(&pin.ID, &pin.Picture)
		if err != nil {
			return pins, fmt.Errorf("scan to receive %d pins after %d: %w", count, afterPinID, err)
		}
		pins = append(pins, pin)
	}

	return pins, nil
}

func (p *pinRepoPG) AddNewPin(ctx context.Context, pin *entity.Pin) error {
	titles := fetchTitles(pin.Tags)

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction for add new pin: %w", err)
	}

	err = p.addTags(ctx, tx, titles)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("add tags: %w", err)
	}

	tagIds, err := p.getTagIdsByTitles(ctx, tx, titles)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("get tag ids by titles: %w", err)
	}

	pinID, err := p.addPin(ctx, tx, pin)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("add pin: %w", err)
	}

	err = p.addTagsOnPin(ctx, tx, tagIds, pinID)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("link of the tag to the picture: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction for add new pin: %w", err)
	}
	return nil
}

func (p *pinRepoPG) GetAuthorPin(ctx context.Context, pinID int) (*user.User, error) {
	return nil, errors.New("unimplemented")
}

func (p *pinRepoPG) addPin(ctx context.Context, tx pgx.Tx, pin *entity.Pin) (int, error) {
	sqlRow, args, err := p.sqlBuilder.Insert("pin").
		Columns("author", "title", "description", "picture", "public").
		Values(pin.AuthorID, pin.Title, pin.Description, pin.Picture, pin.Public).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build sql query row for insert pin: %w", err)
	}

	row := tx.QueryRow(ctx, sqlRow, args...)
	pinID := 0
	err = row.Scan(&pinID)
	if err != nil {
		return 0, fmt.Errorf("scan the result of the insert query to add pin: %w", err)
	}
	return pinID, nil
}
