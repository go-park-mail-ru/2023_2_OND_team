package realtime

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/api/realtime"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const _numWorkers = 64

type Channel struct {
	Name  string
	Topic string
}

type SubscriberHub map[Channel]map[string]*Client

type Node struct {
	subscriptions [_numWorkers]SubscriberHub
	broker        Broker
	numWorkers    int
	mu            sync.RWMutex
}

func NewNode() (*Node, error) {
	node := &Node{}
	broker, err := NewKafkaBroker(node, KafkaConfig{
		Addres:            []string{"localhost:9092"},
		PartitionsOnTopic: _numWorkers,
		MaxNumTopic:       10,
	})
	if err != nil {
		return nil, fmt.Errorf("new kafka broker: %w", err)
	}
	node.broker = broker
	for ind := range node.subscriptions {
		node.subscriptions[ind] = make(SubscriberHub)
	}
	node.numWorkers = _numWorkers
	node.mu = sync.RWMutex{}
	return node, nil
}

func (n *Node) SendOut(channel Channel, message []byte) {
	clients := []*Client{}
	ind := index(channel.Name)
	n.mu.RLock()

	if m, ok := n.subscriptions[ind][channel]; ok {
		for _, client := range m {
			clients = append(clients, client)
		}
	}
	n.mu.RUnlock()
	for _, client := range clients {
		mes := &rt.Message{}
		err := proto.Unmarshal(message, mes)
		if err != nil {
			fmt.Println(err)
		}
		client.Transport.Send(mes)
	}
	mes := &rt.Message{}
	err := proto.Unmarshal(message, mes)
	if err != nil {
		fmt.Println(err)
	}
}

type Client struct {
	ID        uuid.UUID
	Transport rt.RealTime_SubscribeServer
}

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
	if err := s.node.broker.Publish(pm.Channel.Topic, pm.Channel.Name, message); err != nil {
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
		ID:        id,
		Transport: ss,
	}
	channel := Channel{
		Name:  c.GetName(),
		Topic: c.GetTopic(),
	}
	ind := index(channel.Name)

	s.node.mu.Lock()
	subscribeChannel, ok := s.node.subscriptions[ind][channel]
	if !ok {
		subscribeChannel = make(map[string]*Client)
		s.node.subscriptions[ind][channel] = subscribeChannel
	}
	subscribeChannel[id.String()] = client
	s.node.mu.Unlock()
	select {}
}

func index(nameChannel string) uint32 {
	h := fnv.New32()
	h.Write([]byte(nameChannel))
	return h.Sum32() % _numWorkers
}
