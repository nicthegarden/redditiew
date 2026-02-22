# QuickStart Guide - RedditView

Get RedditView up and running in minutes! Choose your operating system below.

## Table of Contents
- [Windows QuickStart](#windows-quickstart) ⏱️ 5 minutes
- [Linux QuickStart](#linux-quickstart) ⏱️ 5 minutes
- [Verify Installation](#verify-installation)
- [Troubleshooting](#troubleshooting)

---

## Windows QuickStart

### Step 1: Install Prerequisites

#### Option A: Using Chocolatey (Recommended)
If you have Chocolatey installed:
```powershell
choco install golang nodejs
```

#### Option B: Manual Installation
1. **Go**
   - Download from https://golang.org/dl
   - Download `go1.21.x.windows-amd64.msi` (or newer)
   - Run the installer, follow prompts
   - Choose default installation path (`C:\Program Files\Go`)

2. **Node.js**
   - Download from https://nodejs.org
   - Download the **LTS version** (Recommended)
   - Run the installer, choose default settings
   - Restart your computer

#### Verify Installation
Open **PowerShell** or **Command Prompt** and run:
```powershell
go version
node --version
npm --version
```

Expected output:
```
go version go1.21.x windows/amd64
v18.x.x or higher
9.x.x or higher
```

### Step 2: Clone the Repository

```powershell
# Open PowerShell as Administrator (optional but recommended)

# Navigate to where you want the project
cd C:\Users\YourUsername\Projects

# Clone the repository
git clone https://github.com/yourusername/redditiew-local.git

# Enter the directory
cd redditiew-local
```

### Step 3: Install Dependencies

```powershell
# Install Node.js dependencies
npm install

# This installs:
# - Express.js (API server)
# - Axios (HTTP client)
# - And other required packages
```

⏳ This may take 1-2 minutes. Wait for it to complete.

### Step 4: Build the Application

```powershell
# Build the TUI application
cd apps\tui
go build -o redditview.exe .

# Verify the build succeeded
if (Test-Path redditview.exe) { 
    Write-Host "✓ Build successful!" -ForegroundColor Green 
} else { 
    Write-Host "✗ Build failed!" -ForegroundColor Red 
}

# Return to project root
cd ..\..
```

### Step 5: Run the Application

#### Terminal 1: Start the API Server
```powershell
# From project root
npm start
```

Expected output:
```
API server listening on port 3002
Ready to serve requests
```

#### Terminal 2: Start the TUI
```powershell
# Open a NEW PowerShell window
# Navigate to the project directory
cd C:\Users\YourUsername\Projects\redditiew-local

# Run the TUI
.\apps\tui\redditview.exe
```

🎉 **Done! You should see the RedditView TUI load!**

### 🖼️ What You'll See

```
╔════════════════════════════════════════════════════════════════╗
║                      REDDITVIEW - SYSADMIN                      ║
╠════════════════════════════════════════════════════════════════╣
║ Post 1/200                                                      ║
║ ⬆ 5.2K  💬 234  [Linux Kernel 6.8 Released]                    ║
║         u/torvalds  ·  3 hours ago                              ║
║ ─────────────────────────────────────────────────────────────  ║
║ [Post details showing here]                                     ║
║                                                                  ║
╚════════════════════════════════════════════════════════════════╝
```

---

## Linux QuickStart

### Step 1: Install Prerequisites

Choose your distribution below:

#### Ubuntu/Debian
```bash
# Update package list
sudo apt update

# Install Go
sudo apt install -y golang-go

# Install Node.js (using NodeSource repository for latest LTS)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Verify installation
go version
node --version
npm --version
```

#### Fedora/RHEL/CentOS
```bash
# Install Go
sudo dnf install -y golang

# Install Node.js
sudo dnf install -y nodejs npm

# Verify installation
go version
node --version
npm --version
```

#### Arch Linux
```bash
# Install Go
sudo pacman -S go

# Install Node.js
sudo pacman -S nodejs npm

# Verify installation
go version
node --version
npm --version
```

### Step 2: Clone the Repository

```bash
# Navigate to where you want the project
mkdir -p ~/Projects
cd ~/Projects

# Clone the repository
git clone https://github.com/yourusername/redditiew-local.git

# Enter the directory
cd redditiew-local
```

### Step 3: Install Dependencies

```bash
# Install Node.js dependencies
npm install

# This installs:
# - Express.js (API server)
# - Axios (HTTP client)
# - And other required packages
```

⏳ This may take 1-2 minutes. Wait for it to complete.

### Step 4: Build the Application

```bash
# Navigate to TUI directory
cd apps/tui

# Build the application
go build -o redditview .

# Verify the build
ls -lh redditview

# Expected output: -rwxr-xr-x redditview (11M)

# Return to project root
cd ../..
```

### Step 5: Run the Application

#### Terminal 1: Start the API Server
```bash
# From project root
npm start
```

Expected output:
```
API server listening on port 3002
Ready to serve requests
```

#### Terminal 2: Start the TUI
```bash
# Open a NEW terminal window
# Navigate to the project directory
cd ~/Projects/redditiew-local

# Run the TUI
./apps/tui/redditview
```

🎉 **Done! You should see the RedditView TUI load!**

---

## Verify Installation

### Check All Components Are Running

**Test 1: API Server is Responding**
```bash
# In a new terminal, run:
curl -s http://localhost:3002/api/r/sysadmin.json | head -20
```

Expected: JSON data with posts (not an error)

**Test 2: TUI Loads Successfully**
```bash
# In the TUI terminal, verify you see:
# - Post list with items
# - Subreddit name "sysadmin"
# - Keybinding footer at bottom
```

**Test 3: Basic Navigation**
- Press `j` to move down one post
- Press `k` to move up one post
- Press `q` to quit (return to terminal)

If all tests pass, installation is complete! 🎉

---

## First Time Usage

### Basic Operations

1. **View a Post**
   - Use arrow keys or `j`/`k` to select a post
   - Press `Enter` to view details

2. **View Comments**
   - While in detail view, press `c`
   - Use arrow keys to scroll through comments
   - Press `Esc` to close comments

3. **Open in Browser**
   - Press `w` to open the post in your default browser

4. **Search Posts**
   - Press `Ctrl+F` to open search
   - Type your search query
   - Press `Escape` to close search

5. **Change Subreddit**
   - Press `s` to open subreddit selector
   - Type the subreddit name (without r/)
   - Press `Enter` to load

### Full Keybinding Reference

See [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md) for complete keyboard shortcuts.

---

## Configuration

You can customize RedditView by editing `config.json`:

```json
{
  "tui": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 200,
    "list_height": 10,
    "max_title_length": 80
  },
  "api": {
    "base_url": "http://localhost:3002/api",
    "timeout_seconds": 10
  }
}
```

**Common Customizations:**

Change default subreddit:
```json
"default_subreddit": "programming"  // or "pics", "news", etc.
```

Adjust posts per page (higher = more posts, slower):
```json
"posts_per_page": 100  // or 50, 300, etc.
```

See [CONFIGURATION.md](CONFIGURATION.md) for all options.

---

## Troubleshooting

### "command not found: go" or "command not found: node"

**Windows:**
- Go to Settings → Environment Variables
- Verify `C:\Program Files\Go\bin` is in PATH
- Restart PowerShell

**Linux:**
```bash
# Check if installed
which go
which node

# If not found, install using your package manager (see Step 1)
```

### "Cannot connect to API server"

**Check if server is running:**
```bash
# Windows
netstat -ano | findstr :3002

# Linux
lsof -i :3002
```

**If not running:**
```bash
# Kill any existing process on port 3002
# Then start the server again
npm start
```

### "TUI looks broken / characters are garbled"

This usually means terminal is too small. 
- **Minimum size:** 80 columns × 24 rows
- **Recommended:** 120 columns × 40 rows

Resize your terminal and try again.

### Build fails with "go: command not found"

Reinstall Go:
- **Windows:** Download installer from https://golang.org/dl and run
- **Linux:** Run the install command for your distro (see Step 1)

### "Port 3002 already in use"

Another process is using the port. Options:
```bash
# Option 1: Kill the process using port 3002
# Windows:
netstat -ano | findstr :3002
taskkill /PID <PID> /F

# Linux:
lsof -i :3002
kill -9 <PID>

# Option 2: Change the port in config.json
# Edit "base_url" in config.json and restart
```

### "npm ERR! code ERESOLVE"

This is a dependency conflict. Try:
```bash
# Clear npm cache
npm cache clean --force

# Delete package lock file
rm package-lock.json

# Reinstall
npm install
```

---

## Next Steps

1. **Explore the TUI:**
   - Try different subreddits with `s` key
   - Read comments with `c` key
   - Search with `Ctrl+F`

2. **Customize Configuration:**
   - Edit `config.json` to your preferences
   - See [CONFIGURATION.md](CONFIGURATION.md) for all options

3. **Learn All Keybindings:**
   - Read [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md)
   - Become a power user!

4. **Report Issues:**
   - Found a bug? Create an issue on GitHub
   - Have a feature idea? Start a discussion

---

## Getting Help

- 📖 **Full Documentation:** See [README.md](README.md)
- ⚙️ **Configuration:** See [CONFIGURATION.md](CONFIGURATION.md)
- ⌨️ **Keybindings:** See [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md)
- 🏗️ **Architecture:** See [ARCHITECTURE.md](ARCHITECTURE.md)
- 👨‍💻 **Development:** See [DEVELOPMENT.md](DEVELOPMENT.md)

---

**Happy browsing! 🚀** [Back to README](README.md)
