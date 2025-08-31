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

	var respData GetTasksResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return respData.Tasks, nil
}

func (api *API) GetTask(id int) (*Task, error) {
	resp, err := api.client.Get("/api/task/" + strconv.Itoa(id))
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

	resp, err := api.client.Post("/api/task", "application/json", jsonData)
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

	if req.Name == "" && req.Note == "" {
		return nil, fmt.Errorf("no fields to update")
	}

	currentTask, err := api.GetTask(id)
	if err != nil {
		return nil, fmt.Errorf("task with id %d not found: %w", id, err)
	}

	currentTask.Name = req.Name
	currentTask.Note = req.Note

	jsonData, err := json.Marshal(currentTask)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := api.client.Patch("/api/task/"+strconv.Itoa(id), "application/json", jsonData)
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

	resp, err := api.client.Delete("/api/task/" + strconv.Itoa(id))
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

type GetProjectsResponse struct {
	Projects []Project `json:"projects"`
}

func (api *API) GetProjects() ([]Project, error) {
	resp, err := api.client.Get("/api/projects")
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get lists: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var projectsResponse GetProjectsResponse
	if err := json.Unmarshal(body, &projectsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return projectsResponse.Projects, nil
}
