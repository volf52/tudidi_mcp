# Tudidi MCP Server

A Model Context Protocol (MCP) server for Tudidi task management. This server provides MCP tools to interact with a self-hosted Tudidi instance, allowing AI assistants to manage tasks and lists through the MCP protocol.

## Features

- **Session-based Authentication**: Authenticates with Tudidi server using username/password and maintains session cookies
- **Readonly Mode**: Optional readonly mode prevents destructive operations (create/update/delete)
- **Complete Task Management**: Full CRUD operations for tasks and lists
- **MCP Protocol**: Standard MCP server implementation using stdio transport

## Available Tools

| Tool | Description | Readonly Safe |
|------|-------------|---------------|
| `list_tasks` | List all tasks | ✅ |
| `get_task` | Get specific task by ID | ✅ |
| `create_task` | Create new task | ❌ |
| `update_task` | Update existing task | ❌ |
| `delete_task` | Delete task | ❌ |
| `list_task_lists` | List all task lists | ✅ |

## Installation

### Prerequisites
- Go 1.25.0 or later
- [mise](https://mise.jdx.dev/) (optional, for task automation)
- Access to a Tudidi server instance

### Build

```bash
# Using mise
mise run build

# Or directly with go
go build -o server
```

## Usage

### Basic Usage

```bash
./server --url <tudidi-server-url> --username <username> --password <password>
```

### Readonly Mode

```bash
./server --url <tudidi-server-url> --username <username> --password <password> --readonly
```

### Command Line Options

- `--url` (required): Tudidi server URL
- `--username` (required): Username for authentication
- `--password` (required): Password for authentication  
- `--readonly` (optional): Enable readonly mode to prevent destructive operations

### Example

```bash
./server --url https://my-tudidi.example.com --username admin --password mypassword --readonly
```

## MCP Integration

This server implements the MCP protocol over stdio. It can be integrated with MCP-compatible clients like:

- Claude Desktop
- Other MCP clients

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "tudidi": {
      "command": "/path/to/tudidi_mcp/server",
      "args": [
        "--url", "https://your-tudidi.com",
        "--username", "your-username", 
        "--password", "your-password",
        "--readonly"
      ]
    }
  }
}
```

## Development

### Project Structure

```
tudidi_mcp/
├── main.go              # Server entry point and initialization
├── auth/
│   └── client.go        # HTTP client with authentication
├── config/
│   ├── config.go        # Configuration and CLI parsing
│   └── config_test.go   # Configuration tests
├── tudidi/
│   └── api.go           # Tudidi API operations
├── tools/
│   └── handlers.go      # MCP tool implementations
├── go.mod               # Go module definition
├── mise.toml            # Task automation
└── AGENTS.md           # Development guidelines
```

### Build Commands

```bash
# Build
mise run build
# or: go build -o server

# Test
go test ./...

# Format
go fmt ./...

# Lint
go vet ./...
```

### Adding New Tools

1. Add API methods to `tudidi/api.go`
2. Create tool handlers in `tools/handlers.go`
3. Register tools in the `RegisterTools` method
4. Update this README

## API Endpoints

The server expects these Tudidi API endpoints:

- `POST /api/login` - Authentication
- `GET /api/tasks` - List tasks
- `GET /api/tasks/{id}` - Get task
- `POST /api/tasks` - Create task
- `PUT /api/tasks/{id}` - Update task
- `DELETE /api/tasks/{id}` - Delete task
- `GET /api/lists` - List task lists

## Error Handling

- Authentication failures are logged and cause server exit
- API errors are returned to the MCP client
- Readonly mode violations return descriptive error messages
- Network timeouts and connection issues are handled gracefully

## Security

- Credentials are passed via command line arguments (consider environment variables for production)
- Session cookies are stored in memory only
- HTTPS is recommended for the Tudidi server URL
- Readonly mode provides safe operations for untrusted scenarios

## License

This project is released under the MIT License.

## Contributing

1. Follow the coding guidelines in `AGENTS.md`
2. Ensure all tests pass: `go test ./...`
3. Format code: `go fmt ./...`
4. Check for issues: `go vet ./...`