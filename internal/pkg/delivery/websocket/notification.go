package websocket

import (
	"context"
	"net/http"

	ws "nhooyr.io/websocket"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
)

func (h *HandlerWebSocket) Notification(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgradeWSConnect(w, r)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","code":"websocket_connect","message":"fail connect"}`))
		return
	}
	defer conn.CloseNow()

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	ctx, cancel := context.WithTimeout(context.Background(), _ctxOnServeConnect)
	defer cancel()

	err = h.subscribeOnNotification(ctx, conn, userID)
	if err != nil && ws.CloseStatus(err) == -1 {
		h.log.Error(err.Error())
		conn.Close(ws.StatusInternalError, "subscribe_fail")
	}
}

// func (h *HandlerWebSocket) handleNotification(ctx context.Context, conn *ws.Conn, userID int) {

// }

func (h *HandlerWebSocket) subscribeOnNotification(ctx context.Context, conn *ws.Conn, userID int) error {
	return nil
}
