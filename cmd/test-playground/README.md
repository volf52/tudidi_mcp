# Test Playground

An interactive command-line tool for manually testing the Tudidi API functionality. This is not a test suite, but a playground where you can experiment with different API operations interactively.

## Features

- 🔧 Interactive command-line interface
- 📋 All CRUD operations (Create, Read, Update, Delete)
- 🔒 Toggle between readonly and writable modes
- 📊 Real-time status display
- 🆘 Built-in help system
- ✅ Error handling and validation
- 🎨 Colorful output with emojis

## Building

### Using mise (recommended)
```bash
mise run build-playground
# or
mise bp
```

### Using go directly
```bash
go build -o test-playground ./cmd/test-playground
```

## Usage

### Basic Usage
```bash
# Using command line arguments
./test-playground --url http://localhost:3002 --email admin@test.com --password secret

# Using environment variables (recommended)
export TUDIDI_URL="http://localhost:3002"
export TUDIDI_USER_EMAIL="admin@example.com"
export TUDIDI_USER_PASSWORD="password"
./test-playground

# Readonly mode (default)
./test-playground --readonly

# Writable mode
./test-playground --readonly=false
```

### Environment Variables
```bash
export TUDIDI_URL="http://localhost:3002"        # Server URL
export TUDIDI_USER_EMAIL="admin@example.com"     # Login email
export TUDIDI_USER_PASSWORD="password"           # Password
export TUDIDI_READONLY="false"                   # Enable/disable readonly mode
export TUDIDI_TRANSPORT="stdio"                  # Transport type (stdio/sse)
export TUDIDI_PORT="8080"                        # SSE port (if using SSE)
```

## Available Commands

### Read Operations (safe in readonly mode)
- **`list-tasks` (lt)** - List all tasks with summary table
- **`get-task` (gt)** - Get detailed information about a specific task
- **`list-lists` (ll)** - List all project lists/containers

### Write Operations (disabled in readonly mode)
- **`create-task` (ct)** - Create a new task with interactive prompts
- **`update-task` (ut)** - Update an existing task
- **`delete-task` (dt)** - Delete a task (with confirmation)

### Utility Commands
- **`toggle-readonly` (tr)** - Switch between readonly and writable modes
- **`status` (s)** - Show current connection and mode status
- **`clear` (c)** - Clear the screen
- **`help` (h)** - Show detailed help
- **`quit` (q)** - Exit the playground

## Interactive Examples

### Example Session
```
🔧 Tudidi API Testing Playground
==================================
🔐 Authenticating with http://localhost:3002...
✅ Authentication successful!
🚀 API ready (READONLY MODE - destructive operations disabled)

📋 Available Commands:
  list-tasks (lt)     - List all tasks
  get-task (gt)       - Get specific task by ID
  create-task (ct)    - Create a new task
  ...

Enter command: list-tasks

📋 Fetching tasks...
✅ Found 3 tasks:
ID   | Name                     | Status | Project ID | Created
-----|--------------------------|--------|------------|--------
1    | Fix login bug           | Active | 1          | 2024-01-15
2    | Update documentation    | Done   | 2          | 2024-01-14
3    | Test API endpoints      | New    | 1          | 2024-01-16

Enter command: get-task
Enter task ID: 1

🔍 Fetching task 1...
✅ Task details:
  ID:          1
  Name:        Fix login bug
  Note:        Critical bug affecting user authentication
  Priority:    2
  Status:      Active (1)
  Project ID:  1
  Created:     2024-01-15T10:30:00Z
  ...

Enter command: toggle-readonly
🔄 Switched to WRITABLE mode

Enter command: create-task
Enter task name: Test New Feature
Enter task description (optional): Testing the new user dashboard
Enter project ID (or press Enter for default): 
ℹ️  Using project ID 1 (Main Project)
🔨 Creating task...
✅ Task created successfully!
  ID:   4
  Name: Test New Feature
  Note: Testing the new user dashboard

Enter command: quit
👋 Goodbye!
```

## Error Handling

The playground handles various error scenarios gracefully:

- **Authentication failures** - Clear error messages with retry suggestions
- **Invalid input** - Validation with helpful prompts
- **Network errors** - Connection status and retry information
- **Readonly violations** - Clear explanation when destructive operations are blocked
- **Missing resources** - Proper 404 handling with user-friendly messages

## Tips

1. **Start in readonly mode** to explore safely
2. **Use short commands** (lt, gt, ct) for faster interaction
3. **Toggle readonly mode** to test destructive operations safely
4. **Use `status`** to check current configuration
5. **Use `help`** for detailed command information
6. **Use `clear`** to clean up the display during long sessions

## Development

The playground is designed to be:
- **Safe** - Readonly mode prevents accidental data loss
- **Interactive** - Real-time feedback and prompts
- **Educational** - Clear error messages and help text
- **Flexible** - Easy to extend with new commands
- **Robust** - Comprehensive error handling

## Troubleshooting

### Common Issues

**Authentication Failed**
```
❌ Authentication failed: login failed with status 403
```
- Check your credentials and server URL
- Ensure the server is running and accessible

**Connection Refused**
```
❌ Authentication failed: login request failed: dial tcp: connection refused
```
- Verify the server URL and port
- Check if the Tudidi server is running

**Readonly Mode Violations**
```
❌ Error creating task: operation not allowed in readonly mode
```
- Use `toggle-readonly` to switch to writable mode
- Or restart with `--readonly=false`