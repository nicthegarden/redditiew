# 📋 RedditView Installation Checklist

Follow this checklist to get RedditView running on Arch Linux.

## ✅ System Requirements

- [ ] Running Arch Linux or Arch-based distro (Manjaro, EndeavourOS, etc.)
- [ ] sudo access (can run sudo commands)
- [ ] Internet connection
- [ ] ~500MB free disk space

## ✅ Part 1: Install System Dependencies

### Node.js & npm

- [ ] Run: `sudo pacman -S nodejs npm`
- [ ] Verify: `node --version` (should be v16+)
- [ ] Verify: `npm --version` (should be v8+)

### Go (Required for TUI)

- [ ] Run: `sudo pacman -S go`
- [ ] Verify: `go version` (should be go1.21+)

### Optional: Build Tools

- [ ] Run: `sudo pacman -S base-devel`
  (contains gcc, make, and other tools)

### Optional: curl (for testing)

- [ ] Run: `sudo pacman -S curl`

## ✅ Part 2: Navigate to Project

- [ ] Open terminal
- [ ] Run: `cd /path/to/redditiew-local`
- [ ] Run: `pwd` to confirm you're in the right directory

## ✅ Part 3: Install npm Dependencies

- [ ] Run: `npm install`
- [ ] Wait for completion (this takes a few minutes)
- [ ] Check: No red ERR! messages
- [ ] Result: `node_modules/` folder created

## ✅ Part 4: Build Core Package

- [ ] Run: `npm run build`
- [ ] Check: Output shows "Successfully compiled"
- [ ] Result: `packages/core/dist/` folder created

## ✅ Part 5: Run the Application

Choose ONE option:

### Option A: Web App (Simplest)

- [ ] Run: `npm run dev`
- [ ] Wait for message: "Local: http://localhost:5173"
- [ ] Open browser: http://localhost:5173
- [ ] Check: Reddit posts load from r/sysadmin
- [ ] Click a post and verify comments load
- [ ] Stop with: `Ctrl+C`

### Option B: Go TUI App

- [ ] Open Terminal 1
- [ ] Run: `npm run dev:api`
- [ ] Wait for: "✓ API running on http://localhost:3002"
- [ ] Test: `curl http://localhost:3002/health` (should return OK)
- [ ] Open Terminal 2
- [ ] Run: `cd apps/tui && go run main.go`
- [ ] Check: TUI loads with Reddit posts
- [ ] Navigate with: ↑/↓ or j/k
- [ ] Quit with: q or Ctrl+C

### Option C: Both (Web + TUI)

- [ ] Open Terminal 1: `npm run dev:api`
- [ ] Open Terminal 2: `npm run dev`
- [ ] Open Terminal 3: `cd apps/tui && go run main.go`
- [ ] Test web app at http://localhost:5173
- [ ] Test TUI in Terminal 3
- [ ] Stop each with Ctrl+C

## ✅ Part 6: Troubleshooting

If something goes wrong, follow these steps:

### Web app won't start

```bash
npm run build
npm run dev
```

### TUI won't connect to API

```bash
# Check if API server is running
curl http://localhost:3002/health

# If not, restart it
npm run dev:api
```

### Port already in use

```bash
# Find what's using the port
lsof -i :5173  # for web app
lsof -i :3002  # for API

# Kill it
kill -9 <PID>
```

### npm install fails

```bash
# Clear cache and retry
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

### Go modules not found

```bash
cd apps/tui
go mod download
go run main.go
```

## ✅ Part 7: Verify Everything Works

Run this verification script:

```bash
echo "=== Checking Node.js ==="
node --version  # Should be v16+

echo "=== Checking npm ==="
npm --version   # Should be v8+

echo "=== Checking Go ==="
go version      # Should be go1.21+

echo "=== Checking project structure ==="
ls packages/core/dist/  # Should show compiled files
ls packages/web/src/    # Should show React source

echo "=== Checking if API responds ==="
npm run dev:api &
sleep 2
curl http://localhost:3002/health  # Should return JSON
kill %1

echo ""
echo "✅ All checks passed!"
```

## 🎯 Next Steps

Once everything works:

1. **Read the documentation:**
   - `SETUP_COMPLETE.md` - Project overview
   - `QUICK_REFERENCE.md` - Common commands
   - `TUI_SETUP_GUIDE.md` - TUI details
   - `ARCHITECTURE.md` - System design

2. **Try the apps:**
   - Web: `npm run dev` → http://localhost:5173
   - TUI: `npm run dev:api` + `go run main.go` in apps/tui/

3. **Explore the code:**
   - Web app: `packages/web/src/`
   - Core package: `packages/core/src/`
   - TUI app: `apps/tui/main.go`

4. **Make changes:**
   - Edit files
   - Rebuild core if needed: `npm run build`
   - Changes auto-reload in web app
   - Restart TUI for changes to take effect

## 🆘 Still Having Issues?

1. Re-read this checklist carefully
2. Check the error messages in detail
3. Try the troubleshooting section above
4. Look at relevant documentation files
5. Make sure all dependencies are installed

## 📞 Quick Command Reference

```bash
# Install everything
sudo pacman -S nodejs npm go base-devel

# Setup project
npm install && npm run build

# Run web app
npm run dev

# Run TUI
npm run dev:api  # Terminal 1
cd apps/tui && go run main.go  # Terminal 2

# Check if services are running
curl http://localhost:5173      # Web
curl http://localhost:3002/health  # API

# Clear and reinstall
rm -rf node_modules package-lock.json
npm install
npm run build
```

---

**When everything is checked and working, you're done!** 🎉
