package board

import (
	"context"
	"fmt"
	stdLog "log"
	"testing"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	mock_board "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board/mock"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/require"
)

type (
	CreateBoard             func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string)
	UpdateBoard             CreateBoard
	GetBoardAuthorByBoardID func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int)
)

var (
	sanitizer = bluemonday.UGCPolicy()
)

func TestBoardUsecase_CreateNewBoard(t *testing.T) {

	tests := []struct {
		name         string
		inCtx        context.Context
		newBoardData dto.BoardData
		CreateBoard  CreateBoard
		expNewID     int
		expErr       error
	}{
		{
			name:  "valid board data",
			inCtx: context.Background(),
			newBoardData: dto.BoardData{
				Title:       "valid title",
				Description: "some description",
				AuthorID:    45,
				Public:      false,
				TagTitles:   []string{"nice", "green"},
			},
			CreateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
				mockRepo.EXPECT().CreateBoard(ctx, newBoardData, tagTitles).Return(1, nil).Times(1)
			},
			expNewID: 1,
			expErr:   nil,
		},
		{
			name:  "invalid board title",
			inCtx: context.Background(),
			newBoardData: dto.BoardData{
				Title:       "~nval$d title~~",
				Description: "some description",
				AuthorID:    45,
				Public:      false,
				TagTitles:   []string{"nice", "green"},
			},
			CreateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expNewID: 0,
			expErr:   ErrInvalidBoardTitle,
		},
		{
			name:  "invalid tag titles: all tags",
			inCtx: context.Background(),
			newBoardData: dto.BoardData{
				Title:       "valid title",
				Description: "some description",
				AuthorID:    45,
				Public:      false,
				TagTitles:   []string{"nic~e", "gr~$een"},
			},
			CreateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expNewID: 0,
			expErr:   fmt.Errorf("%v: %w", []string{"nic~e", "gr~$een"}, ErrInvalidTagTitles),
		},
		{
			name:  "invalid tag titles: some tags",
			inCtx: context.Background(),
			newBoardData: dto.BoardData{
				Title:       "valid title",
				Description: "some description",
				AuthorID:    45,
				Public:      false,
				TagTitles:   []string{"nic~e", "green"},
			},
			CreateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expNewID: 0,
			expErr:   fmt.Errorf("%v: %w", []string{"nic~e"}, ErrInvalidTagTitles),
		},
		{
			name:  "invalid tag titles: too many tags",
			inCtx: context.Background(),
			newBoardData: dto.BoardData{
				Title:       "valid title",
				Description: "some description",
				AuthorID:    45,
				Public:      false,
				TagTitles:   []string{"nice", "green", "a", "b", "c", "d", "e", "f"},
			},
			CreateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expNewID: 0,
			expErr:   fmt.Errorf("too many titles: %w", ErrInvalidTagTitles),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			log, err := logger.New(logger.RFC3339FormatTime())
			if err != nil {
				stdLog.Fatal(err)
			}
			mockBoardRepo := mock_board.NewMockRepository(ctl)
			test.CreateBoard(mockBoardRepo, test.inCtx, entity.Board{
				AuthorID:    test.newBoardData.AuthorID,
				Title:       test.newBoardData.Title,
				Description: test.newBoardData.Description,
				Public:      test.newBoardData.Public,
			}, test.newBoardData.TagTitles)

			boardUsecase := New(log, mockBoardRepo, nil, sanitizer)
			newBoardID, err := boardUsecase.CreateNewBoard(test.inCtx, test.newBoardData)

			if err != nil {
				require.EqualError(t, err, test.expErr.Error())
			}
			require.Equal(t, test.expNewID, newBoardID)
		})
	}
}

func TestBoardUsecase_UpdateBoardInfo(t *testing.T) {

	tests := []struct {
		name                    string
		inCtx                   context.Context
		updatedBoardData        dto.BoardData
		GetBoardAuthorByBoardID GetBoardAuthorByBoardID
		UpdateBoard             UpdateBoard
		expErr                  error
	}{
		{
			name:  "valid data, authenticated, with access",
			inCtx: context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			updatedBoardData: dto.BoardData{
				ID:          25,
				Title:       "valid title",
				Description: "some description",
				Public:      false,
				TagTitles:   []string{"nice", "green"},
			},
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(1, nil).Times(1)
			},
			UpdateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, updatedBoardData entity.Board, tagTitles []string) {
				mockRepo.EXPECT().UpdateBoard(ctx, updatedBoardData, tagTitles).Return(nil).Times(1)
			},
			expErr: nil,
		},
		{
			name:  "valid data, authenticated, no access",
			inCtx: context.WithValue(context.Background(), auth.KeyCurrentUserID, 534),
			updatedBoardData: dto.BoardData{
				ID:          25,
				Title:       "valid title",
				Description: "some description",
				Public:      false,
				TagTitles:   []string{"nice", "green"},
			},
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(1, nil).Times(1)
			},
			UpdateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expErr: ErrNoAccess,
		},
		{
			name:  "valid data, no_auth",
			inCtx: context.Background(),
			updatedBoardData: dto.BoardData{
				ID:          25,
				Title:       "valid title",
				Description: "some description",
				Public:      false,
				TagTitles:   []string{"nice", "green"},
			},
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(1, nil).Times(1)
			},
			UpdateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expErr: ErrNoAccess,
		},
		{
			name:  "invalid board id",
			inCtx: context.Background(),
			updatedBoardData: dto.BoardData{
				ID:          122125,
				Title:       "valid title",
				Description: "some description",
				Public:      false,
				TagTitles:   []string{"nice", "green"},
			},
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(0, repository.ErrNoData).Times(1)
			},
			UpdateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expErr: ErrNoSuchBoard,
		},
		{
			name:  "invalid board title",
			inCtx: context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			updatedBoardData: dto.BoardData{
				ID:          1,
				Title:       "va!@#*^*!&@$*lid title",
				Description: "some description",
				Public:      false,
				TagTitles:   []string{"nice", "green"},
			},
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(1, nil).Times(1)
			},
			UpdateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expErr: ErrInvalidBoardTitle,
		},
		{
			name:  "invalid board tags",
			inCtx: context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			updatedBoardData: dto.BoardData{
				ID:          11,
				Title:       "valid title",
				Description: "some description",
				Public:      false,
				TagTitles:   []string{"ni@#@#%!~~ce", "green"},
			},
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(1, nil).Times(1)
			},
			UpdateBoard: func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string) {
			},
			expErr: fmt.Errorf("%v: %w", []string{"ni@#@#%!~~ce"}, ErrInvalidTagTitles),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			log, err := logger.New(logger.RFC3339FormatTime())
			if err != nil {
				stdLog.Fatal(err)
			}
			mockBoardRepo := mock_board.NewMockRepository(ctl)
			test.GetBoardAuthorByBoardID(mockBoardRepo, test.inCtx, test.updatedBoardData.ID)
			test.UpdateBoard(mockBoardRepo, test.inCtx, entity.Board{
				ID:          test.updatedBoardData.ID,
				Title:       test.updatedBoardData.Title,
				Description: test.updatedBoardData.Description,
				Public:      test.updatedBoardData.Public,
			}, test.updatedBoardData.TagTitles)

			boardUsecase := New(log, mockBoardRepo, nil, sanitizer)
			err = boardUsecase.UpdateBoardInfo(test.inCtx, test.updatedBoardData)

			if err != nil {
				require.EqualError(t, err, test.expErr.Error())
			}
		})
	}
}
