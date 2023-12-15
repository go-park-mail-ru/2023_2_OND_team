package websocket

import (
	"context"
	"fmt"
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

	socket := newSocketJSON(conn)

	err = h.subscribeOnNotificationAndServe(ctx, socket, userID)
	if err != nil && ws.CloseStatus(err) == -1 {
		h.log.Error(err.Error())
		conn.Close(ws.StatusInternalError, "subscribe_fail")
	}
}

func (h *HandlerWebSocket) subscribeOnNotificationAndServe(ctx context.Context, w CtxWriter, userID int) error {
	chanNotify, err := h.notifySub.SubscribeOnAllNotifications(ctx, userID)
	if err != nil {
		return fmt.Errorf("subscribe on Notification")
	}

	for notify := range chanNotify {
		if notify.Err() != nil {
			return notify.Err()
		}

		err = w.Write(ctx, notify)
		if err != nil {
			h.log.Error(err.Error())
		}
	}

	return nil
}
