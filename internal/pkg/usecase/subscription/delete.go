package subscription

import "context"

func (u *subscriptionUsecase) UnsubscribeFromUser(ctx context.Context, from, to int) error {
	if from == to {
		return &ErrSelfUnsubscription{}
	}

	if err := u.userRepo.CheckUserExistence(ctx, to); err != nil {
		return err
	}

	return u.subRepo.DeleteSubscriptionUser(ctx, from, to)
}
