package subscription

import (
	"context"
	"errors"
	"strconv"

	"github.com/Masterminds/squirrel"
	userEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/internal/pgtype"
	subRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/subscription"
	"github.com/jackc/pgx/v5/pgconn"
)

type subscriptionRepoPG struct {
	db         pgtype.PgxPoolIface
	sqlBuilder squirrel.StatementBuilderType
}

func NewSubscriptionRepoPG(db pgtype.PgxPoolIface) subRepo.Repository {
	return &subscriptionRepoPG{db: db, sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func convertErrorPostgres(err error) error {

	switch err {
	case context.DeadlineExceeded:
		return &errPkg.ErrTimeoutExceeded{}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		case strconv.Itoa(23505):
			return &subRepo.ErrSubscriptionAlreadyExist{}
		default:
			return &errPkg.InternalError{Message: err.Error(), Layer: string(errPkg.Repo)}
		}
	}
	return &errPkg.InternalError{Message: err.Error(), Layer: string(errPkg.Repo)}
}

func (r *subscriptionRepoPG) CreateSubscriptionUser(ctx context.Context, from, to int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return convertErrorPostgres(err)
	}

	if _, err = tx.Exec(ctx, CreateSubscriptionUser, from, to); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return convertErrorPostgres(err)
		}
		return convertErrorPostgres(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return convertErrorPostgres(err)
	}
	return nil
}

func (r *subscriptionRepoPG) DeleteSubscriptionUser(ctx context.Context, from, to int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return convertErrorPostgres(err)
	}

	status, err := tx.Exec(ctx, DeleteSubscriptionUser, from, to)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return convertErrorPostgres(err)
		}
		return convertErrorPostgres(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return convertErrorPostgres(err)
	}
	if status.RowsAffected() == 0 {
		return &subRepo.ErrNonExistingSubscription{}
	}
	return nil
}

func (r *subscriptionRepoPG) GetUserSubscriptions(ctx context.Context, userID, count, lastID int, currUserID int) ([]userEntity.SubscriptionUser, error) {

	rows, err := r.db.Query(ctx, GetUserSubscriptions, currUserID, userID, lastID, count)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	defer rows.Close()

	subscriptions := make([]userEntity.SubscriptionUser, 0)
	for rows.Next() {
		var subscription userEntity.SubscriptionUser
		if err = rows.Scan(&subscription.ID, &subscription.Username, &subscription.Avatar, &subscription.HasSubscribeFromCurUser); err != nil {
			return nil, convertErrorPostgres(err)
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

func (r *subscriptionRepoPG) GetUserSubscribers(ctx context.Context, userID, count, lastID int, currUserID int) ([]userEntity.SubscriptionUser, error) {

	rows, err := r.db.Query(ctx, GetUserSubscribers, currUserID, userID, lastID, count)
	if err != nil {
		return nil, convertErrorPostgres(err)
	}
	defer rows.Close()

	subscribers := make([]userEntity.SubscriptionUser, 0)
	for rows.Next() {
		var subscriber userEntity.SubscriptionUser
		if err = rows.Scan(&subscriber.ID, &subscriber.Username, &subscriber.Avatar, &subscriber.HasSubscribeFromCurUser); err != nil {
			return nil, convertErrorPostgres(err)
		}
		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}
