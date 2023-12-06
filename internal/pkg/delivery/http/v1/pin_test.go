package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pashagolub/pgxmock/v2"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	boardRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/board/postgres"
	pinRepo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	boardCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
	pinCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorGetPins(t *testing.T) {

	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	pinCase := pinCase.New(log, nil, ramrepo.NewRamPinRepo(db))
	service := New(log, UsecaseHub{
		PinCase: pinCase,
	})

	rawUrl := "https://domain.test:8080/api/v1/pin"

	badCases := []struct {
		rawURL  string
		expResp JsonErrResponse
	}{
		{
			rawURL: fmt.Sprintf("%s?count=%d&minID=%s", rawUrl, 0, url.PathEscape("not integer")),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "bad url params",
				Code:    "parse_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?lastID=%d", rawUrl, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "bad url params",
				Code:    "parse_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d&deleted=bad", rawUrl, 213123, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "bad url params",
				Code:    "parse_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&minID=%d&maxID=12&liked=bad", rawUrl, 0, -1),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "bad url params",
				Code:    "parse_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=&lastID=%d", rawUrl, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "bad url params",
				Code:    "parse_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?lastID=%d", rawUrl, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "bad url params",
				Code:    "parse_params",
			},
		},
	}

	for _, tCase := range badCases {
		t.Run(fmt.Sprintf("TestGetPins bad: %s", tCase.rawURL), func(t *testing.T) {
			req := httptest.NewRequest("GET", tCase.rawURL, nil)
			w := httptest.NewRecorder()
			service.FeedPins(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var actualResp JsonErrResponse

			json.Unmarshal(body, &actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}
}

func TestViewFeedPin(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal("new mock pool")
	}
	defer pool.Close()

	log, err := logger.New()
	if err != nil {
		t.Fatal("new logger")
	}
	defer log.Sync()

	repoBoard := boardRepo.NewBoardRepoPG(pool)
	usecaseBoard := boardCase.New(log, repoBoard, nil, nil)

	repoPin := pinRepo.NewPinRepoPG(pool)
	usecasePin := pinCase.New(log, nil, repoPin)

	service := New(log, UsecaseHub{
		BoardCase: usecaseBoard,
		PinCase:   usecasePin,
	})

	urlFeed := "https://pinspire.online:8080/api/v1/feed/pin?count=20&boardID=12&userID=23&liked=true"
	r, err := http.NewRequest(http.MethodGet, urlFeed, nil)
	if err != nil {
		t.Fatal("new request for test")
	}
	w := httptest.NewRecorder()

	pool.ExpectQuery("SELECT public FROM board").
		WithArgs(12).
		WillReturnRows(
			pgxmock.NewRows([]string{"public"}).
				AddRow(true),
		)

	pool.ExpectQuery("SELECT pin.id, pin.picture FROM pin").
		WithArgs(true, 23, 12, 0, 0).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "picture"}).
				AddRow(1, "https://pin"),
		)

	service.FeedPins(w, r)
	res := w.Result()
	defer res.Body.Close()

	resBode, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, resBode, []byte(`{"status":"ok","message":"ok","body":{"minID":1,"maxID":1,"pins":[{"id":1,"picture":"https://pin","title":null,"description":null,"public":true,"count_likes":0}]}}`))
}

func TestLikePin(t *testing.T) {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal("new mock pool")
	}
	defer pool.Close()

	log, err := logger.New()
	if err != nil {
		t.Fatal("new logger")
	}
	defer log.Sync()

	repoBoard := boardRepo.NewBoardRepoPG(pool)
	usecaseBoard := boardCase.New(log, repoBoard, nil, nil)

	repoPin := pinRepo.NewPinRepoPG(pool)
	usecasePin := pinCase.New(log, nil, repoPin)

	service := New(log, UsecaseHub{
		BoardCase: usecaseBoard,
		PinCase:   usecasePin,
	})

	urlSetLike := "https://pinspire.online:8080/api/v1/pin/like/set/45"
	r, err := http.NewRequest(http.MethodPost, urlSetLike, nil)
	if err != nil {
		t.Fatal("new request for test")
	}

	r = r.WithContext(context.WithValue(r.Context(), auth.KeyCurrentUserID, 25))
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"pinID"}, Values: []string{"45"}},
	}))
	w := httptest.NewRecorder()

	pool.ExpectQuery("SELECT author, title, description, picture, public, deleted_at FROM pin").
		WithArgs(45).
		WillReturnRows(
			pgxmock.NewRows([]string{"author", "title", "description", "picture", "public", "deleted_at"}).
				AddRow(25, "pin", "no", "https://picture", true, nil),
		)

	pool.ExpectQuery("INSERT INTO like_pin").
		WithArgs(45, 25).
		WillReturnRows(
			pgxmock.NewRows([]string{"count"}).
				AddRow(90),
		)

	service.SetLikePin(w, r)
	res := w.Result()
	defer res.Body.Close()

	resBode, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resBode, []byte(`{"status":"ok","message":"ok","body":{"count_like":91}}`))
}
