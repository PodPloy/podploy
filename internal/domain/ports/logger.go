package ports

import (
	"context"
	"time"
)

type fieldType uint8

// Field type constants identify the kind of value stored in a Field.
const (
	UnknownType  fieldType = iota // UnknownType represents an unrecognized field type.
	StringType                    // StringType represents a string-valued field.
	IntType                       // IntType represents an integer-valued field.
	BoolType                      // BoolType represents a boolean-valued field.
	ErrorType                     // ErrorType represents an error-valued field.
	DurationType                  // DurationType represents a time.Duration-valued field.
)

// Field is a structured key-value pair used to attach metadata to log entries.
// The Type discriminator indicates which value field is populated.
type Field struct {
	Key       string
	Type      fieldType
	StringVal string
	IntVal    int64
	Any       any
}

// String creates a Field that carries a string value.
func String(key, val string) Field {
	return Field{Key: key, Type: StringType, StringVal: val}
}

// Int creates a Field that carries an integer value.
func Int(key string, val int) Field {
	return Field{Key: key, Type: IntType, IntVal: int64(val)}
}

// Bool creates a Field that carries a boolean value.
func Bool(key string, val bool) Field {
	ival := int64(0)
	if val {
		ival = 1
	}
	return Field{Key: key, Type: BoolType, IntVal: ival}
}

// Error creates a Field that carries an error under the key "error".
func Error(err error) Field {
	return Field{Key: "error", Type: ErrorType, Any: err}
}

// Duration creates a Field that carries a time.Duration value.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, IntVal: int64(val)}
}

// ILogger defines a structured, leveled logging interface. Implementations
// must support contextual enrichment via the With* methods, which return a
// new logger instance with the additional metadata attached.
type ILogger interface {
	// Debug logs a message at the Debug level.
	Debug(msg string, fields ...Field)
	// Info logs a message at the Info level.
	Info(msg string, fields ...Field)
	// Warn logs a message at the Warn level.
	Warn(msg string, fields ...Field)
	// Error logs a message at the Error level.
	Error(msg string, fields ...Field)
	// Fatal logs a message at the Fatal level and terminates the process.
	Fatal(msg string, fields ...Field)

	// With returns a logger that includes the given fields in every entry.
	With(fields ...Field) ILogger
	// WithError returns a logger with the given error attached.
	WithError(err error) ILogger
	// WithContext returns a logger enriched with values from the context.
	WithContext(ctx context.Context) ILogger
	// WithUser returns a logger with the given user ID attached.
	WithUser(userID string) ILogger
	// WithRequest returns a logger enriched with HTTP request metadata.
	WithRequest(requestID, method, path, status, ip string, latency time.Duration) ILogger

	// Sync flushes any buffered log entries.
	Sync() error
}
