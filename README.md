# Tudidi MCP Server

A Model Context Protocol (MCP) server for Tudidi task management. This server provides MCP tools to interact with a self-hosted Tudidi instance, allowing AI assistants to manage tasks and lists through the MCP protocol.

## Features

- **Session-based Authentication**: Authenticates with Tudidi server using email/password and maintains session cookies
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

## Usage

### Using Command Line Arguments

```bash
./server --url <tudidi-server-url> --email <email> --password <password>
```

### Using Environment Variables

```bash
export TUDIDI_URL="https://my-tudidi.example.com"
export TUDIDI_USER_EMAIL="admin@example.com"
export TUDIDI_USER_PASSWORD="mypassword"
export TUDIDI_READONLY="true"  # optional, for readonly mode

./server
```

### Mixed Usage (Environment + CLI)

Environment variables take precedence over CLI flags:

```bash
export TUDIDI_URL="https://my-tudidi.example.com"
export TUDIDI_USER_EMAIL="admin@example.com"

# Only need to specify password via CLI if not in environment
./server --password mypassword --readonly
```

### Basic Usage

```bash
./server --url <tudidi-server-url> --email <email> --password <password>
```

### Readonly Mode

```bash
./server --url <tudidi-server-url> --email <email> --password <password> --readonly
```

### Command Line Options

- `--url` (required): Tudidi server URL
- `--email` (required): Email for authentication
- `--password` (required): Password for authentication  
- `--readonly` (optional): Enable readonly mode to prevent destructive operations

### Environment Variables

- `TUDIDI_URL`: Tudidi server URL
- `TUDIDI_USER_EMAIL`: Email for authentication
- `TUDIDI_USER_PASSWORD`: Password for authentication
- `TUDIDI_READONLY`: Set to "true" for readonly mode

**Note**: Environment variables take precedence over command line flags.

### Example

```bash
./server --url https://my-tudidi.example.com --email admin@example.com --password mypassword --readonly
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
        "--email", "your-email@example.com", 
        "--password", "your-password",
        "--readonly"
      ]
    }
  }
}
```

Or using environment variables for better security:

```json
{
  "mcpServers": {
    "tudidi": {
      "command": "/path/to/tudidi_mcp/server",
      "env": {
        "TUDIDI_URL": "https://your-tudidi.com",
        "TUDIDI_USER_EMAIL": "your-email@example.com",
        "TUDIDI_USER_PASSWORD": "your-password",
        "TUDIDI_READONLY": "true"
      }
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

- Credentials can be provided via command line arguments or environment variables (environment variables recommended for production)
- Session cookies are stored in memory only
- HTTPS is recommended for the Tudidi server URL
- Readonly mode provides safe operations for untrusted scenarios
- Environment variables help avoid exposing credentials in process lists

## License

This project is released under the MIT License.

## Contributing

1. Follow the coding guidelines in `AGENTS.md`
2. Ensure all tests pass: `go test ./...`
3. Format code: `go fmt ./...`
4. Check for issues: `go vet ./...`