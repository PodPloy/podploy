package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wrapper around zap.SugaredLogger with dynamic log level control.
// It provides structured logging with context fields and supports both
// development and production configurations.
type Logger struct {
	*zap.SugaredLogger
	atomicLevel zap.AtomicLevel
	path  		string
}

type Config struct {
	Level       string
	OutputPath  string
	Development bool
	MaxSize     int
	MaxBackups  int
	MaxAge      int
}

type contextKey string

const defaultLevel = zapcore.InfoLevel

const RequestIDKey contextKey = "request_id"

func (c *Config) validate() error {
    if c.OutputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if _, err := parseLevel(c.Level); err != nil {
		return fmt.Errorf("invalid log level %q: %w", c.Level, err)
	}
	return nil
}

// New creates a new Logger instance with the provided configuration.
// It automatically creates the log directory if it doesn't exist.
// In development mode, logs are written to both stdout and the specified file.
// The log level can be changed dynamically at runtime using SetLevel().
func New(cfg *Config) (*Logger, error) {
	if cfg == nil {
		cfg = &Config{
			Level:      "info",
			OutputPath: "~/podploy/podploy.log",
		}
	}

	if err := cfg.validate(); err != nil {
    	return nil, fmt.Errorf("invalid config: %w", err)
	}

	logPath := cfg.OutputPath
	if strings.HasPrefix(logPath, "~/") {
    	dirname, err := os.UserHomeDir()
    	if err != nil {
        	return nil, fmt.Errorf("failed to get user home directory: %w", err)
    	}

    	logPath = filepath.Join(dirname, logPath[2:])
	}


	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	writeSyncer := getWriteSyncer(logPath, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Development)

	initialLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		initialLevel = defaultLevel
	}
	atomicLevel := zap.NewAtomicLevelAt(initialLevel)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
	}

	var encoder zapcore.Encoder
	if cfg.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		atomicLevel, 
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	zapLogger := zap.New(core, options...)

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
		atomicLevel:   atomicLevel,
		path:          logPath,
	}, nil
}

func getWriteSyncer(path string, maxSize, maxBks, maxAge int, development bool) zapcore.WriteSyncer {
    lumberjackLogger := &lumberjack.Logger{
        Filename:   path,
        MaxSize:    maxSize, // MB
        MaxBackups: maxBks,
        MaxAge:     maxAge, // days
        Compress:   true,
    }

	var fileSyncer zapcore.WriteSyncer = zapcore.AddSync(lumberjackLogger)
    
    if development {
        return zapcore.NewMultiWriteSyncer(
            zapcore.AddSync(os.Stdout),
            fileSyncer,
        )
    }
    
    return fileSyncer
}

func (l *Logger) With(fields ...any) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(fields...),
		atomicLevel:   l.atomicLevel,
		path:          l.path,
	}
}

func (l *Logger) WithUser(userID string) *Logger {
	return l.With("user_id", userID)
}

func (l *Logger) WithRequest(requestID, method, path string) *Logger {
	return l.With(
		"request_id", requestID,
		"method", method,
		"path", path,
	)
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
    logger := l

	if deadline, ok := ctx.Deadline(); ok {
		logger = logger.With("deadline", deadline)
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		logger = logger.With("request_id", requestID)
	} else if requestID, ok := ctx.Value("request_id").(string); ok {
		logger = logger.With("request_id", requestID)
	}

	return logger
}

func (l *Logger) WithComponent(component string) *Logger {
	return l.With("component", component)
}

func (l *Logger) WithError(err error) *Logger {
	return l.With("error", err.Error())
}

func (l *Logger) Sync() error {
    err := l.SugaredLogger.Sync()
    if err != nil {
        if strings.Contains(err.Error(), "inappropriate ioctl") || strings.Contains(err.Error(), "bad file descriptor") {
            return nil
        }
    }

    return err
}

func (l *Logger) GetLogPath() string {
	return l.path
}

func (l *Logger) SetLevel(level string) error {
	newLevel, err := parseLevel(level)
	if err != nil {
		return fmt.Errorf("unknown log level: %s", level)
	}
	l.atomicLevel.SetLevel(newLevel)
	return nil
}

func parseLevel(level string) (zapcore.Level, error) {
    l, err := zapcore.ParseLevel(level)
	if err == nil {
		return l, nil
	}
    
    switch strings.ToLower(level) {
    case "warn", "warning":
        return zapcore.WarnLevel, nil
    case "err", "error":
        return zapcore.ErrorLevel, nil
    default:
        return defaultLevel, fmt.Errorf("unknown level: %s", level)
    }
}

func DefaultLogger() (*Logger, error) {
	return New(&Config{
		Level:       "debug",
		OutputPath:  "~/podploy/podploy.log",
		Development: true,
	})
}

func ProductionLogger(path string) (*Logger, error) {
	if path == "" {
		path = "~/podploy/podploy.log"
	}
	
	return New(&Config{
		Level:       "info",
		OutputPath:  path,
		Development: false,
	})
}