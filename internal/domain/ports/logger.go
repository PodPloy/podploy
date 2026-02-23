package ports

import (
	"context"
	"time"
)

type fieldType uint8

const (
	UnknownType fieldType = iota
	StringType
	IntType
	BoolType
	ErrorType
	DurationType
)

type Field struct {
	Key       string
	Type      fieldType
	StringVal string
	IntVal    int64
	Any       any
}

func String(key, val string) Field {
	return Field{Key: key, Type: StringType, StringVal: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Type: IntType, IntVal: int64(val)}
}

func Bool(key string, val bool) Field {
	ival := int64(0)
	if val {
		ival = 1
	}
	return Field{Key: key, Type: BoolType, IntVal: ival}
}

func Error(err error) Field {
	return Field{Key: "error", Type: ErrorType, Any: err}
}

func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, IntVal: int64(val)}
}

type ILogger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	With(fields ...Field) ILogger
	WithError(err error) ILogger
	WithContext(ctx context.Context) ILogger
	WithUser(userID string) ILogger
	WithRequest(requestID, method, path, status, ip string, latency time.Duration) ILogger

	Sync() error
}
