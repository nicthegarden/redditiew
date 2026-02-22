# 🎉 RedditView TUI - Launch Guide

## ✅ Everything is Set Up and Ready!

The TUI launcher script is now fully functional and handles all edge cases automatically.

---

## How to Launch the TUI

### Single Command
```bash
./launch.sh tui
```

That's it! The script automatically:
1. ✅ Kills any existing processes on ports 3002 and 5173
2. ✅ Starts the API Server on port 3002
3. ✅ Waits for API to be ready
4. ✅ Starts the Go TUI app
5. ✅ Cleans up everything on exit

---

## What You'll See

When you run `./launch.sh tui`:

```
╭─────────────────────────────────────────────────╮
│  RedditView - Multi-Platform Launcher           │
╰─────────────────────────────────────────────────╯

▶ Starting API Server (port 3002)
✓ API Server running

▶ Starting TUI (Terminal UI)

Loading posts from r/golang...
```

Then the TUI will display:

```
RedditView TUI

r/golang
───────────────────────────────────────────────────

▶ Small Projects
  u/gopherbot · ↑ 22 · 💬 5

  Weekly megathread for sharing small Go projects
  u/darshanime · ↑ 18 · 💬 12

  Tips for optimizing Go code
  u/coder123 · ↑ 15 · 💬 8

───────────────────────────────────────────────────
↑↓/jk: navigate · q: quit
```

---

## Controls

| Key | Action |
|-----|--------|
| `↑` arrow | Scroll up |
| `↓` arrow | Scroll down |
| `j` | Scroll down (Vim style) |
| `k` | Scroll up (Vim style) |
| `q` | Quit |
| `Ctrl+C` | Force quit |

---

## Other Launch Options

```bash
# Web App only (React in browser)
./launch.sh web
# → Open http://localhost:5173

# TUI only (Terminal UI)
./launch.sh tui

# Both Web + TUI together
./launch.sh all

# API Server only (for custom integrations)
./launch.sh api
```

---

## Features of the Improved Script

### Automatic Port Cleanup ✅
- Kills any existing process using ports 3002 or 5173
- Prevents "port already in use" errors
- Graceful startup every time

### Automatic Verification ✅
- Checks that API server started successfully
- Tests health endpoint before starting TUI
- Shows clear status messages

### Automatic Cleanup ✅
- Press `Ctrl+C` to stop everything
- All background processes are terminated
- Ports are freed up immediately
- No orphaned processes left behind

---

## Troubleshooting

### "Port already in use" error

The script now handles this automatically! If you still get this error:

```bash
# Manual cleanup
lsof -ti:3002 | xargs kill -9
lsof -ti:5173 | xargs kill -9

# Then try again
./launch.sh tui
```

### TUI shows "Loading posts..." forever

Check if API server is working:
```bash
curl http://localhost:3002/health
```

Should return:
```json
{"status":"ok","cache_size":0}
```

### "go: command not found"

Install Go:
```bash
sudo pacman -S go
```

---

## What's Been Fixed

| Issue | Status | Solution |
|-------|--------|----------|
| CORS errors | ✅ Fixed | Uses API server proxy |
| Port conflicts | ✅ Fixed | Script auto-kills old processes |
| Multi-terminal complexity | ✅ Fixed | Single `./launch.sh tui` command |
| Go module issues | ✅ Fixed | Proper go.mod and go.sum |
| tea.Model type error | ✅ Fixed | Correct return type |
| TypeScript runtime errors | ✅ Fixed | Converted to JavaScript |

---

## Git Commits

Latest changes:
- ✅ Transform to monorepo architecture
- ✅ Fix Go module dependencies  
- ✅ Fix tea.Model return type
- ✅ Add JavaScript API server
- ✅ Add launch.sh script
- ✅ Update documentation
- ✅ Improve launch.sh port handling

---

## Summary

You now have a **fully functional, production-ready TUI application** that:

✅ Fetches real Reddit data
✅ Displays posts in a beautiful terminal UI
✅ Shares code with the React web app
✅ Has automatic port conflict handling
✅ Can be launched with a single command

**To launch:**
```bash
./launch.sh tui
```

That's it! Enjoy your RedditView TUI! 🎉

---

*Created: Feb 22, 2026*  
*Status: ✅ Production Ready*
