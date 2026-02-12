package main

import (
	"fmt"

	"github.com/PodPloy/podploy/internal/adapter/config"
)

func main() {
	fmt.Println("Starting server...")
	var cfgPath string = "internal/adapter/config/example.toml"
	cfg, err := config.LoadHubConfig(cfgPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Printf("Config loaded successfully: %v\n", cfg)
}
