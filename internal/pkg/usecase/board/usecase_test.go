package board

import (
	"context"
	"testing"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	uEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository"
	mock_board "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board/mock"
	mock_user "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user/mock"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/require"
)

type (
	CreateBoard              func(mockRepo *mock_board.MockRepository, ctx context.Context, newBoardData entity.Board, tagTitles []string)
	UpdateBoard              CreateBoard
	GetBoardAuthorByBoardID  func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int)
	GetContributorBoardsIDs  func(mockRepo *mock_board.MockRepository, ctx context.Context, contributorID int)
	GetBoardsByUserID        func(mockRepo *mock_board.MockRepository, ctx context.Context, userID int, isAuthor bool, accessableBoardsIDs []int)
	GetContributorsByBoardID func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int)
	GetBoardByID             func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int, hasAccess bool)
	GetUserIdByUsername      func(mockRepo *mock_user.MockRepository, ctx context.Context, username string)
	DeleteBoardByID          func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int)
)

var (
	sanitizer = bluemonday.UGCPolicy()
)

func TestBoardUsecase_CreateNewBoard(t *testing.T) {

	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		t.Fatalf("test: log init - %s", err.Error())
	}

	tests := []struct {
		name         string
		inCtx        context.Context
		newBoardData dto.BoardData
		CreateBoard  CreateBoard
		expNewID     int
		wantErr      bool
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockBoardRepo := mock_board.NewMockRepository(ctl)
			test.CreateBoard(mockBoardRepo, test.inCtx, entity.Board{
				AuthorID:    test.newBoardData.AuthorID,
				Title:       test.newBoardData.Title,
				Description: test.newBoardData.Description,
				Public:      test.newBoardData.Public,
			}, test.newBoardData.TagTitles)

			boardUsecase := New(log, mockBoardRepo, nil, sanitizer)
			newBoardID, err := boardUsecase.CreateNewBoard(test.inCtx, test.newBoardData)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expNewID, newBoardID)
			}
		})
	}
}

func TestBoardUsecase_UpdateBoardInfo(t *testing.T) {

	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		t.Fatalf("test: log init - %s", err.Error())
	}

	tests := []struct {
		name                    string
		inCtx                   context.Context
		updatedBoardData        dto.BoardData
		GetBoardAuthorByBoardID GetBoardAuthorByBoardID
		UpdateBoard             UpdateBoard
		wantErr                 bool
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
			wantErr: true,
			expErr:  ErrNoAccess,
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
			wantErr: true,
			expErr:  ErrNoAccess,
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
			wantErr: true,
			expErr:  ErrNoSuchBoard,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockBoardRepo := mock_board.NewMockRepository(ctl)
			test.GetBoardAuthorByBoardID(mockBoardRepo, test.inCtx, test.updatedBoardData.ID)
			test.UpdateBoard(mockBoardRepo, test.inCtx, entity.Board{
				ID:          test.updatedBoardData.ID,
				Title:       test.updatedBoardData.Title,
				Description: test.updatedBoardData.Description,
				Public:      test.updatedBoardData.Public,
			}, test.updatedBoardData.TagTitles)

			boardUsecase := New(log, mockBoardRepo, nil, sanitizer)
			err := boardUsecase.UpdateBoardInfo(test.inCtx, test.updatedBoardData)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBoardUsecase_GetBoardsByUsername(t *testing.T) {

	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		t.Fatalf("test: log init - %s", err.Error())
	}

	tests := []struct {
		name                    string
		inCtx                   context.Context
		username                string
		GetUserIdByUsername     GetUserIdByUsername
		GetContributorBoardsIDs GetContributorBoardsIDs
		GetBoardsByUserID       GetBoardsByUserID
		expBoards               []dto.UserBoard
		wantErr                 bool
		expErr                  error
	}{
		{
			name:     "exisitng user with valid username",
			inCtx:    context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			username: "validGuy",
			GetUserIdByUsername: func(mockRepo *mock_user.MockRepository, ctx context.Context, username string) {
				mockRepo.EXPECT().GetUserIdByUsername(ctx, username).Return(3, nil).Times(1)
			},
			GetContributorBoardsIDs: func(mockRepo *mock_board.MockRepository, ctx context.Context, contributorID int) {
				mockRepo.EXPECT().GetContributorBoardsIDs(ctx, contributorID).Return([]int{1, 2, 3}, nil).Times(1)
			},
			GetBoardsByUserID: func(mockRepo *mock_board.MockRepository, ctx context.Context, userID int, isAuthor bool, accessableBoardsIDs []int) {
				mockRepo.EXPECT().GetBoardsByUserID(ctx, userID, isAuthor, accessableBoardsIDs).Return(
					[]dto.UserBoard{
						{
							BoardID:    23,
							Title:      "title",
							CreatedAt:  "25:10:2022",
							PinsNumber: 2,
							Pins:       []string{"/pic1", "/pic2"},
						},
						{
							BoardID:    21,
							Title:      "title21",
							CreatedAt:  "25:10:2012",
							PinsNumber: 0,
							Pins:       []string{},
						},
					}, nil).Times(1)
			},
			expBoards: []dto.UserBoard{
				{
					BoardID:    23,
					Title:      "title",
					CreatedAt:  "25:10:2022",
					PinsNumber: 2,
					Pins:       []string{"/pic1", "/pic2"},
				},
				{
					BoardID:    21,
					Title:      "title21",
					CreatedAt:  "25:10:2012",
					PinsNumber: 0,
					Pins:       []string{},
				},
			},
		},
		{
			name:     "non-exisitng user with valid username",
			inCtx:    context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			username: "validGuy",
			GetUserIdByUsername: func(mockRepo *mock_user.MockRepository, ctx context.Context, username string) {
				mockRepo.EXPECT().GetUserIdByUsername(ctx, username).Return(0, repository.ErrNoData).Times(1)
			},
			GetContributorBoardsIDs: func(mockRepo *mock_board.MockRepository, ctx context.Context, contributorID int) {
			},
			GetBoardsByUserID: func(mockRepo *mock_board.MockRepository, ctx context.Context, userID int, isAuthor bool, accessableBoardsIDs []int) {
			},
			expBoards: nil,
			wantErr:   true,
			expErr:    ErrInvalidUsername,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockBoardRepo := mock_board.NewMockRepository(ctl)
			mockUserRepo := mock_user.NewMockRepository(ctl)

			test.GetUserIdByUsername(mockUserRepo, test.inCtx, test.username)
			test.GetContributorBoardsIDs(mockBoardRepo, test.inCtx, 1)
			test.GetBoardsByUserID(mockBoardRepo, test.inCtx, 3, *new(bool), []int{1, 2, 3})

			boardUsecase := New(log, mockBoardRepo, mockUserRepo, sanitizer)
			userBoards, err := boardUsecase.GetBoardsByUsername(test.inCtx, test.username)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expBoards, userBoards)
			}
		})
	}
}

func TestBoardUsecase_GetCertainBoard(t *testing.T) {

	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		t.Fatalf("test: log init - %s", err.Error())
	}

	tests := []struct {
		name                     string
		inCtx                    context.Context
		boardID                  int
		GetBoardAuthorByBoardID  GetBoardAuthorByBoardID
		GetContributorsByBoardID GetContributorsByBoardID
		GetBoardByID             GetBoardByID
		hasAccess                bool
		expBoard                 dto.UserBoard
		wantErr                  bool
		expErr                   error
	}{
		{
			name:    "private board, valid board id, request from author",
			inCtx:   context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			boardID: 22,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(1, nil).Times(1)
			},
			GetContributorsByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetContributorsByBoardID(ctx, boardID).Return([]uEntity.User{}, nil).Times(1)
			},
			GetBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int, hasAccess bool) {
				mockRepo.EXPECT().GetBoardByID(ctx, boardID, hasAccess).Return(dto.UserBoard{
					BoardID:     boardID,
					Title:       "title",
					Description: "description",
					CreatedAt:   "10:10:2020",
					PinsNumber:  1,
					Pins:        []string{"/pic1"},
					TagTitles:   []string{"good", "bad"},
				}, nil).Times(1)
			},
			hasAccess: true,
			expBoard: dto.UserBoard{
				BoardID:     22,
				Title:       "title",
				Description: "description",
				CreatedAt:   "10:10:2020",
				PinsNumber:  1,
				Pins:        []string{"/pic1"},
				TagTitles:   []string{"good", "bad"},
			},
		},
		{
			name:    "private board, valid board id, request from contributor",
			inCtx:   context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			boardID: 22,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(2, nil).Times(1)
			},
			GetContributorsByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetContributorsByBoardID(ctx, boardID).Return([]uEntity.User{{ID: 1}}, nil).Times(1)
			},
			GetBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int, hasAccess bool) {
				mockRepo.EXPECT().GetBoardByID(ctx, boardID, hasAccess).Return(dto.UserBoard{
					BoardID:     boardID,
					Title:       "title",
					Description: "description",
					CreatedAt:   "10:10:2020",
					PinsNumber:  1,
					Pins:        []string{"/pic1"},
					TagTitles:   []string{"good", "bad"},
				}, nil).Times(1)
			},
			hasAccess: true,
			expBoard: dto.UserBoard{
				BoardID:     22,
				Title:       "title",
				Description: "description",
				CreatedAt:   "10:10:2020",
				PinsNumber:  1,
				Pins:        []string{"/pic1"},
				TagTitles:   []string{"good", "bad"},
			},
		},
		{
			name:    "private board, valid board id, request from not author, not contributor",
			inCtx:   context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			boardID: 22,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(2, nil).Times(1)
			},
			GetContributorsByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetContributorsByBoardID(ctx, boardID).Return([]uEntity.User{{ID: 123}}, nil).Times(1)
			},
			GetBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int, hasAccess bool) {
				mockRepo.EXPECT().GetBoardByID(ctx, boardID, hasAccess).Return(dto.UserBoard{}, repository.ErrNoData).Times(1)
			},
			hasAccess: false,
			expBoard:  dto.UserBoard{},
			wantErr:   true,
			expErr:    ErrNoSuchBoard,
		},
		{
			name:    "private board, valid board id, request from unauthorized",
			inCtx:   context.Background(),
			boardID: 22,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(2, nil).Times(1)
			},
			GetContributorsByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetContributorsByBoardID(ctx, boardID).Return([]uEntity.User{{ID: 123}}, nil).Times(1)
			},
			GetBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int, hasAccess bool) {
				mockRepo.EXPECT().GetBoardByID(ctx, boardID, hasAccess).Return(dto.UserBoard{}, repository.ErrNoData).Times(1)
			},
			hasAccess: false,
			expBoard:  dto.UserBoard{},
			wantErr:   true,
			expErr:    ErrNoSuchBoard,
		},
		{
			name:    "public board, valid board id, request from unauthorized",
			inCtx:   context.Background(),
			boardID: 22,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(2, nil).Times(1)
			},
			GetContributorsByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetContributorsByBoardID(ctx, boardID).Return([]uEntity.User{{ID: 123}}, nil).Times(1)
			},
			GetBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int, hasAccess bool) {
				mockRepo.EXPECT().GetBoardByID(ctx, boardID, hasAccess).Return(dto.UserBoard{
					BoardID:     boardID,
					Title:       "title",
					Description: "description",
					CreatedAt:   "10:10:2020",
					PinsNumber:  1,
					Pins:        []string{"/pic1"},
					TagTitles:   []string{"good", "bad"},
				}, nil).Times(1)
			},
			hasAccess: false,
			expBoard: dto.UserBoard{
				BoardID:     22,
				Title:       "title",
				Description: "description",
				CreatedAt:   "10:10:2020",
				PinsNumber:  1,
				Pins:        []string{"/pic1"},
				TagTitles:   []string{"good", "bad"},
			},
		},
		{
			name:    "invalid board id",
			inCtx:   context.Background(),
			boardID: 1222,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(0, repository.ErrNoData).Times(1)
			},
			GetContributorsByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
			},
			GetBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int, hasAccess bool) {
			},
			expBoard: dto.UserBoard{},
			wantErr:  true,
			expErr:   ErrNoSuchBoard,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockBoardRepo := mock_board.NewMockRepository(ctl)
			test.GetBoardAuthorByBoardID(mockBoardRepo, test.inCtx, test.boardID)
			test.GetContributorsByBoardID(mockBoardRepo, test.inCtx, test.boardID)
			test.GetBoardByID(mockBoardRepo, test.inCtx, test.boardID, test.hasAccess)

			boardUsecase := New(log, mockBoardRepo, nil, sanitizer)
			board, err := boardUsecase.GetCertainBoard(test.inCtx, test.boardID)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expBoard, board)
			}
		})
	}
}

func TestBoardUsecase_DeleteCertainBoard(t *testing.T) {

	log, err := logger.New(logger.RFC3339FormatTime())
	if err != nil {
		t.Fatalf("test: log init - %s", err.Error())
	}

	tests := []struct {
		name                    string
		inCtx                   context.Context
		boardID                 int
		GetBoardAuthorByBoardID GetBoardAuthorByBoardID
		DeleteBoardByID         DeleteBoardByID
		wantErr                 bool
		expErr                  error
	}{
		{
			name:    "valid board id, deletion by author",
			inCtx:   context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			boardID: 23,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(1, nil).Times(1)
			},
			DeleteBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().DeleteBoardByID(ctx, boardID).Return(nil).Times(1)
			},
		},
		{
			name:    "valid board id, deletion by another user",
			inCtx:   context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			boardID: 23,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(2, nil).Times(1)
			},
			DeleteBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
			},
			wantErr: true,
			expErr:  ErrNoAccess,
		},
		{
			name:    "valid board id, deletion by unauthorized user",
			inCtx:   context.Background(),
			boardID: 23,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(2, nil).Times(1)
			},
			DeleteBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
			},
			wantErr: true,
			expErr:  ErrNoAccess,
		},
		{
			name:    "invalid board id",
			inCtx:   context.WithValue(context.Background(), auth.KeyCurrentUserID, 1),
			boardID: 1221323,
			GetBoardAuthorByBoardID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
				mockRepo.EXPECT().GetBoardAuthorByBoardID(ctx, boardID).Return(0, repository.ErrNoData).Times(1)
			},
			DeleteBoardByID: func(mockRepo *mock_board.MockRepository, ctx context.Context, boardID int) {
			},
			wantErr: true,
			expErr:  ErrNoSuchBoard,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockBoardRepo := mock_board.NewMockRepository(ctl)
			test.GetBoardAuthorByBoardID(mockBoardRepo, test.inCtx, test.boardID)
			test.DeleteBoardByID(mockBoardRepo, test.inCtx, test.boardID)

			boardUsecase := New(log, mockBoardRepo, nil, sanitizer)
			err = boardUsecase.DeleteCertainBoard(test.inCtx, test.boardID)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
