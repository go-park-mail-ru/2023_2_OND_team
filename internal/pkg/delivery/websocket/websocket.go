package websocket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/message"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type HandlerWebSocket struct {
	originPatterns []string
	log            *log.Logger
	messageCase    usecase.Usecase
	client         rt.RealTimeClient
}

type Option func(h *HandlerWebSocket)

const _ctxOnServeConnect = 24 * time.Hour

func SetOriginPatterns(patterns []string) Option {
	return func(h *HandlerWebSocket) {
		h.originPatterns = patterns
	}
}

func New(log *log.Logger, mesCase usecase.Usecase, opts ...Option) *HandlerWebSocket {
	gRPCConn, err := grpc.Dial("localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(fmt.Errorf("grpc dial: %w", err).Error())
	}

	client := rt.NewRealTimeClient(gRPCConn)
	handlerWS := &HandlerWebSocket{log: log, messageCase: mesCase, client: client}
	for _, opt := range opts {
		opt(handlerWS)
	}

	return handlerWS
}

func (h *HandlerWebSocket) WebSocketConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Accept(w, r, &ws.AcceptOptions{OriginPatterns: h.originPatterns})
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

	err = h.serveWebSocketConn(ctx, conn, userID)
	if err != nil {
		h.log.Error(err.Error())
	}
}

func (h *HandlerWebSocket) serveWebSocketConn(ctx context.Context, conn *ws.Conn, userID int) error {
	request := &Request{}
	var err error
	for {
		err = wsjson.Read(ctx, conn, request)
		if err != nil {
			h.log.Error(err.Error())
			return fmt.Errorf("read message: %w", err)
		}
		switch request.Action {
		case "Publish":
			switch request.Message.Type {
			case "create":
				mesCopy := &message.Message{}
				*mesCopy = request.Message.Message
				mesCopy.From = userID
				id, err := h.messageCase.SendMessage(ctx, mesCopy)
				if err != nil {
					h.log.Warn(err.Error())
					continue
				}
				wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "ok", "", "publish success", map[string]int{"id": id}))
				_, err = h.client.Publish(ctx, &rt.PublishMessage{
					Channel: &rt.Channel{
						Name:  request.Channel.Name,
						Topic: request.Channel.Topic,
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
				wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "ok", "", "publish success", nil))
				_, err = h.client.Publish(ctx, &rt.PublishMessage{
					Channel: &rt.Channel{
						Name:  request.Channel.Name,
						Topic: request.Channel.Topic,
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
				wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "ok", "", "publish success", nil))
				_, err = h.client.Publish(ctx, &rt.PublishMessage{
					Channel: &rt.Channel{
						Name:  request.Channel.Name,
						Topic: request.Channel.Topic,
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
		case "Subscribe":
			err = h.subscribe(ctx, h.client, request, conn)
			if err != nil {
				h.log.Warn(err.Error())
				wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "error", "subscribe_fail", "failed to subscribe to the channel", nil))
				continue
			}
			wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "ok", "", "you have successfully subscribed to the channel", nil))
		default:
			wsjson.Write(ctx, conn, newResponseOnRequest(request.ID, "error", "unsupported", "unsupported action", nil))
		}
	}
}

func (h *HandlerWebSocket) subscribe(ctx context.Context, client rt.RealTimeClient, req *Request, conn *ws.Conn) error {
	sc, err := client.Subscribe(ctx, &rt.Channel{
		Name:  req.Channel.Name,
		Topic: req.Channel.Topic,
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
					msg, err = h.messageCase.GetMessage(ctx, int(mes.Object.Id))
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
				err = wsjson.Write(ctx, conn, newMessageFromChannel(req.Channel, "ok", "", Object{
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
