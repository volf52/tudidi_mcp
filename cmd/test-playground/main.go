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

// Command represents a playground command
type Command struct {
	Name        string
	Alias       string
	Description string
	Handler     func(*PlaygroundContext, *bufio.Scanner)
	RequiresAPI bool
}

// PlaygroundContext holds the shared state
type PlaygroundContext struct {
	API    *tudidi.API
	Config *config.Config
	Client *auth.Client
}

// Commands defines all available playground commands
var Commands = []Command{
	{"help", "h", "Show detailed help", cmdShowHelp, false},
	{"quit", "q", "Exit playground", cmdQuit, false},
	{"status", "s", "Show current status", cmdShowStatus, false},
	{"clear", "c", "Clear screen", cmdClearScreen, false},
	{"toggle-readonly", "tr", "Toggle readonly mode", cmdToggleReadonly, false},
	{"list-tasks", "lt", "List all tasks", cmdListTasks, true},
	{"get-task", "gt", "Get specific task by ID", cmdGetTask, true},
	{"create-task", "ct", "Create a new task", cmdCreateTask, true},
	{"update-task", "ut", "Update existing task", cmdUpdateTask, true},
	{"delete-task", "dt", "Delete a task", cmdDeleteTask, true},
	{"list-projects", "lp", "List all projects", cmdListProjects, true},
	{"search-projects", "sp", "Search projects by name", cmdSearchProjects, true},
}

func main() {
	fmt.Println("🔧 Tudidi API Testing Playground")
	fmt.Println("==================================")

	cfg, err := initializeConfig()
	if err != nil {
		os.Exit(1)
	}

	client, err := initializeClient(cfg)
	if err != nil {
		log.Fatalf("❌ Client initialization failed: %v", err)
	}

	api := tudidi.NewAPI(client, cfg.Readonly)
	ctx := &PlaygroundContext{
		API:    api,
		Config: cfg,
		Client: client,
	}

	showStartupMessage(cfg)
	runInteractiveLoop(ctx)
}

func initializeConfig() (*config.Config, error) {
	cfg, err := config.ParseArgs()
	if err != nil {
		fmt.Printf("❌ Configuration error: %v\n", err)
		fmt.Println("\nUsage examples:")
		fmt.Println("  ./test-playground --url http://localhost:3002 --email admin@test.com --password secret")
		fmt.Println("  TUDIDI_URL=http://localhost:3002 TUDIDI_USER_EMAIL=admin@test.com TUDIDI_USER_PASSWORD=secret ./test-playground")
		return nil, err
	}
	return cfg, nil
}

func initializeClient(cfg *config.Config) (*auth.Client, error) {
	client, err := auth.NewClient(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	fmt.Printf("🔐 Authenticating with %s...\n", cfg.URL)
	if err := client.Login(cfg.Email, cfg.Password); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Println("✅ Authentication successful!")
	return client, nil
}

func showStartupMessage(cfg *config.Config) {
	readonlyStatus := ""
	if cfg.Readonly {
		readonlyStatus = " (READONLY MODE - destructive operations disabled)"
	}
	fmt.Printf("🚀 API ready%s\n\n", readonlyStatus)
}

func runInteractiveLoop(ctx *PlaygroundContext) {
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

		handleCommand(ctx, command, scanner)
		fmt.Println()
	}
}

func handleCommand(ctx *PlaygroundContext, input string, scanner *bufio.Scanner) {
	// Find matching command
	for _, cmd := range Commands {
		if input == cmd.Name || input == cmd.Alias {
			if cmd.RequiresAPI && ctx.API == nil {
				fmt.Println("❌ API not available")
				return
			}
			cmd.Handler(ctx, scanner)
			return
		}
	}

	fmt.Printf("❌ Unknown command: %s\n", input)
	fmt.Println("Type 'help' for available commands")
}

func showMenu() {
	fmt.Println("📋 Available Commands:")
	for _, cmd := range Commands {
		fmt.Printf("  %-20s (%s) - %s\n", cmd.Name, cmd.Alias, cmd.Description)
	}
}

// Command Handlers

func cmdShowHelp(ctx *PlaygroundContext, scanner *bufio.Scanner) {
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
	fmt.Println("  list-projects, lp")
	fmt.Println("    Lists all project lists/containers")
	fmt.Println("    Safe to use in readonly mode")
	fmt.Println()
	fmt.Println("  search-projects, sp")
	fmt.Println("    Search projects by name (case-insensitive)")
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

func cmdQuit(ctx *PlaygroundContext, scanner *bufio.Scanner) {
	fmt.Println("👋 Goodbye!")
	os.Exit(0)
}

func cmdShowStatus(ctx *PlaygroundContext, scanner *bufio.Scanner) {
	fmt.Println("📊 Current Status:")
	fmt.Printf("  Server URL:  %s\n", ctx.Config.URL)
	fmt.Printf("  Email:       %s\n", ctx.Config.Email)
	fmt.Printf("  Transport:   %s", ctx.Config.Transport)
	if ctx.Config.Transport == "sse" {
		fmt.Printf(" (port %d)", ctx.Config.Port)
	}
	fmt.Println()

	readonlyStatus := "WRITABLE (create/update/delete enabled)"
	if ctx.Config.Readonly {
		readonlyStatus = "READONLY (create/update/delete disabled)"
	}
	fmt.Printf("  Mode:        %s\n", readonlyStatus)
	fmt.Printf("  Time:        %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func cmdClearScreen(ctx *PlaygroundContext, scanner *bufio.Scanner) {
	fmt.Print("\033[2J\033[H")
	fmt.Println("🔧 Tudidi API Testing Playground")
	fmt.Println("==================================")
}

func cmdToggleReadonly(ctx *PlaygroundContext, scanner *bufio.Scanner) {
	ctx.Config.Readonly = !ctx.Config.Readonly
	ctx.API = tudidi.NewAPI(ctx.Client, ctx.Config.Readonly)

	status := "WRITABLE"
	if ctx.Config.Readonly {
		status = "READONLY"
	}
	fmt.Printf("🔄 Switched to %s mode\n", status)
}

func cmdListTasks(ctx *PlaygroundContext, scanner *bufio.Scanner) {
	fmt.Println("📋 Fetching tasks...")
	tasks, err := ctx.API.GetTasks()
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

func cmdGetTask(ctx *PlaygroundContext, scanner *bufio.Scanner) {
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
	task, err := ctx.API.GetTask(id)
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

func cmdCreateTask(ctx *PlaygroundContext, scanner *bufio.Scanner) {
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
		projects, err := ctx.API.GetProjects()
		if err != nil || len(projects) == 0 {
			fmt.Println("❌ No projects available and no project ID specified")
			return
		}
		projectID = projects[0].ID
		fmt.Printf("ℹ️  Using project ID %d (%s)\n", projectID, projects[0].Name)
	}

	req := tudidi.CreateTaskRequest{
		Name:      name,
		Note:      note,
		ProjectID: projectID,
		Status:    tudidi.NotStarted,
	}

	fmt.Println("🔨 Creating task...")
	task, err := ctx.API.CreateTask(req)
	if err != nil {
		fmt.Printf("❌ Error creating task: %v\n", err)
		return
	}

	fmt.Printf("✅ Task created successfully!\n")
	fmt.Printf("  ID:   %d\n", task.ID)
	fmt.Printf("  Name: %s\n", task.Name)
	fmt.Printf("  Note: %s\n", task.Note)
}

func cmdUpdateTask(ctx *PlaygroundContext, scanner *bufio.Scanner) {
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
	task, err := ctx.API.UpdateTask(id, req)
	if err != nil {
		fmt.Printf("❌ Error updating task: %v\n", err)
		return
	}

	fmt.Printf("✅ Task updated successfully!\n")
	fmt.Printf("  ID:   %d\n", task.ID)
	fmt.Printf("  Name: %s\n", task.Name)
	fmt.Printf("  Note: %s\n", task.Note)
}

func cmdDeleteTask(ctx *PlaygroundContext, scanner *bufio.Scanner) {
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
	task, err := ctx.API.GetTask(id)
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
	err = ctx.API.DeleteTask(id)
	if err != nil {
		fmt.Printf("❌ Error deleting task: %v\n", err)
		return
	}

	fmt.Println("✅ Task deleted successfully!")
}

func cmdListProjects(ctx *PlaygroundContext, scanner *bufio.Scanner) {
	fmt.Println("📁 Fetching project lists...")
	projects, err := ctx.API.GetProjects()
	if err != nil {
		fmt.Printf("❌ Error fetching projects: %v\n", err)
		return
	}

	displayProjects(projects)
}

func cmdSearchProjects(ctx *PlaygroundContext, scanner *bufio.Scanner) {
	fmt.Print("Enter project name to search for: ")
	if !scanner.Scan() {
		return
	}
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Println("❌ Project name cannot be empty")
		return
	}

	fmt.Printf("🔍 Searching projects by name: %s...\n", name)
	projects, err := ctx.API.SearchProjectsByName(name)
	if err != nil {
		fmt.Printf("❌ Error searching projects: %v\n", err)
		return
	}

	displayProjects(projects)
}

// Helper Functions

func displayProjects(projects []tudidi.Project) {
	if len(projects) == 0 {
		fmt.Println("📝 No projects found")
		return
	}

	fmt.Printf("✅ Found %d project(s):\n", len(projects))
	fmt.Println("ID   | Name                     | Active | Description")
	fmt.Println("-----|--------------------------|--------|------------")

	for _, project := range projects {
		active := "No"
		if project.Active {
			active = "Yes"
		}
		name := truncateString(project.Name, 24)
		desc := truncateString(project.Description, 20)
		fmt.Printf("%-4d | %-24s | %-6s | %s\n",
			project.ID, name, active, desc)
	}
}

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
