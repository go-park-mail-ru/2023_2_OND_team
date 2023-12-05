package realtime

import (
	"fmt"
	"os"
	"sync"

	"google.golang.org/protobuf/proto"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
)

type Node struct {
	subscriptions [_numWorkers]SubscriberHub
	broker        Broker
	numWorkers    int
	mu            sync.RWMutex
}

func NewNode() (*Node, error) {
	node := &Node{}

	broker, err := NewKafkaBroker(node, KafkaConfig{
		// Addres:            []string{"localhost:9092"},
		Addres:            []string{os.Getenv("KAFKA_BROKER_ADDRESS") + ":" + os.Getenv("KAFKA_BROKER_PORT")},
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
		client.transport.Send(mes)
	}
}

func (n *Node) AddSubscriber(c *rt.Channel, client *Client) {
	channel := Channel{
		Name:  c.GetName(),
		Topic: c.GetTopic(),
	}

	ind := index(channel.Name)

	n.mu.Lock()
	defer n.mu.Unlock()
	subscribeChannel, ok := n.subscriptions[ind][channel]
	if !ok {
		subscribeChannel = make(map[string]*Client)
		n.subscriptions[ind][channel] = subscribeChannel
	}
	subscribeChannel[client.id.String()] = client
}
