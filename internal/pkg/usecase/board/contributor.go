package board

import (
	"context"
)

func (bCase *boardUsecase) AddContributorsToBoard(ctx context.Context, boardID int, usersId []int, role string) error {
	return bCase.boardRepo.AddContributors(ctx, boardID, usersId, role)
}
