package subscription

import (
	"context"

	userEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	subRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/subscription"
	uRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

//go:generate mockgen -destination=./mock/subscription_mock.go -package=mock -source=usecase.go Usecase
type Usecase interface {
	SubscribeToUser(ctx context.Context, from, to int) error
	UnsubscribeFromUser(ctx context.Context, from, to int) error
	GetSubscriptionInfoForUser(ctx context.Context, subOpts *userEntity.SubscriptionOpts) ([]userEntity.SubscriptionUser, error)
}

type subscriptionUsecase struct {
	subRepo  subRepo.Repository
	userRepo uRepo.Repository
	log      *logger.Logger
}

func New(log *logger.Logger, subRepo subRepo.Repository, uRepo uRepo.Repository) Usecase {
	return &subscriptionUsecase{subRepo: subRepo, userRepo: uRepo, log: log}
}
