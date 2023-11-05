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
	goodCases := []struct {
		rawURL  string
		expResp JsonResponse
	}{
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 1, 2),
			expResp: JsonResponse{
				Status:  "ok",
				Message: "pins received are sorted by id",
				Body: map[string]interface{}{
					"lastID": 3,
					"pins": []interface{}{
						map[string]interface{}{"id": 3},
					},
				},
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 2, 3),
			expResp: JsonResponse{
				Status:  "ok",
				Message: "pins received are sorted by id",
				Body: map[string]interface{}{
					"lastID": 5,
					"pins": []interface{}{
						map[string]interface{}{"id": 4},
						map[string]interface{}{"id": 5},
					},
				},
			},
		},
	}

	for _, tCase := range goodCases {
		t.Run(fmt.Sprintf("TestGetPins good: %s", tCase.rawURL), func(t *testing.T) {
			req := httptest.NewRequest("GET", tCase.rawURL, nil)
			w := httptest.NewRecorder()
			service.GetPins(w, req)

			var actualResp JsonResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)

			require.Equal(t, tCase.expResp.Status, actualResp.Status)
			require.Equal(t, tCase.expResp.Message, actualResp.Message)
			expLastID := tCase.expResp.Body.(map[string]interface{})["lastID"].(int)
			actualLastID := actualResp.Body.(map[string]interface{})["lastID"].(float64)

			expIDs, actualIDs := make([]int, 0), make([]int, 0)
			for _, pin := range tCase.expResp.Body.(map[string]interface{})["pins"].([]interface{}) {
				expIDs = append(expIDs, pin.(map[string]interface{})["id"].(int))
			}
			for _, pin := range actualResp.Body.(map[string]interface{})["pins"].([]interface{}) {
				actualIDs = append(actualIDs, int(pin.(map[string]interface{})["id"].(float64)))
			}

			require.Equal(t, expLastID, int(actualLastID))
			require.Equal(t, expIDs, actualIDs)
		})
	}

	badCases := []struct {
		rawURL  string
		expResp JsonErrResponse
	}{
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, -2, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 213123, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, -1),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?count=&lastID=%d", rawUrl, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
				Code:    "bad_params",
			},
		},
		{
			rawURL: fmt.Sprintf("%s?lastID=%d", rawUrl, 3),
			expResp: JsonErrResponse{
				Status:  "error",
				Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
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
