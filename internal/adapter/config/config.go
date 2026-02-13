package config

import (
	"fmt"
	"strings"

	toml "github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	koanf "github.com/knadh/koanf/v2"
)
type HubConfig struct {
	Port 	    int
	Host 	    string
	LogLevel    string
	Environment string
    OutputPath  string
    MaxAge      int
    MaxSize     int
    MaxBackups  int
}

const (
    defaultHost = "0.0.0.0"
    defaultPort = 8080
    defaultEnv  = "production"
    defaultLogLevel = "info"
    defaultOutputPath = "~/podploy/podploy.log"
    defaultMaxSize = 100 // MB
    defaultMaxBackups = 5
    defaultMaxAge = 30 // days
)

func LoadHubConfig(path string) (*HubConfig, error) {
    k := koanf.New(".")

    if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
        return nil, fmt.Errorf("load config: %w", err)
    }

    conf := &HubConfig{}
    if err := k.Unmarshal("hub", conf); err != nil {
        return nil, fmt.Errorf("unmarshal hub: %w", err)
    }

    if err := normalizeHubConfig(conf); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    return conf, nil
}

func normalizeHubConfig(conf *HubConfig) error {
    if conf.Host == "" {
        conf.Host = defaultHost
    }

    if conf.OutputPath == "" {
        conf.OutputPath = defaultOutputPath
    }

    if conf.MaxSize <= 0 {
        conf.MaxSize = defaultMaxSize
    }

    if conf.MaxBackups < 0 {
        conf.MaxBackups = defaultMaxBackups
    }

    if conf.MaxAge <= 0 {
        conf.MaxAge = defaultMaxAge
    }
    
    if conf.Port == 0 {
        conf.Port = defaultPort
    } else if conf.Port < 1 || conf.Port > 65535 {
        return fmt.Errorf("port %d out of range [1-65535]", conf.Port)
    }
    
    conf.Environment = normalizeString(conf.Environment, defaultEnv)

    logLevel, err := normalizeLogLevel(conf.LogLevel)
    if err != nil {
        return err
    }

    conf.LogLevel = logLevel
    
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

