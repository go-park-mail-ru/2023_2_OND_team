package v1

import (
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/structs"
	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/search"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	userEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
)

type converterHTTP struct {
	sanitizer validation.SanitizerXSS
	censor    validation.ProfanityCensor
}

func NewConverterHTTP(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) converterHTTP {
	return converterHTTP{sanitizer, censor}
}

func (c *converterHTTP) ToCertainBoardFromService(board *entity.BoardWithContent) structs.CertainBoard {
	b := structs.CertainBoard{
		ID:          board.BoardInfo.ID,
		AuthorID:    board.BoardInfo.AuthorID,
		Title:       board.BoardInfo.Title,
		Description: board.BoardInfo.Description,
		CreatedAt:   board.BoardInfo.CreatedAt.Format(TimeFormat),
		PinsNumber:  board.PinsNumber,
		Pins:        board.Pins,
		Tags:        board.TagTitles,
	}
	b.Sanitize(c.sanitizer, c.censor)
	return b
}

func (c *converterHTTP) ToCertainBoardUsernameFromService(board *entity.BoardWithContent, username string) structs.CertainBoardWithUsername {
	b := structs.CertainBoardWithUsername{
		ID:             board.BoardInfo.ID,
		AuthorID:       board.BoardInfo.AuthorID,
		AuthorUsername: username,
		Title:          board.BoardInfo.Title,
		Description:    board.BoardInfo.Description,
		CreatedAt:      board.BoardInfo.CreatedAt.Format(TimeFormat),
		PinsNumber:     board.PinsNumber,
		Pins:           board.Pins,
		Tags:           board.TagTitles,
	}
	b.Sanitize(c.sanitizer, c.censor)
	return b
}

func (c *converterHTTP) ToBoardFromService(board *entity.Board) *entity.Board {
	board.Sanitize(c.sanitizer, c.censor)
	return board
}

func (c *converterHTTP) ToUsersForSearchFromService(users []search.UserForSearch) []search.UserForSearch {
	for id := range users {
		users[id].Sanitize(c.sanitizer, c.censor)
	}
	return users
}

func (c *converterHTTP) ToBoardsForSearchFromService(boards []search.BoardForSearch) []search.BoardForSearch {
	for id := range boards {
		boards[id].Sanitize(c.sanitizer, c.censor)
	}
	return boards
}

func (c *converterHTTP) ToPinsForSearchFromService(pins []search.PinForSearch) []search.PinForSearch {
	for id := range pins {
		pins[id].Sanitize(c.sanitizer, c.censor)
	}
	return pins
}

func (c *converterHTTP) ToSubscriptionUsersFromService(users []user.SubscriptionUser) []user.SubscriptionUser {
	for id := range users {
		users[id].Sanitize(c.sanitizer, c.censor)
	}
	return users
}

func (c *converterHTTP) ToUserInfoFromService(user *userEntity.User, isSubscribed bool, subsCount int) structs.UserInfo {
	u := structs.UserInfo{
		ID:           user.ID,
		Username:     user.Username,
		Avatar:       user.Avatar,
		Name:         user.Name.String,
		Surname:      user.Surname.String,
		About:        user.AboutMe.String,
		IsSubscribed: isSubscribed,
		SubsCount:    subsCount,
	}
	u.Sanitize(c.sanitizer, c.censor)
	return u
}

func (c *converterHTTP) ToProfileInfoFromService(user *userEntity.User, subsCount int) structs.ProfileInfo {
	p := structs.ProfileInfo{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		SubsCount: subsCount,
	}
	p.Sanitize(c.sanitizer, c.censor)
	return p
}

func (c *converterHTTP) ToUserFromService(user *userEntity.User) *userEntity.User {
	user.Sanitize(c.sanitizer, c.censor)
	return user
}

func (c *converterHTTP) ToPinFromService(pin *pin.Pin) *pin.Pin {
	pin.Sanitize(c.sanitizer, c.censor)
	return pin
}

func (c *converterHTTP) ToMessagesFromService(mes []message.Message) []message.Message {
	for id := range mes {
		mes[id].Sanitize(c.sanitizer, c.censor)
	}
	return mes
}

func (c *converterHTTP) ToCommentsFromService(comments []comment.Comment) []comment.Comment {
	for id := range comments {
		comments[id].Sanitize(c.sanitizer, c.censor)
	}
	return comments
}
