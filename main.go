package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"tudidi_mcp/auth"
	"tudidi_mcp/config"
	"tudidi_mcp/tools"
	"tudidi_mcp/tudidi"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	cfg, err := config.ParseArgs()
	if err != nil {
		config.PrintUsage()
		log.Fatalf("Configuration error: %v", err)
	}

	// Create HTTP client with authentication
	client, err := auth.NewClient(cfg.URL)
	if err != nil {
		log.Fatalf("Failed to create HTTP client: %v", err)
	}

	// Authenticate with Tudidi server
	if err := client.Login(cfg.Email, cfg.Password); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Create Tudidi API instance
	api := tudidi.NewAPI(client, cfg.Readonly)

	// Create MCP server
	opts := &mcp.ServerOptions{
		Instructions: "Tudidi MCP Server for task management",
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "tudidi",
		Version: "1.0.0",
	}, opts)

	// Register tools
	handlers := tools.NewHandlers(api)
	handlers.RegisterTools(server)

	// Log server status
	readonlyStatus := ""
	if cfg.Readonly {
		readonlyStatus = " (readonly mode)"
	}
	log.Printf("Tudidi MCP server connected to %s%s using %s transport", cfg.URL, readonlyStatus, cfg.Transport)

	// Run based on transport type
	switch cfg.Transport {
	case "stdio":
		transport := &mcp.LoggingTransport{
			Transport: &mcp.StdioTransport{},
			Writer:    os.Stderr,
		}
		if err := server.Run(context.Background(), transport); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	case "sse":
		// Create SSE handler
		handler := mcp.NewSSEHandler(func(req *http.Request) *mcp.Server {
			return server
		})

		addr := fmt.Sprintf(":%d", cfg.Port)
		log.Printf("Starting SSE server on %s", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatalf("SSE server failed: %v", err)
		}
	default:
		log.Fatalf("Unsupported transport: %s", cfg.Transport)
	}
}
