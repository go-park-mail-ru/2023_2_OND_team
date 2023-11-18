package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	fields []F
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

func (log *Logger) WithField(key string, val any) *Logger {
	newFields := make([]F, len(log.fields)+1)
	copy(newFields, log.fields)
	newFields[len(log.fields)] = F{key, val}
	return &Logger{
		Logger: log.Logger,
		fields: newFields,
	}
}

func (log *Logger) Info(msg string, fields ...F) {
	log.multiLevelLog(log.Logger.Info, msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...F) {
	log.multiLevelLog(log.Logger.Warn, msg, fields...)
}

func (log *Logger) Error(msg string, fields ...F) {
	log.multiLevelLog(log.Logger.Error, msg, fields...)
}

func (log *Logger) Infof(template string, args ...any) {
	log.multiLevelSugarLog((*zap.SugaredLogger).Infof, template, args...)
}

func (log *Logger) Warnf(template string, args ...any) {
	log.multiLevelSugarLog((*zap.SugaredLogger).Warnf, template, args...)
}

func (log *Logger) Errorf(template string, args ...any) {
	log.multiLevelSugarLog((*zap.SugaredLogger).Errorf, template, args...)
}

func (log *Logger) InfoMap(msg string, mapFields M) {
	log.Info(msg, mapToSliceFields(mapFields)...)
}

func (log *Logger) makeFields(fields []F) []zap.Field {
	listFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zf := zap.Field{
			Key: field.FieldName,
		}

		switch value := field.Value.(type) {
		case string:
			zf.Type = zapcore.StringType
			zf.String = value
		case int64:
			zf.Type = zapcore.Int64Type
			zf.Integer = value
		case int:
			zf.Type = zapcore.Int64Type
			zf.Integer = int64(value)
		case bool:
			zf.Type = zapcore.BoolType
			zf.Interface = value

		default:
			log.Warn("unknown type field for the logger")
			zf.Type = zapcore.SkipType
		}

		listFields = append(listFields, zf)
	}
	return listFields
}

func (log *Logger) multiLevelLog(logFn zapLogFn, msg string, fields ...F) {
	fields = append(fields, log.fields...)
	logFn(msg, log.makeFields(fields)...)
}

func (log *Logger) multiLevelSugarLog(logFn zapSugarLogFn, template string, args ...any) {
	fieldsSugarLogger := make([]any, 0, 2*len(log.fields))
	for _, field := range log.fields {
		fieldsSugarLogger = append(fieldsSugarLogger, field.FieldName, field.Value)
	}
	logFn(log.Logger.Sugar().With(fieldsSugarLogger...), template, args...)
}
