package security

import (
	"net/http"
	"time"
)

type CSRFConfig struct {
	CookieConfig
	SkipMethods           []string
	HeaderSet             string
	Header                string
	PathToGet             string
	Lifetime              time.Duration
	LenToken              int
	UpdateWithEachRequest bool
}

type CookieConfig struct {
	CookieName     string
	CookieDomain   string
	CookietPath    string
	CookieSecure   bool
	CookieHTTPOnly bool
	CookieSameSite http.SameSite
}

func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		CookieConfig: CookieConfig{
			CookieName:     "_csrf",
			CookietPath:    "/",
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: http.SameSiteStrictMode,
		},
		SkipMethods: []string{http.MethodGet, http.MethodHead, http.MethodOptions},
		HeaderSet:   "X-Set-CSRF-Token",
		Header:      "X-CSRF-Token",
		LenToken:    16,
		Lifetime:    time.Hour,
	}
}
