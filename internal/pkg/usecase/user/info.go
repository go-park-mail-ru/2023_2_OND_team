package user

type ProfileUpdateData struct {
	Username *string
	Email    *string
	Name     *string
	Surname  *string
	AboutMe  *string `json:"about_me`
	Password *string
}
