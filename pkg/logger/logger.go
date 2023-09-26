package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New(options ...ConfigOption) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	for _, opt := range options {
		opt(&cfg)
	}
	log, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("new Logger: %w", err)
	}
	return &Logger{Logger: log}, nil
}
