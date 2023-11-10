package pin

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/internal/pgtype"
)

type S map[string]any

//go:generate mockgen -destination=./mock/pin_mock.go -package=mock -source=repo.go Repository
type Repository interface {
	GetSortedNewNPins(ctx context.Context, count, midID, maxID int) ([]entity.Pin, error)
	GetSortedUserPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, error)
	GetAuthorPin(ctx context.Context, pinID int) (*user.User, error)
	GetPinByID(ctx context.Context, pinID int, revealAuthor bool) (*entity.Pin, error)
	GetBatchPinByID(ctx context.Context, pinID []int) ([]entity.Pin, error)
	AddNewPin(ctx context.Context, pin *entity.Pin) error
	DeletePin(ctx context.Context, pinID, userID int) error
	SetLike(ctx context.Context, pinID, userID int) (int, error)
	IsSetLike(ctx context.Context, pinID, userID int) (bool, error)
	DelLike(ctx context.Context, pinID, userID int) (int, error)
	EditPin(ctx context.Context, pinID int, updateData S, titleTags []string) error
	GetCountLikeByPinID(ctx context.Context, pinID int) (int, error)
	GetTagsByPinID(ctx context.Context, pinID int) ([]entity.Tag, error)
	IsAvailableToUserAsContributorBoard(ctx context.Context, pinID, userID int) (bool, error)

	GetFeedPins(ctx context.Context, cfg entity.FeedPinConfig) (entity.FeedPin, error)
	GetSortedUserLikedPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, error)
}

type pinRepoPG struct {
	db         pgtype.PgxPoolIface
	sqlBuilder sq.StatementBuilderType
}

func NewPinRepoPG(db pgtype.PgxPoolIface) *pinRepoPG {
	return &pinRepoPG{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (p *pinRepoPG) GetFeedPins(ctx context.Context, cfg entity.FeedPinConfig) (entity.FeedPin, error) {
	queryBuild := p.sqlBuilder.Select("pin.id", "pin.picture").
		From("pin").Limit(uint64(cfg.Count))

	pin := entity.Pin{}
	scanFields := []any{&pin.ID, &pin.Picture}

	switch cfg.Protection {
	case entity.FeedAll:
		queryBuild = queryBuild.Columns("pin.public")
		scanFields = append(scanFields, &pin.Public)
	case entity.FeedProtectionPublic:
		queryBuild = queryBuild.Where(sq.Eq{"pin.public": true})
	case entity.FeedProtectionPrivate:
		queryBuild = queryBuild.Where(sq.Eq{"pin.public": false})
	}

	if cfg.Deleted {
		queryBuild = queryBuild.Where(sq.NotEq{"pin.deleted_at": nil})
	} else {
		queryBuild = queryBuild.Where(sq.Eq{"pin.deleted_at": nil})
	}

	if cfg.UserID != 0 && !cfg.Liked {
		queryBuild = queryBuild.Where(sq.Eq{"pin.author": cfg.UserID})
	}

	if cfg.Liked {
		queryBuild = queryBuild.InnerJoin("like_pin ON like_pin.pin_id = pin.id").
			Where(sq.Eq{"like_pin.user_id": cfg.UserID}).
			OrderBy("like_pin.created_at DESC")
	}

	if cfg.BoardID != 0 {
		queryBuild = queryBuild.InnerJoin("membership", "membership.pin_id = pin.id").
			InnerJoin("board", "membership.board_id = board.id").
			Where(sq.Eq{"board.id": cfg.BoardID})
	}

	sqlRow, args, err := queryBuild.
		Where(sq.Or{sq.Lt{"pin.id": cfg.MinID}, sq.Gt{"pin.id": cfg.MaxID}}).
		OrderBy("pin.id DESC").
		ToSql()

	if err != nil {
		return entity.FeedPin{}, fmt.Errorf("query build error: %w", err)
	}
	fmt.Println(sqlRow, args)

	rows, err := p.db.Query(ctx, sqlRow, args...)
	if err != nil {
		return entity.FeedPin{Pins: []entity.Pin{}, Condition: cfg.Condition}, fmt.Errorf("getting pins for feed from storage: %w", err)
	}
	feed := entity.FeedPin{Condition: cfg.Condition}
	for rows.Next() {
		fmt.Println("SCAN")
		err = rows.Scan(scanFields...)
		fmt.Println(pin)
		if err != nil {
			return feed, fmt.Errorf("scan feed pins: %w", err)
		}
		feed.Pins = append(feed.Pins, pin)
	}
	if len(feed.Pins) != 0 && feed.Pins[0].ID > cfg.MaxID {
		feed.MaxID = feed.Pins[0].ID
	}
	if len(feed.Pins) != 0 && (feed.Pins[len(feed.Pins)-1].ID < cfg.MinID || cfg.MinID == 0) {
		feed.MinID = feed.Pins[len(feed.Pins)-1].ID
	}
	fmt.Println(len(feed.Pins))
	return feed, nil
}

func (p *pinRepoPG) GetSortedNewNPins(ctx context.Context, count, minID, maxID int) ([]entity.Pin, error) {
	rows, err := p.db.Query(ctx, SelectWithExcludeLimit, minID, maxID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins: %w", count, err)
	}

	pins := make([]entity.Pin, 0, count)
	pin := entity.Pin{Public: true}
	for rows.Next() {
		err := rows.Scan(&pin.ID, &pin.Picture)
		if err != nil {
			return pins, fmt.Errorf("scan to receive %d pins: %w", count, err)
		}
		pins = append(pins, pin)
	}

	return pins, nil
}

func (p *pinRepoPG) GetSortedUserPins(ctx context.Context, userID, count, minID, maxID int) ([]entity.Pin, error) {
	rows, err := p.db.Query(ctx, SelectUserPinsLimit, userID, minID, maxID, count)
	if err != nil {
		return nil, fmt.Errorf("select to receive %d pins: %w", count, err)
	}

	pins := make([]entity.Pin, 0, count)
	pin := entity.Pin{}
	for rows.Next() {
		err := rows.Scan(&pin.ID, &pin.Picture, &pin.Public)
		if err != nil {
			return pins, fmt.Errorf("scan to receive %d pins: %w", count, err)
		}
		pins = append(pins, pin)
	}

	return pins, nil
}

func (p *pinRepoPG) GetPinByID(ctx context.Context, pinID int, revealAuthor bool) (*entity.Pin, error) {
	pin := &entity.Pin{Author: &user.User{}}
	var err error
	if revealAuthor {
		err = p.getPinByID(ctx, pinID, SelectPinByIDWithAuthor,
			&pin.Author.ID, &pin.Title, &pin.Description,
			&pin.Picture, &pin.Public, &pin.DeletedAt,
			&pin.Author.Username, &pin.Author.Avatar)
	} else {
		err = p.getPinByID(ctx, pinID, SelectPinByID,
			&pin.Author.ID, &pin.Title, &pin.Description,
			&pin.Picture, &pin.Public, &pin.DeletedAt)
	}
	if err != nil {
		return nil, fmt.Errorf("get pin by id from storage: %w", err)
	}

	pin.ID = pinID
	return pin, nil
}

func (p *pinRepoPG) GetBatchPinByID(ctx context.Context, pinID []int) ([]entity.Pin, error) {
	sqlRow, args, err := p.sqlBuilder.Select("id", "author", "public", "deleted_at").
		From("pin").
		Where(sq.Eq{"id": pinID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql query build for get batch pins: %w", err)
	}

	rows, err := p.db.Query(ctx, sqlRow, args...)
	if err != nil {
		return nil, fmt.Errorf("select batch pins: %w", err)
	}

	pin := entity.Pin{Author: &user.User{}}
	pins := make([]entity.Pin, 0, len(pinID))
	for rows.Next() {
		err = rows.Scan(&pin.ID, &pin.Author.ID, &pin.Public, &pin.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("scan result select batch pins: %w", err)
		}
		pins = append(pins, pin)
	}

	if len(pins) != len(pinID) {
		return nil, ErrNumberSelectRows
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
		Values(pin.Author.ID, pin.Title, pin.Description, pin.Picture, pin.Public).
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

func (p *pinRepoPG) getPinByID(ctx context.Context, pinID int, query string, dest ...any) error {
	row := p.db.QueryRow(ctx, query, pinID)
	return row.Scan(dest...)
}
