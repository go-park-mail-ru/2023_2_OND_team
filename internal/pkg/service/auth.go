package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/usecases/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
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
	s.log.Info("request on check login", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	cookie, err := r.Cookie("session_key")
	if err != nil {
		s.log.Info("no cookie", log.F{"error", err.Error()})
		err = responseError(w, "no_auth", "the user is not logged in")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	userID, err := s.sm.GetUserIDBySessionKey(r.Context(), cookie.Value)
	if err != nil {
		err = responseError(w, "no_auth", "no user session found")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	username, err := s.userCase.FindOutUserName(r.Context(), userID)
	if err != nil {
		s.log.Error(err.Error())
		err = responseError(w, "no_auth", "no user was found for this session")
	} else {
		err = responseOk(w, "user found", map[string]string{"username": username})
	}
	if err != nil {
		s.log.Error(err.Error())
	}
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
	s.log.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	defer r.Body.Close()
	params := usecase.NewCredentials()
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		s.log.Info("failed to parse parameters", log.F{"error", err.Error()})
		err = responseError(w, "parse_body", "the correct username and password are expected to be received in JSON format")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	if !isValidPassword(params.Password) || !isValidUsername(params.Username) {
		s.log.Info(err.Error())
		err = responseError(w, "bad_credentials", "invalid user credentials")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	user, err := s.userCase.Authentication(r.Context(), params)
	if err != nil || !isValidPassword(params.Password) || !isValidUsername(params.Username) {
		s.log.Warn(err.Error())
		err = responseError(w, "bad_credentials", "invalid user credentials")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	session, err := s.sm.CreateNewSessionForUser(r.Context(), user.ID)
	if err != nil {
		s.log.Error(err.Error())
		err = responseError(w, "session", "failed to create a session for the user")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	cookie := &http.Cookie{
		Name:     "session_key",
		Value:    session.Key,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  session.Expire,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	err = responseOk(w, "a new session has been created for the user", nil)
	if err != nil {
		s.log.Error(err.Error())
	}
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
	s.log.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	defer r.Body.Close()
	user := &user.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		s.log.Info("failed to parse parameters", log.F{"error", err.Error()})
		err = responseError(w, "parse_body", "the correct username, email and password are expected to be received in JSON format")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	if !IsValidUserForRegistration(user) {
		fmt.Println(isValidEmail(user.Email), isValidUsername(user.Username), isValidPassword(user.Password))
		s.log.Info("invalid user registration data")
		err = responseError(w, "invalid_params", "invalid data for registration is specified")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	err = s.userCase.Register(r.Context(), user)
	if err != nil {
		s.log.Warn(err.Error())
		err = responseError(w, "uniq_fields", "username")
	} else {
		err = responseOk(w, "the user has been successfully registered", nil)
	}
	if err != nil {
		s.log.Error(err.Error())
	}
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
	s.log.Info("request on logout", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	cookie, err := r.Cookie("session_key")
	if err != nil {
		s.log.Info("no cookie", log.F{"error", err.Error()})
		err = responseError(w, "no_auth", "to log out, you must first log in")
		if err != nil {
			s.log.Error(err.Error())
		}
		return
	}

	cookie.Expires = time.Now().UTC().AddDate(0, -1, 0)
	http.SetCookie(w, cookie)

	err = s.sm.DeleteUserSession(r.Context(), cookie.Value)
	if err != nil {
		s.log.Error(err.Error())
		err = responseError(w, "session", "the user logged out, but his session did not end")
	} else {
		err = responseOk(w, "the user has successfully logged out", nil)
	}
	if err != nil {
		s.log.Error(err.Error())
	}
}
