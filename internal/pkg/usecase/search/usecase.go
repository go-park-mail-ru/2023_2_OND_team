package search

import (
	"context"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/search"
	sRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/search"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

//go:generate mockgen -destination=./mock/search_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	GetUsers(ctx context.Context, opts *search.SearchOpts) ([]search.UserForSearch, error)
	GetBoards(ctx context.Context, opts *search.SearchOpts) ([]search.BoardForSearch, error)
	GetPins(ctx context.Context, opts *search.SearchOpts) ([]search.PinForSearch, error)
}

type searchUsecase struct {
	log        *logger.Logger
	searchRepo sRepo.Repository
}

func New(log *logger.Logger, searchRepo sRepo.Repository) Usecase {
	return &searchUsecase{log: log, searchRepo: searchRepo}
}

func (u *searchUsecase) GetUsers(ctx context.Context, opts *search.SearchOpts) ([]search.UserForSearch, error) {
	users, err := u.searchRepo.GetFilteredUsers(ctx, opts)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *searchUsecase) GetBoards(ctx context.Context, opts *search.SearchOpts) ([]search.BoardForSearch, error) {
	boards, err := u.searchRepo.GetFilteredBoards(ctx, opts)
	if err != nil {
		return nil, err
	}
	return boards, nil
}

func (u *searchUsecase) GetPins(ctx context.Context, opts *search.SearchOpts) ([]search.PinForSearch, error) {
	pins, err := u.searchRepo.GetFilteredPins(ctx, opts)
	if err != nil {
		return nil, err
	}
	return pins, nil
}
