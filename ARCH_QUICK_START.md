# 🚀 RedditView on Arch Linux - Quick Setup

## Step 1: Install Missing Dependencies

Open a terminal and run these commands:

### Install Go (Required for TUI)
```bash
sudo pacman -S go
```

### Optional: Install Build Tools
```bash
sudo pacman -S base-devel
```

## Step 2: Verify Installation

Check that everything is installed:
```bash
node --version    # Should show v18+
npm --version     # Should show v8+
go version        # Should show go1.21+
```

## Step 3: Setup RedditView

Navigate to the project directory:
```bash
cd /path/to/redditiew-local
```

Install npm dependencies:
```bash
npm install
```

Build the core package:
```bash
npm run build
```

## Step 4: Run the Application

### Option A: Run Web App (Easiest)

```bash
npm run dev
```

Then open your browser to: **http://localhost:5173**

### Option B: Run TUI App

Open two terminals:

**Terminal 1** - Start the API server:
```bash
npm run dev:api
```

**Terminal 2** - Run the TUI:
```bash
cd apps/tui
go run main.go
```

### Option C: Run Everything

Open three terminals:

**Terminal 1** - API server:
```bash
npm run dev:api
```

**Terminal 2** - Web app:
```bash
npm run dev
# Visit http://localhost:5173
```

**Terminal 3** - TUI app:
```bash
cd apps/tui
go run main.go
```

## Troubleshooting

### "go: command not found"

Go is not installed. Run:
```bash
sudo pacman -S go
```

### "Port 5173 already in use"

Kill the process using the port:
```bash
lsof -i :5173    # Find the process
kill -9 PID      # Replace PID with the number
```

Or install lsof first:
```bash
sudo pacman -S lsof
```

### "npm install" fails

Try clearing npm cache:
```bash
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

### "npm ERR! code EACCES"

Fix npm permissions:
```bash
mkdir -p ~/.npm-global
npm config set prefix '~/.npm-global'
export PATH=~/.npm-global/bin:$PATH
```

Add to your shell profile (~/.bashrc or ~/.zshrc):
```bash
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

## Full Install Script

If you want to automate everything, run:

```bash
bash setup-arch.sh
```

This script will:
1. Update system packages
2. Install Node.js and npm (if not installed)
3. Install Go (if not installed)
4. Install build tools (optional)

## What You Need

| Tool | Version | Package | How to Install |
|------|---------|---------|----------------|
| Node.js | v16+ | nodejs | `sudo pacman -S nodejs` |
| npm | v8+ | npm | `sudo pacman -S npm` |
| Go | v1.21+ | go | `sudo pacman -S go` |
| Python3 | 3.8+ | python | `sudo pacman -S python` (optional) |
| Build Tools | - | base-devel | `sudo pacman -S base-devel` (optional) |

## Verify Everything Works

After installation, run:

```bash
# Check versions
node --version
npm --version
go version

# Navigate to project
cd /path/to/redditiew-local

# Install and build
npm install
npm run build

# Test web app
npm run dev
# Visit http://localhost:5173 - should load without errors
```

## Next: Check the Docs

Once everything is running, read the documentation:

1. **SETUP_COMPLETE.md** - Project overview
2. **QUICK_REFERENCE.md** - Commands and checklists
3. **TUI_SETUP_GUIDE.md** - TUI specific info
4. **ARCHITECTURE.md** - System design

## Common Issues & Fixes

| Issue | Fix |
|-------|-----|
| `go: command not found` | `sudo pacman -S go` |
| Port already in use | `kill -9 PID` (find PID with `lsof -i :PORT`) |
| npm install fails | `npm cache clean --force && npm install` |
| Permission denied | Fix npm permissions (see above) |
| TypeScript errors | `npm run build` again (or `npm install -g typescript`) |

---

**You're all set! Enjoy using RedditView!** 🎉
