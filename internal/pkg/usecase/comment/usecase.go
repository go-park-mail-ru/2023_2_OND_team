package comment

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	commentRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/comment"
)

//go:generate mockgen -destination=./mock/comment_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	PutCommentOnPin(ctx context.Context, userID int, comment *entity.Comment) (int, error)
	GetFeedCommentOnPin(ctx context.Context, userID, pinID, count, lastID int) ([]entity.Comment, int, error)
	DeleteComment(ctx context.Context, userID, commentID int) error
}

type availablePinChecker interface {
	IsAvailablePinForViewingUser(ctx context.Context, userID, pinID int) error
	GetAuthorIdOfThePin(ctx context.Context, pinID int) (int, error)
}

type commentCase struct {
	availablePinChecker

	repo commentRepo.Repository
}

func New(repo commentRepo.Repository, checker availablePinChecker) *commentCase {
	return &commentCase{checker, repo}
}

func (c *commentCase) PutCommentOnPin(ctx context.Context, userID int, comment *entity.Comment) (int, error) {
	err := c.IsAvailablePinForViewingUser(ctx, userID, comment.PinID)
	if err != nil {
		return 0, fmt.Errorf("put comment on not available pin: %w", err)
	}

	comment.Author = &user.User{ID: userID}

	id, err := c.repo.AddComment(ctx, comment)
	if err != nil {
		return 0, fmt.Errorf("put comment on available pin: %w", err)
	}
	return id, nil
}

func (c *commentCase) GetFeedCommentOnPin(ctx context.Context, userID, pinID, count, lastID int) ([]entity.Comment, int, error) {
	err := c.IsAvailablePinForViewingUser(ctx, userID, pinID)
	if err != nil {
		return nil, 0, fmt.Errorf("put comment on not available pin: %w", err)
	}

	feed, err := c.repo.GetCommensToPin(ctx, pinID, lastID, count)
	if err != nil {
		err = fmt.Errorf("get feed comment on pin: %w", err)
	}

	var newLastID int
	if len(feed) > 0 {
		newLastID = feed[len(feed)-1].ID
	}
	return feed, newLastID, err
}

func (c *commentCase) DeleteComment(ctx context.Context, userID, commentID int) error {
	err := c.isAvailableCommentForDelete(ctx, userID, commentID)
	if err != nil {
		return fmt.Errorf("check available delete comment: %w", err)
	}

	err = c.repo.EditStatusCommentOnDeletedByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	return nil
}
