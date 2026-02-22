# RedditView Monorepo Architecture

This is a monorepo containing shared core logic and multiple frontend applications for RedditView.

## Project Structure

```
redditview-monorepo/
├── packages/
│   ├── core/                 # Shared TypeScript business logic
│   │   ├── src/
│   │   │   ├── models/       # Data interfaces (Post, Comment, etc.)
│   │   │   ├── api/          # Reddit API client
│   │   │   ├── cache/        # Caching logic
│   │   │   ├── utils/        # Helper functions (formatTime, formatNum, etc.)
│   │   │   └── index.ts      # Main entry point
│   │   └── package.json
│   └── web/                  # React web application
│       ├── src/
│       ├── public/
│       ├── index.html
│       └── package.json
├── apps/
│   └── tui/                  # Go TUI application (to be created)
│       ├── main.go
│       ├── go.mod
│       └── ...
├── proxy.ts                  # Shared proxy server (handles CORS)
└── root package.json         # Monorepo configuration with workspaces
```

## Core Package (`@redditview/core`)

The core package contains all reusable business logic:

### Models
Data interfaces used across all applications:
- `RedditPost` - Post data structure
- `Comment` - Comment tree structure
- `CacheEntry` - Cache entry format
- `FetchPostsResult`, `FetchCommentsResult`, `SearchResult` - API responses

### API Client (`RedditApiClient`)
Main API client for Reddit interactions:
```typescript
import { RedditApiClient } from '@redditview/core'

const client = new RedditApiClient({
  baseUrl: '/api',
  timeout: 15000
})

// Fetch posts
const result = await client.fetchPosts('javascript')
if (result.error) {
  console.error(result.error)
} else {
  console.log(result.posts)
}

// Fetch comments
const comments = await client.fetchComments('/r/javascript/comments/abc123')

// Search
const search = await client.search('TypeScript tutorial')
```

### Cache
Two cache implementations:

**PostCache** - In-memory cache (Node.js, CLI)
```typescript
import { PostCache } from '@redditview/core'

const cache = new PostCache({
  ttl: 3600000, // 1 hour
  maxSize: 100
})
```

**LocalStorageCache** - Browser storage (React)
```typescript
import { LocalStorageCache } from '@redditview/core'

const cache = new LocalStorageCache()
```

### Utils
Helper functions:
- `formatTime(ts)` - Format timestamp (e.g., "2h")
- `formatTimeAgo(ts)` - Format with "ago" suffix (e.g., "2h ago")
- `formatNum(n)` - Format large numbers (e.g., "1.2K", "5.5M")
- `getSubredditFromInput(input)` - Clean subreddit input
- `isTextPost(post)` - Check if post is text-only
- `getPostThumbnail(post)` - Extract thumbnail URL

## React Web App (`@redditview/web`)

The React frontend uses the core package:

```typescript
import { RedditApiClient, formatTime, formatNum } from '@redditview/core'

const client = new RedditApiClient()
const posts = await client.fetchPosts('javascript')
```

### Proxy Configuration
The web app uses Vite's dev proxy (configured in `vite.config.ts`):
```typescript
proxy: {
  '/api': {
    target: 'https://old.reddit.com',
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api/, '')
  }
}
```

## Go TUI App (To Be Created)

A terminal UI application will be built in Go using Bubble Tea:

```
apps/tui/
├── main.go
├── go.mod
├── models/          # Type definitions matching @redditview/core
├── api/             # HTTP client calls to the Node.js API server
└── ui/              # Bubble Tea TUI components
```

The Go TUI will:
1. Call HTTP endpoints exposed by Node.js on port 3002
2. Display data in a terminal UI using Bubble Tea
3. Share the same data models and API contracts as the core

## Development

### Setup
```bash
npm install
```

This installs dependencies for all workspaces.

### Running the Web App

Development mode with Vite proxy:
```bash
npm run dev
```

The app runs on `http://localhost:5173` and proxies `/api` requests through Vite.

### Building the Core

Build the TypeScript core:
```bash
npm run build
```

Watch mode for development:
```bash
npm run dev:core
```

### Running the Proxy

The standalone proxy (for non-Vite development):
```bash
npm run dev:proxy
```

Runs on port 3001.

## Architecture Decisions

### Why a Monorepo?
- **Single source of truth** for shared logic
- **Easy to keep code in sync** across projects
- **Workspace management** via npm workspaces
- **Simplified development** - change core once, use everywhere

### Why Separate the Core?
- **Language agnostic** - Core can be used from JavaScript, TypeScript, Go, Python, etc.
- **Reusability** - Both React and Go TUI use the same business logic
- **Testability** - Core logic can be tested independently
- **Maintainability** - Clear separation of concerns

### Why HTTP for Go TUI?
- **Simplicity** - Go's HTTP client is excellent
- **Decoupling** - TUI doesn't need TypeScript dependencies
- **Deployment** - Can run TUI and backend on different machines
- **Future proof** - Could add mobile, CLI, desktop clients easily

## API Endpoints (Node.js Port 3002)

When the Node.js API server runs, it will expose:

```
GET /api/r/:subreddit.json          # Fetch posts
GET /api/r/:subreddit/comments/...  # Fetch comments
GET /api/search.json?q=:query       # Search posts
```

## Next Steps

1. ✅ Create core package with shared logic
2. ✅ Update React app to use @redditview/core
3. Create Node.js API server (port 3002)
4. Create Go TUI app with Bubble Tea
5. Add tests for core package
6. Add CI/CD pipeline

## Contributing

When adding features:
1. Add types/models to `packages/core/src/models/`
2. Add logic to `packages/core/src/api/` or other modules
3. Build: `npm run build`
4. Use in React: `import { ... } from '@redditview/core'`
5. Use in Go: Call HTTP API endpoints
