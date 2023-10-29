package security

import (
	"fmt"
	"net/http"
	"time"

	mw "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/crypto"
)

func CSRF(cfg CSRFConfig) mw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == cfg.PathToGet && r.Method == http.MethodGet {
				setToken(w, &cfg)
				return
			}

			skip := isSkipMethod(r.Method, &cfg)
			tokenHeader := r.Header.Get(cfg.Header)
			cookie, err := r.Cookie(cfg.CookieName)

			if err != nil && !skip {
				if len(tokenHeader) == 0 {
					responseCSRFErr(http.StatusBadRequest, w, "missing csrf token in request header")
				} else {
					responseCSRFErr(http.StatusForbidden, w, "invalid csrf token")
				}
				return
			}

			if err != nil {
				setToken(w, &cfg)
				next.ServeHTTP(w, r)
				return
			}

			if len(tokenHeader) != cfg.LenToken || cookie.Value != tokenHeader {
				if skip {
					setToken(w, &cfg)
					next.ServeHTTP(w, r)
				} else {
					responseCSRFErr(http.StatusForbidden, w, "invalid csrf token")
				}
				return
			}

			if skip || cfg.UpdateWithEachRequest {
				setToken(w, &cfg)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isSkipMethod(method string, cfg *CSRFConfig) bool {
	for _, skipMethod := range cfg.SkipMethods {
		if method == skipMethod {
			return true
		}
	}
	return false
}

func setToken(w http.ResponseWriter, cfg *CSRFConfig) error {
	token, err := crypto.NewRandomString(cfg.LenToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("set new csrf token: %w", err)
	}

	cookie := &http.Cookie{
		Name:     cfg.CookieName,
		Value:    token,
		Path:     cfg.CookietPath,
		Domain:   cfg.CookieDomain,
		Secure:   cfg.CookieSecure,
		HttpOnly: cfg.CookieHTTPOnly,
		SameSite: cfg.CookieSameSite,
		Expires:  time.Now().UTC().Add(cfg.Lifetime),
	}
	http.SetCookie(w, cookie)
	w.Header().Set(cfg.HeaderSet, token)
	return nil
}

func responseCSRFErr(status int, w http.ResponseWriter, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "error", "code": "csrf", "message": "` + message + `"}`))
}
