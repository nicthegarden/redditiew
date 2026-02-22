# RedditView Architecture Diagram

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        REDDITVIEW MONOREPO                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐   │
│  │  React Web App   │  │   Go TUI App     │  │  Other Clients│   │
│  │  (localhost:     │  │  (Terminal)      │  │  (Mobile,    │   │
│  │   5173)          │  │                  │  │   Desktop)   │   │
│  └────────┬─────────┘  └────────┬─────────┘  └─────┬────────┘   │
│           │                     │                   │             │
│           │ /api requests       │ HTTP requests    │             │
│           │ (Vite proxy)        │ (port 3002)      │             │
│           └─────────┬───────────┴───────────┬──────┘             │
│                     │                       │                    │
│                     ▼                       ▼                    │
│           ┌──────────────────────────────────────┐              │
│           │   API Server (api-server.ts)         │              │
│           │   Port: 3002                         │              │
│           │   ┌────────────────────────────────┐ │              │
│           │   │ Proxies & Caches Data          │ │              │
│           │   │ • /api/r/:subreddit            │ │              │
│           │   │ • /api/r/:sub/comments/:id     │ │              │
│           │   │ • /api/search.json             │ │              │
│           │   └────────────────────────────────┘ │              │
│           └──────────────────┬───────────────────┘              │
│                              │                                   │
│                   Uses @redditview/core                          │
│                              │                                   │
│           ┌──────────────────▼───────────────────┐              │
│           │   packages/core (TypeScript)         │              │
│           │   ┌──────────────────────────────┐   │              │
│           │   │ API Client                   │   │              │
│           │   │ • RedditApiClient            │   │              │
│           │   │ • fetchPosts()               │   │              │
│           │   │ • fetchComments()            │   │              │
│           │   │ • search()                   │   │              │
│           │   └──────────────────────────────┘   │              │
│           │   ┌──────────────────────────────┐   │              │
│           │   │ Data Models                  │   │              │
│           │   │ • RedditPost                 │   │              │
│           │   │ • Comment                    │   │              │
│           │   │ • CacheEntry                 │   │              │
│           │   └──────────────────────────────┘   │              │
│           │   ┌──────────────────────────────┐   │              │
│           │   │ Cache                        │   │              │
│           │   │ • PostCache (in-memory)      │   │              │
│           │   │ • LocalStorageCache (browser)│   │              │
│           │   └──────────────────────────────┘   │              │
│           │   ┌──────────────────────────────┐   │              │
│           │   │ Utilities                    │   │              │
│           │   │ • formatTime()               │   │              │
│           │   │ • formatNum()                │   │              │
│           │   └──────────────────────────────┘   │              │
│           └──────────────────┬───────────────────┘              │
│                              │                                   │
│                   All HTTPS requests                            │
│                              │                                   │
│           ┌──────────────────▼───────────────────┐              │
│           │      REDDIT API                      │              │
│           │   https://www.reddit.com/            │              │
│           │   https://old.reddit.com/            │              │
│           └────────────────────────────────────────┘              │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

### Web App Request Flow
```
User clicks on subreddit
        ↓
React App calls: fetch('/api/r/golang.json')
        ↓
Vite Dev Proxy intercepts /api/
        ↓
Vite rewrites to: https://old.reddit.com/r/golang.json
        ↓
Reddit returns JSON
        ↓
Vite proxy returns response
        ↓
React renders posts
```

### TUI App Request Flow
```
User navigates menu
        ↓
Go TUI calls: http.Get("http://localhost:3002/api/r/golang.json")
        ↓
API Server (api-server.ts) receives request
        ↓
Checks cache, if miss:
        ↓
Fetches from: https://www.reddit.com/r/golang.json
        ↓
Caches response
        ↓
Returns JSON to Go TUI
        ↓
Bubble Tea renders posts
```

## Monorepo Benefits

```
┌─────────────────────────────────────────────────────┐
│  Single Codebase                                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ✅ Shared Core Package                            │
│     All clients use same business logic             │
│                                                     │
│  ✅ npm Workspaces                                 │
│     Dependencies managed together                  │
│     npm install installs all                       │
│                                                     │
│  ✅ Type Safety                                    │
│     TypeScript interfaces shared                   │
│     Go models auto-generated from schemas          │
│                                                     │
│  ✅ One Git Repo                                   │
│     Easy to keep code in sync                      │
│     Single commit for related changes              │
│                                                     │
│  ✅ Consistent API Contracts                       │
│     Web and TUI call same endpoints                │
│     Share error handling & caching                 │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## Technology Stack

| Layer | Web | TUI | Core |
|-------|-----|-----|------|
| **Frontend** | React 19 | Bubble Tea | - |
| **Styling** | CSS | Lipgloss | - |
| **Routing** | React Router | - | - |
| **Backend** | Vite Proxy | API Server | - |
| **Business Logic** | @redditview/core | @redditview/core | TypeScript |
| **HTTP Client** | Fetch API | http.Client | - |
| **Cache** | LocalStorageCache | PostCache | Core |
| **Data Format** | JSON | JSON | Types |

## Development Setup

### Quick Reference

```bash
# Install everything
npm install

# Build core package
npm run build

# Run web app
npm run dev

# Run TUI (in another terminal)
cd apps/tui && go run main.go

# Run API server (if TUI needs it)
npm run dev:api
```

## File Organization

```
redditiew-monorepo/
├── packages/
│   ├── core/
│   │   ├── src/
│   │   │   ├── api/index.ts         (4️⃣ API Client)
│   │   │   ├── models/index.ts      (1️⃣ Data Models)
│   │   │   ├── cache/index.ts       (3️⃣ Cache Logic)
│   │   │   ├── utils/index.ts       (2️⃣ Utilities)
│   │   │   └── index.ts             (Export all)
│   │   ├── package.json
│   │   └── tsconfig.json
│   └── web/
│       ├── src/
│       │   ├── App.tsx              (Uses core)
│       │   ├── components/
│       │   │   ├── CommentsList.tsx (Uses core)
│       │   │   └── PostDetail.tsx   (Uses core)
│       │   └── index.css
│       ├── index.html
│       ├── vite.config.ts           (Proxy config)
│       └── package.json
├── apps/
│   └── tui/
│       ├── main.go                  (Uses API server)
│       ├── go.mod
│       └── README.md
├── api-server.ts                    (Node.js API server)
├── proxy.ts                         (Old proxy - deprecated)
├── proxy.js                         (Old proxy - deprecated)
├── package.json                     (Monorepo root)
├── MONOREPO_ARCHITECTURE.md
├── TUI_SETUP_GUIDE.md
└── ARCHITECTURE.md                  (This file)
```

## Next Steps

1. ✅ Core package created with all shared logic
2. ✅ Web app refactored to use core
3. ✅ API server created (port 3002)
4. ✅ Go TUI scaffold created
5. 🔄 Add more features to TUI
6. 🔄 Add tests to core
7. 🔄 Setup CI/CD pipeline
8. 🔄 Deploy to production

## Contributing

When adding features:
1. **Start with core** - Add types/logic to `packages/core`
2. **Web implementation** - Use in React `packages/web`
3. **TUI implementation** - Create endpoints if needed, call from Go
4. **Test all** - Verify works in web and TUI

## Performance Optimization

- **Caching** - Core caches results, API server caches results
- **Request Deduplication** - Same request doesn't hit Reddit twice
- **Lazy Loading** - Components load data on demand
- **Browser Storage** - Web app uses LocalStorage for persistence
- **API Compression** - Server supports gzip responses
