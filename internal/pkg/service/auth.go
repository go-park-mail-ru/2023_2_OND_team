package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	repo "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/repository/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

// Login godoc
//
//	@Description	User login, check authentication, get user info
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	JsonResponse{body=user.User}
//	@Failure		400	{object}	JsonErrResponse
//	@Failure		404	{object}	JsonErrResponse
//	@Failure		500	{object}	JsonErrResponse
//	@Router			/api/v1/auth/login [get]
func (s *Service) CheckLogin(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked CheckLogin")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}

// Login godoc
//
//	@Description	User login, creating new session
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			username	body		string	true	"Username"
//	@Param			password	body		string	true	"Password"
//	@Success		200			{object}	JsonResponse{body=Empty}
//	@Failure		400			{object}	JsonErrResponse
//	@Failure		404			{object}	JsonErrResponse
//	@Failure		500			{object}	JsonErrResponse
//	@Header			200			{string}	session_key	"Auth cookie with new valid session id"
//	@Router			/api/v1/auth/login [post]
func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	s.log.Info("request on signup", logger.F{"method", r.Method}, logger.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	defer r.Body.Close()
	params := repo.UserCredentials{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		s.log.Info("failed to parse parameters", logger.F{"error", err.Error()})
		resBody, err := json.Marshal(map[string]any{
			"status": "error",
			"code":   "bad_params",
		})
		if err != nil {
			s.log.Error(err.Error())
		}
		w.Write(resBody)
		return
	}

	user, err := s.userCase.Authentication(r.Context(), params)
	if err != nil {
		s.log.Warn(err.Error())
		resBody, err := json.Marshal(map[string]string{
			"status": "error",
			"code":   "user_authentication",
		})
		if err != nil {
			s.log.Error(err.Error())
		}
		w.Write(resBody)
		return
	}

	session, err := s.sm.CreateNewSessionForUser(r.Context(), user.ID)
	if err != nil {
		s.log.Error(err.Error())
		resBody, err := json.Marshal(map[string]string{
			"status": "error",
			"code":   "create_session",
		})
		if err != nil {
			s.log.Error(err.Error())
		}
		w.Write(resBody)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_key",
		Value:    session.Key,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	resBody, err := json.Marshal(map[string]any{
		"status":  "ok",
		"comment": "set cookie",
		"body":    map[string]any{"user": user},
	})
	if err != nil {
		s.log.Error(err.Error())
	}
	w.Write(resBody)
}

// SignUp godoc
//
//	@Description	User registration
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			username	body		string	true	"Username"
//	@Param			email		body		string	true	"Email"
//	@Param			password	body		string	true	"Password"
//	@Success		200			{object}	JsonResponse{body=Empty}
//	@Failure		400			{object}	JsonErrResponse
//	@Failure		404			{object}	JsonErrResponse
//	@Failure		500			{object}	JsonErrResponse
//	@Router			/api/v1/auth/signup [post]
func (s *Service) Signup(w http.ResponseWriter, r *http.Request) {
	s.log.Info("request on signup", logger.F{"method", r.Method}, logger.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	defer r.Body.Close()
	user := &user.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		s.log.Info("failed to parse parameters", logger.F{"error", err.Error()})
		resBody, err := json.Marshal(map[string]string{
			"status": "error",
			"code":   "bad_params",
		})
		if err != nil {
			s.log.Error(err.Error())
		}
		w.Write(resBody)
		return
	}

	err = s.userCase.Register(r.Context(), user)
	if err != nil {
		s.log.Warn(err.Error())
		resBody, err := json.Marshal(map[string]string{
			"status": "error",
			"code":   "register",
		})
		if err != nil {
			s.log.Error(err.Error())
		}
		w.Write(resBody)
		return
	}

	resBody, err := json.Marshal(map[string]string{
		"status":  "ok",
		"comment": "the user is registered",
	})
	if err != nil {
		s.log.Error(err.Error())
	}
	w.Write(resBody)
}

// Logout godoc
//
//	@Description	User logout, session deletion
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	JsonResponse{body=Empty}
//	@Failure		400	{object}	JsonErrResponse
//	@Failure		404	{object}	JsonErrResponse
//	@Failure		500	{object}	JsonErrResponse
//	@Header			200	{string}	Session-id	"Auth cookie with expired session id"
//	@Router			/api/v1/auth/logout [delete]
func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	s.log.Info("request on signup", logger.F{"method", r.Method}, logger.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	cookie, err := r.Cookie("session_key")
	if err != nil {
		s.log.Info("no cookie", logger.F{"error", err.Error()})
		resBody, err := json.Marshal(map[string]string{
			"status": "error",
			"code":   "no_cookie",
		})
		if err != nil {
			s.log.Error(err.Error())
		}
		w.Write(resBody)
		return
	}

	err = s.sm.DeleteUserSession(r.Context(), cookie.Value)
	if err != nil {
		s.log.Error(err.Error())
	}

	cookie.Expires.AddDate(0, -1, 0)
	http.SetCookie(w, cookie)
	resBody, err := json.Marshal(map[string]string{
		"status":  "ok",
		"comment": "cookie delete",
	})
	if err != nil {
		s.log.Error(err.Error())
	}
	w.Write(resBody)
}
