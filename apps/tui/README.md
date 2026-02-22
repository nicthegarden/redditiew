# RedditView TUI - Go Terminal User Interface

A terminal UI for browsing Reddit built with Bubble Tea, sharing the same data models and API as the React web app.

## Architecture

```
Go TUI App
    ↓
HTTP Requests (localhost:3002)
    ↓
Node.js API Server (@redditview/core)
    ↓
Reddit API
```

## Prerequisites

- Go 1.21 or later
- Node.js 16+ (for running the API server)

## Setup

### 1. Install Go Dependencies

```bash
go mod download
```

### 2. Start the Node.js API Server

From the monorepo root:

```bash
# Build the core package first
npm run build

# Start the API server on port 3002
npm run dev:api
```

Or manually:
```bash
node --loader ts-node/esm api-server.ts
```

The server will be available at `http://localhost:3002`.

### 3. Run the TUI App

```bash
go run main.go
```

## Controls

- **↑/↓** or **k/j** - Navigate posts
- **Enter** - View post details (coming soon)
- **q** or **Ctrl+C** - Quit

## Data Models

The Go TUI uses the same data structures as `@redditview/core`:

```go
type RedditPostData struct {
    ID       string
    Title    string
    Author   string
    Score    int
    Created  int64
    Comments int
    SelfText string
    URL      string
    SubName  string
}
```

These match the TypeScript interfaces in `packages/core/src/models/`.

## API Endpoints

The TUI calls these endpoints on the API server:

### Get Posts
```
GET http://localhost:3002/api/r/:subreddit.json?limit=20
```

Response format matches Reddit API structure.

### Search (Coming Soon)
```
GET http://localhost:3002/api/search.json?q=:query&limit=50
```

### Health Check
```
GET http://localhost:3002/health
```

Returns `{"status": "ok", "cache_size": 0}`

## Roadmap

- [ ] Display post details on Enter
- [ ] Show comments
- [ ] Search functionality
- [ ] Subreddit switcher
- [ ] Favorites/bookmarks
- [ ] Dark/light theme
- [ ] Pagination (load more posts)
- [ ] Responsive layout

## Development

### Building

```bash
go build -o redditview main.go
```

### Running with Hot Reload

Install `air`:
```bash
go install github.com/cosmtrek/air@latest
```

Then:
```bash
air
```

### Testing

```bash
go test ./...
```

## Project Structure

```
apps/tui/
├── main.go           # Entry point and TUI model
├── go.mod            # Go module definition
├── go.sum            # Go dependencies lock file
├── models.go         # (future) Data structures
├── api/              # (future) API client
│   └── client.go
├── ui/               # (future) UI components
│   ├── list.go
│   └── detail.go
└── README.md
```

## Contributing

When adding features:
1. Keep data models in sync with `@redditview/core`
2. Use the API server for all data fetching
3. Follow Bubble Tea patterns for state management
4. Test thoroughly with different terminal sizes

## Related Projects

- **React Web App**: `packages/web/`
- **Shared Core**: `packages/core/`
- **API Server**: `api-server.ts` (root)
- **Monorepo Docs**: `MONOREPO_ARCHITECTURE.md`

## License

Same as parent RedditView project
