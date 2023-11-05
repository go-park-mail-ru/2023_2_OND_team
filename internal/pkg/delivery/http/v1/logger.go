package v1

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) getRequestLogger(r *http.Request) *logger.Logger {
	if log, ok := r.Context().Value(logger.KeyLogger).(*logger.Logger); ok {
		return log
	}
	return h.log
}
