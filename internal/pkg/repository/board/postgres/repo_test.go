package board

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/board"
	dto "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board/dto"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"
)

func getTx(mockDB pgxmock.PgxPoolIface) pgx.Tx {
	mockDB.ExpectBegin()
	tx, err := mockDB.Begin(context.Background())
	if err != nil {
		log.Fatalf("can't get tx for test: %s", err.Error())
	}
	return tx
}

func TestBoardRepo_insertBoard(t *testing.T) {
	mockDB, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("test: get new mock pool - %s", err.Error())
	}
	boardRepo := NewBoardRepoPG(mockDB)

	cases := []struct {
		name    string
		tx      pgx.Tx
		board   entity.Board
		setMock func()
		wantErr bool
		expErr  error
		expID   int
	}{
		{
			name: "insert valid board",
			tx:   getTx(mockDB),
			board: entity.Board{
				AuthorID:    1,
				Title:       "title",
				Description: "desc",
				Public:      false,
			},
			setMock: func() {
				row := mockDB.NewRows([]string{"id"}).AddRow(55)
				mockDB.ExpectQuery("INSERT INTO board").WithArgs(1, "title", "desc", false).WillReturnRows(row)
			},
			expID: 55,
		},
		{
			name: "invalid authorID",
			tx:   getTx(mockDB),
			board: entity.Board{
				AuthorID:    -1,
				Title:       "title",
				Description: "desc",
				Public:      false,
			},
			setMock: func() {
				mockDB.ExpectQuery("INSERT INTO board").WithArgs(-1, "title", "desc", false).WillReturnError(pgx.ErrNoRows)
			},
			wantErr: true,
			expErr:  fmt.Errorf("scan result of insterting new board: %w", pgx.ErrNoRows),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			test.setMock()
			boardID, err := boardRepo.insertBoard(context.Background(), test.tx, test.board)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expID, boardID)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestBoardRepo_CreateBoard(t *testing.T) {
	mockDB, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("test: get new mock pool - %s", err.Error())
	}

	boardRepo := NewBoardRepoPG(mockDB)

	cases := []struct {
		name      string
		board     entity.Board
		tagTitles []string
		setMock   func()
		wantErr   bool
		expErr    error
		expID     int
	}{
		{
			name: "valid board",
			board: entity.Board{
				AuthorID:    1,
				Title:       "title",
				Description: "desc",
				Public:      false,
			},
			tagTitles: []string{"cool", "view"},
			setMock: func() {
				mockDB.ExpectBegin()

				row := mockDB.NewRows([]string{"id"}).AddRow(25)
				mockDB.ExpectQuery("INSERT INTO board").WithArgs(1, "title", "desc", false).WillReturnRows(row)

				mockDB.ExpectExec("INSERT INTO tag").WithArgs("cool", "view").WillReturnResult(pgxmock.NewResult("INSERT", 2))
				mockDB.ExpectExec("INSERT INTO board_tag").WithArgs("cool", "view").WillReturnResult(pgxmock.NewResult("INSERT", 2))

				mockDB.ExpectCommit()
			},
			expID: 25,
		},
		{
			name: "invalid author id",
			board: entity.Board{
				AuthorID:    -1231,
				Title:       "title",
				Description: "desc",
				Public:      false,
			},
			tagTitles: []string{"cool", "view"},
			setMock: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery("INSERT INTO board").WithArgs(-1231, "title", "desc", false).WillReturnError(fmt.Errorf("scan result of insterting new board: %w", pgx.ErrNoRows))
				mockDB.ExpectRollback()
			},
			wantErr: true,
			expErr:  errors.New("inserting board within transaction: scan result of insterting new board: scan result of insterting new board: no rows in result set"),
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			test.setMock()
			boardID, err := boardRepo.CreateBoard(context.Background(), test.board, test.tagTitles)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expID, boardID)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}

func TestBoardRepo_GetBoardsByUserID(t *testing.T) {
	mockDB, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("test: get new mock pool - %s", err.Error())
	}

	boardRepo := NewBoardRepoPG(mockDB)

	cases := []struct {
		name                string
		authorID            int
		isAuthor            bool
		accessableBoardsIDs []int
		setMock             func()
		wantErr             bool
		expErr              error
		expBoards           []dto.UserBoard
	}{
		{
			name:                "valid user ID, author",
			authorID:            2,
			isAuthor:            true,
			accessableBoardsIDs: []int{1, 2},
			setMock: func() {

				rows := mockDB.NewRows([]string{"board.id", "board.title", "board.description", "board.created_at", "pins_number", "pins", "tags"}).
					AddRow(4, "title", "desc", "12:12:2022", 1, []string{"/pic1"}, []string{"blue"}).
					AddRow(5, "title_", "desc", "12:11:2022", 0, []string{}, []string{})

				mockDB.ExpectQuery(
					`SELECT (.+) FROM board LEFT JOIN membership ON (.+) WHERE (.+) GROUP BY (.+) ORDER BY (.+)`,
				).WithArgs(2).WillReturnRows(rows)

			},
			expBoards: []dto.UserBoard{
				{
					BoardID:     4,
					Title:       "title",
					Description: "desc",
					CreatedAt:   "12:12:2022",
					PinsNumber:  1,
					Pins:        []string{"/pic1"},
					TagTitles:   []string{"blue"},
				},
				{
					BoardID:     5,
					Title:       "title_",
					Description: "desc",
					CreatedAt:   "12:11:2022",
					PinsNumber:  0,
					Pins:        []string{},
					TagTitles:   []string{},
				},
			},
		},
		{
			name:                "valid user ID, contributor",
			authorID:            3,
			isAuthor:            false,
			accessableBoardsIDs: []int{3, 4},
			setMock: func() {

				rows := mockDB.NewRows([]string{"board.id", "board.title", "board.description", "board.created_at", "pins_number", "pins", "tags"}).
					AddRow(4, "title", "desc", "12:12:2022", 1, []string{"/pic1"}, []string{"sun"})

				mockDB.ExpectQuery(
					`SELECT (.+) FROM board LEFT JOIN membership ON (.+) WHERE (.+) GROUP BY (.+) ORDER BY (.+)`,
				).WithArgs(3, true, 3, 4).WillReturnRows(rows)

			},
			expBoards: []dto.UserBoard{
				{
					BoardID:     4,
					Title:       "title",
					Description: "desc",
					CreatedAt:   "12:12:2022",
					PinsNumber:  1,
					Pins:        []string{"/pic1"},
					TagTitles:   []string{"sun"},
				},
			},
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			test.setMock()
			boards, err := boardRepo.GetBoardsByUserID(context.Background(), test.authorID, test.isAuthor, test.accessableBoardsIDs)

			if test.wantErr {
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expBoards, boards)
			}

			require.NoError(t, mockDB.ExpectationsWereMet())
		})
	}
}
