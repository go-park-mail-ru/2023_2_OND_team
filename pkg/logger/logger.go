package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New() (*Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("new Logger: %w", err)
	}
	return &Logger{Logger: log}, nil
}
