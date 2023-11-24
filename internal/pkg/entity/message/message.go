package message

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Chat [2]int

type Message struct {
	ID      int
	From    int         `json:"from"`
	To      int         `json:"to"`
	Content pgtype.Text `json:"content"`
}

func (m Message) WhatChat() Chat {
	return Chat{m.From, m.To}
}

type FeedUserChats []ChatWithUser

type ChatWithUser struct {
	MessageLastID int       `json:"-"`
	WichWhomChat  user.User `json:"user"`
}
