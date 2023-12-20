package share

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/share"
	pool "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/internal/pgtype"
)

type Repository interface {
	AddSharedLink(ctx context.Context, link *share.SharedLink) (int, error)
	GetSharedLink(ctx context.Context, linkID int) (*share.SharedLink, error)
}

type querer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type shareRepoPG struct {
	db         pool.PgxPoolIface
	sqlBuilder sq.StatementBuilderType
}

func NewShareRepoPG(db pool.PgxPoolIface) *shareRepoPG {
	return &shareRepoPG{
		db:         db,
		sqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *shareRepoPG) AddSharedLink(ctx context.Context, link *share.SharedLink) (int, error) {
	if link.IsDistributedAll {
		id, err := addNewSharedLink(ctx, s.db, link)
		if err != nil {
			return 0, fmt.Errorf("add shared link in storage: %w", err)
		}
		return id, nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction for adding shared link: %w", err)
	}

	id, err := addNewSharedLink(ctx, tx, link)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("add shared link in storage: %w", err)
	}

	err = s.distributeAccess(ctx, tx, id, link.Users)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("distribute access: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit transaction for adding shared link: %w", err)
	}

	return id, nil
}

func (s *shareRepoPG) GetSharedLink(ctx context.Context, linkID int) (*share.SharedLink, error) {
	rows, err := s.db.Query(ctx, SelectLink, linkID)
	if err != nil {
		return nil, fmt.Errorf("select shared link: %w", err)
	}
	defer rows.Close()

	link := &share.SharedLink{}

	rows.Next()

	var userIdPgType pgtype.Int8
	var role int
	err = rows.Scan(&link.BoardID, &role, &userIdPgType)
	if err != nil {
		return nil, fmt.Errorf("scan link fog getting shared link: %w", err)
	}

	link.Role = getStringContributorRoleFromInt(role)

	if !userIdPgType.Valid {
		link.IsDistributedAll = true
		return link, nil
	}

	link.Users = append(link.Users, int(userIdPgType.Int64))

	var userID int
	for rows.Next() {
		err = rows.Scan(nil, nil, &userID)
		if err != nil {
			return nil, fmt.Errorf("scan distributed users for get shared link: %w", err)
		}

		link.Users = append(link.Users, userID)
	}

	return link, nil
}

func (s *shareRepoPG) distributeAccess(ctx context.Context, tx pgx.Tx, linkID int, users []int) error {
	insertBuilder := s.sqlBuilder.Insert("access_link").Columns("link_id", "user_id")

	for _, user := range users {
		insertBuilder = insertBuilder.Values(linkID, user)
	}

	sqlRow, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build sql query on distribute access: %w", err)
	}

	_, err = tx.Exec(ctx, sqlRow, args...)
	if err != nil {
		return fmt.Errorf("insert rows: %w", err)
	}
	return nil
}

func addNewSharedLink(ctx context.Context, q querer, link *share.SharedLink) (int, error) {
	role, ok := getContributorRoleFromString(link.Role)
	if !ok {
		return 0, ErrUnknownRole
	}

	var newLinkId int
	err := q.QueryRow(ctx, InsertNewLink, link.BoardID, role).Scan(&newLinkId)
	if err != nil {
		return 0, fmt.Errorf("insert link: %w", err)
	}

	return newLinkId, nil
}
