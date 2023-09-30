package user

import "time"

type User struct {
	ID       int       `json:"id" example:"123"`
	Username string    `json:"username" example:"Green"`
	Name     string    `json:"name" example:"Peter"`
	Surname  string    `json:"surname" example:"Green"`
	Email    string    `json:"email" example:"digital@gmail.com"`
	Avatar   string    `json:"avatar" example:"pinspire.online/avatars/avatar.jpg"`
	Password string    `json:"password" example:"pass123"`
	Birthday time.Time `json:"birthday"`
} // @name User
