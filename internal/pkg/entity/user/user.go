package user

type User struct {
	ID       int    `json:"-" example:"123"`
	Username string `json:"username" example:"Green"`
	Name     string `json:"-" example:"Peter"`
	Surname  string `json:"-" example:"Green"`
	Email    string `json:"email" example:"digital@gmail.com"`
	Avatar   string `json:"-" example:"pinspire.online/avatars/avatar.jpg"`
	Password string `json:"password,omitempty" example:"pass123"`
} // @name User
