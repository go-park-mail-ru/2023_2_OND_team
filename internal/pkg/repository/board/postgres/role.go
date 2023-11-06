package board

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
)

func (repo *boardRepoPG) RoleUserHaveOnThisBoard(ctx context.Context, boardID int, userID int) (board.UserRole, error) {
	row := repo.db.QueryRow(ctx, SelectAuthorOrContributorRole, userID, boardID)
	var (
		author int
		role   pgtype.Text
	)
	err := row.Scan(&author, &role)
	if err != nil {
		return 0, fmt.Errorf("scan select row for getting user role: %w", err)
	}
	if userID == author {
		return board.Author, nil
	}
	return getUserRole(role), nil
}

func getUserRole(role pgtype.Text) board.UserRole {
	if !role.Valid {
		return board.RegularUser
	}

	switch role.String {
	case "read-write":
		return board.ContributorForReading | board.ContributorForAdding
	case "read-only":
		return board.ContributorForReading
	default:
		return board.RegularUser
	}
}
