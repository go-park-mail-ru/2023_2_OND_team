package pin

import (
	sq "github.com/Masterminds/squirrel"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

func addFilters(queryBuild sq.SelectBuilder, cfg entity.FeedPinConfig, pin *entity.Pin,
	fields []any) (sq.SelectBuilder, []any) {

	queryBuild, fields = protection(queryBuild, cfg, pin, fields)
	queryBuild = addFilterDeleted(queryBuild, cfg.Deleted)
	queryBuild, ok := addFilterLiked(queryBuild, cfg)
	if !ok {
		queryBuild = addFilterUser(queryBuild, cfg)
	}
	queryBuild = addFilterBoard(queryBuild, cfg)
	return queryBuild, fields
}

func protection(queryBuild sq.SelectBuilder, cfg entity.FeedPinConfig, pin *entity.Pin,
	fields []any) (sq.SelectBuilder, []any) {

	switch cfg.Protection {
	case entity.FeedAll:
		queryBuild = queryBuild.Columns("pin.public")
		fields = append(fields, &pin.Public)
	case entity.FeedProtectionPublic:
		queryBuild = queryBuild.Where(sq.Eq{"pin.public": true})
		pin.Public = true
	case entity.FeedProtectionPrivate:
		queryBuild = queryBuild.Where(sq.Eq{"pin.public": false})
	}
	return queryBuild, fields
}

func addFilterDeleted(queryBuild sq.SelectBuilder, deleted bool) sq.SelectBuilder {
	if deleted {
		queryBuild = queryBuild.Where(sq.NotEq{"pin.deleted_at": nil})
	} else {
		queryBuild = queryBuild.Where(sq.Eq{"pin.deleted_at": nil})
	}
	return queryBuild
}

func addFilterBoard(queryBuild sq.SelectBuilder, cfg entity.FeedPinConfig) sq.SelectBuilder {
	if boardID, ok := cfg.Board(); ok {
		queryBuild = queryBuild.InnerJoin("membership ON membership.pin_id = pin.id").
			InnerJoin("board ON membership.board_id = board.id").
			Where(sq.Eq{"board.id": boardID})
	}
	return queryBuild
}

func addFilterLiked(queryBuild sq.SelectBuilder, cfg entity.FeedPinConfig) (sq.SelectBuilder, bool) {
	if userID, ok := cfg.User(); ok && cfg.Liked {
		queryBuild = queryBuild.InnerJoin("like_pin ON like_pin.pin_id = pin.id").
			Where(sq.Eq{"like_pin.user_id": userID}).
			OrderBy("like_pin.created_at DESC")
		return queryBuild, true
	}
	return queryBuild, false
}

func addFilterUser(queryBuild sq.SelectBuilder, cfg entity.FeedPinConfig) sq.SelectBuilder {
	if userID, ok := cfg.User(); ok {
		queryBuild = queryBuild.Where(sq.Eq{"pin.author": userID})
	}
	return queryBuild
}
