# Getting Started: RedditView with TUI

This guide walks you through setting up and running both the React web app and the Go TUI app.

## Project Overview

RedditView is now structured as a monorepo with:
1. **@redditview/core** - Shared TypeScript business logic (models, API client, utilities)
2. **@redditview/web** - React web UI
3. **Go TUI** - Terminal UI application

All three share the same data models and API contracts.

## Quick Start

### Prerequisites
- Node.js 16+ (for web app and core)
- Go 1.21+ (for TUI)
- npm or yarn (for Node.js packages)

### Option A: Run the React Web App (as before)

```bash
# Install dependencies
npm install

# Build the core package
npm run build

# Start the web app
npm run dev
```

Visit `http://localhost:5173`

### Option B: Run the Go TUI App

```bash
# Terminal 1: Start the API server (serves data to TUI)
npm run build
npm run dev:api

# Terminal 2: Go to the TUI app
cd apps/tui
go run main.go
```

### Option C: Run Both Web and TUI Together

```bash
# Terminal 1: Start the API server
npm run dev:api

# Terminal 2: Start the web app (uses Vite's built-in proxy)
npm run dev

# Terminal 3: Run the TUI app
cd apps/tui
go run main.go
```

## How It Works

### React Web App Flow
```
Browser (localhost:5173)
    ↓
Vite Dev Proxy (localhost:5173)
    ↓
Reddit API
```

The web app uses Vite's built-in proxy configured in `vite.config.ts`.

### Go TUI App Flow
```
Terminal TUI (localhost)
    ↓
API Server (localhost:3002)
    ↓
Reddit API
```

The TUI app calls the Node.js API server, which proxies requests to Reddit.

## Project Structure

```
redditiew-monorepo/
├── packages/
│   ├── core/              # Shared TypeScript package
│   │   ├── src/
│   │   │   ├── api/       # Reddit API client
│   │   │   ├── models/    # Data types
│   │   │   ├── cache/     # Caching logic
│   │   │   ├── utils/     # Helper functions
│   │   │   └── index.ts
│   │   └── package.json
│   └── web/               # React web app
│       ├── src/
│       ├── index.html
│       └── vite.config.ts
├── apps/
│   └── tui/               # Go TUI app
│       ├── main.go        # Entry point
│       ├── go.mod
│       └── README.md
├── api-server.ts          # Node.js API server (port 3002)
├── proxy.ts               # Old proxy server (not used)
└── MONOREPO_ARCHITECTURE.md
```

## Available Commands

### Root Level
```bash
npm run dev          # Start React web app (dev only)
npm run dev:api      # Start API server (for TUI and web to use)
npm run dev:core     # Watch and rebuild core package
npm run build        # Build everything for production
```

### TUI App
```bash
cd apps/tui
go run main.go       # Run the TUI
go build             # Build executable
go test              # Run tests
```

## Shared Core Package

The `@redditview/core` package exports:

```typescript
// Models
export { RedditPost, Comment, CacheEntry, ... }

// API Client
export { RedditApiClient, apiClient }

// Cache
export { PostCache, LocalStorageCache }

// Utils
export { formatTime, formatNum, ... }
```

### Using in React
```typescript
import { RedditApiClient, formatTime } from '@redditview/core'

const client = new RedditApiClient()
const result = await client.fetchPosts('javascript')
```

### Using in Go
```go
// Models are defined locally in Go matching the TypeScript definitions
type RedditPostData struct {
    Title   string
    Author  string
    Score   int
    // ...
}

// API calls go to the Node.js server
resp, _ := http.Get("http://localhost:3002/api/r/golang.json")
```

## Development Workflow

1. **Update shared logic** → Edit `packages/core/src/`
2. **Build** → `npm run build`
3. **Use in React** → Automatic (via npm workspaces)
4. **Use in TUI** → Restart the API server to pick up changes
5. **Commit** → All changes in one monorepo

## Next Steps

- [ ] Add post detail view to TUI
- [ ] Display comments in TUI
- [ ] Add search functionality
- [ ] Implement infinite scroll / pagination
- [ ] Add tests for core package
- [ ] Add CI/CD pipeline
- [ ] Deploy to production

## Troubleshooting

### "Cannot find module '@redditview/core'"
```bash
npm install
npm run build
```

### API Server not responding
```bash
npm run dev:api
# Check if it's running on http://localhost:3002/health
```

### TUI shows "Loading posts..." forever
```bash
# Make sure API server is running
npm run dev:api

# Check the server is accessible
curl http://localhost:3002/health
```

### Port 3002 already in use
```bash
# Find process using port 3002
lsof -i :3002

# Kill it (if needed)
kill -9 <PID>
```

## Architecture Benefits

✅ **Code Reuse** - Business logic shared between web and TUI  
✅ **Consistency** - Same data models everywhere  
✅ **Maintainability** - Single source of truth for core logic  
✅ **Scalability** - Easy to add mobile, desktop, or other clients  
✅ **Type Safety** - TypeScript definitions guide API contracts  

## Further Reading

- [Monorepo Architecture](./MONOREPO_ARCHITECTURE.md)
- [Go TUI README](./apps/tui/README.md)
- [React Web App](./packages/web/)
- [Core Package](./packages/core/)

## Support

If you encounter issues:
1. Check the relevant README (web, TUI, core)
2. Review the error message carefully
3. Try rebuilding everything: `npm run build`
4. Check ports are available (5173 for web, 3002 for API)
