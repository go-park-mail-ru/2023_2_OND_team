package session

import "time"

type Session struct {
	Key    string
	UserID int
	Expire time.Time
}
