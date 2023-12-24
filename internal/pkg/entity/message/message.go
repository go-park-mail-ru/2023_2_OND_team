package message

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/validation"
)

//go:generate easyjson message.go

type Chat [2]int

//easyjson:json
type Message struct {
	ID      int         `json:"id,omitempty"`
	From    int         `json:"from"`
	To      int         `json:"to"`
	Content pgtype.Text `json:"content"`
}

func (m *Message) Sanitize(sanitizer validation.SanitizerXSS, censor validation.ProfanityCensor) {
	if m != nil {
		m.Content = pgtype.Text{
			String: sanitizer.Sanitize(m.Content.String),
			Valid:  m.Content.Valid,
		}
	}
}

func (m Message) WhatChat() Chat {
	return Chat{m.From, m.To}
}

type FeedUserChats []ChatWithUser

type ChatWithUser struct {
	MessageLastID int       `json:"-"`
	WichWhomChat  user.User `json:"user"`
}
