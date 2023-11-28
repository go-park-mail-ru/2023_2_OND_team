package realtime

import (
	"hash/fnv"

	"github.com/google/uuid"

	rt "github.com/go-park-mail-ru/2023_2_OND_team/internal/api/realtime"
)

type Channel struct {
	Name  string
	Topic string
}

type SubscriberHub map[Channel]map[string]*Client

type Client struct {
	id        uuid.UUID
	transport rt.RealTime_SubscribeServer
}

func index(nameChannel string) uint32 {
	h := fnv.New32()
	h.Write([]byte(nameChannel))
	return h.Sum32() % _numWorkers
}
