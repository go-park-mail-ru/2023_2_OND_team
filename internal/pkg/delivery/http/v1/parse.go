package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func fetchURLParamInt(r *http.Request, param string) (int, error) {
	paramInt64, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("fetch integer param from query request %s: %w", r.URL.RawQuery, err)
	}
	return int(paramInt64), nil
}

func decodeBody(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("decode body: %w", err)
	}
	return nil
}
