package auth

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/session"
)

type authContextValueKey string

const (
	KeyCurrentUserID authContextValueKey = "userID"

	SessionCookieName string = "session_key"
)

type authMiddleware struct {
	sm session.SessionManager
}

func NewAuthMiddleware(sm session.SessionManager) authMiddleware {
	return authMiddleware{sm}
}

func (am authMiddleware) ContextWithUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(SessionCookieName); err == nil {
			if userID, err := am.sm.GetUserIDBySessionKey(r.Context(), cookie.Value); err == nil {
				r = r.WithContext(context.WithValue(r.Context(), KeyCurrentUserID, userID))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(KeyCurrentUserID).(int)
		if ok {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"status": "error", "code": "no_auth", "message": "authentication required"}`))
		}
	})
}
