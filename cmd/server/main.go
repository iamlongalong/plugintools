package main

import (
	"flag"
	"fmt"
	"log"

	"gay/plugintools/internal/config"
	"gay/plugintools/internal/core"
	"gay/plugintools/internal/server"
	"gay/plugintools/internal/tools"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create tool registry
	registry := core.NewRegistry()

	// Register tools
	if err := registerTools(registry); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// Create and start server
	srv := server.NewServer(registry)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := srv.Start(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// registerTools 注册所有工具
func registerTools(registry core.ToolRegistry) error {
	tools := []core.Tool{
		tools.NewFileManager(),
		tools.NewShellExecutor(),
		tools.NewScheduler(),
	}

	for _, tool := range tools {
		if err := registry.Register(tool); err != nil {
			return err
		}
		log.Printf("Registered tool: %s", tool.GetInfo().Name)
	}

	return nil
}
