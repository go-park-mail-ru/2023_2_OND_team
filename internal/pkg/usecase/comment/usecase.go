package comment

import (
	"context"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	commentRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/realtime/notification"
)

//go:generate mockgen -destination=./mock/comment_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	PutCommentOnPin(ctx context.Context, userID int, comment *entity.Comment) (int, error)
	GetFeedCommentOnPin(ctx context.Context, userID, pinID, count, lastID int) ([]entity.Comment, int, error)
	DeleteComment(ctx context.Context, userID, commentID int) error
	GetCommentWithAuthor(ctx context.Context, commentID int) (*entity.Comment, error)
}

type availablePinChecker interface {
	IsAvailablePinForViewingUser(ctx context.Context, userID, pinID int) error
	GetAuthorIdOfThePin(ctx context.Context, pinID int) (int, error)
}

type commentCase struct {
	availablePinChecker

	notifyCase notification.Usecase
	repo       commentRepo.Repository

	notifyIsEnable bool
}

func New(repo commentRepo.Repository, checker availablePinChecker, notifyCase notification.Usecase) *commentCase {
	comCase := &commentCase{
		availablePinChecker: checker,
		repo:                repo,
		notifyCase:          notifyCase,
	}

	if notifyCase != nil {
		comCase.notifyIsEnable = true
	}
	return comCase
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

	ctx = context.Background()

	if c.notifyIsEnable {
		go c.notifyCase.NotifyCommentLeftOnPin(ctx, id)
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

func (c *commentCase) GetCommentWithAuthor(ctx context.Context, commentID int) (*entity.Comment, error) {
	comment, err := c.repo.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("get comment with author: %w", err)
	}

	return comment, nil
}
