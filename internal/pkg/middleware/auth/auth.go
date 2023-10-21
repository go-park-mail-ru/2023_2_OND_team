package auth

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type authContextValueKey string

const (
	KeyCurrentUserID authContextValueKey = "userID"

	SessionCookieName string = "session_key"
)

type authMiddleware struct {
	sm  session.SessionManager
	log *logger.Logger
}

func NewAuthMiddleware(sm session.SessionManager, log *logger.Logger) authMiddleware {
	return authMiddleware{sm, log}
}

func (am authMiddleware) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionCookieName)
		if err == http.ErrNoCookie {
			return
		}
		if err != nil {
			am.log.Sugar().Errorf("auth middleware: %s", err)
		}

		userID, err := am.sm.GetUserIDBySessionKey(r.Context(), cookie.Value)
		if err != nil {
			return
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), KeyCurrentUserID, userID)))
	})
}
