# RedditView - Multi-Interface Reddit Browser

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

### Web UI Features  
- **Modern UI Design** - Clean, intuitive interface
- **Responsive Layout** - Works on desktop and tablet
- **Real-time Updates** - Live comment and post data
- **Customizable Theme** - Light/dark mode support

## 🏗️ Why This Architecture? (Proxy Server Design)

RedditView uses a **local proxy server architecture** instead of direct Reddit API access. Here's why this design is perfect for your niche needs:

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
│  │  (Go Binary) │              │  (Node.js/React Frontend)   │  │
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
│         │   Running on localhost:3002    │                      │
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

### How It Works

**1. Data Flow**
```
Terminal/Browser → Local Proxy → Reddit Public API → Local Proxy → Display
```

**2. Request Example**
```
You: "Show me posts from r/sysadmin"
      ↓
TUI/Web: GET http://localhost:3002/api/r/sysadmin.json
      ↓
Proxy Server:
  - Checks if cached (serve instantly)
  - If not cached, fetches from Reddit
  - Caches result (next request is instant)
  - Returns JSON to your app
      ↓
TUI/Web: Displays posts beautifully
```

**3. Why This Benefits You**
- **Zero Configuration** - No API keys, no OAuth flow
- **Lightning Fast** - Results cached locally
- **Multiple Interfaces** - Same backend, different UIs
- **Offline-ish** - Cached data available without internet
- **Privacy Friendly** - No third-party tracking
- **Self-Contained** - Everything runs locally

### Components

| Component | Language | Purpose |
|-----------|----------|---------|
| **TUI Application** | Go | Terminal interface with Bubble Tea framework |
| **Proxy Server** | Node.js/Express | Local API server, caching, rate limiting |
| **Web UI** | Node.js/React | Browser-based interface (same backend) |

### Key Design Benefits

✅ **Simplicity** - One local server, multiple clients
✅ **Performance** - Built-in caching layer
✅ **Flexibility** - Easy to add new interfaces (CLI, Web, Desktop, etc.)
✅ **Reliability** - No network dependency for cached data
✅ **Scalability** - Can serve multiple instances

---

## 📸 Screenshots

### TUI - Post List View
![TUI Post List](TUI.png)

### TUI - Comments View
![TUI Comments](TUI-Comment.png)

### Web UI
![Web UI](WebUI.png)

## 🚀 Quick Start

### For Linux Users
```bash
# Install dependencies (choose one)
# Ubuntu/Debian
sudo apt-get install golang-go nodejs npm

# Fedora/RHEL
sudo dnf install golang nodejs npm

# Install the application
git clone https://github.com/yourusername/redditiew-local.git
cd redditiew-local
npm install
go build -o apps/tui/redditview ./apps/tui
npm run build

# Start the server
npm start

# In another terminal, run the TUI
./apps/tui/redditview
```

### For Windows Users
```powershell
# Install dependencies from:
# - Go: https://golang.org/dl
# - Node.js: https://nodejs.org

# Clone and setup
git clone https://github.com/yourusername/redditiew-local.git
cd redditiew-local
npm install
go build -o apps/tui/redditview.exe ./apps/tui
npm run build

# Start the server
npm start

# In another PowerShell window, run the TUI
.\apps\tui\redditview.exe
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

## 🎮 Basic Usage

### Terminal UI (TUI)

**Starting the Application**
```bash
# Ensure API server is running
npm start

# In another terminal
./apps/tui/redditview
```

**Basic Navigation**
| Action | Keys |
|--------|------|
| Browse posts | `j`/`k` or `↑`/`↓` |
| View post details | `Enter` |
| View comments | `c` |
| Scroll details/comments | `↑`/`↓` or `Page Up`/`Page Down` |
| Switch subreddits | `s` |
| Search posts | `Ctrl+F` |
| Open post in browser | `w` |
| Back to list | `Esc` or `Tab` |
| Quit | `q` |

See [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md) for complete keybinding documentation.

### Web UI

1. Open http://localhost:3000 in your web browser
2. Select a subreddit from the sidebar
3. Click on any post to view details
4. Click "View Comments" to expand the comment section

## ⚙️ Configuration

Create/edit `config.json` in the project root:

```json
{
  "tui": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 200,
    "list_height": 10,
    "max_title_length": 80
  },
  "web": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 20,
    "theme": "dark"
  },
  "api": {
    "base_url": "http://localhost:3002/api",
    "timeout_seconds": 10
  }
}
```

👉 **See [CONFIGURATION.md](CONFIGURATION.md) for detailed configuration options**

## 🏗️ Project Structure

```
redditiew-local/
├── apps/
│   ├── tui/                    # Terminal User Interface (Go)
│   │   └── main.go            # TUI application
│   └── web/                   # Web interface (placeholder)
├── packages/
│   ├── api/                   # API server implementation
│   └── reddit-scraper/        # Reddit data scraper
├── api-server.js              # API server (JavaScript)
├── config.json                # Configuration file
├── package.json               # Node.js dependencies
├── README.md                  # This file
├── QUICKSTART.md              # Quick start guide
├── CONFIGURATION.md           # Configuration reference
└── ARCHITECTURE.md            # Technical architecture
```

## 🔧 System Requirements

### Minimum Requirements
- **Go** 1.19+ (for TUI)
- **Node.js** 16+ (for API server)
- **npm** 7+
- **Terminal** with 80x24 character minimum

### Recommended
- **Go** 1.21+
- **Node.js** 18+ LTS
- **Terminal** with 120x40 character minimum
- **Modern OS** (Windows 10+, Ubuntu 20.04+, macOS 10.15+)

## 🌐 API Reference

The application uses a local REST API server running on `localhost:3002`.

### Key Endpoints

**Get Posts from Subreddit**
```
GET /api/r/:subreddit.json?limit=200
```

**Get Comments for Post**
```
GET /api/r/:subreddit/comments/:postid
```

See the API server implementation in `api-server.js` for complete details.

## 🐛 Troubleshooting

### TUI Won't Start
```bash
# Check if API server is running
curl http://localhost:3002/api/r/sysadmin.json

# Rebuild the binary
cd apps/tui && go build -o redditview .

# Check Go installation
go version
```

### Comments Not Loading
- Ensure API server is running: `npm start`
- Check your internet connection
- Verify the subreddit name is valid
- Try refreshing with `F5`

### Performance Issues
- Increase terminal window size (minimum 80x24)
- Reduce `posts_per_page` in config.json
- Ensure system has at least 512MB RAM

## 📋 System Support

| OS | Status | Notes |
|----|---------| ------|
| Linux | ✅ Fully Supported | Tested on Ubuntu, Fedora, Arch |
| Windows | ✅ Fully Supported | Windows 10 and newer |
| macOS | ✅ Fully Supported | Intel and Apple Silicon |

## 🔗 Key Dependencies

### Backend
- **Go** - TUI application language
- **Bubble Tea** - TUI framework
- **Node.js/Express** - API server
- **Axios** - HTTP client

### Frontend (Web)
- **React** - UI framework
- **Tailwind CSS** - Styling
- **Axios** - HTTP client

## 📝 Recent Changes

### Latest Features (v0.2.0)
- ✨ Enhanced comment scrolling with proper height calculation
- ✨ Open posts directly in browser with `w` key
- ✨ Increased page scroll distance for faster navigation
- ✨ 200 posts per page by default (up from 50)
- 🐛 Fixed comment scrolling state propagation
- 🐛 Fixed list display bug affecting post visibility

See [git log](https://github.com/yourusername/redditiew-local/commits) for complete history.

## 🤝 Contributing

We welcome contributions! Please see [DEVELOPMENT.md](DEVELOPMENT.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

## 🙋 Support & Feedback

- **Issues**: [GitHub Issues](https://github.com/yourusername/redditiew-local/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/redditiew-local/discussions)
- **Email**: support@example.com

## 🎉 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- Reddit - Data source

---

**Happy browsing! 🚀**

For more information, see the [complete documentation](#-documentation).
