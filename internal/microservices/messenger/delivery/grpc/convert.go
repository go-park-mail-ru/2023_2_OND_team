package grpc

import (
	mess "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/messenger"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
)

func convertFeedChat(feed message.FeedUserChats) []*mess.ChatWithUser {
	res := make([]*mess.ChatWithUser, len(feed))

	for ind := range feed {
		res[ind] = &mess.ChatWithUser{
			LastMessageID: int64(feed[ind].MessageLastID),
			Chat: &mess.WichWhomChat{
				UserID:   int64(feed[ind].WichWhomChat.ID),
				Username: feed[ind].WichWhomChat.Username,
				Avatar:   feed[ind].WichWhomChat.Avatar,
			},
		}
	}
	return res
}

func convertFeedMessage(feed []message.Message) []*mess.Message {
	res := make([]*mess.Message, len(feed))

	for ind := range res {
		res[ind] = &mess.Message{
			Id:       &mess.MsgID{Id: int64(feed[ind].ID)},
			UserFrom: int64(feed[ind].From),
			UserTo:   int64(feed[ind].To),
			Content:  feed[ind].Content.String,
		}
	}
	return res
}
