package user

type profileUpdateData struct {
	Username *string
	Email    *string
	Name     *string
	Surname  *string
	Password *string
}

func NewProfileUpdateData() *profileUpdateData {
	return &profileUpdateData{}
}
