package user

type profileUpdateData struct {
	Username *string
	Email    *string
	Name     *string
	Surname  *string
	AboutMe  *string `json:"about_me`
	Password *string
}

func NewProfileUpdateData() *profileUpdateData {
	return &profileUpdateData{}
}
