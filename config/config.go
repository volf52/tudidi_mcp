package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	URL      string
	Email    string
	Password string
	Readonly bool
}

func ParseArgs() (*Config, error) {
	var config Config

	flag.StringVar(&config.URL, "url", "http://localhost:3002", "Tudidi server URL (required)")
	flag.StringVar(&config.Email, "email", "", "Email for authentication (required)")
	flag.StringVar(&config.Password, "password", "", "Password for authentication (required)")
	flag.BoolVar(&config.Readonly, "readonly", false, "Run in readonly mode (prevents destructive operations)")

	flag.Parse()

	if config.URL == "" {
		return nil, fmt.Errorf("URL is required")
	}
	if config.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	return &config, nil
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s --url <tudidi-url> --email <user> --password <pass> [--readonly]\n", os.Args[0])
	flag.PrintDefaults()
}
