package server

import "go.uber.org/config"

type Config struct {
	Host     string
	Port     string
	CertFile string
	KeyFile  string
}

const ConfigName = "app.server"

func NewConfig(cfg *config.YAML) *Config {
	value := cfg.Get(ConfigName)
	return &Config{
		Host:     value.Get("host").String(),
		Port:     value.Get("port").String(),
		CertFile: value.Get("certFile").String(),
		KeyFile:  value.Get("keyFile").String(),
	}
}
