package user

type userCredentials struct {
	Username string
	Password string
}

func NewCredentials() userCredentials {
	return userCredentials{}
}
