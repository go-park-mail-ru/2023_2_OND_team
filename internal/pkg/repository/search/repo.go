package search

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/search"
)

type Repository interface {
	GetFilteredUsers(ctx context.Context, opts *search.SearchOpts) ([]search.UserForSearch, error)
	GetFilteredPins(ctx context.Context, opts *search.SearchOpts) ([]search.PinForSearch, error)
	GetFilteredBoards(ctx context.Context, opts *search.SearchOpts) ([]search.BoardForSearch, error)
}
