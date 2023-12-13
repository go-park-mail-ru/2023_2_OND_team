package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
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

	err = h.subscribe(ctx, conn, userID)
	if err != nil {
		h.log.Error(err.Error())
		conn.Close(ws.StatusInternalError, "subscribe_fail")
		return
	}

	err = h.serveChat(ctx, conn, userID)
	if err != nil {
		h.log.Error(err.Error())
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

		switch request.Message.Type {
		case "create":
			mesCopy := &message.Message{}
			*mesCopy = request.Message.Message
			mesCopy.From = userID
			id, err := h.messageCase.SendMessage(ctx, userID, mesCopy)
			if err != nil {
				h.log.Warn(err.Error())
				continue
			}
			wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "ok", "", "publish success", map[string]any{"id": id, "eventType": "create"}))
			_, err = h.client.Publish(ctx, &rt.PublishMessage{
				Channel: &rt.Channel{
					Name:  request.Channel.Name,
					Topic: _topicChat,
				},
				Message: &rt.Message{
					Body: &rt.Message_Object{
						Object: &rt.EventObject{
							Type: rt.EventType_EV_CREATE,
							Id:   int64(id),
						},
					},
				},
			})
			if err != nil {
				h.log.Error(err.Error())
			}
		case "update":
			mesCopy := &message.Message{}
			*mesCopy = request.Message.Message
			err = h.messageCase.UpdateContentMessage(ctx, userID, mesCopy)
			if err != nil {
				h.log.Warn(err.Error())
				continue
			}
			wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "ok", "", "publish success", map[string]string{"eventType": "update"}))
			_, err = h.client.Publish(ctx, &rt.PublishMessage{
				Channel: &rt.Channel{
					Name:  request.Channel.Name,
					Topic: _topicChat,
				},
				Message: &rt.Message{
					Body: &rt.Message_Object{
						Object: &rt.EventObject{
							Type: rt.EventType_EV_UPDATE,
							Id:   int64(request.Message.Message.ID),
						},
					},
				},
			})
			if err != nil {
				h.log.Error(err.Error())
			}

		case "delete":
			err = h.messageCase.DeleteMessage(ctx, userID, request.Message.Message.ID)
			if err != nil {
				h.log.Warn(err.Error())
				continue
			}
			wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "ok", "", "publish success", map[string]string{"eventType": "delete"}))
			_, err = h.client.Publish(ctx, &rt.PublishMessage{
				Channel: &rt.Channel{
					Name:  request.Channel.Name,
					Topic: _topicChat,
				},
				Message: &rt.Message{
					Body: &rt.Message_Object{
						Object: &rt.EventObject{
							Type: rt.EventType_EV_DELETE,
							Id:   int64(request.Message.Message.ID),
						},
					},
				},
			})
			if err != nil {
				h.log.Error(err.Error())
			}
		default:
			wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "error", "unsupported", "unsupported eventType", nil))
		}
	}
}

func (h *HandlerWebSocket) subscribe(ctx context.Context, conn *ws.Conn, userID int) error {
	channel := Channel{Name: strconv.Itoa(userID)}

	sc, err := h.client.Subscribe(ctx, &rt.Channel{
		Name:  channel.Name,
		Topic: _topicChat,
	})
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	go func() {
		for {
			obj, err := sc.Recv()
			if err != nil {
				return
			}
			mes, ok := obj.Body.(*rt.Message_Object)
			if ok {
				var msg *message.Message
				if mes.Object.Type == rt.EventType_EV_DELETE {
					msg = &message.Message{ID: int(mes.Object.Id)}
				} else {
					msg, err = h.messageCase.GetMessage(ctx, userID, int(mes.Object.Id))
					if err != nil {
						h.log.Error(err.Error())
						return
					}
				}
				objType := ""
				switch mes.Object.Type {
				case rt.EventType_EV_CREATE:
					objType = "create"
				case rt.EventType_EV_UPDATE:
					objType = "update"
				case rt.EventType_EV_DELETE:
					objType = "delete"
				}
				err = wsjson.Write(ctx, conn, newMessageFromChannel(channel, "ok", "", Object{
					Type:    objType,
					Message: *msg,
				}))
				if err != nil {
					h.log.Error(err.Error())
					return
				}
			}
		}
	}()
	return nil
}
