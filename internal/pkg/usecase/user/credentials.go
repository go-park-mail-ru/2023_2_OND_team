package user

//go:generate easyjson credentials.go

//easyjson:json
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
