package comment

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotAvailableAction = errors.New("action not available for user")

func (c *commentCase) isAvailableCommentForDelete(ctx context.Context, userID, commentID int) error {
	comment, err := c.repo.GetCommentByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("get comment for check available comment for delete: %w", err)
	}

	if comment.Author.ID == userID {
		return nil
	}

	authorPinID, err := c.GetAuthorIdOfThePin(ctx, comment.PinID)
	if err != nil {
		return fmt.Errorf("get author pin for check availabel comment: %w", err)
	}
	if authorPinID != userID {
		return ErrNotAvailableAction
	}
	return nil
}
