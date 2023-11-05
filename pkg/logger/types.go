package logger

import "go.uber.org/zap/zapcore"

type (
	F struct {
		FieldName string
		Value     any
	}

	M map[string]any
)

type zapLogFn func(msg string, fields ...zapcore.Field)

type ctxLoggerKey int8

var KeyLogger ctxLoggerKey = 0

func mapToSliceFields(m M) []F {
	fields := make([]F, 0, len(m))
	for key, val := range m {
		fields = append(fields, F{key, val})
	}
	return fields
}
