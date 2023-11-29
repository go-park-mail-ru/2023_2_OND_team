package realtime

import (
	"context"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const _numWorkers = 64

type Server struct {
	rt.UnimplementedRealTimeServer

	node *Node
}

func NewServer(node *Node) Server {
	return Server{
		UnimplementedRealTimeServer: rt.UnimplementedRealTimeServer{},
		node:                        node,
	}
}

type Broker interface {
	Publish(topic string, channel string, message []byte) error
	Close()
}

func (s Server) Publish(ctx context.Context, pm *rt.PublishMessage) (*empty.Empty, error) {
	message, err := proto.Marshal(pm.Message)
	if err != nil {
		return nil, status.Error(codes.Internal, "marshaling message")
	}
	if err := s.node.broker.Publish(pm.Channel.GetTopic(), pm.Channel.GetName(), message); err != nil {
		return nil, status.Error(codes.Internal, "publish message")
	}
	return &empty.Empty{}, nil
}

func (s Server) Subscribe(c *rt.Channel, ss rt.RealTime_SubscribeServer) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return status.Error(codes.Internal, "generate uuid v4")
	}
	client := &Client{
		id:        id,
		transport: ss,
	}

	s.node.AddSubscriber(c, client)

	<-ss.Context().Done()
	return nil
}
