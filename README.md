# Tudidi MCP Server

A Model Context Protocol (MCP) server for Tudidi task management. This server provides MCP tools to interact with a self-hosted Tudidi instance, allowing AI assistants to manage tasks and lists through the MCP protocol.

## Features

- **Session-based Authentication**: Authenticates with Tudidi server using email/password and maintains session cookies
- **Readonly Mode**: Optional readonly mode prevents destructive operations (create/update/delete) - defaults to true
- **Multiple Transports**: Supports both stdio and SSE (Server-Sent Events) transports
- **Complete Task Management**: Full CRUD operations for tasks and lists
- **MCP Protocol**: Standard MCP server implementation

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
# Default transport (stdio) with readonly mode enabled by default
./server --url <tudidi-server-url> --email <email> --password <password>

# Disable readonly mode
./server --url <tudidi-server-url> --email <email> --password <password> --readonly=false

# Use SSE transport on default port (8080)
./server --url <tudidi-server-url> --email <email> --password <password> --transport sse

# Use SSE transport on custom port
./server --url <tudidi-server-url> --email <email> --password <password> --transport sse --port 3000
```

### Using Environment Variables

```bash
export TUDIDI_URL="https://my-tudidi.example.com"
export TUDIDI_USER_EMAIL="admin@example.com"
export TUDIDI_USER_PASSWORD="mypassword"
export TUDIDI_READONLY="false"  # optional, defaults to true
export TUDIDI_TRANSPORT="sse"   # optional, defaults to stdio
export TUDIDI_PORT="3000"       # optional, defaults to 8080 for SSE

./server
```

### Mixed Usage (Environment + CLI)

Environment variables take precedence over CLI flags:

```bash
export TUDIDI_URL="https://my-tudidi.example.com"
export TUDIDI_USER_EMAIL="admin@example.com"

# Only need to specify password via CLI if not in environment
./server --password mypassword --readonly=false --transport sse --port 3000
```

### Transport Options

#### Stdio Transport (Default)
```bash
./server --url <tudidi-server-url> --email <email> --password <password> --transport stdio
```

#### SSE Transport
```bash
# Default port (8080)
./server --url <tudidi-server-url> --email <email> --password <password> --transport sse

# Custom port
./server --url <tudidi-server-url> --email <email> --password <password> --transport sse --port 3000
```

When using SSE transport, the server starts an HTTP server on the specified port (default 8080) and serves MCP over Server-Sent Events.

### Basic Usage

```bash
# Default: readonly=true, transport=stdio
./server --url <tudidi-server-url> --email <email> --password <password>
```

### Readonly Mode

```bash
# Explicitly enable readonly mode (default)
./server --url <tudidi-server-url> --email <email> --password <password> --readonly=true

# Disable readonly mode to allow destructive operations
./server --url <tudidi-server-url> --email <email> --password <password> --readonly=false
```

### Command Line Options

- `--url` (required): Tudidi server URL
- `--email` (required): Email for authentication
- `--password` (required): Password for authentication  
- `--readonly` (optional): Enable/disable readonly mode to prevent destructive operations (default: true)
- `--transport` (optional): Transport type - 'stdio' or 'sse' (default: stdio)
- `--port` (optional): Port for SSE transport (default: 8080, ignored for stdio)

### Environment Variables

- `TUDIDI_URL`: Tudidi server URL
- `TUDIDI_USER_EMAIL`: Email for authentication
- `TUDIDI_USER_PASSWORD`: Password for authentication
- `TUDIDI_READONLY`: Set to "true" or "false" for readonly mode (default: true)
- `TUDIDI_TRANSPORT`: Transport type - 'stdio' or 'sse' (default: stdio)
- `TUDIDI_PORT`: Port for SSE transport (default: 8080)

**Note**: Environment variables take precedence over command line flags.

### Example

```bash
# Stdio transport with readonly disabled
./server --url https://my-tudidi.example.com --email admin@example.com --password mypassword --readonly=false

# SSE transport on default port (8080) with readonly enabled (default)
./server --url https://my-tudidi.example.com --email admin@example.com --password mypassword --transport sse

# SSE transport on custom port
./server --url https://my-tudidi.example.com --email admin@example.com --password mypassword --transport sse --port 3000
```

## MCP Integration

This server implements the MCP protocol over stdio (default) or SSE transports. It can be integrated with MCP-compatible clients like:

- Claude Desktop (stdio transport)
- Web-based MCP clients (SSE transport)
- Other MCP clients

### Stdio Transport Integration

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
        "--readonly=false"
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
        "TUDIDI_READONLY": "false"
      }
    }
  }
}
```

### SSE Transport Integration

For SSE transport, start the server and connect to the specified port:

```bash
# Default port 8080
./server --url https://your-tudidi.com --email your-email@example.com --password your-password --transport sse

# Custom port 3000
./server --url https://your-tudidi.com --email your-email@example.com --password your-password --transport sse --port 3000
```

Then connect your MCP client to the SSE endpoint at `http://localhost:8080` (or your specified port).

## Development

### Project Structure

```
tudidi_mcp/
├── main.go              # Server entry point and initialization
├── cmd/
│   └── test-playground/ # Interactive testing tool
│       ├── main.go      # Test playground implementation
│       └── README.md    # Playground documentation
├── auth/
│   └── client.go        # HTTP client with authentication
├── config/
│   ├── config.go        # Configuration and CLI parsing
│   └── config_test.go   # Configuration tests
├── tudidi/
│   ├── api.go           # Tudidi API operations
│   ├── api_test.go      # Comprehensive API tests
│   └── README.md        # API testing documentation
├── tools/
│   └── handlers.go      # MCP tool implementations
├── go.mod               # Go module definition
├── mise.toml            # Task automation
└── AGENTS.md           # Development guidelines
```

### Build Commands

```bash
# Build main server
mise run build
# or: go build -o server

# Build test playground
mise run build-playground
# or: go build -o test-playground ./cmd/test-playground

# Test
go test ./...

# Format
go fmt ./...

# Lint
go vet ./...
```

### Testing & Development Tools

#### Interactive Test Playground
For manual testing and API exploration:

```bash
# Build the playground
mise run build-playground

# Run with your server
./test-playground --url http://localhost:3002 --email admin@test.com --password secret

# Or use environment variables
export TUDIDI_URL="http://localhost:3002"
export TUDIDI_USER_EMAIL="admin@example.com"
export TUDIDI_USER_PASSWORD="password"
./test-playground
```

The playground provides an interactive CLI for testing all API operations. See [`cmd/test-playground/README.md`](cmd/test-playground/README.md) for detailed usage.

#### API Integration Tests
For automated testing against a live server:

```bash
# Set test environment variables
export TUDIDI_TEST_URL="http://localhost:3002"
export TUDIDI_TEST_EMAIL="test@example.com" 
export TUDIDI_TEST_PASSWORD="password"

# Run API tests
go test ./tudidi -v
```

See [`tudidi/README.md`](tudidi/README.md) for comprehensive testing documentation.

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