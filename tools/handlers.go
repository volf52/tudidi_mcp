package tools

import (
	"context"
	"fmt"
	"tudidi_mcp/tudidi"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Handlers struct {
	api *tudidi.API
}

func NewHandlers(api *tudidi.API) *Handlers {
	return &Handlers{api: api}
}

func (h *Handlers) RegisterTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_tasks",
		Description: "List all tasks",
	}, h.listTasks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_task",
		Description: "Get a specific task by ID",
	}, h.getTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_task",
		Description: "Create a new task",
	}, h.createTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_task",
		Description: "Update an existing task",
	}, h.updateTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_task",
		Description: "Delete a task",
	}, h.deleteTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_task_lists",
		Description: "List all task lists",
	}, h.listTaskLists)
}

type TaskIDArgs struct {
	ID int `json:"id" jsonschema:"Task ID"`
}

type CreateTaskArgs struct {
	Title       string `json:"title" jsonschema:"Task title"`
	Description string `json:"description,omitempty" jsonschema:"Task description"`
	ListID      int    `json:"list_id,omitempty" jsonschema:"List ID to assign task to"`
}

type UpdateTaskArgs struct {
	ID          int    `json:"id" jsonschema:"Task ID"`
	Title       string `json:"title,omitempty" jsonschema:"New task title"`
	Description string `json:"description,omitempty" jsonschema:"New task description"`
	Completed   *bool  `json:"completed,omitempty" jsonschema:"Task completion status"`
}

func (h *Handlers) listTasks(ctx context.Context, req *mcp.CallToolRequest, args any) (*mcp.CallToolResult, any, error) {
	tasks, err := h.api.GetTasks()
	if err != nil {
		return nil, nil, err
	}

	result := map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Found %d tasks", len(tasks))},
		},
	}, result, nil
}

func (h *Handlers) getTask(ctx context.Context, req *mcp.CallToolRequest, args TaskIDArgs) (*mcp.CallToolResult, any, error) {
	task, err := h.api.GetTask(args.ID)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Task: %s", task.Title)},
		},
	}, task, nil
}

func (h *Handlers) createTask(ctx context.Context, req *mcp.CallToolRequest, args CreateTaskArgs) (*mcp.CallToolResult, any, error) {
	createReq := tudidi.CreateTaskRequest{
		Title:       args.Title,
		Description: args.Description,
		ListID:      args.ListID,
	}

	task, err := h.api.CreateTask(createReq)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Created task: %s", task.Title)},
		},
	}, task, nil
}

func (h *Handlers) updateTask(ctx context.Context, req *mcp.CallToolRequest, args UpdateTaskArgs) (*mcp.CallToolResult, any, error) {
	updateReq := tudidi.UpdateTaskRequest{
		Title:       args.Title,
		Description: args.Description,
		Completed:   args.Completed,
	}

	task, err := h.api.UpdateTask(args.ID, updateReq)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Updated task: %s", task.Title)},
		},
	}, task, nil
}

func (h *Handlers) deleteTask(ctx context.Context, req *mcp.CallToolRequest, args TaskIDArgs) (*mcp.CallToolResult, any, error) {
	err := h.api.DeleteTask(args.ID)
	if err != nil {
		return nil, nil, err
	}

	result := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Task %d deleted successfully", args.ID),
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Deleted task %d", args.ID)},
		},
	}, result, nil
}

func (h *Handlers) listTaskLists(ctx context.Context, req *mcp.CallToolRequest, args any) (*mcp.CallToolResult, any, error) {
	lists, err := h.api.GetLists()
	if err != nil {
		return nil, nil, err
	}

	result := map[string]interface{}{
		"lists": lists,
		"count": len(lists),
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Found %d task lists", len(lists))},
		},
	}, result, nil
}
