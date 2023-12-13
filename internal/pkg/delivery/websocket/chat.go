package websocket

import (
	"context"
	"fmt"
	"net/http"

	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
)

func (h *HandlerWebSocket) Chat(w http.ResponseWriter, r *http.Request) {
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

	err = h.subscribeOnChat(ctx, conn, userID)
	if err != nil {
		h.log.Error(err.Error())
		conn.Close(ws.StatusInternalError, "subscribe_fail")
		return
	}

	err = h.serveChat(ctx, conn, userID)
	if err != nil && ws.CloseStatus(err) == -1 {
		h.log.Error(err.Error())
		conn.Close(ws.StatusInternalError, "serve_chat")
	}
}

func (h *HandlerWebSocket) serveChat(ctx context.Context, conn *ws.Conn, userID int) error {
	request := &PublsihRequest{}
	var err error
	for {
		err = wsjson.Read(ctx, conn, request)
		if err != nil {
			h.log.Error(err.Error())
			return fmt.Errorf("read message: %w", err)
		}

		h.handlePublishRequestMessage(ctx, conn, userID, request)
	}
}

func (h *HandlerWebSocket) handlePublishRequestMessage(ctx context.Context, conn *ws.Conn, userID int, r *PublsihRequest) {
	switch r.Message.Type {
	case "create":
		r.Message.Message.From = userID
		id, err := h.messageCase.SendMessage(ctx, userID, &r.Message.Message)
		if err != nil {
			h.log.Warn(err.Error())
			return
		}
		wsjson.Write(ctx, conn, newResponseOnRequest(r.ID, "ok", "", "publish success", map[string]any{"id": id, "eventType": "create"}))

	case "update":
		err := h.messageCase.UpdateContentMessage(ctx, userID, &r.Message.Message)
		if err != nil {
			h.log.Warn(err.Error())
			return
		}
		wsjson.Write(ctx, conn, newResponseOnRequest(r.ID, "ok", "", "publish success", map[string]string{"eventType": "update"}))

	case "delete":
		err := h.messageCase.DeleteMessage(ctx, userID, &r.Message.Message)
		if err != nil {
			h.log.Warn(err.Error())
			return
		}
		wsjson.Write(ctx, conn, newResponseOnRequest(r.ID, "ok", "", "publish success", map[string]string{"eventType": "delete"}))

	default:
		wsjson.Write(ctx, conn, newResponseOnRequest(r.ID, "error", "unsupported", "unsupported eventType", nil))
	}
}

func (h *HandlerWebSocket) subscribeOnChat(ctx context.Context, conn *ws.Conn, userID int) error {
	chanEvMsg, err := h.messageCase.SubscribeUserToAllChats(ctx, userID)
	if err != nil {
		return fmt.Errorf("subscribe user on chat: %w", err)
	}

	go func() {
		for eventMessage := range chanEvMsg {
			if eventMessage.Err != nil {
				h.log.Error(err.Error())
				return
			}

			err = wsjson.Write(ctx, conn, newMessageFromChannel("ok", "", Object{
				Type:    eventMessage.Type,
				Message: *eventMessage.Message,
			}))
			if err != nil {
				h.log.Error(err.Error())
				return
			}
		}
	}()
	return nil
}
