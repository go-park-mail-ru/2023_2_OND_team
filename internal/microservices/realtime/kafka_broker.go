package realtime

import (
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

var ErrAlreadyMaxNumberTopic = errors.New("the maximum number has already been created topics")

const _timeoutCreateTopic = time.Second

type KafkaConfig struct {
	Addres            []string
	PartitionsOnTopic int
	MaxNumTopic       int
}

type kafkaBroker struct {
	node           *Node
	mainClient     sarama.Client
	existingTopics map[string]struct{}
	cfg            KafkaConfig
	producer       sarama.SyncProducer

	m sync.RWMutex
}

var _ Broker = (*kafkaBroker)(nil)

func NewKafkaBroker(node *Node, cfg KafkaConfig) (*kafkaBroker, error) {
	clientCfg := sarama.NewConfig()
	clientCfg.Producer.Return.Successes = true
	clientCfg.Producer.Partitioner = sarama.NewCustomHashPartitioner(func() hash.Hash32 { return fnv.New32() })

	client, err := sarama.NewClient(cfg.Addres, clientCfg)
	if err != nil {
		return nil, fmt.Errorf("new client for kafka broker: %w", err)
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("new sync producer for kafka broker: %w", err)
	}

	k := &kafkaBroker{
		node:           node,
		cfg:            cfg,
		mainClient:     client,
		producer:       producer,
		m:              sync.RWMutex{},
		existingTopics: make(map[string]struct{}),
	}
	return k, nil
}

func (k *kafkaBroker) Publish(topic string, channel string, message []byte) error {
	created, err := k.checkOrCreateTopic(topic)
	if err != nil {
		return fmt.Errorf("publish to topic %s: %w", topic, err)
	}

	if created {
		err = k.serveTopic(topic)
		if err != nil {
			return fmt.Errorf("serve new topic %s: %w", topic, err)
		}
	}

	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(channel),
		Value: sarama.ByteEncoder(message),
	})
	if err != nil {
		return fmt.Errorf("send message with topic %s and channel %s to kafka server: %w", topic, channel, err)
	}
	return nil
}

func (k *kafkaBroker) Close() {
	k.producer.Close()
	k.mainClient.Close()
}

func (k *kafkaBroker) checkOrCreateTopic(topic string) (bool, error) {
	isExists := false
	var countTopics int
	k.m.RLock()
	if _, ok := k.existingTopics[topic]; ok {
		isExists = true
	} else {
		countTopics = len(k.existingTopics)
	}
	k.m.RUnlock()

	if isExists {
		return false, nil
	}
	if countTopics == k.cfg.MaxNumTopic {
		return false, ErrAlreadyMaxNumberTopic
	}

	detail := &sarama.TopicDetail{
		NumPartitions:     int32(k.cfg.PartitionsOnTopic),
		ReplicationFactor: -1,
	}
	_, err := k.mainClient.LeastLoadedBroker().CreateTopics(&sarama.CreateTopicsRequest{
		ValidateOnly: true,
		Timeout:      _timeoutCreateTopic,
		TopicDetails: map[string]*sarama.TopicDetail{topic: detail},
	})
	if err != nil {
		return false, fmt.Errorf("create topic: %w", err)
	}

	k.m.Lock()
	if _, ok := k.existingTopics[topic]; !ok && len(k.existingTopics) == k.cfg.MaxNumTopic {
		k.m.Unlock()
		go delTopic(topic, k.mainClient)
		return false, ErrAlreadyMaxNumberTopic
	}
	k.existingTopics[topic] = struct{}{}
	k.m.Unlock()

	k.mainClient.RefreshController()
	return true, nil
}

func delTopic(topic string, client sarama.Client) {
	client.LeastLoadedBroker().DeleteTopics(&sarama.DeleteTopicsRequest{
		Topics:  []string{topic},
		Timeout: _timeoutCreateTopic,
	})
}

func (k *kafkaBroker) serveTopic(topic string) error {
	cons, err := sarama.NewConsumer(k.cfg.Addres, sarama.NewConfig())
	if err != nil {
		return fmt.Errorf("serve topic %s: %w", topic, err)
	}

	for i := int32(0); int(i) < k.node.numWorkers; i++ {
		go func(partition int32) {
			offset, err := k.mainClient.GetOffset(topic, partition, -1)
			if err != nil {
				return
			}
			partConsumer, err := cons.ConsumePartition(topic, int32(partition), offset)
			if err != nil {
				return
			}
			for message := range partConsumer.Messages() {
				k.node.SendOut(Channel{Name: string(message.Key), Topic: topic}, message.Value)
			}
		}(i)
	}
	return nil
}
