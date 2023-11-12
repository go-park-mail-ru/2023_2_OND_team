package board

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board"
)

func (b *boardUsecase) CheckAvailabilityFeedPinCfgOnBoard(ctx context.Context, cfg pin.FeedPinConfig,
	userID int, isAuth bool) error {

	boardID, ok := cfg.Board()
	if !ok {
		return nil
	}

	if !isAuth && cfg.Protection != pin.FeedProtectionPublic {
		return ErrNoAccess
	}

	protection, err := b.boardRepo.GerProtectionStatusBoard(ctx, boardID)
	if err != nil {
		return fmt.Errorf("get protection status board for check availability: %w", err)
	}

	if !isAuth && protection != board.ProtectionPublic {
		return ErrNoAccess
	}
	if !isAuth {
		return nil
	}

	role, err := b.boardRepo.RoleUserHaveOnThisBoard(ctx, boardID, userID)
	if err != nil {
		return fmt.Errorf("get user role of the board for check availability: %w", err)
	}

	if role&(board.Author|board.ContributorForAdding|board.ContributorForReading) == 0 &&
		(cfg.Protection != pin.FeedProtectionPublic || protection != board.ProtectionPublic) {

		return ErrNoAccess
	}
	return nil
}
