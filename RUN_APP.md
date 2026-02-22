# 🚀 RedditView - Quick Start Guide

## ✅ Current Status

The API server is **RUNNING NOW** on `http://localhost:3002`

```
✓ API Server (port 3002) - ACTIVE
✓ Reddit data fetching - WORKING
✓ Response caching - ACTIVE
```

---

## How to Start the App

### Option 1: Web App Only (Easiest)

```bash
./launch.sh web
```

Then open your browser: **http://localhost:5173**

The web app will automatically proxy API requests to localhost:3002.

---

### Option 2: TUI (Terminal UI) + API Server

```bash
./launch.sh tui
```

This starts:
1. API Server (port 3002)
2. Go TUI app with full feature parity to web

**Features:**
- Complete post browsing and navigation
- Real-time search and filtering
- Subreddit switching (press `s`)
- Full post details with content display
- Comment viewing support
- Professional color scheme
- Responsive terminal layout

**Controls:**
- `j/k` or `↑↓`: Navigate posts
- `/`: Search posts
- `s`: Switch subreddit
- `Enter`: View post details
- `c`: View comments
- `b`: Go back
- `q`: Quit

---

### Option 3: All Three (Web + TUI + API)

```bash
./launch.sh all
```

This starts all services:
1. API Server (port 3002)
2. Web App (port 5173)
3. TUI Terminal App

---

### Option 4: API Server Only

```bash
./launch.sh api
```

Then make requests to:
- `http://localhost:3002/api/r/golang` - Get posts
- `http://localhost:3002/health` - Health check
- `http://localhost:3002/api/stats` - Cache stats

---

## What's Different Now

### ✅ Fixed Issues

| Issue | Status | Fix |
|-------|--------|-----|
| CORS errors | ✅ Fixed | Removed hardcoded localhost:3001 URLs |
| go.sum missing | ✅ Fixed | Cleaned up go.mod dependencies |
| Update method type error | ✅ Fixed | Correct tea.Model return type |
| TypeScript server issues | ✅ Fixed | Converted to JavaScript |
| Need separate terminals | ✅ Fixed | Created launch.sh script |

### ✅ New Files

- `api-server.js` - Working API server
- `launch.sh` - Multi-platform launcher
- `packages/core/` - Shared TypeScript library
- `packages/web/` - React web app
- `apps/tui/` - Go TUI app

---

## Architecture

```
┌─────────────────────────────────────────────┐
│          RedditView - Full Stack            │
├─────────────────────────────────────────────┤
│                                             │
│  Web Browser        Terminal                │
│       │                  │                  │
│       └─────────┬────────┘                  │
│                 │                           │
│         ┌───────▼────────┐                  │
│         │  API Server    │                  │
│         │  port 3002     │                  │
│         │ (caching)      │                  │
│         └───────┬────────┘                  │
│                 │                           │
│         ┌───────▼──────────┐                │
│         │  Reddit API      │                │
│         │ (www.reddit.com) │                │
│         └──────────────────┘                │
│                                             │
└─────────────────────────────────────────────┘
```

---

## Quick Reference

### Available Commands

```bash
# Install dependencies
npm install && npm run build

# Run web app
./launch.sh web

# Run TUI
./launch.sh tui

# Run both
./launch.sh all

# Run API server only
./launch.sh api

# Stop everything (Ctrl+C)
```

### API Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/r/:subreddit` | GET | Fetch posts from subreddit |
| `/api/r/:sub/comments/:id` | GET | Fetch comments on post |
| `/api/search.json?q=:query` | GET | Search Reddit |
| `/health` | GET | Health check |
| `/api/stats` | GET | Cache statistics |

---

## Troubleshooting

### "Cannot connect to API server"

Make sure the API server is running:
```bash
curl http://localhost:3002/health
```

Should return:
```json
{"status":"ok","cache_size":0}
```

### "Port already in use"

Find and kill the process:
```bash
# Port 3002 (API)
lsof -i :3002 | grep -v COMMAND | awk '{print $2}' | xargs kill -9

# Port 5173 (Web)
lsof -i :5173 | grep -v COMMAND | awk '{print $2}' | xargs kill -9
```

### "go: command not found"

Install Go:
```bash
sudo pacman -S go
```

---

## What's Next?

1. **Try the Web App:**
   ```bash
   ./launch.sh web
   ```
   Open: http://localhost:5173

2. **Try the TUI:**
   ```bash
   ./launch.sh tui
   ```

3. **Explore the Code:**
   - Shared logic: `packages/core/src/`
   - Web app: `packages/web/src/`
   - TUI app: `apps/tui/main.go`

---

## Git Status

Latest commits:
- ✅ Transform to monorepo architecture
- ✅ Fix Go module dependencies
- ✅ Fix tea.Model return type
- ✅ Add JavaScript API server
- ✅ Add launch script
- ✅ Implement split-view TUI with search, filtering, and detailed post view

All changes committed! 🎉
