package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"tudidi_mcp/auth"
	"tudidi_mcp/config"
	"tudidi_mcp/tudidi"
)

func main() {
	fmt.Println("🔧 Tudidi API Testing Playground")
	fmt.Println("==================================")

	// Parse config
	cfg, err := config.ParseArgs()
	if err != nil {
		fmt.Printf("❌ Configuration error: %v\n", err)
		fmt.Println("\nUsage examples:")
		fmt.Println("  ./test-playground --url http://localhost:3002 --email admin@test.com --password secret")
		fmt.Println("  TUDIDI_URL=http://localhost:3002 TUDIDI_USER_EMAIL=admin@test.com TUDIDI_USER_PASSWORD=secret ./test-playground")
		os.Exit(1)
	}

	// Create HTTP client with authentication
	client, err := auth.NewClient(cfg.URL)
	if err != nil {
		log.Fatalf("❌ Failed to create HTTP client: %v", err)
	}

	// Authenticate with Tudidi server
	fmt.Printf("🔐 Authenticating with %s...\n", cfg.URL)
	if err := client.Login(cfg.Email, cfg.Password); err != nil {
		log.Fatalf("❌ Authentication failed: %v", err)
	}
	fmt.Println("✅ Authentication successful!")

	// Create API instance
	api := tudidi.NewAPI(client, cfg.Readonly)

	readonlyStatus := ""
	if cfg.Readonly {
		readonlyStatus = " (READONLY MODE - destructive operations disabled)"
	}
	fmt.Printf("🚀 API ready%s\n\n", readonlyStatus)

	// Interactive mode
	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		fmt.Print("Enter command: ")

		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}

		switch command {
		case "help", "h":
			showHelp()
		case "quit", "q", "exit":
			fmt.Println("👋 Goodbye!")
			return
		case "list-tasks", "lt":
			listTasks(api)
		case "get-task", "gt":
			getTask(api, scanner)
		case "create-task", "ct":
			createTask(api, scanner)
		case "update-task", "ut":
			updateTask(api, scanner)
		case "delete-task", "dt":
			deleteTask(api, scanner)
		case "list-lists", "ll":
			listLists(api)
		case "toggle-readonly", "tr":
			api = toggleReadonly(api, client, !cfg.Readonly)
			cfg.Readonly = !cfg.Readonly
		case "status", "s":
			showStatus(cfg, api)
		case "clear", "c":
			clearScreen()
		default:
			fmt.Printf("❌ Unknown command: %s\n", command)
			fmt.Println("Type 'help' for available commands")
		}

		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("📋 Available Commands:")
	fmt.Println("  list-tasks (lt)     - List all tasks")
	fmt.Println("  get-task (gt)       - Get specific task by ID")
	fmt.Println("  create-task (ct)    - Create a new task")
	fmt.Println("  update-task (ut)    - Update existing task")
	fmt.Println("  delete-task (dt)    - Delete a task")
	fmt.Println("  list-lists (ll)     - List all project lists")
	fmt.Println("  toggle-readonly (tr)- Toggle readonly mode")
	fmt.Println("  status (s)          - Show current status")
	fmt.Println("  clear (c)           - Clear screen")
	fmt.Println("  help (h)            - Show detailed help")
	fmt.Println("  quit (q)            - Exit")
}

func showHelp() {
	fmt.Println("🆘 Detailed Command Help:")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("📋 READ OPERATIONS:")
	fmt.Println("  list-tasks, lt")
	fmt.Println("    Lists all tasks from the server")
	fmt.Println("    Safe to use in readonly mode")
	fmt.Println()
	fmt.Println("  get-task, gt")
	fmt.Println("    Retrieves a specific task by ID")
	fmt.Println("    Will prompt for task ID")
	fmt.Println("    Safe to use in readonly mode")
	fmt.Println()
	fmt.Println("  list-lists, ll")
	fmt.Println("    Lists all project lists/containers")
	fmt.Println("    Safe to use in readonly mode")
	fmt.Println()
	fmt.Println("✍️  WRITE OPERATIONS (disabled in readonly mode):")
	fmt.Println("  create-task, ct")
	fmt.Println("    Creates a new task")
	fmt.Println("    Will prompt for task name, description, and project ID")
	fmt.Println()
	fmt.Println("  update-task, ut")
	fmt.Println("    Updates an existing task")
	fmt.Println("    Will prompt for task ID and new values")
	fmt.Println()
	fmt.Println("  delete-task, dt")
	fmt.Println("    Deletes a task by ID")
	fmt.Println("    Will prompt for confirmation")
	fmt.Println()
	fmt.Println("⚙️  UTILITY COMMANDS:")
	fmt.Println("  toggle-readonly, tr")
	fmt.Println("    Switches between readonly and writable modes")
	fmt.Println()
	fmt.Println("  status, s")
	fmt.Println("    Shows current connection and mode status")
	fmt.Println()
	fmt.Println("  clear, c")
	fmt.Println("    Clears the screen")
}

func listTasks(api *tudidi.API) {
	fmt.Println("📋 Fetching tasks...")
	tasks, err := api.GetTasks()
	if err != nil {
		fmt.Printf("❌ Error fetching tasks: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("📝 No tasks found")
		return
	}

	fmt.Printf("✅ Found %d tasks:\n", len(tasks))
	fmt.Println("ID   | Name                     | Status | Project ID | Created")
	fmt.Println("-----|--------------------------|--------|------------|--------")

	for _, task := range tasks {
		status := getStatusText(task.Status)
		createdAt := formatDate(task.CreatedAt)
		name := truncateString(task.Name, 24)
		fmt.Printf("%-4d | %-24s | %-6s | %-10d | %s\n",
			task.ID, name, status, task.ProjectID, createdAt)
	}
}

func getTask(api *tudidi.API, scanner *bufio.Scanner) {
	fmt.Print("Enter task ID: ")
	if !scanner.Scan() {
		return
	}

	idStr := strings.TrimSpace(scanner.Text())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("❌ Invalid task ID: %s\n", idStr)
		return
	}

	fmt.Printf("🔍 Fetching task %d...\n", id)
	task, err := api.GetTask(id)
	if err != nil {
		fmt.Printf("❌ Error fetching task: %v\n", err)
		return
	}

	fmt.Println("✅ Task details:")
	fmt.Printf("  ID:          %d\n", task.ID)
	fmt.Printf("  UUID:        %s\n", task.UUID)
	fmt.Printf("  Name:        %s\n", task.Name)
	fmt.Printf("  Note:        %s\n", task.Note)
	fmt.Printf("  Due Date:    %s\n", task.DueDate)
	fmt.Printf("  Today:       %t\n", task.Today)
	fmt.Printf("  Priority:    %d\n", task.Priority)
	fmt.Printf("  Status:      %s (%d)\n", getStatusText(task.Status), task.Status)
	fmt.Printf("  Project ID:  %d\n", task.ProjectID)
	fmt.Printf("  User ID:     %d\n", task.UserID)
	fmt.Printf("  Completed:   %s\n", task.CompletedAt)
	fmt.Printf("  Created:     %s\n", task.CreatedAt)
	fmt.Printf("  Updated:     %s\n", task.UpdatedAt)
	fmt.Printf("  Parent Task: %d\n", task.ParentTaskID)
}

func createTask(api *tudidi.API, scanner *bufio.Scanner) {
	fmt.Print("Enter task name: ")
	if !scanner.Scan() {
		return
	}
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Println("❌ Task name cannot be empty")
		return
	}

	fmt.Print("Enter task description (optional): ")
	if !scanner.Scan() {
		return
	}
	note := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter project ID (or press Enter for default): ")
	if !scanner.Scan() {
		return
	}
	projectIDStr := strings.TrimSpace(scanner.Text())

	var projectID int
	if projectIDStr != "" {
		var err error
		projectID, err = strconv.Atoi(projectIDStr)
		if err != nil {
			fmt.Printf("❌ Invalid project ID: %s\n", projectIDStr)
			return
		}
	} else {
		// Try to get first available project
		lists, err := api.GetLists()
		if err != nil || len(lists) == 0 {
			fmt.Println("❌ No projects available and no project ID specified")
			return
		}
		projectID = lists[0].ID
		fmt.Printf("ℹ️  Using project ID %d (%s)\n", projectID, lists[0].Name)
	}

	req := tudidi.CreateTaskRequest{
		Name:      name,
		Note:      note,
		ProjectID: projectID,
		Status:    tudidi.NotStarted,
	}

	fmt.Println("🔨 Creating task...")
	task, err := api.CreateTask(req)
	if err != nil {
		fmt.Printf("❌ Error creating task: %v\n", err)
		return
	}

	fmt.Printf("✅ Task created successfully!\n")
	fmt.Printf("  ID:   %d\n", task.ID)
	fmt.Printf("  Name: %s\n", task.Name)
	fmt.Printf("  Note: %s\n", task.Note)
}

func updateTask(api *tudidi.API, scanner *bufio.Scanner) {
	fmt.Print("Enter task ID to update: ")
	if !scanner.Scan() {
		return
	}

	idStr := strings.TrimSpace(scanner.Text())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("❌ Invalid task ID: %s\n", idStr)
		return
	}

	fmt.Print("Enter new title (or press Enter to skip): ")
	if !scanner.Scan() {
		return
	}
	name := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter new description (or press Enter to skip): ")
	if !scanner.Scan() {
		return
	}
	description := strings.TrimSpace(scanner.Text())

	// fmt.Print("Mark as completed? (y/N): ")
	// if !scanner.Scan() {
	// 	return
	// }
	// completedStr := strings.TrimSpace(strings.ToLower(scanner.Text()))
	//
	// switch completedStr {
	// case "y", "yes":
	// 	val := true
	// 	completed = &val
	// case "n", "no":
	// 	val := false
	// 	completed = &val
	// }

	req := tudidi.UpdateTaskRequest{}
	if name != "" {
		req.Name = name
	}
	if description != "" {
		req.Note = description
	}

	if req.Name == "" && req.Note == "" {
		fmt.Println("❌ No updates specified")
		return
	}

	fmt.Printf("🔄 Updating task %d...\n", id)
	task, err := api.UpdateTask(id, req)
	if err != nil {
		fmt.Printf("❌ Error updating task: %v\n", err)
		return
	}

	fmt.Printf("✅ Task updated successfully!\n")
	fmt.Printf("  ID:   %d\n", task.ID)
	fmt.Printf("  Name: %s\n", task.Name)
	fmt.Printf("  Completed: %t\n", task.CompletedAt != "")
}

func deleteTask(api *tudidi.API, scanner *bufio.Scanner) {
	fmt.Print("Enter task ID to delete: ")
	if !scanner.Scan() {
		return
	}

	idStr := strings.TrimSpace(scanner.Text())
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("❌ Invalid task ID: %s\n", idStr)
		return
	}

	// Get task details first
	task, err := api.GetTask(id)
	if err != nil {
		fmt.Printf("❌ Error fetching task: %v\n", err)
		return
	}

	fmt.Printf("⚠️  About to delete task:\n")
	fmt.Printf("  ID:   %d\n", task.ID)
	fmt.Printf("  Name: %s\n", task.Name)
	fmt.Print("Are you sure? (y/N): ")

	if !scanner.Scan() {
		return
	}

	confirm := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("❌ Deletion cancelled")
		return
	}

	fmt.Printf("🗑️  Deleting task %d...\n", id)
	err = api.DeleteTask(id)
	if err != nil {
		fmt.Printf("❌ Error deleting task: %v\n", err)
		return
	}

	fmt.Println("✅ Task deleted successfully!")
}

func listLists(api *tudidi.API) {
	fmt.Println("📁 Fetching project lists...")
	lists, err := api.GetLists()
	if err != nil {
		fmt.Printf("❌ Error fetching lists: %v\n", err)
		return
	}

	if len(lists) == 0 {
		fmt.Println("📝 No project lists found")
		return
	}

	fmt.Printf("✅ Found %d project lists:\n", len(lists))
	fmt.Println("ID   | Name                     | Active | Description")
	fmt.Println("-----|--------------------------|--------|------------")

	for _, list := range lists {
		active := "No"
		if list.Active {
			active = "Yes"
		}
		name := truncateString(list.Name, 24)
		desc := truncateString(list.Description, 20)
		fmt.Printf("%-4d | %-24s | %-6s | %s\n",
			list.ID, name, active, desc)
	}
}

func toggleReadonly(api *tudidi.API, client *auth.Client, readonly bool) *tudidi.API {
	newAPI := tudidi.NewAPI(client, readonly)
	status := "WRITABLE"
	if readonly {
		status = "READONLY"
	}
	fmt.Printf("🔄 Switched to %s mode\n", status)
	return newAPI
}

func showStatus(cfg *config.Config, api *tudidi.API) {
	fmt.Println("📊 Current Status:")
	fmt.Printf("  Server URL:  %s\n", cfg.URL)
	fmt.Printf("  Email:       %s\n", cfg.Email)
	fmt.Printf("  Transport:   %s", cfg.Transport)
	if cfg.Transport == "sse" {
		fmt.Printf(" (port %d)", cfg.Port)
	}
	fmt.Println()

	readonlyStatus := "WRITABLE (create/update/delete enabled)"
	if cfg.Readonly {
		readonlyStatus = "READONLY (create/update/delete disabled)"
	}
	fmt.Printf("  Mode:        %s\n", readonlyStatus)
	fmt.Printf("  Time:        %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
	fmt.Println("🔧 Tudidi API Testing Playground")
	fmt.Println("==================================")
}

// Helper functions

func getStatusText(status int) string {
	switch status {
	case 0:
		return "New"
	case 1:
		return "Active"
	case 2:
		return "Done"
	default:
		return "Unknown"
	}
}

func formatDate(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}

	// Try to parse and reformat
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t.Format("2006-01-02")
	}

	// If parsing fails, return first 10 characters
	if len(dateStr) > 10 {
		return dateStr[:10]
	}
	return dateStr
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
