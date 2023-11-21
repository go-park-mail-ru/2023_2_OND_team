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
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
)

type subscriptionRepoPG struct {
	db         pgtype.PgxPoolIface
	sqlBuilder squirrel.StatementBuilderType
}

func NewSubscriptionRepoPG(db pgtype.PgxPoolIface) subRepo.Repository {
	return &subscriptionRepoPG{db: db, sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func convertErrorPostgres(ctx context.Context, err error) error {
	logger := logger.GetLoggerFromCtx(ctx)

	if errors.Is(err, context.DeadlineExceeded) {
		return &errPkg.ErrTimeoutExceeded{}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		case strconv.Itoa(23505):
			return &subRepo.ErrSubscriptionAlreadyExist{}
		default:
			logger.Warnf("Unexpected error from subscription repo - postgres: %s\n", err.Error())
			return &errPkg.InternalError{}
		}
	}
	logger.Warnf("Unexpected error from subscription repo: %s\n", err.Error())
	return &errPkg.InternalError{}
}

func (r *subscriptionRepoPG) CreateSubscriptionUser(ctx context.Context, from, to int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return convertErrorPostgres(ctx, err)
	}

	if _, err = tx.Exec(ctx, CreateSubscriptionUser, from, to); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return convertErrorPostgres(ctx, err)
		}
		return convertErrorPostgres(ctx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return convertErrorPostgres(ctx, err)
	}
	return nil
}

func (r *subscriptionRepoPG) DeleteSubscriptionUser(ctx context.Context, from, to int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return convertErrorPostgres(ctx, err)
	}

	status, err := tx.Exec(ctx, DeleteSubscriptionUser, from, to)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return convertErrorPostgres(ctx, err)
		}
		return convertErrorPostgres(ctx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return convertErrorPostgres(ctx, err)
	}
	if status.RowsAffected() == 0 {
		return &subRepo.ErrNonExistingSubscription{}
	}
	return nil
}

func (r *subscriptionRepoPG) GetUserSubscriptions(ctx context.Context, userID, count, offset int, currUserID int) ([]userEntity.SubscriptionUser, error) {

	var getUserSubscriptions squirrel.SelectBuilder
	if currUserID != 0 {
		getUserSubscriptions = r.sqlBuilder.Select(
			"p.username, p.avatar, s.who IS NOT NULL AS is_subscribed",
		)
	} else {
		getUserSubscriptions = r.sqlBuilder.Select(
			"p.username, p.avatar, false AS is_subscribed",
		)
	}
	getUserSubscriptions = getUserSubscriptions.
		From("subscription_user f").
		LeftJoin("profile p ON f.whom = p.id").
		LeftJoin("subscription_user s ON f.whom = s.whom AND s.who = $1", currUserID).
		Where("f.who = $2", userID).
		Where("p.deleted_at IS NULL").
		OrderBy("f.whom ASC").
		Limit(uint64(count)).
		Offset(uint64(offset))

	sqlRow, args, err := getUserSubscriptions.ToSql()
	if err != nil {
		return nil, convertErrorPostgres(ctx, err)
	}

	rows, err := r.db.Query(ctx, sqlRow, args...)
	defer rows.Close()

	subscriptions := make([]userEntity.SubscriptionUser, 0)
	for rows.Next() {
		var subscription userEntity.SubscriptionUser
		if err = rows.Scan(&subscription.Username, &subscription.Avatar, &subscription.HasSubscribeFromCurUser); err != nil {
			return nil, convertErrorPostgres(ctx, err)
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

func (r *subscriptionRepoPG) GetUserSubscribers(ctx context.Context, userID, count, offset int, currUserID int) ([]userEntity.SubscriptionUser, error) {

	var getUserSubscribers squirrel.SelectBuilder
	if currUserID != 0 {
		getUserSubscribers = r.sqlBuilder.Select(
			"p.username, p.avatar, s.who IS NOT NULL AS is_subscribed",
		)
	} else {
		getUserSubscribers = r.sqlBuilder.Select(
			"p.username, p.avatar, false AS is_subscribed",
		)
	}
	getUserSubscribers = getUserSubscribers.
		From("subscription_user f").
		LeftJoin("profile p ON f.who = p.id").
		LeftJoin("subscription_user s ON f.who = s.whom AND s.who = $1", currUserID).
		Where("f.whom = $2", userID).
		Where("p.deleted_at IS NULL").
		OrderBy("f.who ASC").
		Limit(uint64(count)).
		Offset(uint64(offset))

	sqlRow, args, err := getUserSubscribers.ToSql()
	if err != nil {
		return nil, convertErrorPostgres(ctx, err)
	}

	rows, err := r.db.Query(ctx, sqlRow, args...)
	if err != nil {
		return nil, convertErrorPostgres(ctx, err)
	}
	defer rows.Close()

	subscribers := make([]userEntity.SubscriptionUser, 0)
	for rows.Next() {
		var subscriber userEntity.SubscriptionUser
		if err = rows.Scan(&subscriber.Username, &subscriber.Avatar, &subscriber.HasSubscribeFromCurUser); err != nil {
			return nil, convertErrorPostgres(ctx, err)
		}
		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}
