package message

import (
	"github.com/jackc/pgx/v5/pgtype"

	mess "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

func convertFeedMessage(feed *mess.FeedMessage) []entity.Message {
	res := make([]entity.Message, len(feed.Messages))
	for ind := range feed.Messages {
		res[ind] = entity.Message{
			ID:   int(feed.Messages[ind].GetId().Id),
			From: int(feed.Messages[ind].UserFrom),
			To:   int(feed.Messages[ind].UserTo),
			Content: pgtype.Text{
				String: feed.Messages[ind].GetContent(),
				Valid:  true,
			},
		}
	}
	return res
}

func convertFeedChat(feed *mess.FeedChat) entity.FeedUserChats {
	res := make(entity.FeedUserChats, len(feed.Chats))

	for ind := range feed.Chats {
		res[ind] = entity.ChatWithUser{
			MessageLastID: int(feed.Chats[ind].GetLastMessageID()),
			WichWhomChat: user.User{
				ID:       int(feed.Chats[ind].GetChat().GetUserID()),
				Username: feed.Chats[ind].GetChat().GetUsername(),
				Avatar:   feed.Chats[ind].GetChat().GetAvatar(),
			},
		}
	}
	return res
}
