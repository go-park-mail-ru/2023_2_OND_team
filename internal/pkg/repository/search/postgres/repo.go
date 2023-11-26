package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/search"
	errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/internal/pgtype"
	searchRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/search"
	"github.com/jackc/pgx/v5/pgconn"
)

type defaultSearchTemplate string

func (t defaultSearchTemplate) GetTempl() string {
	return fmt.Sprintf("%%%s%%", t)
}

type searchRepoPG struct {
	db         pgtype.PgxPoolIface
	sqlBuilder squirrel.StatementBuilderType
}

func NewSearchRepoPG(db pgtype.PgxPoolIface) *searchRepoPG {
	return &searchRepoPG{db, squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func convertErrorPostgres(err error) error {

	switch err {
	case context.DeadlineExceeded:
		return &errPkg.ErrTimeoutExceeded{}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		}
	}

	return &errPkg.InternalError{Message: err.Error(), Layer: string(errPkg.Repo)}
}

func (r *searchRepoPG) GetFilteredUsers(ctx context.Context, opts *search.SearchOpts) ([]search.UserForSearch, error) {
	sqlRow, args, err := r.SelectUsersForSearch(opts)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	rows, err := r.db.Query(ctx, sqlRow, args...)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	defer rows.Close()

	users := make([]search.UserForSearch, 0)
	for rows.Next() {
		user := search.UserForSearch{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Avatar, &user.SubsCount, &user.HasSubscribeFromCurUser); err != nil {
			return nil, convertErrorPostgres(err)
		}
		users = append(users, user)
	}

	if len(users) == 0 && opts.General.Offset == 0 {
		return nil, &searchRepo.ErrNoUsers{}
	}

	return users, nil
}

func (r *searchRepoPG) GetFilteredBoards(ctx context.Context, opts *search.SearchOpts) ([]search.BoardForSearch, error) {

	sqlRow, args, err := r.SelectBoardsForSearch(opts)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	rows, err := r.db.Query(ctx, sqlRow, args...)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	defer rows.Close()

	boards := make([]search.BoardForSearch, 0)
	for rows.Next() {
		board := search.BoardForSearch{}
		if err := rows.Scan(&board.BoardHeader.ID, &board.BoardHeader.Title, &board.BoardHeader.CreatedAt, &board.PinsNumber, &board.PreviewPins); err != nil {
			return nil, convertErrorPostgres(err)
		}
		boards = append(boards, board)
	}

	if len(boards) == 0 && opts.General.Offset == 0 {
		return nil, &searchRepo.ErrNoBoards{}
	}

	return boards, nil
}

func (r *searchRepoPG) GetFilteredPins(ctx context.Context, opts *search.SearchOpts) ([]search.PinForSearch, error) {

	sqlRow, args, err := r.SelectPinsForSearch(opts)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	rows, err := r.db.Query(ctx, sqlRow, args...)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	defer rows.Close()

	pins := make([]search.PinForSearch, 0)
	for rows.Next() {
		pin := search.PinForSearch{}
		if err := rows.Scan(&pin.ID, &pin.Title, &pin.Picture, &pin.Likes); err != nil {
			return nil, convertErrorPostgres(err)
		}
		pins = append(pins, pin)
	}

	if len(pins) == 0 && opts.General.Offset == 0 {
		return nil, &searchRepo.ErrNoPins{}
	}

	return pins, nil

}
