package search

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/search"
)

func SetUserSortType(sq squirrel.SelectBuilder, opts *search.SearchOpts) squirrel.SelectBuilder {
	switch opts.SortBy {
	case "subscribers":
		return sq.OrderBy(fmt.Sprintf("subscribers %s", opts.General.SortOrder))
	default:
		return sq.OrderBy(fmt.Sprintf("p1.id %s", opts.General.SortOrder))
	}
}

func SetBoardSortType(sq squirrel.SelectBuilder, opts *search.SearchOpts) squirrel.SelectBuilder {
	switch opts.SortBy {
	case "pins":
		return sq.OrderBy(fmt.Sprintf("pins_number %s", opts.General.SortOrder))
	default:
		return sq.OrderBy(fmt.Sprintf("board.id %s", opts.General.SortOrder))
	}
}

func SetPinSortType(sq squirrel.SelectBuilder, opts *search.SearchOpts) squirrel.SelectBuilder {
	switch opts.SortBy {
	case "likes":
		return sq.OrderBy(fmt.Sprintf("likes %s", opts.General.SortOrder))
	default:
		return sq.OrderBy(fmt.Sprintf("p.id %s", opts.General.SortOrder))
	}
}

func (r *searchRepoPG) SelectBoardsForSearch(opts *search.SearchOpts) (string, []interface{}, error) {
	SelectBoardsForSearch := r.sqlBuilder.Select(
		"board.id",
		"board.title",
		"board.created_at",
		"COUNT(DISTINCT pin.id) FILTER (WHERE pin.deleted_at IS NULL) AS pins_number",
		"COALESCE((ARRAY_AGG(DISTINCT pin.picture) FILTER (WHERE pin.deleted_at IS NULL AND pin.picture IS NOT NULL))[:3], ARRAY[]::TEXT[]) AS pins",
	).From(
		"board",
	).LeftJoin(
		"membership ON board.id = membership.board_id",
	).LeftJoin(
		"pin ON membership.pin_id = pin.id",
	).Where(
		squirrel.Eq{"board.deleted_at": nil},
	).Where(
		squirrel.ILike{"board.title": defaultSearchTemplate(opts.General.Template).GetTempl()},
	).Where(
		fmt.Sprintf("(board.public OR board.author = %d OR %d IN (SELECT user_id FROM contributor WHERE board_id = board.id))", opts.General.CurrUserID, opts.General.CurrUserID),
	).GroupBy(
		"board.id",
		"board.title",
		"board.created_at",
	)

	SelectBoardsForSearch = SetBoardSortType(SelectBoardsForSearch, opts)
	SelectBoardsForSearch = SelectBoardsForSearch.Limit(uint64(opts.General.Count)).Offset(uint64(opts.General.Offset))

	return SelectBoardsForSearch.ToSql()
}

func (r *searchRepoPG) SelectUsersForSearch(opts *search.SearchOpts) (string, []interface{}, error) {
	SelectUsersForSearch := r.sqlBuilder.Select(
		"p1.id",
		"p1.username",
		"p1.avatar",
		"COUNT(s1.who) AS subscribers",
		"s2.who IS NOT NULL AS is_subscribed",
	).From(
		"profile p1",
	).LeftJoin(
		"subscription_user s1 ON p1.id = s1.whom",
	).LeftJoin(
		"profile p2 ON s1.who = p2.id",
	).LeftJoin(
		fmt.Sprintf("subscription_user s2 ON p1.id = s2.whom AND s2.who = %d", opts.General.CurrUserID),
	).Where(
		squirrel.And{
			squirrel.Eq{"p1.deleted_at": nil},
			squirrel.Eq{"p2.deleted_at": nil},
			squirrel.ILike{"p1.username": defaultSearchTemplate(opts.General.Template).GetTempl()},
		},
	).GroupBy(
		"p1.id",
		"p1.username",
		"p1.avatar",
		"s2.who IS NOT NULL",
	)

	SelectUsersForSearch = SetUserSortType(SelectUsersForSearch, opts)
	SelectUsersForSearch = SelectUsersForSearch.Limit(uint64(opts.General.Count)).Offset(uint64(opts.General.Offset))

	return SelectUsersForSearch.ToSql()
}

func (r *searchRepoPG) SelectPinsForSearch(opts *search.SearchOpts) (string, []interface{}, error) {
	SelectPinsForSearch := r.sqlBuilder.Select(
		"p.id",
		"p.title",
		"p.picture",
		"COUNT(*) AS likes",
	).From(
		"pin p",
	).LeftJoin(
		"like_pin lp ON p.id = lp.pin_id",
	).Where(
		squirrel.And{
			squirrel.Eq{"p.deleted_at": nil},
			squirrel.Or{
				squirrel.Eq{"p.public": true},
				squirrel.Eq{"p.author": opts.General.CurrUserID},
			},
			squirrel.ILike{"p.title": defaultSearchTemplate(opts.General.Template).GetTempl()},
		},
	).GroupBy(
		"p.id",
		"p.title",
		"p.picture",
	)

	SelectPinsForSearch = SetPinSortType(SelectPinsForSearch, opts)
	SelectPinsForSearch = SelectPinsForSearch.Limit(uint64(opts.General.Count)).Offset(uint64(opts.General.Offset))

	return SelectPinsForSearch.ToSql()
}
