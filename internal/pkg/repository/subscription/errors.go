package subscription

import errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"

type ErrSubscriptionAlreadyExist struct{}

func (e *ErrSubscriptionAlreadyExist) Error() string {
	return "subscription on that user already exists"
}

func (e *ErrSubscriptionAlreadyExist) Type() errPkg.Type {
	return errPkg.ErrAlreadyExists
}

type ErrNonExistingSubscription struct{}

func (e *ErrNonExistingSubscription) Error() string {
	return "such subscription doesn't exist"
}

func (e *ErrNonExistingSubscription) Type() errPkg.Type {
	return errPkg.ErrNotFound
}
