package server

import (
	"fmt"
	"strconv"

	"go.uber.org/config"
)

const ConfigName = "app.server"

type Config struct {
	Host     string
	Port     string
	CertFile string
	KeyFile  string
	https    bool
}

func NewConfig(filename string) (*Config, error) {
	cfg, err := config.NewYAML(config.File(filename))
	if err != nil {
		return nil, fmt.Errorf("new YAML provider: %w", err)
	}

	value := cfg.Get(ConfigName)
	c := &Config{
		Host:     value.Get("host").String(),
		Port:     value.Get("port").String(),
		CertFile: value.Get("certFile").String(),
		KeyFile:  value.Get("keyFile").String(),
	}

	https, err := strconv.ParseBool(value.Get("https").String())
	if err != nil {
		return nil, fmt.Errorf("parse param https in server.Config: %w", err)
	}
	c.https = https
	return c, nil
}
