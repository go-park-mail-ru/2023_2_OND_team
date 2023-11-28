package search

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/search"
	sRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/search"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/microcosm-cc/bluemonday"
)

//ctx context.Context, template string, currUserID, lastID, count int

type Usecase interface {
	GetUsers(ctx context.Context, opts *search.SearchOpts) ([]search.UserForSearch, error)
	GetBoards(ctx context.Context, opts *search.SearchOpts) ([]search.BoardForSearch, error)
	GetPins(ctx context.Context, opts *search.SearchOpts) ([]search.PinForSearch, error)
}

type searchUsecase struct {
	log        *logger.Logger
	searchRepo sRepo.Repository
	sanitizer  *bluemonday.Policy
}

func New(log *logger.Logger, searchRepo sRepo.Repository, sanitizer *bluemonday.Policy) Usecase {
	return &searchUsecase{log: log, searchRepo: searchRepo, sanitizer: sanitizer}
}

func (u *searchUsecase) GetUsers(ctx context.Context, opts *search.SearchOpts) ([]search.UserForSearch, error) {
	users, err := u.searchRepo.GetFilteredUsers(ctx, opts)
	if err != nil {
		return nil, err
	}

	for id := range users {
		users[id].Sanitize(u.sanitizer)
	}

	return users, nil
}

func (u *searchUsecase) GetBoards(ctx context.Context, opts *search.SearchOpts) ([]search.BoardForSearch, error) {
	boards, err := u.searchRepo.GetFilteredBoards(ctx, opts)
	if err != nil {
		return nil, err
	}

	for id := range boards {
		boards[id].Sanitize(u.sanitizer)
	}

	return boards, nil
}

func (u *searchUsecase) GetPins(ctx context.Context, opts *search.SearchOpts) ([]search.PinForSearch, error) {
	pins, err := u.searchRepo.GetFilteredPins(ctx, opts)
	if err != nil {
		return nil, err
	}

	for id := range pins {
		pins[id].Sanitize(u.sanitizer)
	}

	return pins, nil
}
