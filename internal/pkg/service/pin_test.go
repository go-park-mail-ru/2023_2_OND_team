package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	pinCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestGetPins(t *testing.T) {

	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB()
	defer db.Close()

	pinCase := pinCase.New(log, ramrepo.NewRamPinRepo(db))

	service := New(log, nil, nil, pinCase)

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
						map[string]interface{}{"id": 2, "picture": "/2"},
						map[string]interface{}{"id": 3, "picture": "/3"},
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
						map[string]interface{}{"id": 4, "picture": "/4"},
						map[string]interface{}{"id": 5, "picture": "/5"},
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

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			var actualResp JsonResponse

			json.Unmarshal(body, &actualResp)
			fmt.Println(tCase.expResp)
			fmt.Println(actualResp)
			require.True(t, reflect.DeepEqual(tCase.expResp, actualResp))
			// require.EqualExportedValues(t, tCase.expResp, actualResp)
		})
	}

	// badCases := []struct {
	// 	rawURL  string
	// 	expResp JsonErrResponse
	// }{
	// 	{
	// 		rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, 3),
	// 		expResp: JsonErrResponse{
	// 			Status:  "error",
	// 			Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
	// 			Code:    "bad_params",
	// 		},
	// 	},
	// 	{
	// 		rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, -2, 3),
	// 		expResp: JsonErrResponse{
	// 			Status:  "error",
	// 			Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
	// 			Code:    "bad_params",
	// 		},
	// 	},
	// 	{
	// 		rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 213123, 3),
	// 		expResp: JsonErrResponse{
	// 			Status:  "error",
	// 			Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
	// 			Code:    "bad_params",
	// 		},
	// 	},
	// 	{
	// 		rawURL: fmt.Sprintf("%s?count=%d&lastID=%d", rawUrl, 0, -1),
	// 		expResp: JsonErrResponse{
	// 			Status:  "error",
	// 			Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
	// 			Code:    "bad_params",
	// 		},
	// 	},
	// 	{
	// 		rawURL: fmt.Sprintf("%s?count=&lastID=%d", rawUrl, 3),
	// 		expResp: JsonErrResponse{
	// 			Status:  "error",
	// 			Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
	// 			Code:    "bad_params",
	// 		},
	// 	},
	// 	{
	// 		rawURL: fmt.Sprintf("%s?lastID=%d", rawUrl, 3),
	// 		expResp: JsonErrResponse{
	// 			Status:  "error",
	// 			Message: "expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)",
	// 			Code:    "bad_params",
	// 		},
	// 	},
	// }

	// for _, tCase := range badCases {
	// 	t.Run("TestGetPins bad", func(t *testing.T) {
	// 		req := httptest.NewRequest("GET", tCase.rawURL, nil)
	// 		w := httptest.NewRecorder()
	// 		service.GetPins(w, req)

	// 		resp := w.Result()
	// 		body, _ := io.ReadAll(resp.Body)
	// 		var actualResp JsonResponse

	// 		json.Unmarshal(body, actualResp)
	// 		require.Equal(t, tCase.expResp, actualResp)
	// 	})
	// }

}
