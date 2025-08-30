package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewClient(baseURL string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Jar: jar,
		},
		baseURL: baseURL,
	}, nil
}

func (c *Client) Login(email, password string) error {
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	loginURL := c.baseURL + "/api/login"
	resp, err := c.httpClient.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	fullURL := c.baseURL + endpoint
	return c.httpClient.Get(fullURL)
}

func (c *Client) Post(endpoint string, contentType string, body []byte) (*http.Response, error) {
	fullURL := c.baseURL + endpoint
	return c.httpClient.Post(fullURL, contentType, bytes.NewBuffer(body))
}

func (c *Client) Put(endpoint string, contentType string, body []byte) (*http.Response, error) {
	fullURL := c.baseURL + endpoint
	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.httpClient.Do(req)
}

func (c *Client) Delete(endpoint string) (*http.Response, error) {
	fullURL := c.baseURL + endpoint
	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}
