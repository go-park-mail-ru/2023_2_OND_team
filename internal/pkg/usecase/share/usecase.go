package share

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/share"
	shareRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/share"
)

var (
	ErrInvalidSharedLinkParam = errors.New("invalid link parameters")
	ErrNoAccess               = errors.New("there is no access")
	ErrExceededLimit          = errors.New("the limit of allowed users has been exceeded")
)

const _limitUserDdistributeAccess = 100

type Usecase interface {
	CreateSharedLinkForAddingContributor(ctx context.Context, userID int, link *entity.SharedLink) (int, error)
	CheckLinkAvailability(ctx context.Context, userID, linkID int) (string, error)
}

type boardRepository interface {
	GetBoardAuthorByBoardID(ctx context.Context, boardID int) (int, error)
}

type shareCase struct {
	repo      shareRepo.Repository
	boardRepo boardRepository
}

func New(s shareRepo.Repository, b boardRepository) *shareCase {
	return &shareCase{
		repo:      s,
		boardRepo: b,
	}
}

func (s *shareCase) CreateSharedLinkForAddingContributor(ctx context.Context, userID int, link *entity.SharedLink) (int, error) {
	if link.IsDistributedAll && len(link.Users) != 0 {
		return 0, ErrInvalidSharedLinkParam
	}

	if len(link.Users) > _limitUserDdistributeAccess {
		return 0, ErrExceededLimit
	}

	authorID, err := s.boardRepo.GetBoardAuthorByBoardID(ctx, link.BoardID)
	if err != nil {
		return 0, fmt.Errorf("get board author for check access creating shared link: %w", err)
	}
	if authorID != userID {
		return 0, ErrNoAccess
	}

	id, err := s.repo.AddSharedLink(ctx, link)
	if err != nil {
		return 0, fmt.Errorf("create shared link for adding contributor: %w", err)
	}

	return id, nil
}

func (s *shareCase) CheckLinkAvailability(ctx context.Context, userID, linkID int) (string, error) {
	link, err := s.repo.GetSharedLink(ctx, linkID)
	if err != nil {
		return "", fmt.Errorf("check link availability: %w", err)
	}

	if link.IsDistributedAll {
		return link.Role, nil
	}

	var inAccessList bool
	for _, user := range link.Users {
		if userID == user {
			inAccessList = true
		}
	}

	if !inAccessList {
		return "", ErrNoAccess
	}
	return link.Role, nil
}
