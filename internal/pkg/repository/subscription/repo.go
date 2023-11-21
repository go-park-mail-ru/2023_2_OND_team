package subscription

import (
	"context"

	userEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

type Repository interface {
	CreateSubscriptionUser(ctx context.Context, from, to int) error
	DeleteSubscriptionUser(ctx context.Context, from, to int) error
	GetUserSubscriptions(ctx context.Context, userID, count, offset int, currUserID int) ([]userEntity.SubscriptionUser, error)
	GetUserSubscribers(ctx context.Context, userID, count, offset int, currUserID int) ([]userEntity.SubscriptionUser, error)
}
