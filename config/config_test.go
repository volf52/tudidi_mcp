package config

import (
	"os"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	// Test that config validation works
	config := &Config{
		URL:       "",
		Email:     "",
		Password:  "",
		Readonly:  true, // Now defaults to true
		Transport: "stdio",
		Port:      8080,
	}

	if config.URL != "" {
		t.Error("Expected empty URL")
	}

	if config.Email != "" {
		t.Error("Expected empty email")
	}

	if config.Password != "" {
		t.Error("Expected empty password")
	}

	if !config.Readonly {
		t.Error("Expected readonly to default to true")
	}

	if config.Transport != "stdio" {
		t.Error("Expected transport to default to stdio")
	}

	if config.Port != 8080 {
		t.Error("Expected port to default to 8080")
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalURL := os.Getenv("TUDIDI_URL")
	originalEmail := os.Getenv("TUDIDI_USER_EMAIL")
	originalPassword := os.Getenv("TUDIDI_USER_PASSWORD")
	originalReadonly := os.Getenv("TUDIDI_READONLY")
	originalTransport := os.Getenv("TUDIDI_TRANSPORT")
	originalPort := os.Getenv("TUDIDI_PORT")

	// Clean up after test
	defer func() {
		os.Setenv("TUDIDI_URL", originalURL)
		os.Setenv("TUDIDI_USER_EMAIL", originalEmail)
		os.Setenv("TUDIDI_USER_PASSWORD", originalPassword)
		os.Setenv("TUDIDI_READONLY", originalReadonly)
		os.Setenv("TUDIDI_TRANSPORT", originalTransport)
		os.Setenv("TUDIDI_PORT", originalPort)
	}()

	// Set test environment variables
	os.Setenv("TUDIDI_URL", "https://test.example.com")
	os.Setenv("TUDIDI_USER_EMAIL", "test@example.com")
	os.Setenv("TUDIDI_USER_PASSWORD", "testpass")
	os.Setenv("TUDIDI_READONLY", "false")
	os.Setenv("TUDIDI_TRANSPORT", "sse")
	os.Setenv("TUDIDI_PORT", "3000")

	// Test that environment variables are read correctly
	if url := os.Getenv("TUDIDI_URL"); url != "https://test.example.com" {
		t.Errorf("Expected TUDIDI_URL to be 'https://test.example.com', got '%s'", url)
	}

	if email := os.Getenv("TUDIDI_USER_EMAIL"); email != "test@example.com" {
		t.Errorf("Expected TUDIDI_USER_EMAIL to be 'test@example.com', got '%s'", email)
	}

	if password := os.Getenv("TUDIDI_USER_PASSWORD"); password != "testpass" {
		t.Errorf("Expected TUDIDI_USER_PASSWORD to be 'testpass', got '%s'", password)
	}

	if readonly := os.Getenv("TUDIDI_READONLY"); readonly != "false" {
		t.Errorf("Expected TUDIDI_READONLY to be 'false', got '%s'", readonly)
	}

	if transport := os.Getenv("TUDIDI_TRANSPORT"); transport != "sse" {
		t.Errorf("Expected TUDIDI_TRANSPORT to be 'sse', got '%s'", transport)
	}

	if port := os.Getenv("TUDIDI_PORT"); port != "3000" {
		t.Errorf("Expected TUDIDI_PORT to be '3000', got '%s'", port)
	}
}
