package tudidi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"tudidi_mcp/auth"
)

type API struct {
	client   *auth.Client
	readonly bool
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
	ListID      int    `json:"list_id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type TaskList struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	ListID      int    `json:"list_id,omitempty"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Completed   *bool  `json:"completed,omitempty"`
}

func NewAPI(client *auth.Client, readonly bool) *API {
	return &API{
		client:   client,
		readonly: readonly,
	}
}

func (api *API) GetTasks() ([]Task, error) {
	resp, err := api.client.Get("/api/tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tasks: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return tasks, nil
}

func (api *API) GetTask(id int) (*Task, error) {
	resp, err := api.client.Get("/api/tasks/" + strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get task: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}

	return &task, nil
}

func (api *API) CreateTask(req CreateTaskRequest) (*Task, error) {
	if api.readonly {
		return nil, fmt.Errorf("operation not allowed in readonly mode")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := api.client.Post("/api/tasks", "application/json", jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create task: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}

	return &task, nil
}

func (api *API) UpdateTask(id int, req UpdateTaskRequest) (*Task, error) {
	if api.readonly {
		return nil, fmt.Errorf("operation not allowed in readonly mode")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := api.client.Put("/api/tasks/"+strconv.Itoa(id), "application/json", jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update task: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}

	return &task, nil
}

func (api *API) DeleteTask(id int) error {
	if api.readonly {
		return fmt.Errorf("operation not allowed in readonly mode")
	}

	resp, err := api.client.Delete("/api/tasks/" + strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("task not found")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete task: status %d", resp.StatusCode)
	}

	return nil
}

func (api *API) GetLists() ([]TaskList, error) {
	resp, err := api.client.Get("/api/lists")
	if err != nil {
		return nil, fmt.Errorf("failed to get lists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get lists: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var lists []TaskList
	if err := json.Unmarshal(body, &lists); err != nil {
		return nil, fmt.Errorf("failed to parse lists: %w", err)
	}

	return lists, nil
}
