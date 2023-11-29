package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/microservices/realtime"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

const _address = "localhost:8090"

func main() {
	log, err := logger.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := RealTimeRun(log, _address); err != nil {
		log.Error(err.Error())
	}
}

func RealTimeRun(log *logger.Logger, addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen tcp %s: %w", addr, err)
	}

	node, err := realtime.NewNode()
	if err != nil {
		return fmt.Errorf("new server node: %w", err)
	}

	serv := grpc.NewServer()
	rt.RegisterRealTimeServer(serv, realtime.NewServer(node))

	log.Info("start realtime server", logger.F{"network", "tcp"}, logger.F{"addr", addr})
	return serv.Serve(l)
}