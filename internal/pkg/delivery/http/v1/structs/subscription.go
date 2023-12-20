package structs

import errHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/errors"

//go:generate easyjson subscription.go

//easyjson:json
type SubscriptionAction struct {
	To *int `json:"to" example:"2"`
}

func (s *SubscriptionAction) Validate() error {
	if s.To == nil {
		return &errHTTP.ErrMissingBodyParams{Params: []string{"to"}}
	}
	return nil
}
