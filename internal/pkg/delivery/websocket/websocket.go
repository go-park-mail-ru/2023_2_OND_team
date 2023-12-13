package websocket

import (
	"fmt"
	"net/http"
	"time"

	ws "nhooyr.io/websocket"

	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type HandlerWebSocket struct {
	originPatterns []string
	log            *log.Logger
	messageCase    usecase.Usecase
}

type Option func(h *HandlerWebSocket)

const _ctxOnServeConnect = 24 * time.Hour

func SetOriginPatterns(patterns []string) Option {
	return func(h *HandlerWebSocket) {
		h.originPatterns = patterns
	}
}

func New(log *log.Logger, mesCase usecase.Usecase, opts ...Option) *HandlerWebSocket {
	handlerWS := &HandlerWebSocket{log: log, messageCase: mesCase}
	for _, opt := range opts {
		opt(handlerWS)
	}

	return handlerWS
}

func (h *HandlerWebSocket) upgradeWSConnect(w http.ResponseWriter, r *http.Request) (*ws.Conn, error) {
	conn, err := ws.Accept(w, r, &ws.AcceptOptions{OriginPatterns: h.originPatterns})
	if err != nil {
		return nil, fmt.Errorf("upgrade to websocket connect: %w", err)
	}
	return conn, nil
}
