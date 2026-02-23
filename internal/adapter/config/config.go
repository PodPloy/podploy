package config

import (
	"fmt"
	"strings"

	toml "github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	koanf "github.com/knadh/koanf/v2"
)

type rawServerConfig struct {
	Host    string `koanf:"host"`
	Port    uint   `koanf:"port"`
	Origins string `koanf:"origins"`
}

type rawLoggerConfig struct {
	Level      string `koanf:"level"`
	Env        string `koanf:"environment"`
	OutputPath string `koanf:"output_path"`
	MaxAge     uint   `koanf:"max_age"`
	MaxSize    uint   `koanf:"max_size"`
	MaxBackups uint   `koanf:"max_backups"`
}

type rawHubConfig struct {
	Server rawServerConfig `koanf:"server"`
	Logger rawLoggerConfig `koanf:"logs"`
}

type ServerConfig struct {
	port    uint
	host    string
	origins []string
}

type LoggerConfig struct {
	level      string
	env        string
	outputPath string
	maxAge     uint
	maxSize    uint
	maxBackups uint
}

type HubConfig struct {
	server ServerConfig
	logger LoggerConfig
}

const (
	defaultHost       = "0.0.0.0"
	defaultPort       = 8080
	defaultEnv        = "production"
	defaultLogLevel   = "info"
	defaultOutputPath = "~/.config/podploy/hub.log"
	defaultMaxSize    = 100 // MB
	defaultMaxBackups = 5
	defaultMaxAge     = 30 // days
	defaultOrigins    = "*"
)

func LoadHubConfig(path string) (*HubConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var raw rawHubConfig
	if err := k.Unmarshal("", &raw); err != nil {
		return nil, fmt.Errorf("unmarshal hub: %w", err)
	}
	fmt.Println(raw)

	conf := &HubConfig{
		server: ServerConfig{
			host:    raw.Server.Host,
			port:    raw.Server.Port,
			origins: strings.Split(raw.Server.Origins, ","),
		},
		logger: LoggerConfig{
			level:      raw.Logger.Level,
			env:        raw.Logger.Env,
			outputPath: raw.Logger.OutputPath,
			maxAge:     raw.Logger.MaxAge,
			maxSize:    raw.Logger.MaxSize,
			maxBackups: raw.Logger.MaxBackups,
		},
	}

	if err := normalizeHubConfig(conf); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return conf, nil
}

func normalizeHubConfig(conf *HubConfig) error {
	if conf.server.host == "" {
		conf.server.host = defaultHost
	}

	if conf.server.port == 0 {
		conf.server.port = defaultPort
	} else if conf.server.port > 65535 {
		return fmt.Errorf("port %d out of range [1-65535]", conf.server.port)
	}

	if len(conf.server.origins) == 0 || conf.server.origins[0] == "" {
		conf.server.origins = []string{defaultOrigins}
	}

	if conf.logger.outputPath == "" {
		conf.logger.outputPath = defaultOutputPath
	}

	if conf.logger.maxSize == 0 {
		conf.logger.maxSize = defaultMaxSize
	}

	if conf.logger.maxBackups == 0 {
		conf.logger.maxBackups = defaultMaxBackups
	}

	if conf.logger.maxAge <= 0 {
		conf.logger.maxAge = defaultMaxAge
	}

	conf.logger.env = normalizeString(conf.logger.env, defaultEnv)

	logLevel, err := normalizeLogLevel(conf.logger.level)
	if err != nil {
		return err
	}

	conf.logger.level = logLevel

	return nil
}

func normalizeString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func normalizeLogLevel(level string) (string, error) {
	if level == "" {
		return defaultLogLevel, nil
	}

	nLevel := strings.ToLower(level)

	switch nLevel {
	case "debug", "info", "warn", "error", "fatal":
		return nLevel, nil
	default:
		return "", fmt.Errorf("unknown log level '%s'. Supported: debug, info, warn, error, fatal", level)
	}
}

func (h *HubConfig) Server() ServerConfig {
	return h.server
}

func (h *HubConfig) Logger() LoggerConfig {
	return h.logger
}

func (s *ServerConfig) Host() string {
	return s.host
}

func (s *ServerConfig) Origins() []string {
	return s.origins
}

func (s *ServerConfig) Port() uint {
	return s.port
}

func (l *LoggerConfig) Level() string {
	return l.level
}

func (l *LoggerConfig) Env() string {
	return l.env
}

func (l *LoggerConfig) OutputPath() string {
	return l.outputPath
}

func (l *LoggerConfig) MaxAge() uint {
	return l.maxAge
}

func (l *LoggerConfig) MaxSize() uint {
	return l.maxSize
}

func (l *LoggerConfig) MaxBackups() uint {
	return l.maxBackups
}
