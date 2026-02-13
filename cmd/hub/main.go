package main

import (
	"fmt"

	"github.com/PodPloy/podploy/internal/adapter/config"
	"github.com/PodPloy/podploy/internal/adapter/logger"
)

func main() {
	fmt.Println("Starting server...")
	var cfgPath string = "example.toml"
	cfg, err := config.LoadHubConfig(cfgPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	log, err := logger.New(&logger.Config{
		Level:      cfg.LogLevel,
		OutputPath: cfg.OutputPath,
		Development: cfg.Environment != "production",
		MaxSize:    cfg.MaxSize, // MB
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge, // days
	})
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}
	defer log.Sync()

	log.Info("Server configuration loaded ", "host ", cfg.Host, " port ", cfg.Port, " environment ", cfg.Environment)
	log.Info("Logger initialized with level ", cfg.LogLevel, " and output path ", cfg.OutputPath)
	log.Info("Max log file size: ", cfg.MaxSize, " MB, Max backups: ", cfg.MaxBackups, " Max age: ", cfg.MaxAge, " days")
}
