# RedditView - Multi-Interface Reddit Browser

**Repository**: [github.com/nicthegarden/redditiew](https://github.com/nicthegarden/redditiew)

**[🎬 Watch Demo Video](https://www.youtube.com/watch?v=upnz0AzVMn8)**

A modern, feature-rich Reddit client available in both Terminal User Interface (TUI) and Web UI formats. Browse Reddit posts, read comments, and open threads directly from your terminal or web browser.

## 🎯 Features

### Core Functionality
- **Browse Reddit Posts** - Navigate posts from any subreddit with smooth pagination
- **View Comments** - Read threaded discussions with scrollable comment panels
- **Search Posts** - Full-text search across loaded posts
- **Subreddit Switching** - Quickly switch between subreddits without restarting
- **Open in Browser** - Launch post URLs directly in your default browser
- **Real-time Stats** - See post scores, comment counts, and author information

### TUI Application Features
- **Keyboard-Driven Navigation** - Efficient keybindings for power users
- **Split-View Display** - See post list, details, and comments simultaneously
- **Smooth Scrolling** - Arrow keys, Page Up/Down for precise navigation
- **Responsive Design** - Adapts to any terminal size
- **Error Handling** - Graceful error messages and recovery
- **Sort Control** - Toggle between hot and new posts with instant refresh
- **Quick Subreddit Shortcuts** - Configurable 1-9 keyboard shortcuts for favorite subreddits
- **Smart Comment Navigation** - Auto-close comments when switching posts
- **Warning System** - Visual alerts when navigating at comment boundaries

### Web UI Features  
- **Modern UI Design** - Clean, intuitive interface
- **Responsive Layout** - Works on desktop and tablet
- **Real-time Updates** - Live comment and post data
- **Customizable Theme** - Light/dark mode support
- **Sort Toggle** - Switch between hot and new posts
- **Keyboard Navigation** - Arrow keys, PageDown, Spacebar for efficient browsing
- **Touch Gestures** - Swipe left/right for post navigation, up/down for scrolling
- **Auto-load Comments** - Comments load automatically when post is selected
- **Dual Search Modes** - Filter local posts or search all of Reddit

### Web UI Search Modes

The search functionality in the Web UI offers two distinct modes:

#### 1. **Local Filtering** (Default)
- **How:** Type in the filter bar (middle input)
- **Scope:** Filters posts from the current subreddit
- **Speed:** Instant (client-side)
- **Coverage:** Searches post titles only
- **Result:** Narrows down already-loaded posts
- **No network request**

#### 2. **Subreddit-Scoped Search**
- **How:** Type in filter bar + Press `Enter` (while on a subreddit)
- **Scope:** Searches **within the current subreddit only**
- **Speed:** Network request required (~1-2 seconds)
- **Coverage:** Reddit's search algorithm for the subreddit
- **Result:** Loads new search results within the subreddit

**Example Workflow:**
1. Load r/sysadmin (shows 50 posts)
2. Type "python" → Filters to show only Python-related posts from r/sysadmin (instant, 10 posts shown)
3. Press Enter → Searches r/sysadmin for "python" (network request, loads 50 new results from r/sysadmin)
4. Switch to r/programming, type "python", press Enter → Now searches r/programming for "python"

**Key Difference:**
- Pressing Enter searches **within the current subreddit**, not all of Reddit
- This keeps you focused on the community you're browsing
- Change subreddit to search in a different community

## 🏗️ Architecture

RedditView uses a **local proxy server architecture** instead of direct Reddit API access.

### The Problem with Direct Reddit API Access
- ❌ Reddit's OAuth2 requires user authentication credentials
- ❌ No "read-only" mode - can't safely embed credentials
- ❌ Rate limiting per endpoint (60 requests/hour per user)
- ❌ Complex authentication workflows in terminal
- ❌ Browser dependencies for OAuth flow

### The Solution: Local Proxy Server
RedditView solves this by running a **lightweight proxy server** that:
- ✅ Fetches Reddit data server-side (bypasses client-side auth)
- ✅ Uses public endpoints (no credentials needed)
- ✅ Centralizes rate limiting and caching
- ✅ Works in terminal without browser
- ✅ Feeds your niche needs: read content without full Reddit UI

### System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        YOU (Your Computer)                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐              ┌─────────────────────────────┐  │
│  │   TUI App    │              │    Web Browser/UI           │  │
│  │  (Go Binary) │              │  (React/Vite Frontend)      │  │
│  │              │              │                             │  │
│  │  • Terminal  │              │  • Modern responsive UI     │  │
│  │  • Keyboard  │              │  • Mouse friendly           │  │
│  │  • Full TTY  │              │  • Same data backend        │  │
│  └──────┬───────┘              └────────────┬────────────────┘  │
│         │                                    │                   │
│         └────────────────┬───────────────────┘                   │
│                          │                                       │
│                          ▼                                       │
│         ┌────────────────────────────────┐                      │
│         │   LOCAL PROXY SERVER (Node.js) │                      │
│         │   Running on localhost:8765    │                      │
│         │                                │                      │
│         │  • Caches Reddit data          │                      │
│         │  • Handles rate limiting       │                      │
│         │  • Provides REST API           │                      │
│         │  • No credentials needed       │                      │
│         └───────────┬────────────────────┘                      │
│                     │                                           │
└─────────────────────┼───────────────────────────────────────────┘
                      │
                      │ (Your Internet Connection)
                      │
                      ▼
            ┌──────────────────┐
            │  reddit.com API  │
            │  Public Endpoints│
            │  (No Auth)       │
            └──────────────────┘
```

### Components

| Component | Language | Purpose |
|-----------|----------|---------|
| **TUI Application** | Go | Terminal interface with Bubble Tea framework |
| **Proxy Server** | Node.js (native http) | Local API server, caching, rate limiting |
| **Web UI** | React/Vite | Browser-based interface (same backend) |

### Key Design Benefits

- **Zero Configuration** - No API keys, no OAuth flow
- **Lightning Fast** - Results cached locally (1 minute TTL)
- **Multiple Interfaces** - Same backend, different UIs
- **Privacy Friendly** - No third-party tracking
- **Self-Contained** - Everything runs locally

---

## 📸 Screenshots

### TUI - Post List View
![TUI Post List](TUI.png)

### TUI - Comments View
![TUI Comments](TUI-Comment.png)

### Web UI
![Web UI](WebUI.png)

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21+ (for TUI)
- **Node.js** 18+ LTS
- **npm** 7+

### Installation

```bash
# Clone the repository
git clone https://github.com/nicthegarden/redditiew.git
cd redditiew

# Install dependencies
npm install

# Build the TUI binary
go build -o apps/tui/redditview ./apps/tui

# Build the web UI
npm run build

# Start the API server (port 8765)
npm start

# In another terminal, run the TUI
./apps/tui/redditview
```

### Web UI
```bash
# Start the development server (port 5173)
npm run dev

# Open http://localhost:5173 in your browser
```

👉 **See [QUICKSTART.md](QUICKSTART.md) for detailed step-by-step instructions**

## 📖 Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started in 5 minutes (Windows & Linux)
- **[INSTALLATION.md](INSTALLATION.md)** - Detailed installation and build instructions
- **[SYSTEMD_SETUP.md](SYSTEMD_SETUP.md)** - Run as systemd service with auto-start
- **[CONFIGURATION.md](CONFIGURATION.md)** - Configure the application to your preferences
- **[TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md)** - Complete keyboard shortcut reference
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Technical architecture and design decisions
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Contributing and development guide

## 🎮 Complete Keybindings Reference

### Terminal UI (TUI) - Full Keyboard Reference

#### Navigation & Browsing
| Key | Action | Context |
|-----|--------|---------|
| `↑` / `k` | Move cursor up / Scroll up | List view: previous post, Details/Comments: scroll up |
| `↓` / `j` | Move cursor down / Scroll down | List view: next post, Details/Comments: scroll down |
| `Enter` | View post details | List view → Details view |
| `Esc` or `Tab` | Back to list | Details/Comments view → List view |
| `h` / `Left` | Previous post | Details/Comments view |
| `l` / `Right` | Next post | Details/Comments view |
| `Page Up` | Page up | Scroll within any view |
| `Page Down` | Page down | Scroll within any view |

#### Comment Navigation
| Key | Action | Details |
|-----|--------|---------|
| `c` | Toggle comments | Details view: show/hide comments panel |
| `↑` at top | Navigate to previous post | Comments: auto-closes comments and loads previous post |
| `↓` at bottom | Navigate to next post | Comments: auto-closes comments and loads next post |
| `↑`/`↓` in middle | Scroll comments | Normal scrolling when not at boundary |
| **⚠️ Visual Warning** | Shows orange alert | When at comment boundary with next/previous post available |

#### Subreddit & Search
| Key | Action |
|-----|--------|
| `Ctrl+R` | Change subreddit |
| `Ctrl+F` | Search posts |
| `1-9` | Quick jump to favorite subreddit |
| `Esc` | Cancel search or subreddit input |

#### Sorting & Refresh
| Key | Action |
|-----|--------|
| `t` | Toggle sort (Hot ↔ New) |
| `F5` | Refresh posts |

#### Utility
| Key | Action |
|-----|--------|
| `w` | Open post in browser |
| `q` / `Ctrl+C` | Quit application |

---

### TUI - View States Explained

#### 📋 List View (Default)
Browse and navigate through posts from the subreddit.
- Show post number, score, comments count
- Shows current sort status: 📊 Hot or 🆕 New
- Navigate with ↑/↓ keys

**Footer shows**: `Post X/Y [Sort Status] • Enter: view • 1-9: subreddit • t: toggle sort • F5: refresh • q: quit`

#### 📄 Details View (Post Content)
View full post content, score, and metadata.
- Scroll through longer post content
- Press `c` to view comments
- See author and subreddit info
- Open URL with `w` key

**Footer shows**: `↑↓: scroll details • h/l: switch posts • w: open URL • Esc/Tab: back to list • c: view comments • t: toggle sort • q: quit`

#### 💬 Comments View (Comments Panel)
Read comments and threaded discussions.
- Scroll through comments
- **Smart boundary warnings** when at top/bottom:
  - ⚠️ Orange warning when at top: `⚠️ Next ↑ will load previous post`
  - ⚠️ Orange warning when at bottom: `⚠️ Next ↓ will load next post`
- Normal footer when scrolling in middle
- Up/Down arrows navigate to next/previous post at boundaries

**Footer shows** (normal): `↑↓: scroll comments • h/l: switch posts • w: open URL • Esc: close comments • Ctrl+F: search • t: toggle sort • q: quit`

**Footer shows** (warning): Orange text indicating boundary navigation

---

## ⚙️ Configuration

Create/edit `config.json` in the project root:

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
- **`default_subreddit`** - Subreddit to load on startup (default: "sysadmin")
- **`default_sort`** - Default sort order: "hot" or "new" (default: "hot")
- **`posts_per_page`** - Number of posts to fetch (default: 200)
- **`list_height`** - Height of post list in split view (default: 10)
- **`max_title_length`** - Truncate titles longer than this (default: 80)
- **`subreddit_shortcuts`** - Map 1-9 keys to favorite subreddits

#### Web Settings
- **`default_subreddit`** - Subreddit to load in browser (default: "sysadmin")
- **`default_sort`** - Default sort order: "hot" or "new" (default: "hot")
- **`posts_per_page`** - Posts displayed per page (default: 20)
- **`theme`** - UI theme: "light" or "dark" (default: "dark")

#### API Settings
- **`base_url`** - API server URL (default: "http://localhost:8765/api")
- **`timeout_seconds`** - HTTP request timeout (default: 10)

👉 **See [CONFIGURATION.md](CONFIGURATION.md) for detailed configuration options**

## 🏗️ Project Structure

```
redditiew/
├── apps/
│   ├── tui/                    # Terminal User Interface (Go)
│   │   ├── main.go             # TUI application with full keybindings
│   │   └── redditview          # Compiled binary
│   └── web/                   # Web interface (placeholder)
├── packages/
│   ├── core/                   # Core shared logic
│   │   ├── src/api/            # Reddit API client
│   │   ├── src/cache/          # Caching utilities
│   │   └── src/models/         # Data models
│   └── web/                   # React web UI package
├── src/                        # Web UI source (root workspace)
├── api-server.ts               # Node.js API server (pure http module)
├── proxy.ts                    # Vite proxy configuration
├── vite.config.ts              # Vite bundler configuration
├── config.json                 # Configuration file
├── package.json                # Node.js dependencies
├── launch.sh                   # Launch script
└── setup.sh                    # Setup script
```

## 🔧 System Requirements

### Minimum Requirements
- **Go** 1.21+ (for TUI)
- **Node.js** 18+ LTS
- **npm** 7+
- **Terminal** with 80x24 character minimum

### Recommended
- **Go** 1.24+
- **Node.js** 20+ LTS
- **Terminal** with 120x40 character minimum for best experience
- **Modern OS** (Windows 10+, Ubuntu 20.04+, macOS 10.15+)

## 🌐 API Reference

The application uses a local REST API server running on `localhost:8765`.

### Key Endpoints

**Get Posts from Subreddit**
```
GET /api/r/:subreddit
GET /api/r/:subreddit/hot
GET /api/r/:subreddit/new
```

**Get Comments for Post**
```
GET /api/r/:subreddit/comments/:id
```

**Search All of Reddit**
```
GET /api/search.json?q=query&type=link&limit=50
```

**Search Within a Subreddit (New in v1.1.3)**
```
GET /api/r/:subreddit/search.json?q=query&type=link&limit=50
```
Returns search results restricted to the specified subreddit with 60-second caching.

**Health Check**
```
GET /health
```

**Stats**
```
GET /api/stats
```

See the API server implementation in `api-server.js` for complete details.

## 🐛 Troubleshooting

### TUI Won't Start
```bash
# Check if API server is running
curl http://localhost:8765/api/r/sysadmin

# Rebuild the binary
cd apps/tui && go build -o redditview .

# Check Go installation
go version
```

### TUI Hangs on Quit
- **Fixed in v0.4.0** - HTTP requests now timeout after 30 seconds
- Application will exit cleanly even with pending requests
- Previous versions may have hung indefinitely

### Comments Not Loading
- Ensure API server is running: `npm start`
- Check your internet connection
- Verify the subreddit name is valid
- Try refreshing with `F5`
- Check terminal window height (must be at least 24 lines)

### Sort Toggle Not Working
- Ensure config.json has valid `default_sort` value ("hot" or "new")
- API server must be running
- Check that `/health` endpoint is accessible

### Performance Issues
- Increase terminal window size (minimum 80x24, recommended 120x40)
- Reduce `posts_per_page` in config.json
- Ensure system has at least 512MB RAM

### API Server Port Conflict
```bash
# If port 8765 is in use, kill the existing process
lsof -i :8765  # macOS/Linux
netstat -ano | findstr :8765  # Windows
```

## 📋 System Support

| OS | Status | Notes |
|----|--------|-------|
| Linux | ✅ Fully Supported | Tested on Ubuntu, Fedora, Arch |
| Windows | ✅ Fully Supported | Windows 10 and newer |
| macOS | ✅ Fully Supported | Intel and Apple Silicon |

## 🔗 Key Dependencies

### Backend (TUI)
- **Go** - TUI application language
- **Bubble Tea** - TUI framework for elegant terminal UIs
- **Lipgloss** - Terminal styling and rendering

### Backend (API Server)
- **Node.js** - Runtime
- **TypeScript** - Language

### Frontend (Web)
- **React** 19 - UI framework
- **Vite** - Build tool
- **React Router** - Client-side routing

## 📝 Recent Changes

### Latest Release (v1.1.4)
- ✨ **NEW: Mobile Reader Mode** - Cleaner reading experience on mobile/tablet
  - Hide thread listing to focus on reading selected post
  - Reader Mode button next to version (mobile only)
  - Hide search bar in Reader Mode for maximum space
  - Toggle button to switch between modes instantly
  - Desktop view unaffected (button hidden on desktop)

### Previous Release (v1.1.3)
- ✨ **NEW: Subreddit-Scoped Search Endpoint** - `/api/r/:subreddit/search.json?q=:query`
  - Search results are now restricted to the current subreddit
  - 60-second caching to prevent rate limiting
  - Validates query parameter (1+ characters required)
  - Properly handles special characters in search queries
  - Returns `404` for invalid subreddits
- 🐛 Fixed: Web UI now properly calls subreddit search instead of global search

### Previous Release (v1.1.2)
- ✨ Updated Web UI to use subreddit-scoped search when pressing Enter
- ✨ "Subreddit search: query" header indicator

### Previous Release (v1.1.1)
- 🐛 Fixed missing `onKeyDown` handler for filter input
- 🐛 Fixed HTTP error handling in search submissions
- 🐛 Fixed `isRedditSearch` flag reset behavior

### Previous Release (v1.1.0)
- ✨ Keyboard navigation (arrow keys, PageDown, Spacebar)
- ✨ Touch/swipe gestures for mobile (left/right, up/down)
- ✨ Auto-fetch comments when post selected

### v0.4.0 and Earlier
- ✨ Warning system for comment boundary navigation with orange highlighting
- 🐛 Fixed TUI hang on quit with 30-second HTTP timeout
- ✨ Smart arrow key navigation in comments view
- ✨ Sort toggle functionality (hot/new posts)
- ✨ Keyboard shortcuts (1-9) for favorite subreddits
- ✨ Auto-close comments when navigating posts
- ✨ Configuration support for defaults

See [git log](https://github.com/nicthegarden/redditiew/commits) for complete history.

## 🤝 Contributing

We welcome contributions! Please see [DEVELOPMENT.md](DEVELOPMENT.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

## 🙋 Support & Feedback

- **Repository**: [GitHub - nicthegarden/redditiew](https://github.com/nicthegarden/redditiew)
- **Issues**: [GitHub Issues](https://github.com/nicthegarden/redditiew/issues)
- **Discussions**: [GitHub Discussions](https://github.com/nicthegarden/redditiew/discussions)
- **Feedback**: Report issues or suggest features on GitHub

## 🎉 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework for Go
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling library
- Reddit - Data source (public API)
- Community - Contributors and testers

---

**Happy browsing! 🚀**
