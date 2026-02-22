# RedditView Architecture

## Why This Architecture? (Proxy Server Design)

RedditView uses a **local proxy server architecture** instead of direct Reddit API access. This section explains the design rationale and how it solves real-world limitations.

### The Problem: Direct Reddit API Access Doesn't Work for TUI

Reddit's public API has fundamental limitations that make direct access problematic for a terminal application:

#### 1. **OAuth2 Authentication Required**
- Reddit API requires OAuth2 authentication
- Most endpoints need `client_id`, `client_secret`, and `refresh_token`
- **Problem**: Can't safely embed credentials in a distributed application
- **Problem**: Users would need to run OAuth flow in browser, then paste tokens

#### 2. **No "Read-Only" Credentials**
- Reddit API doesn't offer simple API keys for public data
- All authentication requires full OAuth2 flow
- **Problem**: Even for reading public posts, you need user authentication
- **Problem**: Creates privacy concerns (app asks "which subreddits do you read?")

#### 3. **Strict Rate Limiting**
- Rate limit: **60 requests per hour per endpoint per user**
- When you hit the limit, Reddit blocks requests for 1 hour
- **Problem**: A TUI with multiple users on same machine = shared rate limit
- **Problem**: Browsing through subreddits quickly hits the limit

#### 4. **No Multi-Client Support**
- Each authenticated user has separate rate limits
- **Problem**: If 2 people use TUI on same machine, each uses their auth token
- **Problem**: Can't share a single credential (privacy, complexity)

#### 5. **Complex OAuth Workflow in Terminal**
- OAuth requires browser interaction
- **Problem**: Terminal app can't open browser automatically (security)
- **Problem**: Users would see: "Visit http://... and paste the code"
- **Problem**: Bad user experience compared to native apps

### The Solution: Local Proxy Server

RedditView solves all these problems by running a lightweight **local proxy server** that:

#### ✅ **No Authentication Required**
- Uses Reddit's **public JSON endpoints** (no OAuth)
- Example: `https://www.reddit.com/r/golang.json` returns data without authentication
- Server fetches data server-side, bypassing authentication completely

#### ✅ **Centralized Rate Limiting**
- Single proxy server = single rate limit bucket with Reddit
- All clients (TUI, Web UI, future clients) share the same rate limit
- Server implements **intelligent caching** to minimize requests

#### ✅ **Instant Caching**
- First request to `r/golang`: Fetches from Reddit (1 request)
- Second request to `r/golang` within 30 seconds: Instant from cache
- Users can browse seamlessly without hitting Reddit repeatedly

#### ✅ **Multiple Interfaces**
- Same backend works with TUI, Web UI, mobile, etc.
- All interfaces benefit from centralized caching
- Easy to add new clients in the future

#### ✅ **Zero Configuration**
- No API keys, no authentication tokens
- Users don't need Reddit accounts (for read-only access)
- Just run the server and start browsing

### Architecture Comparison

```
┌─────────────────────────────────────────────────────────────────┐
│                     DIRECT API ACCESS (❌ Doesn't Work)          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  User's Computer                          Reddit API             │
│  ┌──────────────┐                    ┌──────────────────┐        │
│  │  TUI App     │ ──[OAuth Flow]───> │                  │        │
│  │  (Go Binary) │ <──[Auth Token]──  │  Requires:       │        │
│  └──────────────┘                    │  • OAuth2 setup  │        │
│                                      │  • User browser  │        │
│  Problems:                           │  • Rate limit:   │        │
│  ❌ Needs OAuth                     │    60/hour/user  │        │
│  ❌ Can't embed credentials          │  • 🔒 Secure     │        │
│  ❌ Hits rate limit quickly          │  • No caching    │        │
│  ❌ Complex setup                    │                  │        │
│  ❌ Bad UX (paste tokens)           └──────────────────┘        │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│              LOCAL PROXY SERVER ARCHITECTURE (✅ Works!)         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Your Computer                                                   │
│  ┌──────────────┐         ┌────────────────────┐                │
│  │  TUI App     │         │  Web Browser       │                │
│  │  (Go Binary) │         │  (React UI)        │                │
│  └──────┬───────┘         └────────┬───────────┘                │
│         │                          │                             │
│         └──────────┬───────────────┘                             │
│                    │                                             │
│                    ▼                                             │
│         ┌────────────────────────┐                              │
│         │  LOCAL PROXY SERVER    │                              │
│         │  (Node.js/Express)     │                              │
│         │  Port: 3002            │                              │
│         │                        │                              │
│         │  ✅ No Auth needed     │                              │
│         │  ✅ Built-in cache     │                              │
│         │  ✅ Single rate limit  │                              │
│         │  ✅ Zero config        │                              │
│         └───────────┬────────────┘                              │
│                     │                                            │
│                     │ (Uses public endpoints)                   │
│                     │                                            │
│                     ▼                                            │
│             ┌──────────────────┐                                │
│             │  reddit.com      │                                │
│             │  (Public JSON)   │                                │
│             │  /r/subreddit    │                                │
│             └──────────────────┘                                │
│                                                                   │
│  Benefits:                                                       │
│  ✅ No authentication setup                                      │
│  ✅ Instant caching (repeat requests are instant)              │
│  ✅ Multiple clients share same backend                        │
│  ✅ Better rate limiting (centralized)                         │
│  ✅ Perfect for niche use case: read Reddit in terminal       │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Comparison Table

| Aspect | Direct API | Proxy Server |
|--------|-----------|--------------|
| **Authentication** | ❌ Requires OAuth2 | ✅ None needed |
| **Setup Complexity** | ❌ High (OAuth flow) | ✅ Simple (run server) |
| **Rate Limiting** | ❌ 60/hour per user | ✅ Shared, w/ caching |
| **Multi-Client** | ❌ Complex (per-user) | ✅ Seamless |
| **Caching** | ❌ Client-side only | ✅ Server-side + client |
| **Data Freshness** | ⚠️ Always fresh | ⚠️ Cache TTL (30 sec) |
| **Configuration** | ❌ Complex | ✅ Zero config |
| **User Experience** | ❌ Token pasting | ✅ Just run TUI |
| **Privacy** | ⚠️ Must authenticate | ✅ Anonymous read-only |

### Why This Is Perfect for the Niche

RedditView's proxy server design is specifically optimized for:

1. **Reading Reddit Content in Terminal** - No need for full Reddit UI or authentication
2. **Offline Exploration** - Browse cached content without internet
3. **Batch Processing** - Cache enables efficient exploration of multiple subreddits
4. **Zero Setup** - No API keys, no authentication, just run it
5. **Resource Efficient** - Lightweight proxy and TUI consume minimal resources

---

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

### Request/Response Cycle

```
┌─────────────────────────────────────────────────────────────────┐
│                    TYPICAL REQUEST FLOW                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Step 1: User Action                                            │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ User types 's' in TUI or clicks subreddit in Web UI  │       │
│  │ Request: "Show me posts from r/golang"              │       │
│  └──────────────────────────────────────────────────────┘       │
│                          ↓                                       │
│  Step 2: Client Makes Request                                   │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ TUI:    GET http://localhost:3002/api/r/golang.json │       │
│  │ WEB:    GET http://localhost:3002/api/r/golang.json │       │
│  │ (Both use same endpoint)                             │       │
│  └──────────────────────────────────────────────────────┘       │
│                          ↓                                       │
│  Step 3: Proxy Server Checks Cache                              │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ Server receives request                              │       │
│  │ Looks up: "r/golang" in cache                       │       │
│  │                                                      │       │
│  │ If cached (and fresh):                              │       │
│  │   → Jump to Step 5 (instant response!)              │       │
│  │                                                      │       │
│  │ If not cached or stale:                             │       │
│  │   → Continue to Step 4                              │       │
│  └──────────────────────────────────────────────────────┘       │
│                          ↓                                       │
│  Step 4: Fetch from Reddit (Cache Miss)                         │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ Server makes request:                                │       │
│  │ GET https://www.reddit.com/r/golang.json            │       │
│  │                                                      │       │
│  │ Reddit API returns:                                  │       │
│  │ {                                                    │       │
│  │   "data": {                                          │       │
│  │     "children": [                                    │       │
│  │       { "data": { "title": "...", "score": ... } }  │       │
│  │     ]                                                │       │
│  │   }                                                  │       │
│  │ }                                                    │       │
│  │                                                      │       │
│  │ Server caches response (expires in 30 sec)          │       │
│  └──────────────────────────────────────────────────────┘       │
│                          ↓                                       │
│  Step 5: Return to Client                                       │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ Server sends JSON response to client                │       │
│  │ Time taken:                                          │       │
│  │  - First request:  ~500ms (fetch from Reddit)       │       │
│  │  - Repeat within 30s: ~50ms (from cache)            │       │
│  └──────────────────────────────────────────────────────┘       │
│                          ↓                                       │
│  Step 6: Render in Client                                       │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ TUI: Parse JSON, render posts with Bubble Tea       │       │
│  │ WEB: Parse JSON, render posts with React            │       │
│  │                                                      │       │
│  │ User sees posts instantly!                          │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### Caching Strategy

```
┌─────────────────────────────────────────────────────────────────┐
│                    CACHE BEHAVIOR                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Request Timeline:                                              │
│                                                                   │
│  Time   Request              Cache Status      Time to Response │
│  ─────────────────────────────────────────────────────────────  │
│  00s    GET r/golang         MISS               ~500ms          │
│          └─> Fetch from Reddit, store in cache                 │
│                                                                   │
│  02s    GET r/golang         HIT                ~50ms           │
│          └─> Serve from cache (28s remaining)                  │
│                                                                   │
│  15s    GET r/sysadmin       MISS               ~500ms          │
│          └─> Different subreddit, not cached                   │
│                                                                   │
│  25s    GET r/golang         HIT                ~50ms           │
│          └─> Still cached (5s remaining)                       │
│                                                                   │
│  32s    GET r/golang         MISS (EXPIRED)     ~500ms          │
│          └─> Cache expired, fetch fresh data from Reddit       │
│                                                                   │
│  Benefits:                                                       │
│  • User browses smoothly (no waiting on Reddit)                │
│  • Repeat requests are instant                                  │
│  • Reduces load on Reddit API                                  │
│  • Single rate limit shared across all clients                 │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### TUI App Request Flow (Detailed)
```
User navigates menu
        ↓
Go TUI calls: http.Get("http://localhost:3002/api/r/golang.json")
        ↓
API Server (api-server.js) receives request
        ↓
Checks cache:
  │
  ├─ [HIT] Serve from cache (instant)
  │         ↓
  │       Return JSON → Go TUI → Bubble Tea renders
  │
  └─ [MISS] Fetch from Reddit
           ↓
         GET https://www.reddit.com/r/golang.json
           ↓
         Reddit API responds with ~25 posts
           ↓
         Server stores in in-memory cache
           ↓
         Return JSON → Go TUI → Bubble Tea renders
```

### Web App Request Flow (Detailed)
```
User clicks on subreddit in React app
        ↓
React calls: fetch('http://localhost:3002/api/r/golang.json')
        ↓
API Server (api-server.js) receives request
        ↓
Checks cache:
  │
  ├─ [HIT] Serve from cache (instant)
  │         ↓
  │       Return JSON → React state update → Re-render
  │
  └─ [MISS] Fetch from Reddit
           ↓
         GET https://www.reddit.com/r/golang.json
           ↓
         Reddit API responds with ~25 posts
           ↓
         Server stores in in-memory cache
           ↓
         Return JSON → React state update → Re-render
```

### Real-World Request/Response Examples

**Example 1: First Request to r/golang (Cache Miss)**
```bash
# Client Request
GET http://localhost:3002/api/r/golang.json

# Server logs:
[12:34:56] GET /api/r/golang.json - Cache MISS
[12:34:56] Fetching from https://www.reddit.com/r/golang.json
[12:34:57] Response: 25 posts cached for 30 seconds
[12:34:57] Returning 200 OK (450KB)

# Total time: ~500ms (waiting for Reddit)
```

**Example 2: Second Request to r/golang (Cache Hit)**
```bash
# Client Request (5 seconds later)
GET http://localhost:3002/api/r/golang.json

# Server logs:
[12:35:01] GET /api/r/golang.json - Cache HIT (25s remaining)
[12:35:01] Returning 200 OK (450KB) from cache

# Total time: ~50ms (served from memory)
```

**Example 3: Request to Different Subreddit (Cache Miss)**
```bash
# Client Request
GET http://localhost:3002/api/r/sysadmin.json

# Server logs:
[12:35:02] GET /api/r/sysadmin.json - Cache MISS (different subreddit)
[12:35:02] Fetching from https://www.reddit.com/r/sysadmin.json
[12:35:03] Response: 25 posts cached for 30 seconds
[12:35:03] Returning 200 OK (520KB)

# Total time: ~500ms
```

**Example 4: Rate Limit Protection**
```
Time:     Request              Rate Limit Usage    Reddit Response
─────────────────────────────────────────────────────────────────
12:00     GET /api/r/golang    Used: 1/60         ✅ 200 OK
12:01     GET /api/r/golang    Used: 1/60 (cache) ✅ 200 OK (instant)
12:02     GET /api/r/sysadmin  Used: 2/60         ✅ 200 OK
12:03     GET /api/r/golang    Used: 2/60 (cache) ✅ 200 OK (instant)

# With caching, 4 requests only use 2 rate limit slots!
# Without caching and multiple users, would use 4 slots per person
```

---

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

## Complete Architecture Diagrams

### Deployment Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    DEPLOYMENT OPTIONS                        │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│ Option 1: DEVELOPMENT (Local Machine)                        │
│ ┌────────────────────────────────────────────────────────┐   │
│ │ Your Laptop/Desktop                                    │   │
│ │ ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │   │
│ │ │ TUI App  │  │ Web UI   │  │ API Server (3002)    │  │   │
│ │ └──────────┘  └──────────┘  └──────────────────────┘  │   │
│ │                                                         │   │
│ │ Setup: npm install && npm run build                   │   │
│ │ Usage: ./redditview (TUI) or localhost:3000 (Web)     │   │
│ └────────────────────────────────────────────────────────┘   │
│                                                                │
│ Option 2: SYSTEMD SERVICE (Linux)                            │
│ ┌────────────────────────────────────────────────────────┐   │
│ │ Linux Machine (CachyOS, Ubuntu, Fedora, etc.)         │   │
│ │                                                         │   │
│ │ User-Level (~/.config/systemd/user/)                  │   │
│ │ ├─ redditview-api.service (port 3002)                 │   │
│ │ ├─ redditview-tui.service (in tmux)                   │   │
│ │ └─ redditview-web.service (port 3000)                 │   │
│ │                                                         │   │
│ │ System-Level (/etc/systemd/system/) [sudo]            │   │
│ │ └─ Same services, system-wide installation            │   │
│ │                                                         │   │
│ │ Setup: ./setup.sh (interactive or automated)          │   │
│ │ Status: systemctl --user status redditview-*          │   │
│ └────────────────────────────────────────────────────────┘   │
│                                                                │
└──────────────────────────────────────────────────────────────┘
```

### Service Communication Architecture

```
┌──────────────────────────────────────────────────────────────┐
│          HOW SERVICES COMMUNICATE WITH EACH OTHER            │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│  TUI Application (Go)        Web UI (React)                   │
│  ┌──────────────┐            ┌──────────────┐                │
│  │ • Keyboard   │            │ • Mouse      │                │
│  │   input      │            │   clicks     │                │
│  │ • Terminal   │            │ • Browser    │                │
│  │   rendering  │            │   rendering  │                │
│  └──────┬───────┘            └──────┬───────┘                │
│         │                           │                        │
│         │ HTTP Requests             │ HTTP Requests          │
│         │ (localhost:3002)          │ (localhost:3002)       │
│         │                           │                        │
│         └───────────────┬───────────┘                        │
│                         │                                    │
│                         ▼                                    │
│         ┌────────────────────────────┐                      │
│         │   API Server (Node.js)     │                      │
│         │   Port: 3002               │                      │
│         │                            │                      │
│         │  GET /api/r/subreddit      │                      │
│         │  GET /api/comments/:id     │                      │
│         │  POST /api/search          │                      │
│         │  GET /api/health (probe)   │                      │
│         │                            │                      │
│         │  ✅ In-Memory Cache        │                      │
│         │  ✅ Rate Limit Manager     │                      │
│         │  ✅ Reddit API Gateway     │                      │
│         └────────────────┬───────────┘                      │
│                          │                                   │
│                          │ HTTPS                            │
│                          │ (to Reddit)                      │
│                          │                                   │
│                          ▼                                   │
│              ┌────────────────────────┐                     │
│              │   Reddit Public API    │                     │
│              │   www.reddit.com       │                     │
│              │   old.reddit.com       │                     │
│              │                        │                     │
│              │ Returns: JSON (posts,  │                     │
│              │ comments, metadata)    │                     │
│              └────────────────────────┘                     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Request Path Analysis

```
┌────────────────────────────────────────────────────────────────┐
│              WHERE DOES YOUR REQUEST GO?                       │
├────────────────────────────────────────────────────────────────┤
│                                                                  │
│ Scenario 1: TUI User Browses r/golang                          │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │ You:       Press 's', type "golang", press Enter        │   │
│ │            ↓                                              │   │
│ │ TUI:       http.Get("localhost:3002/api/r/golang.json") │   │
│ │            ↓                                              │   │
│ │ Proxy:     [Check cache] → Cache MISS                   │   │
│ │            ↓                                              │   │
│ │ Reddit:    GET www.reddit.com/r/golang.json             │   │
│ │            ↓                                              │   │
│ │ Proxy:     Store response in memory cache (30s TTL)     │   │
│ │            ↓                                              │   │
│ │ TUI:       Receive JSON, render posts                   │   │
│ │            ↓                                              │   │
│ │ You:       See r/golang posts in TUI!                   │   │
│ │                                                           │   │
│ │ Latency: ~500ms (first request)                         │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│ Scenario 2: TUI User Re-loads r/golang (5 sec later)           │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │ You:       Press 'r' to refresh                          │   │
│ │            ↓                                              │   │
│ │ TUI:       http.Get("localhost:3002/api/r/golang.json") │   │
│ │            ↓                                              │   │
│ │ Proxy:     [Check cache] → Cache HIT (25s remaining)    │   │
│ │            ↓                                              │   │
│ │ TUI:       Receive JSON from memory instantly            │   │
│ │            ↓                                              │   │
│ │ You:       See r/golang posts again (instant!)          │   │
│ │                                                           │   │
│ │ Latency: ~50ms (cached, no Reddit request!)             │   │
│ │ Rate Limit Usage: 0 (still have 59/60)                  │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
│ Scenario 3: Both TUI and Web UI Browse Simultaneously          │
│ ┌──────────────────────────────────────────────────────────┐   │
│ │ You (TUI):    Browse r/sysadmin                          │   │
│ │ Friend (Web): Browse r/sysadmin                          │   │
│ │                                                           │   │
│ │ TUI:       GET /api/r/sysadmin.json                      │   │
│ │ Web:       GET /api/r/sysadmin.json                      │   │
│ │            ↓                                              │   │
│ │ Proxy:     First request → Fetch from Reddit, cache it │   │
│ │            Second request → Both get same cached data    │   │
│ │                                                           │   │
│ │ Reddit Rate Limit: Only 1 request used (not 2!)         │   │
│ │ Benefit: Share rate limit, both get fast response       │   │
│ └──────────────────────────────────────────────────────────┘   │
│                                                                  │
└────────────────────────────────────────────────────────────────┘
```

---

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
