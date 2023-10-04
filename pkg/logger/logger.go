package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type F struct {
	FieldName string
	Value     string
}

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

func (log *Logger) Info(msg string, fields ...F) {
	listFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		listFields = append(listFields, zap.Field{
			Key:    field.FieldName,
			Type:   zapcore.StringType,
			String: field.Value,
		})
	}
	log.Logger.Info(msg, listFields...)
}
