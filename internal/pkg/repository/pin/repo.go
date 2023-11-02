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

type S map[string]any
type Repository interface {
	GetSortedNewNPins(ctx context.Context, count, midID, maxID int) ([]entity.Pin, error)
	GetAuthorPin(ctx context.Context, pinID int) (*user.User, error)
	GetPinByID(ctx context.Context, pinID int) (*entity.Pin, error)
	AddNewPin(ctx context.Context, pin *entity.Pin) error
	DeletePin(ctx context.Context, pinID, userID int) error
	SetLike(ctx context.Context, pinID, userID int) error
	DelLike(ctx context.Context, pinID, userID int) error
	EditPin(ctx context.Context, pinID int, updateData S, titleTags []string) error
	GetCountLikeByPinID(ctx context.Context, pinID int) (int, error)
	GetTagsByPinID(ctx context.Context, pinID int) ([]entity.Tag, error)
	IsAvailableToUserAsContributorBoard(ctx context.Context, pinID, userID int) (bool, error)
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

func (p *pinRepoPG) GetSortedNewNPins(ctx context.Context, count, minID, maxID int) ([]entity.Pin, error) {
	rows, err := p.db.Query(ctx, SelectWithExcludeLimit, minID, maxID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins: %w", count, err)
	}

	pins := make([]entity.Pin, 0, count)
	pin := entity.Pin{}
	for rows.Next() {
		err := rows.Scan(&pin.ID, &pin.Picture)
		if err != nil {
			return pins, fmt.Errorf("scan to receive %d pins: %w", count, err)
		}
		pins = append(pins, pin)
	}

	return pins, nil
}

func (p *pinRepoPG) GetPinByID(ctx context.Context, pinID int) (*entity.Pin, error) {
	row := p.db.QueryRow(ctx, SelectPinByID, pinID)
	pin := &entity.Pin{}
	err := row.Scan(&pin.AuthorID, &pin.Title, &pin.Description,
		&pin.Picture, &pin.Public, &pin.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("get pin by id from storage: %w", err)
	}
	pin.ID = pinID
	return pin, nil
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

	pinID, err := p.addPin(ctx, tx, pin)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("add pin: %w", err)
	}

	err = p.addTagsByTitleOnPin(ctx, tx, titles, pinID, true)
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

func (p *pinRepoPG) DeletePin(ctx context.Context, pinID, userID int) error {
	_, err := p.db.Exec(ctx, UpdatePinSetStatusDelete, pinID, userID)
	if err != nil {
		return fmt.Errorf("set pin deleted at now: %w", err)
	}
	return nil
}

func (p *pinRepoPG) EditPin(ctx context.Context, pinID int, updateData S, titleTags []string) error {
	if len(updateData) == 0 && titleTags == nil {
		return nil
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction for edit pin: %w", err)
	}

	if len(updateData) != 0 {
		err = p.updateHeaderPin(ctx, tx, pinID, updateData)
	}
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("edit pin header: %w", err)
	}

	if titleTags != nil {
		err = p.updateSetOfTagsInPin(ctx, tx, pinID, titleTags)
	}
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("edit tags on pin: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction for edit pin: %w", err)
	}
	return nil
}

func (p *pinRepoPG) updateHeaderPin(ctx context.Context, tx pgx.Tx, pinID int, newHeader S) error {
	sqlRow, args, err := p.sqlBuilder.Update("pin").
		SetMap(newHeader).
		Where(sq.Eq{"id": pinID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql row for update header pin: %w", err)
	}

	_, err = tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("update header pin in storage: %w", err)
	}
	return nil
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
