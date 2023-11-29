package app

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/config"
)

type ConfigFiles struct {
	ServerConfigFile string
	AddrAuthServer   string
}

type redisConfig struct {
	Password string
	Addr     string
}

func NewConfig(filename string) (redisConfig, error) {
	cfg, err := config.ParseConfig(filename)
	if err != nil {
		return redisConfig{}, fmt.Errorf("new redis config: %w", err)
	}

	return redisConfig{
		Password: cfg.Get("requirepass"),
		Addr:     cfg.Get("host") + ":" + cfg.Get("port"),
	}, nil
}
