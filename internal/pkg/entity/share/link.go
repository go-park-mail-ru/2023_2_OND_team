package share

//go:generate easyjson
//easyjson:json
type SharedLink struct {
	BoardID          int    `json:"-"`
	Role             string `json:"role" enum:"read-write,read-only"`
	Users            []int  `json:"users_distributed"`
	IsDistributedAll bool   `json:"is_distributed_all"`
}
