package config

import (
	"fmt"
	"strings"

	"github.com/PodPloy/podploy/internal/domain/models"
	toml "github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	koanf "github.com/knadh/koanf/v2"
)

const (
    defaultHost = "0.0.0.0"
    defaultPort = 8080
    defaultEnv  = "production"
    defaultLogLevel = "info"
)

func LoadHubConfig(path string) (*models.HubConfig, error) {
    k := koanf.New(".")

    if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
        return nil, fmt.Errorf("load config: %w", err)
    }

    conf := &models.HubConfig{}
    if err := k.Unmarshal("hub", conf); err != nil {
        return nil, fmt.Errorf("unmarshal hub: %w", err)
    }

    if err := normalizeHubConfig(conf); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    return conf, nil
}

func normalizeHubConfig(conf *models.HubConfig) error {
    if conf.Host == "" {
        conf.Host = defaultHost
    }
    
    if conf.Port == 0 {
        conf.Port = defaultPort
    } else if conf.Port < 1 || conf.Port > 65535 {
        return fmt.Errorf("port %d out of range [1-65535]", conf.Port)
    }
    
    conf.Environment = normalizeString(conf.Environment, defaultEnv)
    
    conf.LogLevel = normalizeLogLevel(conf.LogLevel)
    
    return nil
}

func normalizeString(value, defaultValue string) string {
    if value == "" {
        return defaultValue
    }
    return value
}

func normalizeLogLevel(level string) string {
    if level == "" {
        return defaultLogLevel
    }

	nLevel := strings.ToLower(level)
    
    switch nLevel {
    case "debug", "info", "warn", "error", "fatal":
        return nLevel
    default:
        return defaultLogLevel
    }
}