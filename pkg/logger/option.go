package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConfigOption func(cfg *zap.Config)

func SetFormatTime(layout string) ConfigOption {
	return func(cfg *zap.Config) {
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(layout)
	}
}

func RFC3339FormatTime() ConfigOption {
	return SetFormatTime(time.RFC3339)
}

func SetOutputPaths(files ...string) ConfigOption {
	return func(cfg *zap.Config) {
		cfg.OutputPaths = files
	}
}

func SetErrorOutputPaths(files ...string) ConfigOption {
	return func(cfg *zap.Config) {
		cfg.ErrorOutputPaths = files
	}
}
