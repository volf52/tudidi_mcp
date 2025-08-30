package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	URL       string
	Email     string
	Password  string
	Readonly  bool
	Transport string
	Port      int
}

func ParseArgs() (*Config, error) {
	var config Config

	flag.StringVar(&config.URL, "url", "http://localhost:3002", "Tudidi server URL (required)")
	flag.StringVar(&config.Email, "email", "", "Email for authentication (required)")
	flag.StringVar(&config.Password, "password", "", "Password for authentication (required)")
	flag.BoolVar(&config.Readonly, "readonly", true, "Run in readonly mode (prevents destructive operations)")
	flag.StringVar(&config.Transport, "transport", "stdio", "Transport type: 'stdio' or 'sse'")
	flag.IntVar(&config.Port, "port", 8080, "Port for SSE transport (ignored for stdio)")

	flag.Parse()

	// Override with environment variables if available
	if envURL := os.Getenv("TUDIDI_URL"); envURL != "" {
		config.URL = envURL
	}
	if envEmail := os.Getenv("TUDIDI_USER_EMAIL"); envEmail != "" {
		config.Email = envEmail
	}
	if envPassword := os.Getenv("TUDIDI_USER_PASSWORD"); envPassword != "" {
		config.Password = envPassword
	}
	if envReadonly := os.Getenv("TUDIDI_READONLY"); envReadonly != "" {
		config.Readonly = envReadonly == "true"
	}
	if envTransport := os.Getenv("TUDIDI_TRANSPORT"); envTransport != "" {
		config.Transport = envTransport
	}
	if envPort := os.Getenv("TUDIDI_PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil {
			config.Port = port
		}
	}

	if config.URL == "" {
		return nil, fmt.Errorf("URL is required (use --url flag or TUDIDI_URL environment variable)")
	}
	if config.Email == "" {
		return nil, fmt.Errorf("email is required (use --email flag or TUDIDI_USER_EMAIL environment variable)")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("password is required (use --password flag or TUDIDI_USER_PASSWORD environment variable)")
	}
	if config.Transport != "stdio" && config.Transport != "sse" {
		return nil, fmt.Errorf("transport must be 'stdio' or 'sse', got: %s", config.Transport)
	}
	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535, got: %d", config.Port)
	}

	return &config, nil
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s --url <tudidi-url> --email <user> --password <pass> [--readonly] [--transport <stdio|sse>] [--port <port>]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
	fmt.Fprintf(os.Stderr, "  TUDIDI_URL          Tudidi server URL\n")
	fmt.Fprintf(os.Stderr, "  TUDIDI_USER_EMAIL   Email for authentication\n")
	fmt.Fprintf(os.Stderr, "  TUDIDI_USER_PASSWORD Password for authentication\n")
	fmt.Fprintf(os.Stderr, "  TUDIDI_READONLY     Set to 'true' or 'false' for readonly mode (default: true)\n")
	fmt.Fprintf(os.Stderr, "  TUDIDI_TRANSPORT    Transport type: 'stdio' or 'sse' (default: stdio)\n")
	fmt.Fprintf(os.Stderr, "  TUDIDI_PORT         Port for SSE transport (default: 8080)\n")
	fmt.Fprintf(os.Stderr, "\nCommand Line Flags:\n")
	flag.PrintDefaults()
}
