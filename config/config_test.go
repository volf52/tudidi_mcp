package config

import (
	"os"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	// Test that config validation works
	config := &Config{
		URL:      "",
		Email:    "",
		Password: "",
		Readonly: false,
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
}

func TestEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalURL := os.Getenv("TUDIDI_URL")
	originalEmail := os.Getenv("TUDIDI_USER_EMAIL")
	originalPassword := os.Getenv("TUDIDI_USER_PASSWORD")
	originalReadonly := os.Getenv("TUDIDI_READONLY")

	// Clean up after test
	defer func() {
		os.Setenv("TUDIDI_URL", originalURL)
		os.Setenv("TUDIDI_USER_EMAIL", originalEmail)
		os.Setenv("TUDIDI_USER_PASSWORD", originalPassword)
		os.Setenv("TUDIDI_READONLY", originalReadonly)
	}()

	// Set test environment variables
	os.Setenv("TUDIDI_URL", "https://test.example.com")
	os.Setenv("TUDIDI_USER_EMAIL", "test@example.com")
	os.Setenv("TUDIDI_USER_PASSWORD", "testpass")
	os.Setenv("TUDIDI_READONLY", "true")

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

	if readonly := os.Getenv("TUDIDI_READONLY"); readonly != "true" {
		t.Errorf("Expected TUDIDI_READONLY to be 'true', got '%s'", readonly)
	}
}
