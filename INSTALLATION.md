# Installation Guide

Detailed technical installation instructions for RedditView.

## 🖼️ Application Preview

See what RedditView looks like once installed:

**TUI Application - Post List View:**
![RedditView TUI](TUI.png)

**TUI Application - Comments View:**
![RedditView Comments](TUI-Comment.png)

**Web Interface:**
![RedditView Web UI](WebUI.png)

---

## Table of Contents
- [Prerequisites](#prerequisites)
- [Windows Installation](#windows-installation)
- [Linux Installation](#linux-installation)
- [macOS Installation](#macos-installation)
- [Docker Installation](#docker-installation)
- [Building from Source](#building-from-source)
- [Systemd Service Setup](#systemd-service-setup)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### System Requirements
- **Minimum RAM:** 512 MB
- **Minimum Disk Space:** 500 MB
- **Terminal:** 80×24 characters minimum
- **Network:** Internet connection for Reddit API access

### Software Requirements
- **Go:** 1.19+ (1.21+ recommended)
- **Node.js:** 16+ (18 LTS recommended)
- **npm:** 7+ (9+ recommended)
- **Git:** Any recent version

### Optional
- **Docker:** For containerized deployment
- **Make:** For build automation (optional)

---

## Windows Installation

### Step 1: Install Go

1. Visit https://golang.org/dl
2. Download **go1.21.x.windows-amd64.msi** (or latest)
3. Run the installer
4. Accept default path: `C:\Program Files\Go`
5. Click "Finish"

**Verify Installation:**
```powershell
go version
# Output: go version go1.21.x windows/amd64
```

### Step 2: Install Node.js

1. Visit https://nodejs.org
2. Download **LTS version** (recommended over Current)
3. Run the installer
4. Accept default settings
5. Allow PATH modification when asked
6. Restart your computer

**Verify Installation:**
```powershell
node --version
npm --version
```

### Step 3: Clone Repository

```powershell
# Choose a location for the project
cd $env:USERPROFILE\Projects
mkdir redditiew-local
cd redditiew-local

# Clone the repository
git clone https://github.com/yourusername/redditiew-local.git .
```

### Step 4: Install Dependencies

```powershell
# Install Node.js packages
npm install

# This installs:
# - express (API server)
# - axios (HTTP client)
# - And other dependencies
```

### Step 5: Build TUI

```powershell
# Navigate to TUI directory
cd apps\tui

# Build the executable
go build -o redditview.exe .

# Verify build succeeded
if (Test-Path redditview.exe) { 
    Get-Item redditview.exe | Select-Object Name, Length
}

# Expected output:
# Name           Length
# ----           ------
# redditview.exe 11000000

# Return to project root
cd ..\..
```

### Step 6: Configure

```powershell
# Edit configuration (optional)
# Open config.json in your editor and customize:
notepad config.json
```

### Step 7: Run Application

**Terminal 1: Start API Server**
```powershell
npm start
```

**Terminal 2: Start TUI**
```powershell
.\apps\tui\redditview.exe
```

---

## Linux Installation

### Ubuntu/Debian

#### Install Dependencies
```bash
# Update package manager
sudo apt update

# Install Go
sudo apt install -y golang-go

# Install Node.js from NodeSource repository
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Verify installation
go version
node --version
npm --version
```

#### Clone and Build
```bash
# Clone repository
mkdir -p ~/projects
cd ~/projects
git clone https://github.com/yourusername/redditiew-local.git
cd redditiew-local

# Install dependencies
npm install

# Build TUI
cd apps/tui
go build -o redditview .

# Verify binary
ls -lh redditview

# Return to root
cd ../..
```

#### Run
```bash
# Terminal 1: API Server
npm start

# Terminal 2: TUI
./apps/tui/redditview
```

### Fedora/RHEL/CentOS

#### Install Dependencies
```bash
# Install Go
sudo dnf install -y golang

# Install Node.js
sudo dnf install -y nodejs npm

# Verify
go version
node --version
```

#### Build and Run
(Same as Ubuntu - see above)

### Arch Linux

#### Install Dependencies
```bash
# Install packages
sudo pacman -S go nodejs npm

# Verify
go version
node --version
```

#### Build and Run
(Same as Ubuntu - see above)

---

## macOS Installation

### Using Homebrew (Recommended)

#### Install Dependencies
```bash
# Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Go and Node.js
brew install go node

# Verify
go version
node --version
npm --version
```

### Manual Installation

#### Install Go
1. Visit https://golang.org/dl
2. Download `go1.21.x.darwin-amd64.pkg`
3. Run installer and follow prompts

#### Install Node.js
1. Visit https://nodejs.org
2. Download LTS version
3. Run installer

#### Clone and Build
```bash
# Clone repository
mkdir -p ~/projects
cd ~/projects
git clone https://github.com/yourusername/redditiew-local.git
cd redditiew-local

# Install dependencies
npm install

# Build TUI
cd apps/tui
go build -o redditview .

# Return to root
cd ../..
```

#### Run
```bash
# Terminal 1: API Server
npm start

# Terminal 2: TUI (may need to allow execution)
chmod +x ./apps/tui/redditview
./apps/tui/redditview
```

---

## Docker Installation

### Using Docker Compose

#### Prerequisites
- Docker installed (https://www.docker.com/products/docker-desktop)
- Docker Compose installed

#### Create docker-compose.yml
```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "3002:3002"
    environment:
      - NODE_ENV=production
    command: npm start

  tui:
    build:
      context: .
      dockerfile: Dockerfile.tui
    depends_on:
      - api
    environment:
      - API_BASE_URL=http://api:3002/api
    stdin_open: true
    tty: true
```

#### Build and Run
```bash
# Build images
docker-compose build

# Run services
docker-compose up

# Access TUI
docker-compose exec tui /app/redditview
```

### Using Docker CLI

#### Build Images
```bash
# API server
docker build -t redditview-api .

# TUI
docker build -f Dockerfile.tui -t redditview-tui .
```

#### Run Containers
```bash
# Start API
docker run -p 3002:3002 redditview-api

# Start TUI (in another terminal)
docker run -it --network host redditview-tui
```

---

## Building from Source

### Understanding the Build Process

**Project Structure:**
```
redditiew-local/
├── apps/tui/main.go         # TUI source code
├── api-server.js            # API server
├── config.json              # Configuration
└── package.json             # Dependencies
```

### Build Steps

#### 1. Verify Prerequisites
```bash
go version          # Should be 1.19+
node --version      # Should be 16+
npm --version       # Should be 7+
```

#### 2. Clone Repository
```bash
git clone https://github.com/yourusername/redditiew-local.git
cd redditiew-local
```

#### 3. Install Node Dependencies
```bash
npm install
```

This installs:
- `express` - Web framework for API server
- `axios` - HTTP client for making requests
- `cors` - CORS middleware for API
- Development dependencies for building

#### 4. Build TUI

**Linux/macOS:**
```bash
cd apps/tui
go build -o redditview .
```

**Windows:**
```powershell
cd apps\tui
go build -o redditview.exe .
```

**With custom flags:**
```bash
# For smaller binary
go build -ldflags="-s -w" -o redditview apps/tui/main.go

# With version info
go build -ldflags="-X main.Version=1.0.0" -o redditview apps/tui/main.go
```

#### 5. Configure (Optional)
```bash
# Create custom config if needed
cp config.json config.json.backup
# Edit config.json with your settings
```

#### 6. Run

**Option A: Separate Terminals**
```bash
# Terminal 1: API Server
npm start

# Terminal 2: TUI
./apps/tui/redditview
```

**Option B: Background Process (Linux/macOS)**
```bash
# Start API in background
npm start &

# Get the PID
API_PID=$!

# Start TUI
./apps/tui/redditview

# Stop API when done
kill $API_PID
```

---

## Verification

### Test Installation

**1. Check Go Installation**
```bash
go env GOROOT
go env GOPATH
```

**2. Check Node Installation**
```bash
npm list -g --depth=0
```

**3. Test API Server**
```bash
# In project root:
npm start

# In another terminal:
curl http://localhost:3002/api/r/sysadmin.json | jq '. | keys' | head -10
```

**4. Test TUI**
```bash
# Run the binary
./apps/tui/redditview

# If it loads and shows posts, installation is successful!
# Press 'q' to exit
```

---

## Systemd Service Setup

For production deployments and automatic startup on boot, RedditView can be installed as a systemd service.

**See [SYSTEMD_SETUP.md](SYSTEMD_SETUP.md) for complete systemd integration guide:**

- Run as user-level systemd service
- Automatic startup on system boot
- Three deployment modes (API+TUI, API-only, Web-only)
- Auto-restart on crashes with retry logic
- Centralized logging via journalctl
- Easy service management with systemctl

**Quick setup:**
```bash
cd /path/to/redditiew-local
./setup.sh --mode both --enable --start
```

This will install, enable, and start RedditView as systemd services in one command.

---

## Troubleshooting

### "go: command not found"

**Windows:**
- Add `C:\Program Files\Go\bin` to PATH
- Restart PowerShell

**Linux/macOS:**
```bash
# Install Go via package manager
# Ubuntu: sudo apt install golang-go
# macOS: brew install go
# Arch: sudo pacman -S go

# Or add to PATH if installed manually:
export PATH=$PATH:/usr/local/go/bin
```

### "npm: command not found"

**Solution:**
1. Verify Node.js is installed: `node --version`
2. Reinstall Node.js from https://nodejs.org
3. Restart terminal after installation

### Build fails with module errors

```bash
# Clear cache
go clean -cache -modcache

# Tidy dependencies
go mod tidy

# Rebuild
go build -o redditview .
```

### "Port 3002 already in use"

**Windows:**
```powershell
# Find process using port 3002
netstat -ano | findstr :3002

# Kill the process (replace PID with actual number)
taskkill /PID <PID> /F
```

**Linux/macOS:**
```bash
# Find process
lsof -i :3002

# Kill the process
kill -9 <PID>
```

### API server won't start

```bash
# Check Node.js installation
node --version

# Verify npm packages installed
npm list

# Try clearing npm cache
npm cache clean --force

# Reinstall
rm -rf node_modules package-lock.json
npm install

# Try starting again
npm start
```

### TUI won't load

**Check API is running:**
```bash
curl http://localhost:3002/api/r/sysadmin.json
```

If this fails:
1. Start API server: `npm start`
2. Wait 2-3 seconds for server to start
3. Try curl again
4. Then start TUI

### Terminal display issues

**Symptoms:** Garbled characters, broken layout

**Solutions:**
1. Maximize terminal window (minimum 80×24)
2. Try different terminal application
3. Check terminal supports colors: `echo $TERM` (should show xterm-256color or similar)
4. Rebuild binary with recent Go version

### Permission denied (Linux/macOS)

```bash
# Make binary executable
chmod +x apps/tui/redditview

# Run it
./apps/tui/redditview
```

---

## Advanced Configuration

### Custom Build Output

**Optimize for size:**
```bash
go build -ldflags="-s -w" -o redditview .
# Result: ~6MB instead of 11MB
```

**Static binary (no dependencies):**
```bash
CGO_ENABLED=0 go build -o redditview .
```

**For specific OS:**
```bash
# Build for Linux on macOS
GOOS=linux GOARCH=amd64 go build -o redditview-linux .

# Build for Windows on Linux
GOOS=windows GOARCH=amd64 go build -o redditview.exe .
```

### Development Build

```bash
# With debug symbols (larger but better debugging)
go build -gcflags="all=-N -l" -o redditview .

# Run with race detector
go run -race ./apps/tui/main.go
```

---

## Performance Optimization

### Compile-Time Optimization

```bash
# Release build (optimized)
go build -ldflags="-s -w" -trimpath -o redditview .
```

### Runtime Optimization

In `config.json`:
```json
{
  "tui": {
    "posts_per_page": 100  # Lower = faster load
  },
  "api": {
    "timeout_seconds": 8   # Lower = faster feedback
  }
}
```

---

## See Also

- 🚀 [QUICKSTART.md](QUICKSTART.md) - Fast setup guide
- ⚙️ [CONFIGURATION.md](CONFIGURATION.md) - Configuration options
- ⌨️ [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md) - Keyboard shortcuts
- 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture

---

**Installation complete! See [QUICKSTART.md](QUICKSTART.md) for next steps.**
