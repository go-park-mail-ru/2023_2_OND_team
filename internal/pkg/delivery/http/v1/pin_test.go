package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	pinCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestGetPins(t *testing.T) {

	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	pinCase := pinCase.New(log, nil, ramrepo.NewRamPinRepo(db))
	service := New(log, nil, nil, pinCase, nil)

	rawUrl := "https://domain.test:8080/api/v1/pin"

	badCases := []struct {
		rawURL  string
		expResp JsonErrResponse
	}{
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), maxID, minID(positive integers, the absence of these parameters is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, -2, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), maxID, minID(positive integers, the absence of these parameters is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 213123, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), maxID, minID(positive integers, the absence of these parameters is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, -1),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), maxID, minID(positive integers, the absence of these parameters is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=&lastID=%d", rawUrl, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), maxID, minID(positive integers, the absence of these parameters is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?lastID=%d", rawUrl, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), maxID, minID(positive integers, the absence of these parameters is equal to the value 0)",
				Code:    "bad_params",
			},
		},
	}

	for _, tCase := range badCases {
		t.Run(fmt.Sprintf("TestGetPins bad: %s", tCase.rawURL), func(t *testing.T) {
			req := httptest.NewRequest("GET", tCase.rawURL, nil)
			w := httptest.NewRecorder()
			service.GetPins(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			var actualResp JsonErrResponse

			json.Unmarshal(body, &actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}
}
