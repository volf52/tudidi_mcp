package tudidi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"tudidi_mcp/auth"
)

// Mock HTTP Client for testing helper methods
type mockHTTPClient struct {
	responses map[string]*http.Response
	lastURL   string
	lastBody  []byte
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	m.lastURL = url
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader("Not Found")),
	}, nil
}

func (m *mockHTTPClient) Post(url, contentType string, body []byte) (*http.Response, error) {
	m.lastURL = url
	m.lastBody = body
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(strings.NewReader("{}")),
	}, nil
}

func (m *mockHTTPClient) Patch(url, contentType string, body []byte) (*http.Response, error) {
	m.lastURL = url
	m.lastBody = body
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("{}")),
	}, nil
}

func (m *mockHTTPClient) Delete(url string) (*http.Response, error) {
	m.lastURL = url
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func createMockAPI(readonly bool, responses map[string]*http.Response) *API {
	authClient := &auth.Client{} // We'll use reflection or expose the HTTP client interface

	// For testing purposes, we'll create a simplified API
	api := &API{
		client:   authClient,
		readonly: readonly,
	}

	// Note: mockClient would be used if we exposed the HTTP client interface
	// For now, we test the helper methods directly
	return api
}

func TestHandleResponse(t *testing.T) {
	api := &API{readonly: false}

	tests := []struct {
		name           string
		statusCode     int
		body           string
		expectedStatus []int
		result         interface{}
		expectError    bool
		errorContains  string
	}{
		{
			name:           "Success with JSON response",
			statusCode:     http.StatusOK,
			body:           `{"name": "test", "id": 1}`,
			expectedStatus: []int{http.StatusOK},
			result:         &Task{},
			expectError:    false,
		},
		{
			name:           "Not Found",
			statusCode:     http.StatusNotFound,
			body:           "",
			expectedStatus: []int{http.StatusOK},
			result:         &Task{},
			expectError:    true,
			errorContains:  "resource not found",
		},
		{
			name:           "Unexpected status",
			statusCode:     http.StatusInternalServerError,
			body:           "",
			expectedStatus: []int{http.StatusOK},
			result:         &Task{},
			expectError:    true,
			errorContains:  "unexpected status: 500",
		},
		{
			name:           "Success with no result expected",
			statusCode:     http.StatusNoContent,
			body:           "",
			expectedStatus: []int{http.StatusOK, http.StatusNoContent},
			result:         nil,
			expectError:    false,
		},
		{
			name:           "Invalid JSON",
			statusCode:     http.StatusOK,
			body:           `{"invalid": json}`,
			expectedStatus: []int{http.StatusOK},
			result:         &Task{},
			expectError:    true,
			errorContains:  "failed to parse response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}

			err := api.handleResponse(resp, tt.result, tt.expectedStatus...)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDoMutatingRequest_ReadonlyMode(t *testing.T) {
	api := &API{readonly: true}

	tests := []struct {
		name   string
		method string
	}{
		{"POST request", "POST"},
		{"PATCH request", "PATCH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.doMutatingRequest(tt.method, "/test", map[string]string{"test": "data"}, nil, http.StatusOK)

			if err == nil {
				t.Error("Expected readonly error, got nil")
			}

			if !strings.Contains(err.Error(), "readonly mode") {
				t.Errorf("Expected readonly error, got: %v", err)
			}
		})
	}
}

func TestDoMutatingRequest_UnsupportedMethod(t *testing.T) {
	api := &API{readonly: false}

	err := api.doMutatingRequest("PUT", "/test", nil, nil, http.StatusOK)

	if err == nil {
		t.Error("Expected unsupported method error, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported method") {
		t.Errorf("Expected unsupported method error, got: %v", err)
	}
}

func TestDoDelete_ReadonlyMode(t *testing.T) {
	api := &API{readonly: true}

	err := api.doDelete("/test")

	if err == nil {
		t.Error("Expected readonly error, got nil")
	}

	if !strings.Contains(err.Error(), "readonly mode") {
		t.Errorf("Expected readonly error, got: %v", err)
	}
}

// Test JSON marshaling errors
func TestDoMutatingRequest_MarshalError(t *testing.T) {
	api := &API{readonly: false}

	// Create a payload that can't be marshaled (channel)
	invalidPayload := make(chan int)

	err := api.doMutatingRequest("POST", "/test", invalidPayload, nil, http.StatusCreated)

	if err == nil {
		t.Error("Expected marshal error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to marshal request") {
		t.Errorf("Expected marshal error, got: %v", err)
	}
}

// Integration test with mock to test full flow
func TestHelperMethodIntegration(t *testing.T) {
	// This test would require exposing the HTTP client interface
	// For now, we'll keep the existing integration tests in api_test.go
	// which test the full stack with a real server

	t.Skip("Integration testing is handled by existing api_test.go with real server")
}

// Benchmark the helper methods
func BenchmarkHandleResponse(b *testing.B) {
	api := &API{readonly: false}
	task := &Task{}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"name": "test task", "id": 1, "note": "test note"}`)),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset body for each iteration
		resp.Body = io.NopCloser(strings.NewReader(`{"name": "test task", "id": 1, "note": "test note"}`))
		err := api.handleResponse(resp, task, http.StatusOK)
		if err != nil {
			b.Fatalf("handleResponse failed: %v", err)
		}
	}
}

// Test helper to create test responses
func createJSONResponse(statusCode int, data interface{}) *http.Response {
	jsonData, _ := json.Marshal(data)
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(jsonData)),
	}
}

func TestSearchProjectsByName(t *testing.T) {
	api := &API{readonly: false}

	tests := []struct {
		name          string
		searchName    string
		expectError   bool
		errorContains string
	}{
		{
			name:          "Empty search name",
			searchName:    "",
			expectError:   true,
			errorContains: "name cannot be empty",
		},
		{
			name:        "Valid search name",
			searchName:  "test",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test would require a mock HTTP client for complete testing
			// For now, we test the validation logic
			if tt.searchName == "" {
				_, err := api.SearchProjectsByName(tt.searchName)
				if !tt.expectError {
					t.Errorf("Expected no error, got %v", err)
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
			}
		})
	}
}

func TestCreateJSONResponse_Helper(t *testing.T) {
	task := Task{ID: 1, Name: "Test Task"}
	resp := createJSONResponse(http.StatusOK, task)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var parsedTask Task
	err = json.Unmarshal(body, &parsedTask)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsedTask.ID != task.ID || parsedTask.Name != task.Name {
		t.Errorf("Expected task %+v, got %+v", task, parsedTask)
	}
}
