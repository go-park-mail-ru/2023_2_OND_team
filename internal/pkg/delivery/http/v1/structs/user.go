package structs

//go:generate easyjson user.go

//easyjson:json
type UserInfo struct {
	ID           int    `json:"id" example:"123"`
	Username     string `json:"username" example:"Snapshot"`
	Avatar       string `json:"avatar" example:"/pic1"`
	Name         string `json:"name" example:"Bob"`
	Surname      string `json:"surname" example:"Dylan"`
	About        string `json:"about" example:"Cool guy"`
	IsSubscribed bool   `json:"is_subscribed" example:"true"`
	SubsCount    int    `json:"subscribers" example:"23"`
}

//easyjson:json
type ProfileInfo struct {
	ID        int    `json:"id" example:"1"`
	Username  string `json:"username" example:"baobab"`
	Avatar    string `json:"avatar" example:"/pic1"`
	SubsCount int    `json:"subscribers" example:"12"`
}
