package subscription

import (
	"context"

	userEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
)

func (u *subscriptionUsecase) GetSubscriptionInfoForUser(ctx context.Context, subOpts *userEntity.SubscriptionOpts) ([]userEntity.SubscriptionUser, error) {
	if err := u.userRepo.CheckUserExistence(ctx, subOpts.UserID); err != nil {
		return nil, err
	}

	var (
		currUserID, _ = ctx.Value(auth.KeyCurrentUserID).(int)
		users         = make([]userEntity.SubscriptionUser, 0)
		err           error
	)

	switch subOpts.Filter {
	case "subscriptions":
		users, err = u.subRepo.GetUserSubscriptions(ctx, subOpts.UserID, subOpts.Count, subOpts.LastID, currUserID)
	case "subscribers":
		users, err = u.subRepo.GetUserSubscribers(ctx, subOpts.UserID, subOpts.Count, subOpts.LastID, currUserID)
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}
