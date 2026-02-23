package main

import (
	"fmt"

	"github.com/PodPloy/podploy/internal/adapter/config"
	"github.com/PodPloy/podploy/internal/adapter/http"
	"github.com/PodPloy/podploy/internal/adapter/logger"
)

func main() {
	fmt.Println("Starting server...")
	cfgPath := "hub.example.toml"
	cfg, err := config.LoadHubConfig(cfgPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	serverConf := cfg.Server()
	logConf := cfg.Logger()

	log, err := logger.New(&logger.Config{
		Level:       logConf.Level(),
		OutputPath:  logConf.OutputPath(),
		Development: logConf.Env() == "development",
		MaxSize:     logConf.MaxSize(),
		MaxBackups:  logConf.MaxBackups(),
		MaxAge:      logConf.MaxAge(),
	})
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	server, err := http.New(&http.Config{
		Host:    serverConf.Host(),
		Port:    serverConf.Port(),
		Origins: serverConf.Origins(),
	}, log)
	if err != nil {
		log.WithError(err)
		return
	}

	err = server.Start()
	if err != nil {
		log.WithError(err)
		return
	}
}
