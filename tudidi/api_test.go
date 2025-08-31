package tudidi

import (
	"fmt"
	"os"
	"testing"
	"time"
	"tudidi_mcp/auth"
)

// Test configuration - requires environment variables to be set
var (
	testURL      = os.Getenv("TUDIDI_TEST_URL")      // e.g., "http://localhost:3002"
	testEmail    = os.Getenv("TUDIDI_TEST_EMAIL")    // e.g., "test@example.com"
	testPassword = os.Getenv("TUDIDI_TEST_PASSWORD") // e.g., "password"
)

func setupTestAPI(t *testing.T, readonly bool) *API {
	if testURL == "" || testEmail == "" || testPassword == "" {
		t.Skip("Skipping API tests - set TUDIDI_TEST_URL, TUDIDI_TEST_EMAIL, and TUDIDI_TEST_PASSWORD environment variables")
	}

	client, err := auth.NewClient(testURL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Login(testEmail, testPassword)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	return NewAPI(client, readonly)
}

func TestNewAPI(t *testing.T) {
	if testURL == "" || testEmail == "" || testPassword == "" {
		t.Skip("Skipping API tests - set TUDIDI_TEST_URL, TUDIDI_TEST_EMAIL, and TUDIDI_TEST_PASSWORD environment variables")
	}

	client, err := auth.NewClient(testURL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test readonly API
	api := NewAPI(client, true)
	if api.client != client {
		t.Error("Expected client to be set")
	}
	if !api.readonly {
		t.Error("Expected readonly to be true")
	}

	// Test writable API
	api = NewAPI(client, false)
	if api.readonly {
		t.Error("Expected readonly to be false")
	}
}

func TestGetTasks(t *testing.T) {
	api := setupTestAPI(t, false) // readonly doesn't matter for GET

	tasks, err := api.GetTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

	// Tasks list can be empty, that's okay
	t.Logf("Retrieved %d tasks", len(tasks))

	// If we have tasks, validate the structure
	for i, task := range tasks {
		if task.ID == 0 {
			t.Errorf("Task %d has invalid ID: %d", i, task.ID)
		}
		if task.Name == "" {
			t.Errorf("Task %d has empty name", i)
		}
		// UUID, dates, and other fields can be empty, that's valid
	}
}

func TestGetTasksReadonly(t *testing.T) {
	api := setupTestAPI(t, true)

	tasks, err := api.GetTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks in readonly mode: %v", err)
	}

	t.Logf("Retrieved %d tasks in readonly mode", len(tasks))
}

func TestGetProjects(t *testing.T) {
	api := setupTestAPI(t, false)

	projects, err := api.GetProjects()
	if err != nil {
		t.Fatalf("GetProjects failed: %v", err)
	}

	if projects == nil {
		t.Error("Expected projects slice, got nil")
	}

	t.Logf("Found %d projects", len(projects))
	for i, project := range projects {
		t.Logf("Project %d: %+v", i, project)
	}
}

func TestGetProjectsReadonly(t *testing.T) {
	api := setupTestAPI(t, true)

	projects, err := api.GetProjects()
	if err != nil {
		t.Fatalf("GetProjects failed: %v", err)
	}

	if projects == nil {
		t.Error("Expected projects slice, got nil")
	}

	t.Logf("Found %d projects (readonly)", len(projects))
}

func TestTaskCRUDOperations(t *testing.T) {
	api := setupTestAPI(t, false)

	// First, get projects to ensure we have a valid project ID
	projects, err := api.GetProjects()
	if err != nil {
		t.Fatalf("Failed to get projects: %v", err)
	}

	var projectID int
	if len(projects) > 0 {
		projectID = projects[0].ID
	} else {
		t.Skip("No projects available for testing - cannot create tasks without a project")
	}

	// Test Create Task
	createReq := CreateTaskRequest{
		Name:      fmt.Sprintf("Test Task %d", time.Now().Unix()),
		Note:      "This is a test task created by automated tests",
		ProjectID: projectID,
		Status:    NotStarted,
	}

	createdTask, err := api.CreateTask(createReq)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if createdTask.ID == 0 {
		t.Error("Created task has invalid ID")
	}
	if createdTask.Name != createReq.Name {
		t.Errorf("Expected task name %s, got %s", createReq.Name, createdTask.Name)
	}
	if createdTask.ProjectID != projectID {
		t.Errorf("Expected project ID %d, got %d", projectID, createdTask.ProjectID)
	}

	t.Logf("Created task with ID: %d", createdTask.ID)

	// Test Get Task
	retrievedTask, err := api.GetTask(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if retrievedTask.ID != createdTask.ID {
		t.Errorf("Expected task ID %d, got %d", createdTask.ID, retrievedTask.ID)
	}
	if retrievedTask.Name != createdTask.Name {
		t.Errorf("Expected task name %s, got %s", createdTask.Name, retrievedTask.Name)
	}

	// Test Update Task
	updateReq := UpdateTaskRequest{
		Name: "Updated Test Task",
		Note: "Updated description",
	}

	updatedTask, err := api.UpdateTask(createdTask.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	if updatedTask.ID != createdTask.ID {
		t.Errorf("Expected task ID %d, got %d", createdTask.ID, updatedTask.ID)
	}

	t.Logf("Updated task with ID: %d", updatedTask.ID)

	// Test Delete Task
	err = api.DeleteTask(createdTask.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	t.Logf("Deleted task with ID: %d", createdTask.ID)

	// Verify task is deleted
	_, err = api.GetTask(createdTask.ID)
	if err == nil {
		t.Error("Expected error when getting deleted task")
	}
	if err.Error() != "task not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestGetNonExistentTask(t *testing.T) {
	api := setupTestAPI(t, false)

	// Try to get a task with a very high ID that likely doesn't exist
	nonExistentID := 999999
	_, err := api.GetTask(nonExistentID)
	if err == nil {
		t.Error("Expected error when getting non-existent task")
	}
	if err.Error() != "task not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestReadonlyModeEnforcement(t *testing.T) {
	api := setupTestAPI(t, true)

	// Test Create Task in readonly mode
	createReq := CreateTaskRequest{
		Name:      "Should Not Be Created",
		Note:      "This should fail",
		ProjectID: 1,
		Status:    NotStarted,
	}

	_, err := api.CreateTask(createReq)
	if err == nil {
		t.Error("Expected error when creating task in readonly mode")
	}
	if err.Error() != "operation not allowed in readonly mode" {
		t.Errorf("Expected readonly error, got: %v", err)
	}

	// Test Update Task in readonly mode
	updateReq := UpdateTaskRequest{
		Name: "Should Not Be Updated",
	}

	_, err = api.UpdateTask(1, updateReq)
	if err == nil {
		t.Error("Expected error when updating task in readonly mode")
	}
	if err.Error() != "operation not allowed in readonly mode" {
		t.Errorf("Expected readonly error, got: %v", err)
	}

	// Test Delete Task in readonly mode
	err = api.DeleteTask(1)
	if err == nil {
		t.Error("Expected error when deleting task in readonly mode")
	}
	if err.Error() != "operation not allowed in readonly mode" {
		t.Errorf("Expected readonly error, got: %v", err)
	}
}

func TestUpdateNonExistentTask(t *testing.T) {
	api := setupTestAPI(t, false)

	updateReq := UpdateTaskRequest{
		Name: "Should Not Work",
	}

	// Try to update a task with a very high ID that likely doesn't exist
	nonExistentID := 999999
	_, err := api.UpdateTask(nonExistentID, updateReq)
	if err == nil {
		t.Error("Expected error when updating non-existent task")
	}
	if err.Error() != "task not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestDeleteNonExistentTask(t *testing.T) {
	api := setupTestAPI(t, false)

	// Try to delete a task with a very high ID that likely doesn't exist
	nonExistentID := 999999
	err := api.DeleteTask(nonExistentID)
	if err == nil {
		t.Error("Expected error when deleting non-existent task")
	}
	if err.Error() != "task not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}
}

func TestCreateTaskValidation(t *testing.T) {
	api := setupTestAPI(t, false)

	// Test creating task with invalid/missing project ID
	createReq := CreateTaskRequest{
		Name:      "Test Task",
		Note:      "Test note",
		ProjectID: 999999, // Non-existent project ID
		Status:    NotStarted,
	}

	_, err := api.CreateTask(createReq)
	if err == nil {
		t.Error("Expected error when creating task with invalid project ID")
	}
	// The exact error message will depend on the server implementation
	t.Logf("Expected error for invalid project ID: %v", err)
}

// Benchmark tests for performance
func BenchmarkGetTasks(b *testing.B) {
	if testURL == "" || testEmail == "" || testPassword == "" {
		b.Skip("Skipping benchmark - set TUDIDI_TEST_* environment variables")
	}

	client, err := auth.NewClient(testURL)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	err = client.Login(testEmail, testPassword)
	if err != nil {
		b.Fatalf("Failed to login: %v", err)
	}

	api := NewAPI(client, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := api.GetTasks()
		if err != nil {
			b.Fatalf("Failed to get tasks: %v", err)
		}
	}
}

func BenchmarkGetProjects(b *testing.B) {
	if testURL == "" || testEmail == "" || testPassword == "" {
		b.Skip("Skipping benchmark - set TUDIDI_TEST_URL, TUDIDI_TEST_EMAIL, and TUDIDI_TEST_PASSWORD environment variables")
	}

	api := setupTestAPI(&testing.T{}, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := api.GetProjects()
		if err != nil {
			b.Errorf("GetProjects failed: %v", err)
		}
	}
}
