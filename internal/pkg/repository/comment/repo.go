package comment

import (
	"context"
	"errors"
	"fmt"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/internal/pgtype"
)

//go:generate mockgen -destination=./mock/comment_mock.go -package=mock -source=repo.go Repository
type Repository interface {
	AddComment(ctx context.Context, comment *entity.Comment) (int, error)
	GetCommentByID(ctx context.Context, id int) (*entity.Comment, error)
	EditStatusCommentOnDeletedByID(ctx context.Context, id int) error
	GetCommensToPin(ctx context.Context, pinID, lastID, count int) ([]entity.Comment, error)
}

var ErrUserRequired = errors.New("the comment does not have its author specified")

type commentRepoPG struct {
	db pgtype.PgxPoolIface
}

func NewCommentRepoPG(db pgtype.PgxPoolIface) *commentRepoPG {
	return &commentRepoPG{db}
}

func (c *commentRepoPG) AddComment(ctx context.Context, comment *entity.Comment) (int, error) {
	if comment.Author == nil {
		return 0, ErrUserRequired
	}

	var idInsertedComment int
	err := c.db.QueryRow(ctx, InsertNewComment, comment.Author.ID, comment.PinID, comment.Content).
		Scan(&idInsertedComment)
	if err != nil {
		return 0, fmt.Errorf("add comment in storage: %w", err)
	}
	return idInsertedComment, nil
}

func (c *commentRepoPG) GetCommentByID(ctx context.Context, id int) (*entity.Comment, error) {
	comment := &entity.Comment{ID: id, Author: &user.User{}}

	err := c.db.QueryRow(ctx, SelectCommentByID, id).
		Scan(&comment.Author.ID, &comment.Author.Username, &comment.Author.Avatar, &comment.PinID, &comment.Content)
	if err != nil {
		return nil, fmt.Errorf("get comment by id from storage: %w", err)
	}

	return comment, nil
}

func (c *commentRepoPG) EditStatusCommentOnDeletedByID(ctx context.Context, id int) error {
	if _, err := c.db.Exec(ctx, UpdateCommentOnDeleted, id); err != nil {
		return fmt.Errorf("edit status comment on deleted comment by id from storage: %w", err)
	}
	return nil
}

func (c *commentRepoPG) GetCommensToPin(ctx context.Context, pinID, lastID, count int) ([]entity.Comment, error) {
	rows, err := c.db.Query(ctx, SelectCommentsByPinID, pinID, lastID, count)
	if err != nil {
		return nil, fmt.Errorf("get comments to pin from storage: %w", err)
	}
	defer rows.Close()

	cmts := make([]entity.Comment, 0, count)
	cmt := entity.Comment{
		Author: &user.User{},
		PinID:  pinID,
	}

	for rows.Next() {
		cmt.Author = &user.User{}
		err = rows.Scan(&cmt.ID, &cmt.Author.ID, &cmt.Author.Username, &cmt.Author.Avatar, &cmt.Content)
		if err != nil {
			return cmts, fmt.Errorf("scan a comment when getting comments on a pin: %w", err)
		}

		cmts = append(cmts, cmt)
	}
	return cmts, nil
}
