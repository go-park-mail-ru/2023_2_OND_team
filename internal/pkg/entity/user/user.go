package user

import "time"

type User struct {
	ID       int
	Username string
	Name     string
	Surname  string
	Email    string
	Avatar   string
	Password string
	Birthday time.Time
}
