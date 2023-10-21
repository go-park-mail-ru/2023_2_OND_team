package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

// Login godoc
//
//	@Description	User login, check authentication, get user info
//	@Tags			Auth
//	@Produce		json
//	@Param			session_key	header		string	false	"Auth session id"	example(senjs7rvdnrgkjdr)
//	@Success		200			{object}	JsonResponse{body=user.User}
//	@Failure		400			{object}	JsonErrResponse
//	@Failure		404			{object}	JsonErrResponse
//	@Failure		500			{object}	JsonErrResponse
//	@Router			/api/v1/auth/login [get]
func (h *HandlerHTTP) CheckLogin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on check login", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	userID, _ := r.Context().Value(auth.KeyCurrentUserID).(int)

	username, avatar, err := h.userCase.FindOutUsernameAndAvatar(r.Context(), userID)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "no_auth", "no user was found for this session")
	} else {
		err = responseOk(w, "user found", map[string]string{"username": username, "avatar": avatar})
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}

// Login godoc
//
//	@Description	User login, creating new session
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			username	body		string	true	"Username"	example(clicker123)
//	@Param			password	body		string	true	"Password"	example(safe_pass)
//	@Success		200			{object}	JsonResponse
//	@Failure		400			{object}	JsonErrResponse
//	@Failure		404			{object}	JsonErrResponse
//	@Failure		500			{object}	JsonErrResponse
//	@Header			200			{string}	session_key	"Auth cookie with new valid session id"
//	@Router			/api/v1/auth/login [post]
func (h *HandlerHTTP) Login(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	params := usecase.NewCredentials()
	err := json.NewDecoder(r.Body).Decode(&params)
	defer r.Body.Close()
	if err != nil {
		h.log.Info("failed to parse parameters", log.F{"error", err.Error()})
		err = responseError(w, "parse_body", "the correct username and password are expected to be received in JSON format")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	if !isValidPassword(params.Password) || !isValidUsername(params.Username) {
		h.log.Info("invalid credentials")
		err = responseError(w, "invalid_credentials", "invalid user credentials")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	user, err := h.userCase.Authentication(r.Context(), params)
	if err != nil {
		h.log.Warn(err.Error())
		err = responseError(w, "bad_credentials", "incorrect user credentials")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	session, err := h.sm.CreateNewSessionForUser(r.Context(), user.ID)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "session", "failed to create a session for the user")
		if err != nil {
			h.log.Error(err.Error())
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
		h.log.Error(err.Error())
	}
}

// SignUp godoc
//
//	@Description	User registration
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			username	body		string	true	"Username"	example(clicker123)
//	@Param			email		body		string	true	"Email"		example(clickkk@gmail.com)
//	@Param			password	body		string	true	"Password"	example(safe_pass)
//	@Success		200			{object}	JsonResponse
//	@Failure		400			{object}	JsonErrResponse
//	@Failure		404			{object}	JsonErrResponse
//	@Failure		500			{object}	JsonErrResponse
//	@Router			/api/v1/auth/signup [post]
func (h *HandlerHTTP) Signup(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	user := &user.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	defer r.Body.Close()
	if err != nil {
		h.log.Info("failed to parse parameters", log.F{"error", err.Error()})
		err = responseError(w, "parse_body", "the correct username, email and password are expected to be received in JSON format")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	if err := IsValidUserForRegistration(user); err != nil {
		h.log.Info("invalid user registration data")
		err = responseError(w, "invalid_params", err.Error())
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	err = h.userCase.Register(r.Context(), user)
	if err != nil {
		h.log.Warn(err.Error())
		err = responseError(w, "uniq_fields", "there is already an account with this username or password")
	} else {
		err = responseOk(w, "the user has been successfully registered", nil)
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}

// Logout godoc
//
//	@Description	User logout, session deletion
//	@Tags			Auth
//	@Produce		json
//	@Param			session_key	header		string	false	"Auth session id"	example(senjs7rvdnrgkjdr)
//
//	@Success		200			{object}	JsonResponse
//	@Failure		400			{object}	JsonErrResponse
//	@Failure		404			{object}	JsonErrResponse
//	@Failure		500			{object}	JsonErrResponse
//	@Header			200			{string}	Session-id	"Auth cookie with expired session id"
//	@Router			/api/v1/auth/logout [delete]
func (h *HandlerHTTP) Logout(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on logout", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	cookie, err := r.Cookie("session_key")
	if err != nil {
		h.log.Info("no cookie", log.F{"error", err.Error()})
		err = responseError(w, "no_auth", "to log out, you must first log in")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	cookie.Expires = time.Now().UTC().AddDate(0, -1, 0)
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	err = h.sm.DeleteUserSession(r.Context(), cookie.Value)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "session", "the user logged out, but his session did not end")
	} else {
		err = responseOk(w, "the user has successfully logged out", nil)
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}
