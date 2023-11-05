package v1

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/ramrepo"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/stretchr/testify/require"
)

func checkAuthCookie(cookies []*http.Cookie) bool {
	if cookies == nil {
		return false
	}
	for _, cookie := range cookies {
		if cookie.Name == "session_key" {
			return true
		}
	}
	return false
}

func TestCheckLogin(t *testing.T) {
	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	sm := session.New(log, ramrepo.NewRamSessionRepo(db))
	userCase := user.New(log, ramrepo.NewRamUserRepo(db))
	service := New(log, sm, userCase, nil, nil, nil)

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
				Body:    map[string]interface{}{"username": "dogsLover", "avatar": "https://pinspire.online:8081/upload/avatars/default-avatar.png"},
			},
		},
	}

	for _, tCase := range goodCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(tCase.cookie)
			w := httptest.NewRecorder()

			service.CheckLogin(w, req)

			var actualResp JsonResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
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

			var actualResp JsonErrResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}

}

func TestLogin(t *testing.T) {
	url := "https://domain.test:8080/api/v1/login"
	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	sm := session.New(log, ramrepo.NewRamSessionRepo(db))
	userCase := user.New(log, ramrepo.NewRamUserRepo(db))
	service := New(log, sm, userCase, nil, nil, nil)

	goodCases := []struct {
		name    string
		rawBody string
		expResp JsonResponse
	}{
		{
			"providing correct and valid user credentials",
			`{"username":"dogsLover", "password":"big_string"}`,
			JsonResponse{
				Status:  "ok",
				Message: "a new session has been created for the user",
				Body:    nil,
			},
		},
	}

	for _, tCase := range goodCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(strings.NewReader(tCase.rawBody)))
			w := httptest.NewRecorder()

			service.Login(w, req)

			var actualResp JsonResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
			require.Equal(t, tCase.expResp, actualResp)
			require.True(t, checkAuthCookie(w.Result().Cookies()))
		})
	}

	badCases := []struct {
		name    string
		rawBody string
		expResp JsonErrResponse
	}{
		{
			"providing invalid credentials - broken body",
			"{'username': 'dogsLover', 'password': 'big_string'",
			JsonErrResponse{
				Status:  "error",
				Message: "the correct username and password are expected to be received in JSON format",
				Code:    "parse_body",
			},
		},
		{
			"providing invalid credentials - no username",
			`{"password":"big_string"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "invalid user credentials",
				Code:    "invalid_credentials",
			},
		},
		{
			"providing invalid credentials - no password",
			`{"username":"dogsLover"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "invalid user credentials",
				Code:    "invalid_credentials",
			},
		},
		{
			"providing invalid credentials - short username",
			`{"username":"do", "password":"big_string"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "invalid user credentials",
				Code:    "invalid_credentials",
			},
		},
		{
			"providing invalid credentials - long username",
			`{"username":"dojsbrjfbdrjhbhjldrbgbdrhjgbdjrbgjdhbgjhdbrghbdhj,gbdhjrbgjhdbvkvghkevfghjdvrfhvdhrvbjdfgdrgdr","password":"big_string"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "invalid user credentials",
				Code:    "invalid_credentials",
			},
		},
		{
			"providing invalid credentials - short password",
			`{"username":"dogsLover","password":"bi"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "invalid user credentials",
				Code:    "invalid_credentials",
			},
		},
		{
			"providing invalid credentials - long password",
			`{"username":"dogsLover","password":"biyugsgrusgubskhvfhkdgvfgvdvrjgbsjhgjkshzkljfskfwjkhkfjisuidgoquakflsjuzeofiow3i"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "invalid user credentials",
				Code:    "invalid_credentials",
			},
		},
		{
			"providing incorrect credentials - no user with such credentials",
			`{"username":"dogsLover", "password":"doesn't_exist"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "incorrect user credentials",
				Code:    "bad_credentials",
			},
		},
	}

	for _, tCase := range badCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(strings.NewReader(tCase.rawBody)))
			w := httptest.NewRecorder()

			service.Login(w, req)

			var actualResp JsonErrResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
			require.Equal(t, tCase.expResp, actualResp)
			require.False(t, checkAuthCookie(w.Result().Cookies()))
		})
	}
}

func TestSignUp(t *testing.T) {
	url := "https://domain.test:8080/api/v1/signup"
	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	sm := session.New(log, ramrepo.NewRamSessionRepo(db))
	userCase := user.New(log, ramrepo.NewRamUserRepo(db))
	service := New(log, sm, userCase, nil, nil, nil)

	goodCases := []struct {
		name    string
		rawBody string
		expResp JsonResponse
	}{
		{
			"providing correct and valid data for signup",
			`{"username":"newbie", "password":"getHigh123", "email":"world@uandex.ru"}`,
			JsonResponse{
				Status:  "ok",
				Message: "the user has been successfully registered",
				Body:    nil,
			},
		},
	}

	for _, tCase := range goodCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(strings.NewReader(tCase.rawBody)))
			w := httptest.NewRecorder()

			service.Signup(w, req)

			var actualResp JsonResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}

	badCases := []struct {
		name    string
		rawBody string
		expResp JsonErrResponse
	}{
		{
			"user with such data already exists",
			`{"username":"dogsLover", "password":"big_string", "email":"dogslove@gmail.com"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "there is already an account with this username or password",
				Code:    "uniq_fields",
			},
		},
		{
			"invalid data - broken body",
			`{"username":"dogsLover", "password":"big_string", "email":"dogslove@gmail.com"`,
			JsonErrResponse{
				Status:  "error",
				Message: "the correct username, email and password are expected to be received in JSON format",
				Code:    "parse_body",
			},
		},
		{
			"invalid data - no username",
			`{"password":"big_string", "email":"dogslove@gmail.com"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "username",
				Code:    "invalid_params",
			},
		},
		{
			"invalid data - no username, password",
			`{"email":"dogslove@gmail.com"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "password,username",
				Code:    "invalid_params",
			},
		},
		{
			"invalid data - short username",
			`{"username":"sh", "password":"big_string", "email":"dogslove@gmail.com"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "username",
				Code:    "invalid_params",
			},
		},
		{
			"invalid data - incorrect email",
			`{"username":"sh", "password":"big_string", "email":"dog"}`,
			JsonErrResponse{
				Status:  "error",
				Message: "email,username",
				Code:    "invalid_params",
			},
		},
	}

	for _, tCase := range badCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, url, io.NopCloser(strings.NewReader(tCase.rawBody)))
			w := httptest.NewRecorder()

			service.Signup(w, req)

			var actualResp JsonErrResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}
}

func TestLogout(t *testing.T) {
	url := "https://domain.test:8080/api/v1/logout"
	log, _ := logger.New(logger.RFC3339FormatTime())
	defer log.Sync()

	db, _ := ramrepo.OpenDB(strconv.FormatInt(int64(rand.Int()), 10))
	defer db.Close()

	sm := session.New(log, ramrepo.NewRamSessionRepo(db))
	userCase := user.New(log, ramrepo.NewRamUserRepo(db))
	service := New(log, sm, userCase, nil, nil, nil)

	goodCases := []struct {
		name    string
		cookie  *http.Cookie
		expResp JsonResponse
	}{
		{
			"user is logged in - providing valid cookie",
			&http.Cookie{
				Name:  "session_key",
				Value: "461afabf38b3147c",
			},
			JsonResponse{
				Status:  "ok",
				Message: "the user has successfully logged out",
				Body:    nil,
			},
		},
	}

	for _, tCase := range goodCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(tCase.cookie)
			w := httptest.NewRecorder()

			service.Logout(w, req)

			var actualResp JsonResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}

	badCases := []struct {
		name    string
		cookie  *http.Cookie
		expResp JsonErrResponse
	}{
		{
			"user isn't logged in - providing invalid cookie",
			&http.Cookie{
				Name:  "not_auth_cookie",
				Value: "blablalba",
			},
			JsonErrResponse{
				Status:  "error",
				Message: "to log out, you must first log in",
				Code:    "no_auth",
			},
		},
	}

	for _, tCase := range badCases {
		t.Run(tCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(tCase.cookie)
			w := httptest.NewRecorder()

			service.Logout(w, req)

			var actualResp JsonErrResponse
			json.NewDecoder(w.Result().Body).Decode(&actualResp)
			require.Equal(t, tCase.expResp, actualResp)
		})
	}

}
