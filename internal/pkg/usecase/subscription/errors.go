package subscription

import errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"

type ErrSelfSubscription struct{}

func (e *ErrSelfSubscription) Error() string {
	return "can't subscribe on yourself"
}

func (e *ErrSelfSubscription) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrSelfUnsubscription struct{}

func (e *ErrSelfUnsubscription) Error() string {
	return "can't unsubscribe from yourself"
}

func (e *ErrSelfUnsubscription) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}
