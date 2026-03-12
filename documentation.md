# RedditView Feature Documentation

This document provides detailed documentation of all features in RedditView.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [API Server](#api-server)
3. [TUI Application](#tui-application)
4. [Web UI](#web-ui)
5. [Configuration](#configuration)
6. [External Access](#external-access)
7. [Systemd Services](#systemd-services)

---

## Architecture Overview

RedditView is a multi-interface Reddit browser consisting of three main components:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Your Computer                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐              ┌─────────────────────────────┐  │
│  │   TUI App    │              │    Web Browser/UI           │  │
│  │  (Go Binary) │              │  (React/Vite Frontend)      │  │
│  └──────┬───────┘              └────────────┬────────────────┘  │
│         │                                    │                   │
│         └────────────────┬───────────────────┘                   │
│                          ▼                                       │
│         ┌────────────────────────────────┐                      │
│         │   LOCAL PROXY SERVER (Node.js)  │                      │
│         │   Port: 8765                     │                      │
│         └────────────────────────────────┘                      │
│                          │                                       │
└──────────────────────────┼───────────────────────────────────────┘
                           │
                           ▼
                  ┌──────────────────┐
                  │  reddit.com API  │
                  │  Public Endpoints│
                  └──────────────────┘
```

### Technology Stack

| Component | Technology | Version |
|-----------|------------|---------|
| TUI | Go + Bubble Tea | Go 1.21+ |
| API Server | Node.js (native http) | Node 18+ |
| Web UI | React + Vite | React 19, Vite 7 |
| TypeScript | TypeScript | 5.x |

---

## API Server

The API server (`api-server.ts`) is a lightweight proxy that fetches data from Reddit's public API.

### Features

- **Caching**: 1-minute TTL cache for all Reddit requests
- **Rate Limiting**: Handled implicitly through caching
- **CORS**: Enabled for cross-origin requests
- **No Authentication**: Uses Reddit's public endpoints

### Ports

| Port | Service | Description |
|------|---------|-------------|
| 8765 | API Server | REST API for TUI and Web |

### Endpoints

#### Health Check
```
GET /health
```
Returns server status and cache statistics.

**Response:**
```json
{
  "status": "ok",
  "cache_size": 15,
  "uptime": 3600
}
```

#### Server Stats
```
GET /api/stats
```
Returns detailed server statistics.

**Response:**
```json
{
  "cache_size": 15,
  "cache_ttl": 60000,
  "uptime": 3600.5
}
```

#### Get Posts
```
GET /api/r/:subreddit
GET /api/r/:subreddit/hot
GET /api/r/:subreddit/new
```

**Parameters:**
- `subreddit` - Subreddit name (e.g., "sysadmin", "programming")
- Sort options: `hot`, `new` (default: hot)
- `limit` - Number of posts (default: 50, max: 100)
- `after` - Pagination token

**Example:**
```bash
curl http://localhost:8765/api/r/sysadmin
curl http://localhost:8765/api/r/programming/new?limit=100
```

#### Get Comments
```
GET /api/r/:subreddit/comments/:id
```

**Example:**
```bash
curl http://localhost:8765/api/r/sysadmin/comments/abc123
```

#### Search
```
GET /api/search.json?q=:query&type=link&limit=50
```

**Parameters:**
- `q` - Search query (required)
- `type` - Search type (default: link)
- `limit` - Results limit (default: 50)

**Example:**
```bash
curl "http://localhost:8765/api/search.json?q=kubernetes"
```

---

## TUI Application

The Terminal User Interface (TUI) is built with Go using the Bubble Tea framework.

### View Modes

#### 1. List View (Default)
The main post browsing interface showing:
- Post number
- Score
- Comment count
- Post title
- Current sort indicator (📊 Hot or 🆕 New)

**Footer:** `Post X/Y [Sort Status] • Enter: view • 1-9: subreddit • t: toggle sort • F5: refresh • q: quit`

#### 2. Details View
View full post content including:
- Full post title
- Post body/content
- Author information
- Subreddit name
- Score and timestamps

**Footer:** `↑↓: scroll • h/l: prev/next • w: open URL • Esc: back • c: comments • t: sort • q: quit`

#### 3. Comments View
Threaded comment display with:
- Comment author
- Comment score
- Comment body (collapsible for nested replies)
- Visual hierarchy for threaded comments

**Smart Navigation:**
- At top: Press ↑ to load previous post
- At bottom: Press ↓ to load next post
- Warning displayed at boundaries

### Keybindings

#### Navigation
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up / Scroll up |
| `↓` / `j` | Move down / Scroll down |
| `Enter` | View post details |
| `Esc` / `Tab` | Back to list |
| `h` / `←` | Previous post |
| `l` / `→` | Next post |
| `Page Up` | Page up |
| `Page Down` | Page down |

#### Comments
| Key | Action |
|-----|--------|
| `c` | Toggle comments panel |

#### Subreddit & Search
| Key | Action |
|-----|--------|
| `Ctrl+R` | Change subreddit |
| `Ctrl+F` | Search posts |
| `1-9` | Quick jump to favorite subreddit |
| `Esc` | Cancel input |

#### Sorting & Refresh
| Key | Action |
|-----|--------|
| `t` | Toggle sort (Hot ↔ New) |
| `F5` | Refresh posts |

#### Utility
| Key | Action |
|-----|--------|
| `w` | Open post in browser |
| `q` / `Ctrl+C` | Quit |

### Features

- **Split-View Display**: Post list and details shown simultaneously
- **Smooth Scrolling**: Arrow keys and Page Up/Down for navigation
- **Responsive Design**: Adapts to terminal size
- **Error Handling**: Graceful error messages and recovery
- **Sort Toggle**: Instant switch between hot and new posts
- **Subreddit Shortcuts**: 1-9 keys for favorite subreddits
- **Smart Comment Navigation**: Auto-close comments when switching posts
- **Warning System**: Visual alerts at comment boundaries

---

## Web UI

The web interface is built with React 19 and Vite.

### Features

- **Modern Responsive Design**: Works on desktop and tablet
- **Mouse Support**: Full point-and-click navigation
- **Keyboard Navigation**: Full keyboard control for efficient browsing
- **Touch/Swipe Support**: Swipe gestures for mobile devices
- **Theme Support**: Light/dark mode toggle
- **Real-time Updates**: Live post and comment data
- **Sort Toggle**: Switch between hot and new posts
- **Comment Auto-load**: Comments load automatically when post is selected

### Keyboard Controls

#### Left Pane (Posts List)

| Key | Action |
|-----|--------|
| `↑` Up Arrow | Select previous post |
| `↓` Down Arrow | Select next post |
| `Enter` | Open selected post in right pane |
| `Ctrl+F` | Focus on browser find (in right pane) |

#### Right Pane (Post Details & Comments)

| Key | Action |
|-----|--------|
| `→` Right Arrow | Scroll comments down (150px) |
| `←` Left Arrow | Scroll comments up (150px) |
| `PageDown` | Scroll comments down by page (400px) |
| `PageUp` | Scroll comments up by page (400px) |
| `Spacebar` | Scroll comments down by page, auto-advance to next post when at bottom |

#### Tab Navigation

| Key | Action |
|-----|--------|
| `Tab` | Move focus forward (Search → Filter → Posts List) |
| `Shift+Tab` | Move focus backward (Posts List → Filter → Search) |

### Touch/Swipe Gestures (Mobile)

| Gesture | Action |
|---------|--------|
| **Swipe Right** (← direction) | Go to previous post |
| **Swipe Left** (→ direction) | Go to next post |
| **Swipe Up** | Scroll comments down in right pane |
| **Swipe Down** | Scroll comments up in right pane |

**Note:** Swipes require at least 50px movement and must complete within 500ms.

### Navigation Workflow

1. **Browse Posts**: Use `↑↓` arrow keys to navigate the left pane
2. **Select Post**: Press `Enter` or click to open in right pane
3. **Read Details**: Post content and comments auto-load in right pane
4. **Scroll Comments**: Use arrow keys, PageDown/PageUp, or spacebar
5. **Auto-Advance**: When scrolling reaches the bottom of comments with spacebar, automatically move to next post

### Ports

| Port | Service | Description |
|------|---------|-------------|
| 5173 | Vite Dev Server | Development (default) |
| 5173+ | Vite Dev Server | Auto-selected if port in use |

### Access

**Development:**
```bash
npm run dev
# Open http://localhost:5173
```

**Production:**
```bash
npm run build
npm run preview
```

---

## Configuration

### config.json

The main configuration file located in the project root.

```json
{
  "tui": {
    "default_subreddit": "sysadmin",
    "default_sort": "hot",
    "posts_per_page": 200,
    "list_height": 10,
    "max_title_length": 80,
    "subreddit_shortcuts": {
      "1": "sysadmin",
      "2": "programming",
      "3": "linuxadmin",
      "4": "homelab",
      "5": "devops"
    }
  },
  "web": {
    "default_subreddit": "sysadmin",
    "default_sort": "hot",
    "posts_per_page": 20,
    "theme": "dark"
  },
  "api": {
    "base_url": "http://localhost:8765/api",
    "timeout_seconds": 10
  }
}
```

### Configuration Options

#### TUI Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `default_subreddit` | string | "sysadmin" | Subreddit to load on startup |
| `default_sort` | string | "hot" | Sort order: "hot" or "new" |
| `posts_per_page` | number | 200 | Posts to fetch per request |
| `list_height` | number | 10 | Height of post list in split view |
| `max_title_length` | number | 80 | Maximum title display length |
| `subreddit_shortcuts` | object | {...} | 1-9 key mappings |

#### Web Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `default_subreddit` | string | "sysadmin" | Default subreddit |
| `default_sort` | string | "hot" | Sort order |
| `posts_per_page` | number | 20 | Posts per page |
| `theme` | string | "dark" | UI theme: "light" or "dark" |

#### API Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `base_url` | string | "http://localhost:8765/api" | API server URL |
| `timeout_seconds` | number | 10 | HTTP request timeout |

---

## External Access

RedditView can be configured for network access, allowing other computers on your LAN to use the API and Web UI.

### Configuration

1. **API Server**: Edit `api-server.ts` to listen on `0.0.0.0` instead of `localhost`
2. **Vite Dev Server**: Edit `vite.config.ts` to listen on `0.0.0.0`

### Network Ports

| Service | Port | Purpose |
|---------|------|---------|
| API Server | 8765 | REST API access |
| Web UI | 5173 | Browser interface |

### Example Access

```bash
# From another computer on the network
curl http://192.168.1.104:8765/health
curl http://192.168.1.104:8765/api/r/sysadmin

# Open in browser
http://192.168.1.104:5173
```

### Firewall Configuration

**Ubuntu/Debian:**
```bash
sudo ufw allow 8765/tcp
sudo ufw allow 5173/tcp
```

**Fedora/RHEL:**
```bash
sudo firewall-cmd --permanent --add-port=8765/tcp
sudo firewall-cmd --permanent --add-port=5173/tcp
sudo firewall-cmd --reload
```

---

## Systemd Services

RedditView can run as systemd services for auto-start on boot.

### Service Files

Three systemd service files are provided:

1. **redditview-api.service** - API server (port 8765)
2. **redditview-web.service** - Web UI (port 5173)
3. **redditview-tui.service** - TUI application

### Installation

**User-level services (recommended):**
```bash
# Create directories
mkdir -p ~/.config/systemd/user

# Copy service files
cp systemd-templates/redditview-api.service ~/.config/systemd/user/
cp systemd-templates/redditview-web.service ~/.config/systemd/user/

# Enable and start
systemctl --user enable redditview-api
systemctl --user start redditview-api
```

**System-level services:**
```bash
# Copy to system directory
sudo cp systemd-templates/redditview-api.service /etc/systemd/system/

# Enable and start
sudo systemctl enable redditview-api
sudo systemctl start redditview-api
```

### Management Commands

```bash
# Start services
systemctl --user start redditview-api
systemctl --user start redditview-web

# Check status
systemctl --user status redditview-api

# View logs
journalctl --user -u redditview-api -f

# Restart on changes
systemctl --user daemon-reload
```

---

## Troubleshooting

### TUI Won't Start

1. **Check API server is running:**
   ```bash
   curl http://localhost:8765/health
   ```

2. **Verify Go is installed:**
   ```bash
   go version
   ```

3. **Rebuild the binary:**
   ```bash
   cd apps/tui && go build -o redditview .
   ```

### TUI Hangs on Quit

Fixed in v0.4.0 - HTTP requests now timeout after 30 seconds.

### Comments Not Loading

1. Ensure API server is running: `npm start`
2. Check internet connection
3. Verify subreddit name is valid
4. Try refreshing with `F5`

### API Server Port Conflict

```bash
# Check what's using the port
lsof -i :8765

# Kill the process if needed
kill <PID>
```

---

## Recent Changes

### v0.4.0
- Warning system for comment boundary navigation
- 30-second HTTP timeout for graceful quit
- All HTTP requests timeout gracefully

### v0.3.0
- Smart arrow key navigation in comments
- Sort toggle (hot/new)
- Keyboard shortcuts (1-9)
- Auto-close comments when navigating posts

### v0.2.0
- Enhanced comment scrolling
- Open posts in browser (`w` key)
- 200 posts per page default

---

## License

MIT License - See LICENSE file for details.
