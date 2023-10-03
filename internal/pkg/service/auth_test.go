package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestCheckLogin(t *testing.T) {
	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB()
	defer db.Close()

	sm := session.New(log, ramrepo.NewRamSessionRepo(db))
	userCase := user.New(log, ramrepo.NewRamUserRepo(db))
	service := New(log, sm, userCase, nil)

	url := "https://domain.test:8080/api/v1/login"
	goodCases := []struct {
		name    string
		cookie  *http.Cookie
		expResp JsonResponse
	}{
		{
			"sending valid session_key",
			&http.Cookie{
				Name:  "session_key",
				Value: "461afabf38b3147c",
			},
			JsonResponse{
				Status:  "ok",
				Message: "user found",
				Body:    map[string]interface{}{"username": "dogsLover", "avatar": "https://cdn-icons-png.flaticon.com/512/149/149071.png"},
			},
		},
	}

	for _, tCase := range goodCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(tCase.cookie)
			w := httptest.NewRecorder()

			service.CheckLogin(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			var actualResp JsonResponse
			json.Unmarshal(body, &actualResp)
			actualResp.Body = actualResp.Body.(map[string]interface{})
			require.Equal(t, tCase.expResp, actualResp)
		})
	}

	badCases := []struct {
		name    string
		cookie  *http.Cookie
		expResp JsonErrResponse
	}{
		{
			"sending empty cookie",
			&http.Cookie{
				Name:  "",
				Value: "",
			},
			JsonErrResponse{
				Status:  "error",
				Message: "the user is not logged in",
				Code:    "no_auth",
			},
		},
		{
			"sending invalid cookie",
			&http.Cookie{
				Name:  "session_key",
				Value: "doesn't exist",
			},
			JsonErrResponse{
				Status:  "error",
				Message: "no user session found",
				Code:    "no_auth",
			},
		},
		{
			"sending cookie with invald user",
			&http.Cookie{
				Name:  "session_key",
				Value: "f4280a941b664d02",
			},
			JsonErrResponse{
				Status:  "error",
				Message: "no user was found for this session",
				Code:    "no_auth",
			},
		},
	}

	for _, tCase := range badCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(tCase.cookie)
			w := httptest.NewRecorder()

			service.CheckLogin(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			var actualResp JsonErrResponse
			json.Unmarshal(body, &actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}

}
