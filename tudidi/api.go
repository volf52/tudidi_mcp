package tudidi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"tudidi_mcp/auth"
)

type API struct {
	client   *auth.Client
	readonly bool
}

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Status string

const (
	NotStarted Status = "not_started"
	InProgress Status = "in_progress"
	Completed  Status = "completed"
)

type Tag struct{}

type Task struct {
	ID           int    `json:"id"`
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Note         string `json:"note,omitempty"`
	DueDate      string `json:"due_date,omitempty"`
	Today        bool   `json:"today,omitempty"`
	Priority     int    `json:"priority"`
	Status       int    `json:"status,omitempty"`
	ProjectID    int    `json:"project_id,omitempty"`
	UserID       int    `json:"user_id,omitempty"`
	CompletedAt  string `json:"completed_at,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	Tags         []Tag  `json:"tags"`
	ParentTaskID int    `json:"parent_task_id,omitempty"`
}

type Project struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description,omitempty"`
	Active            bool     `json:"active,omitempty"`
	PinToSidebar      bool     `json:"pin_to_sidebar,omitempty"`
	Priority          Priority `json:"priority,omitempty"`
	DueDateAt         string   `json:"due_date_at,omitempty"`
	UserID            int      `json:"user_id,omitempty"`
	AreaID            int      `json:"area_id,omitempty"`
	TaskShowCompleted bool     `json:"task_show_completed,omitempty"`
	TaskSortOrder     string   `json:"task_sort_order,omitempty"`
	CreatedAt         string   `json:"created_at,omitempty"`
	UpdatedAt         string   `json:"updated_at,omitempty"`
}

type Projects struct {
	Projects []Project `json:"projects"`
}

type CreateTaskRequest struct {
	Name      string `json:"name"`
	Note      string `json:"note,omitempty"`
	ProjectID int    `json:"project_id"`
	Status    Status `json:"status"`
}

type UpdateTaskRequest struct {
	Name string `json:"name,omitempty"`
	Note string `json:"note,omitempty"`
}

func NewAPI(client *auth.Client, readonly bool) *API {
	return &API{
		client:   client,
		readonly: readonly,
	}
}

type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func (api *API) doGet(endpoint string, result interface{}) error {
	resp, err := api.client.Get(endpoint)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return api.handleResponse(resp, result, http.StatusOK)
}

func (api *API) doPost(endpoint string, payload interface{}, result interface{}) error {
	return api.doMutatingRequest("POST", endpoint, payload, result, http.StatusCreated)
}

func (api *API) doPatch(endpoint string, payload interface{}, result interface{}) error {
	return api.doMutatingRequest("PATCH", endpoint, payload, result, http.StatusOK)
}

func (api *API) doDelete(endpoint string) error {
	if api.readonly {
		return fmt.Errorf("operation not allowed in readonly mode")
	}

	resp, err := api.client.Delete(endpoint)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return api.handleResponse(resp, nil, http.StatusOK, http.StatusNoContent)
}

func (api *API) doMutatingRequest(method, endpoint string, payload interface{}, result interface{}, expectedStatus int) error {
	if api.readonly {
		return fmt.Errorf("operation not allowed in readonly mode")
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	var resp *http.Response
	switch method {
	case "POST":
		resp, err = api.client.Post(endpoint, "application/json", jsonData)
	case "PATCH":
		resp, err = api.client.Patch(endpoint, "application/json", jsonData)
	default:
		return fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return api.handleResponse(resp, result, expectedStatus)
}

func (api *API) handleResponse(resp *http.Response, result interface{}, expectedStatuses ...int) error {
	statusOK := false
	for _, status := range expectedStatuses {
		if resp.StatusCode == status {
			statusOK = true
			break
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found")
	}

	if !statusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	if result == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

func (api *API) GetTasks() ([]Task, error) {
	var resp GetTasksResponse
	if err := api.doGet("/api/tasks", &resp); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return resp.Tasks, nil
}

func (api *API) GetTask(id int) (*Task, error) {
	var task Task
	endpoint := "/api/task/" + strconv.Itoa(id)
	if err := api.doGet(endpoint, &task); err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

func (api *API) CreateTask(req CreateTaskRequest) (*Task, error) {
	var task Task
	if err := api.doPost("/api/task", req, &task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return &task, nil
}

func (api *API) UpdateTask(id int, req UpdateTaskRequest) (*Task, error) {
	if req.Name == "" && req.Note == "" {
		return nil, fmt.Errorf("no fields to update")
	}

	currentTask, err := api.GetTask(id)
	if err != nil {
		return nil, fmt.Errorf("task with id %d not found: %w", id, err)
	}

	currentTask.Name = req.Name
	currentTask.Note = req.Note

	var updatedTask Task
	endpoint := "/api/task/" + strconv.Itoa(id)
	if err := api.doPatch(endpoint, currentTask, &updatedTask); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return &updatedTask, nil
}

func (api *API) DeleteTask(id int) error {
	endpoint := "/api/task/" + strconv.Itoa(id)
	if err := api.doDelete(endpoint); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

type GetProjectsResponse struct {
	Projects []Project `json:"projects"`
}

func (api *API) GetProjects() ([]Project, error) {
	var resp GetProjectsResponse
	if err := api.doGet("/api/projects", &resp); err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	return resp.Projects, nil
}

func (api *API) SearchProjectsByName(name string) ([]Project, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	var resp GetProjectsResponse
	if err := api.doGet("/api/projects", &resp); err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	var filtered []Project
	searchLower := strings.ToLower(name)

	for _, project := range resp.Projects {
		if strings.Contains(strings.ToLower(project.Name), searchLower) {
			filtered = append(filtered, project)
		}
	}

	return filtered, nil
}
