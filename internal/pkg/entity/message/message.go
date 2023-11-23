package message

import "github.com/jackc/pgx/v5/pgtype"

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
