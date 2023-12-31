package v1

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/mailru/easyjson"
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
	logger := h.getRequestLogger(r)
	userID, _ := r.Context().Value(auth.KeyCurrentUserID).(int)

	username, avatar, err := h.userCase.FindOutUsernameAndAvatar(r.Context(), userID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "no_auth", "no user was found for this session")
	} else {
		err = responseOk(http.StatusOK, w, "user found", map[string]any{"username": username, "avatar": avatar, "id": userID})
	}
	if err != nil {
		logger.Error(err.Error())
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
	logger := h.getRequestLogger(r)

	params := &usecase.UserCredentials{}
	err := easyjson.UnmarshalFromReader(r.Body, params)
	defer r.Body.Close()
	if err != nil {
		logger.Info("failed to parse parameters", log.F{"error", err.Error()})
		err = responseError(w, "parse_body", "the correct username and password are expected to be received in JSON format")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	if !isValidPassword(params.Password) || !isValidUsername(params.Username) {
		logger.Info("invalid credentials")
		err = responseError(w, "invalid_credentials", "invalid user credentials")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	session, err := h.authCase.Login(r.Context(), params.Username, params.Password)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "session", "failed to create a session for the user")
		if err != nil {
			logger.Error(err.Error())
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

	err = responseOk(http.StatusCreated, w, "a new session has been created for the user", nil)
	if err != nil {
		logger.Error(err.Error())
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
	logger := h.getRequestLogger(r)

	user := &user.User{}
	err := easyjson.UnmarshalFromReader(r.Body, user)
	defer r.Body.Close()
	if err != nil {
		logger.Info("failed to parse parameters", log.F{"error", err.Error()})
		err = responseError(w, "parse_body", "the correct username, email and password are expected to be received in JSON format")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	if err := IsValidUserForRegistration(user); err != nil {
		logger.Info("invalid user registration data")
		err = responseError(w, "invalid_params", err.Error())
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = h.authCase.Register(r.Context(), user)
	if err != nil {
		logger.Warn(err.Error())
		err = responseError(w, "uniq_fields", "there is already an account with this username or email")
	} else {
		err = responseOk(http.StatusCreated, w, "the user has been successfully registered", nil)
	}
	if err != nil {
		logger.Error(err.Error())
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
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	cookie, err := r.Cookie("session_key")
	if err != nil {
		logger.Info("no cookie", log.F{"error", err.Error()})
		err = responseError(w, "no_auth", "to log out, you must first log in")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	cookie.Expires = time.Now().UTC().AddDate(0, -1, 0)
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	err = h.authCase.Logout(r.Context(), &session.Session{
		Key:    cookie.Value,
		UserID: userID,
		Expire: cookie.Expires,
	})
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "session", "the user logged out, but his session did not end")
	} else {
		err = responseOk(http.StatusOK, w, "the user has successfully logged out", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}
