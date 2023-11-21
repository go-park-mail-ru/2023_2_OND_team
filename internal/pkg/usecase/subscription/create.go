package subscription

import "context"

func (u *subscriptionUsecase) SubscribeToUser(ctx context.Context, from, to int) error {
	if from == to {
		return &ErrSelfSubscription{}
	}

	if err := u.userRepo.CheckUserExistence(ctx, to); err != nil {
		return err
	}

	return u.subRepo.CreateSubscriptionUser(ctx, from, to)
}
