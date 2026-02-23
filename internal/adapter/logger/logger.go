package logger

import (
	"context"
	"math"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/PodPloy/podploy/internal/domain/ports"
)

// Config holds the parameters for initializing a Logger, including log level,
// output destination, and log rotation settings.
type Config struct {
	Level       string
	OutputPath  string
	Development bool
	MaxSize     uint
	MaxBackups  uint
	MaxAge      uint
}

// Logger is a structured logger backed by zap that implements the
// ports.ILogger interface.
type Logger struct {
	logger *zap.Logger
	cfg    *Config
}

var _ ports.ILogger = (*Logger)(nil)

// New creates a new Logger from the given Config. If cfg is nil, sensible
// defaults are used (info level, stdout, development mode). It configures log
// encoding, output rotation via lumberjack, and returns a ports.ILogger.
func New(cfg *Config) (ports.ILogger, error) {
	if cfg == nil {
		cfg = &Config{Level: "info", OutputPath: "stdout", Development: true}
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var writer zapcore.WriteSyncer
	if cfg.OutputPath == "stdout" {
		writer = zapcore.AddSync(os.Stdout)
	} else {
		lumber := &lumberjack.Logger{
			Filename:   cfg.OutputPath,
			MaxSize:    safeParseUintInteger(cfg.MaxSize),
			MaxBackups: safeParseUintInteger(cfg.MaxBackups),
			MaxAge:     safeParseUintInteger(cfg.MaxAge),
			Compress:   true,
		}
		writer = zapcore.AddSync(lumber)
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	core := zapcore.NewCore(encoder, writer, atomicLevel)

	z := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		logger: z,
		cfg:    cfg,
	}, nil
}

func safeParseUintInteger(number uint) int {
	if number > math.MaxInt {
		return math.MaxInt
	}

	return int(number)
}

func (l *Logger) mapFields(fields []ports.Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}
	zf := make([]zap.Field, len(fields))

	for i, f := range fields {
		switch f.Type {
		case ports.StringType:
			zf[i] = zap.String(f.Key, f.StringVal)
		case ports.IntType:
			zf[i] = zap.Int64(f.Key, f.IntVal)
		case ports.BoolType:
			zf[i] = zap.Bool(f.Key, f.IntVal == 1)
		case ports.DurationType:
			zf[i] = zap.Duration(f.Key, time.Duration(f.IntVal))
		case ports.ErrorType:
			if err, ok := f.Any.(error); ok {
				zf[i] = zap.Error(err)
			} else {
				zf[i] = zap.Any(f.Key, f.Any)
			}
		default:
			zf[i] = zap.Any(f.Key, f.Any)
		}
	}
	return zf
}

// Info logs a message at the Info level with optional structured fields.
func (l *Logger) Info(msg string, fields ...ports.Field) {
	l.logger.Info(msg, l.mapFields(fields)...)
}

// Error logs a message at the Error level with optional structured fields.
func (l *Logger) Error(msg string, fields ...ports.Field) {
	l.logger.Error(msg, l.mapFields(fields)...)
}

// Debug logs a message at the Debug level with optional structured fields.
func (l *Logger) Debug(msg string, fields ...ports.Field) {
	l.logger.Debug(msg, l.mapFields(fields)...)
}

// Warn logs a message at the Warn level with optional structured fields.
func (l *Logger) Warn(msg string, fields ...ports.Field) {
	l.logger.Warn(msg, l.mapFields(fields)...)
}

// Fatal logs a message at the Fatal level with optional structured fields
// and then calls os.Exit(1).
func (l *Logger) Fatal(msg string, fields ...ports.Field) {
	l.logger.Fatal(msg, l.mapFields(fields)...)
}

// With returns a new Logger that includes the given fields in every subsequent
// log entry.
func (l *Logger) With(fields ...ports.Field) ports.ILogger {
	return &Logger{
		logger: l.logger.With(l.mapFields(fields)...),
		cfg:    l.cfg,
	}
}

// WithError returns a new Logger that attaches the given error to every
// subsequent log entry under the "error" key.
func (l *Logger) WithError(err error) ports.ILogger {
	return &Logger{
		logger: l.logger.With(zap.Error(err)),
		cfg:    l.cfg,
	}
}

// WithUser returns a new Logger that attaches the given user ID to every
// subsequent log entry under the "user_id" key.
func (l *Logger) WithUser(userID string) ports.ILogger {
	return &Logger{
		logger: l.logger.With(zap.String("user_id", userID)),
		cfg:    l.cfg,
	}
}

// WithRequest returns a new Logger enriched with HTTP request metadata such as
// request ID, method, path, status, client IP, and response latency.
func (l *Logger) WithRequest(requestID, method, path, status, ip string, latency time.Duration) ports.ILogger {
	return &Logger{
		logger: l.logger.With(
			zap.String("method", method),
			zap.String("path", path),
			zap.String("request_id", requestID),
			zap.String("status", status),
			zap.String("ip", ip),
			zap.Duration("latency", latency),
		),
		cfg: l.cfg,
	}
}

// WithContext returns a new Logger that extracts contextual values (e.g.
// request ID) from the given context and attaches them to log entries.
func (l *Logger) WithContext(ctx context.Context) ports.ILogger {
	newLogger := l.logger
	if reqID, ok := ctx.Value(ports.RequestIDKey).(string); ok {
		newLogger = newLogger.With(zap.String("request_id", reqID))
	}
	return &Logger{logger: newLogger, cfg: l.cfg}
}

// Sync flushes any buffered log entries. It should be called before the
// application exits.
func (l *Logger) Sync() error {
	return l.logger.Sync()
}
