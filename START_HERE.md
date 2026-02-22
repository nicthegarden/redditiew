# 🚀 START HERE - RedditView on Arch Linux

## Your Current Status

```
✓ Node.js v25.6.1 installed
✓ npm 11.10.1 installed
✗ Go NOT installed (REQUIRED)
```

## What You Need to Do (3 Steps)

### Step 1: Install Go (5 minutes)

Open a terminal and run:

```bash
sudo pacman -S go base-devel curl
```

Verify it worked:
```bash
go version
```

Should show something like: `go version go1.21.x linux/amd64`

### Step 2: Install Project Dependencies (5 minutes)

```bash
cd /path/to/redditiew-local
npm install
npm run build
```

### Step 3: Run the Application (choose one)

**Web App (easiest):**
```bash
npm run dev
```
Then open: http://localhost:5173

**Go TUI App:**
```bash
# Terminal 1
npm run dev:api

# Terminal 2
cd apps/tui
go run main.go
```

**Both Web + TUI:**
```bash
# Terminal 1
npm run dev:api

# Terminal 2
npm run dev

# Terminal 3
cd apps/tui && go run main.go
```

## Documentation Files

Read these in order:

1. **ARCH_QUICK_START.md** ← Start here for Arch Linux
2. **CHECKLIST.md** ← Follow this checklist
3. **QUICK_REFERENCE.md** ← Command reference
4. **SETUP_COMPLETE.md** ← Project overview

## Quick Command Reference

```bash
# Install all system packages
sudo pacman -S go base-devel curl

# Install and build project
npm install && npm run build

# Run web app
npm run dev
# → Open http://localhost:5173

# Run TUI
npm run dev:api      # Terminal 1
cd apps/tui && go run main.go  # Terminal 2

# Run both
npm run dev:api      # Terminal 1
npm run dev          # Terminal 2
cd apps/tui && go run main.go  # Terminal 3
```

## Troubleshooting

### "go: command not found"
```bash
sudo pacman -S go
```

### "npm install" fails
```bash
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

### Port already in use
```bash
# Find what's using the port
lsof -i :5173  # or :3002 for API

# Kill it
kill -9 <PID>
```

## Next Steps

1. Install Go: `sudo pacman -S go`
2. Navigate: `cd /path/to/redditiew-local`
3. Install: `npm install && npm run build`
4. Run: `npm run dev`
5. Open browser: `http://localhost:5173`

## That's It!

Once Go is installed, everything else is automatic. 

Happy coding! 🎉
