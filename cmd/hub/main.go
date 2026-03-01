package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PodPloy/podploy/internal/adapter/api"
	"github.com/PodPloy/podploy/internal/adapter/config"
	"github.com/PodPloy/podploy/internal/adapter/logger"
	"github.com/PodPloy/podploy/internal/domain/ports"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var cfgPath string
	flag.StringVar(&cfgPath, "config", "/etc/podploy/hub.toml", "path config file podploy hub")
	flag.Parse()

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
		os.Exit(1)
	}

	server, err := api.New(&api.Config{
		Host:    serverConf.Host(),
		Port:    serverConf.Port(),
		Origins: serverConf.Origins(),
	}, log)
	if err != nil {
		log.Fatal("Error instance server HTTP ", ports.Error(err))
		return
	}

	err = server.Start(ctx)
	if err != nil {
		log.Fatal("Server Crashed for ", ports.Error(err))
		return
	}

	if err := log.Sync(); err != nil {
		fmt.Printf("Error Sync logger: %v\n", err)
	}

	log.Info("Safe Server Close Successfully")
}
