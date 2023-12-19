package user

//go:generate easyjson info.go

//easyjson:json
type ProfileUpdateData struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	AboutMe  *string `json:"about_me"`
	Password *string `json:"password"`
}
