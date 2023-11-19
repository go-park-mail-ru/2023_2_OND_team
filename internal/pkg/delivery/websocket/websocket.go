package websocket

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type HandlerWebSocket struct {
	log *log.Logger
}

func New(log *log.Logger) *HandlerWebSocket {
	return &HandlerWebSocket{log}
}

type Event struct {
	Event       string `json:"event"`
	ObjectID    int    `json:"objID"`
	Content     string `json:"content"`
	PublisherID int    `json:"publisherID"`
}

var mes [3]Event = [3]Event{
	{Event: "del", ObjectID: 12, PublisherID: 2332},
	{Event: "new", ObjectID: 12, Content: "some text", PublisherID: 2332},
	{Event: "edit", ObjectID: 12, Content: "new some text", PublisherID: 2332},
}

func (h *HandlerWebSocket) WebSocketConnect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conn, err := ws.Accept(w, r, &ws.AcceptOptions{})
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error","code":"websocket_connect","message":"fail connect"}`))
		return
	}
	defer conn.CloseNow()

	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
		err = wsjson.Write(ctx, conn, mes[rand.Int31n(3)])
		if err != nil {
			closeStatus := ws.CloseStatus(err)
			if closeStatus != ws.StatusNormalClosure {
				fmt.Println(closeStatus)
				h.log.Error(err.Error())
			}
			if closeStatus == -1 {
				conn.Close(ws.StatusAbnormalClosure, "error write")
			}
			return
		}
	}
}
