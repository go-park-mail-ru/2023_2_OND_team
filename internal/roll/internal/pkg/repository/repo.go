package roll

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	roll "github.com/go-park-mail-ru/2023_2_OND_team/internal/roll/internal/pkg/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	InsertRollAnswer(ctx context.Context, ans []roll.RollAnswer) error
	CheckUserFilledRoll(ctx context.Context, userID, rollID int) (bool, error)
	GetHistStat(ctx context.Context, rollID, questionID int) ([]roll.HistStatObj, error)
}

type rollRepoPG struct {
	db *pgxpool.Pool
	sq squirrel.StatementBuilderType
}

func NewRollRepoPG(db *pgxpool.Pool) *rollRepoPG {
	return &rollRepoPG{db: db, sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

/*
var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		case strconv.Itoa(23505):
			return &subRepo.ErrSubscriptionAlreadyExist{}
		}
	}
*/

/*
insertTagsQuery := repo.sqlBuilder.
		Insert("tag").
		Columns("title")
	for _, title := range titles {
		insertTagsQuery = insertTagsQuery.Values(title)
	}
*/

func (r *rollRepoPG) InsertRollAnswer(ctx context.Context, ans []roll.RollAnswer) error {
	insertAnsQuery := r.sq.Insert("roll").Columns("id", "user_id", "question_id", "answer")

	for _, answer := range ans {
		insertAnsQuery = insertAnsQuery.Values(answer.RollID, answer.UserID, answer.QuestionID, answer.Answer)
	}

	sqlRow, args, err := insertAnsQuery.ToSql()
	if err != nil {
		return fmt.Errorf("insert poll - build query: %w", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("insert roll - begin tx: %w", err)
	}

	if _, err = tx.Exec(ctx, sqlRow, args...); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("insert roll - rollback tx: %w", err)
		}
		return fmt.Errorf("insert roll: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("insert roll - coomit tx: %w", err)
	}

	return nil
}

func (r *rollRepoPG) InsertRollAnswer_(ctx context.Context, rollID, currUserID, questionID int, answer string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("insert roll - begin tx: %w", err)
	}

	if _, err = tx.Exec(ctx, InsertRollAnswer, rollID, currUserID, questionID, answer); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("insert roll - rollback tx: %w", err)
		}
		return fmt.Errorf("insert roll: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("insert roll - coomit tx: %w", err)
	}

	return nil
}

func (r *rollRepoPG) CheckUserFilledRoll(ctx context.Context, userID, rollID int) (bool, error) {
	var dummy int
	err := r.db.QueryRow(ctx, CheckUserFilledRoll, rollID, userID).Scan(&dummy)
	switch err {
	case pgx.ErrNoRows:
		return false, nil
	default:
		return false, fmt.Errorf("check filled poll - %w", err)
	}
}

func (r *rollRepoPG) GetHistStat(ctx context.Context, rollID, questionID int) ([]roll.HistStatObj, error) {
	rows, err := r.db.Query(ctx, SelectHistStat, rollID, questionID)
	if err != nil {
		return nil, fmt.Errorf("select stat - %w", err)
	}
	defer rows.Close()

	stats := make([]roll.HistStatObj, 0)
	for rows.Next() {
		var stat roll.HistStatObj
		if err := rows.Scan(&stat); err != nil {
			return nil, fmt.Errorf("scan hist stat: %w", err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}
