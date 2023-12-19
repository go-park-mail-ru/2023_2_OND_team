package comment

import (
	"context"
	"fmt"
	"strconv"

	comm "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/comment"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/notification"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/notification"
)

type commentGetter interface {
	GetCommentWithAuthor(ctx context.Context, commentID int) (*comm.Comment, error)
}

type pinGetter interface {
	GetPinWithAuthor(ctx context.Context, pinID int) (*pin.Pin, error)
}

type commentNotify struct {
	notification.NotifyBuilder

	com commentGetter
	pin pinGetter
}

func NewCommentNotify(builder notification.NotifyBuilder, com commentGetter, pin pinGetter) commentNotify {
	return commentNotify{builder, com, pin}
}

func (c commentNotify) Type() entity.NotifyType {
	return c.NotifyBuilder.Type()
}

func (c commentNotify) MessageNotify(data notification.M) (*entity.NotifyMessage, error) {
	return c.NotifyBuilder.BuildNotifyMessage(data)
}

func (c commentNotify) ChannelsNameForSubscribe(_ context.Context, userID int) ([]string, error) {
	return []string{strconv.Itoa(userID)}, nil
}

func (c commentNotify) ChannelNameForPublishWithData(ctx context.Context, commentID int) (string, notification.M, error) {
	com, err := c.com.GetCommentWithAuthor(ctx, commentID)
	if err != nil {
		return "", nil, fmt.Errorf("get comment for receive channel name on publish: %w", err)
	}

	pin, err := c.pin.GetPinWithAuthor(ctx, com.PinID)
	if err != nil {
		return "", nil, fmt.Errorf("get pin for receive channel name on publish: %w", err)
	}

	return strconv.Itoa(pin.Author.ID), notification.M{"Username": com.Author.Username, "TitlePin": pin.Title.String}, nil
}
