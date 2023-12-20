package websocket

import (
	"context"

	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type CtxReader interface {
	Read(ctx context.Context, v any) error
}

type CtxWriter interface {
	Write(ctx context.Context, v any) error
}

type CtxReadWriter interface {
	CtxReader
	CtxWriter
}

type socketJSON struct {
	*ws.Conn
}

func newSocketJSON(conn *ws.Conn) socketJSON {
	return socketJSON{conn}
}

func (s socketJSON) Write(ctx context.Context, v any) error {
	return wsjson.Write(ctx, s.Conn, v)
}

func (s socketJSON) Read(ctx context.Context, v any) error {
	return wsjson.Read(ctx, s.Conn, v)
}
