# RedditView - Arch Linux Installation Guide

This guide will help you install all dependencies and run RedditView on Arch Linux.

## System Requirements

- Arch Linux (or Arch-based distro like Manjaro, EndeavourOS)
- pacman package manager
- sudo access

## Prerequisites Installation

### 1. Update System (Recommended)

```bash
sudo pacman -Syu
```

### 2. Install Node.js and npm (if not already installed)

Check if installed:
```bash
node --version
npm --version
```

If not installed:
```bash
sudo pacman -S nodejs npm
```

Verify:
```bash
node --version  # Should be v16 or higher
npm --version   # Should be v8 or higher
```

### 3. Install Go (Required for TUI app)

```bash
sudo pacman -S go
```

Verify:
```bash
go version  # Should show Go 1.21 or higher
```

### 4. Install TypeScript (Optional but recommended)

```bash
sudo pacman -S typescript
```

Or via npm:
```bash
npm install -g typescript
```

### 5. Install Build Tools (Optional but recommended)

```bash
sudo pacman -S base-devel
```

This includes essential tools like `make`, `gcc`, etc.

## Complete Installation in One Command

If you want to install everything at once:

```bash
sudo pacman -S nodejs npm go typescript base-devel
```

## RedditView Setup

### Step 1: Navigate to Project

```bash
cd /path/to/redditiew-local
```

### Step 2: Install npm Dependencies

```bash
npm install
```

This installs:
- React and React DOM
- Vite (dev server)
- TypeScript
- ESLint
- And all other npm packages

### Step 3: Build Core Package

```bash
npm run build
```

This compiles the TypeScript core package to JavaScript.

### Step 4: Run the Application

#### Option A: Run Web App Only

```bash
npm run dev
```

Then open: http://localhost:5173

#### Option B: Run TUI Only

Terminal 1 - Start API server:
```bash
npm run dev:api
```

Terminal 2 - Start TUI app:
```bash
cd apps/tui
go run main.go
```

#### Option C: Run Everything (Web + TUI)

Terminal 1 - API server:
```bash
npm run dev:api
```

Terminal 2 - Web app:
```bash
npm run dev
# Visit http://localhost:5173
```

Terminal 3 - TUI app:
```bash
cd apps/tui
go run main.go
```

## Verification Checklist

- [ ] `node --version` shows v16+
- [ ] `npm --version` shows v8+
- [ ] `go version` shows go1.21+
- [ ] `npm install` completes without errors
- [ ] `npm run build` completes without errors
- [ ] `npm run dev` starts on port 5173
- [ ] Browser can access http://localhost:5173
- [ ] `npm run dev:api` starts on port 3002
- [ ] `cd apps/tui && go run main.go` runs without errors

## Troubleshooting

### "command not found: node"

Node.js not installed:
```bash
sudo pacman -S nodejs npm
```

### "command not found: go"

Go not installed:
```bash
sudo pacman -S go
```

### "npm install" fails with permission errors

```bash
npm config set prefix ~/.npm-global
export PATH=~/.npm-global/bin:$PATH
```

### Port 3000, 5173, or 3002 already in use

Find what's using the port:
```bash
lsof -i :PORT_NUMBER  # May need to install lsof first

# If not installed:
sudo pacman -S lsof
```

Kill the process:
```bash
kill -9 PID
```

Or use different ports by editing config files.

### "EACCES: permission denied" when installing globally

```bash
sudo pacman -S npm
npm config set prefix ~/.npm-global
```

### Build fails with TypeScript errors

Make sure TypeScript is installed:
```bash
npm install -g typescript
# or
sudo pacman -S typescript
```

Then try again:
```bash
npm run build
```

### Go dependencies missing when running TUI

```bash
cd apps/tui
go mod download
go run main.go
```

## Optional: Install Additional Tools

### For better development experience:

```bash
# Install curl (useful for testing API)
sudo pacman -S curl

# Install git (if not already installed)
sudo pacman -S git

# Install code editor (optional, pick one)
sudo pacman -S code          # Visual Studio Code
# or
sudo pacman -S neovim        # Neovim
# or
sudo pacman -S vim           # Vim
```

### For building native modules (if needed):

```bash
sudo pacman -S python3
```

## Environment Setup

### Set up npm global packages directory (Recommended)

```bash
mkdir -p ~/.npm-global
npm config set prefix ~/.npm-global
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

Or for zsh:
```bash
mkdir -p ~/.npm-global
npm config set prefix ~/.npm-global
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.zshrc
source ~/.zshrc
```

## Next Steps

1. Install all prerequisites: `sudo pacman -S nodejs npm go`
2. Navigate to project: `cd /path/to/redditiew-local`
3. Install dependencies: `npm install`
4. Build core: `npm run build`
5. Run app: `npm run dev` (web) or `npm run dev:api && go run main.go` (TUI)

## Getting Help

If you encounter any issues:

1. Check the error message carefully
2. Try rebuilding: `npm run build`
3. Clear cache: `rm -rf node_modules && npm install`
4. Check documentation in project root:
   - `SETUP_COMPLETE.md`
   - `QUICK_REFERENCE.md`
   - `TUI_SETUP_GUIDE.md`

## Useful Links

- Node.js: https://nodejs.org/
- npm: https://www.npmjs.com/
- Go: https://golang.org/
- Arch Linux Wiki: https://wiki.archlinux.org/

---

**Once everything is installed, follow the Quick Start commands above!**
