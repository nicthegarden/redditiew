# Systemd Service Setup Guide

Run RedditView as a systemd service with automatic startup and management capabilities. This guide covers installing, configuring, and managing RedditView using systemd with both user-level and system-level installation options, including optional tmux support for the TUI.

## 📋 Table of Contents

1. [Overview](#overview)
2. [Installation Scopes](#installation-scopes)
3. [Quick Start](#quick-start)
4. [Installation Modes](#installation-modes)
5. [Detailed Setup](#detailed-setup)
6. [Service Management](#service-management)
7. [Troubleshooting](#troubleshooting)
8. [Advanced Configuration](#advanced-configuration)
9. [Uninstallation](#uninstallation)

---

## 🎯 Overview

RedditView can be installed as systemd services with two different installation scopes and three service modes:

### Installation Scopes

| Scope | Location | Root Required | Best For | Auto-Start |
|-------|----------|---------------|----------|-----------|
| **User** | `~/.config/systemd/user/` | No | Development, desktops, personal use | User login required |
| **System** | `/etc/systemd/system/` | Yes (sudo) | Servers, production, always-on | System boot (no login needed) |

### Service Modes

| Mode | Components | Best For | Requirements |
|------|-----------|----------|--------------|
| **Both** (Default) | API Server + TUI (with tmux) | Full feature access | tmux, Node.js, Go |
| **API Only** | API Server only | Headless servers, remote access | Node.js |
| **Web Only** | Web interface | Browser-only access | Node.js |

### Benefits of Systemd Integration

- ✅ **Automatic startup** on system boot (configurable by scope)
- ✅ **Auto-restart** on crashes with configurable retry logic
- ✅ **Centralized logging** via journalctl
- ✅ **Resource management** and isolation
- ✅ **Easy service control** with standard systemctl commands
- ✅ **Both user and system-level installation** support
- ✅ **Dependency management** between services

---

## 📌 Installation Scopes

### User-Level Installation

**Characteristics:**
- Services stored in: `~/.config/systemd/user/`
- No root/sudo required
- Services run as current user
- Requires user to be logged in (or user lingering enabled)
- Perfect for: development, desktops, personal machines

**Pros:**
- No root privileges needed
- Isolated to single user
- Easy to manage multiple instances
- No system-wide impact

**Cons:**
- Stops when user logs out (unless lingering enabled)
- Only one user can run services
- Not suitable for always-on servers

**Installation:**
```bash
./setup.sh --scope user  # Uses $HOME/.config/systemd/user/
```

**Management (no sudo needed):**
```bash
systemctl --user status redditview-api
systemctl --user start redditview-api
journalctl --user -u redditview-api -f
```

---

### System-Level Installation

**Characteristics:**
- Services stored in: `/etc/systemd/system/`
- Requires root/sudo for installation
- Services run as specified user (e.g., `redditview`)
- Always available (no login required)
- Perfect for: servers, production, system services

**Pros:**
- Runs on system boot (no user login needed)
- Always available and monitored
- Standard system service management
- Multiple instances possible with different users
- Suitable for production deployments

**Cons:**
- Requires root/sudo for installation
- System-wide installation
- Affects all users

**Installation:**
```bash
sudo ./setup.sh --scope system --user redditview
```

**Management (requires sudo):**
```bash
sudo systemctl status redditview-api
sudo systemctl start redditview-api
sudo journalctl -u redditview-api -f
```

---

## 🚀 Quick Start

### Option 1: User-Level (Recommended for Personal Use)

```bash
cd /path/to/redditiew-local

# Interactive setup
./setup.sh

# Or automated
./setup.sh --scope user --mode both --enable --start
```

### Option 2: System-Level (Recommended for Servers)

```bash
cd /path/to/redditiew-local

# Interactive setup (with sudo)
sudo ./setup.sh

# Or automated (with sudo)
sudo ./setup.sh --scope system --mode both --user redditview --enable
```

### Option 3: Custom Configuration

```bash
# User-level with custom path
./setup.sh --scope user --path /opt/redditiew --enable --start

# System-level with custom user
sudo ./setup.sh --scope system --path /opt/redditiew --user reddit-service --enable

# API-only system service
sudo ./setup.sh --scope system --mode api-only --user redditview --enable

# Verbose output for debugging
./setup.sh --scope user --verbose
```


## 📦 Installation Modes

The following modes can be used with either user-level or system-level scope:

### Mode 1: Both (API Server + TUI)

**Install both the API server and Terminal UI using tmux.**

```bash
# User-level
./setup.sh --scope user --mode both

# System-level
sudo ./setup.sh --scope system --mode both --user redditview
```

**Creates:**
- API server service
- TUI service in tmux session named `redditview`

**Behavior:**
- TUI service depends on API service (auto-start in order)
- TUI runs in persistent tmux session
- Automatic restart on crashes
- Logs available via journalctl

**Access:**
```bash
# Connect to TUI session
tmux attach-session -t redditview

# Detach from session (keep running)
Ctrl+B then D

# View logs (user-level)
journalctl --user -u redditview-tui -f

# View logs (system-level)
sudo journalctl -u redditview-tui -f
```

**Perfect for:**
- Development systems
- Personal desktops with display
- Full feature access
- Testing both components together

---

### Mode 2: API Only

**Install only the API server (headless).**

```bash
# User-level
./setup.sh --scope user --mode api-only

# System-level
sudo ./setup.sh --scope system --mode api-only --user redditview
```

**Creates:**
- API server service only

**Behavior:**
- Lightweight, no GUI components
- Perfect for servers
- Can be accessed via web interface on remote systems
- Standard Node.js service

**Access:**
```bash
# API on http://localhost:3002
# Web UI on http://localhost:3000

# View logs (user-level)
journalctl --user -u redditview-api -f

# View logs (system-level)
sudo journalctl -u redditview-api -f
```

**Perfect for:**
- Server deployments
- Headless systems
- Docker environments
- Remote access via web
- Always-on systems

---

### Mode 3: Web Only

**Install only the web interface.**

```bash
# User-level
./setup.sh --scope user --mode web-only

# System-level
sudo ./setup.sh --scope system --mode web-only --user redditview
```

**Creates:**
- Web service

**Behavior:**
- Same as API-only but explicitly for web-only use
- All functionality via browser interface
- Lightweight footprint

**Access:**
```bash
# Web UI on http://localhost:3000

# View logs (user-level)
journalctl --user -u redditview-web -f

# View logs (system-level)
sudo journalctl -u redditview-web -f
```

**Perfect for:**
- Web-only deployments
- Browser-based access
- Remote servers
- Minimal resource usage

---

### Mode 1: Both (API Server + TUI)

**Install both the API server and Terminal UI using tmux.**

```bash
./setup.sh --mode both
```

**Creates:**
- `~/.config/systemd/user/redditview-api.service` - API server
- `~/.config/systemd/user/redditview-tui.service` - TUI in tmux session

**Behavior:**
- TUI service depends on API service (auto-start in order)
- TUI runs in persistent tmux session named `redditview`
- Automatic restart on crashes
- Logs available via journalctl

**Access:**
```bash
# Connect to TUI session
tmux attach-session -t redditview

# Detach from session (keep running)
Ctrl+B then D

# View logs
journalctl --user -u redditview-tui -f
```

**Perfect for:**
- Development systems
- Personal desktops with display
- Full feature access

---

### Mode 2: API Only

**Install only the API server (headless).**

```bash
./setup.sh --mode api-only
```

**Creates:**
- `~/.config/systemd/user/redditview-api.service` - API server only

**Behavior:**
- Lightweight, no GUI components
- Perfect for servers
- Can be accessed via web interface on remote systems
- Standard Node.js service

**Access:**
```bash
# API on http://localhost:3002
# Web UI on http://localhost:3000

# View logs
journalctl --user -u redditview-api -f
```

**Perfect for:**
- Server deployments
- Headless systems
- Docker environments
- Remote access via web

---

### Mode 3: Web Only

**Install only the web interface.**

```bash
./setup.sh --mode web-only
```

**Creates:**
- `~/.config/systemd/user/redditview-web.service` - Web service

**Behavior:**
- Same as API-only but explicitly for web-only use
- All functionality via browser interface
- Lightweight footprint

**Access:**
```bash
# Web UI on http://localhost:3000

# View logs
journalctl --user -u redditview-web -f
```

**Perfect for:**
- Web-only deployments
- Browser-based access
- Remote servers

---

## 🔧 Detailed Setup

### Step 1: Verify Prerequisites

Before running setup.sh, ensure you have:

```bash
# Check systemd
systemctl --user status

# Check tmux (if using 'both' mode)
tmux --version

# Check Node.js
node --version

# Check Go (if TUI binary not built)
go version

# Check git
git --version
```

### Step 2: Navigate to Installation Directory

```bash
cd /path/to/redditiew-local
```

### Step 3: Choose Installation Scope and Run Setup Script

**Interactive (recommended for first-time users):**
```bash
# User-level (no sudo needed)
./setup.sh

# System-level (with sudo)
sudo ./setup.sh
```

**With specific options:**

User-level:
```bash
./setup.sh \
  --scope user \
  --mode both \
  --path /home/username/redditiew \
  --enable \
  --start
```

System-level:
```bash
sudo ./setup.sh \
  --scope system \
  --mode both \
  --path /opt/redditiew \
  --user redditview \
  --enable
```

### Step 4: Verify Installation

**User-level:**
```bash
# List installed services
systemctl --user list-unit-files | grep redditview

# Check service status
systemctl --user status redditview-api
systemctl --user status redditview-tui      # if using 'both' mode

# Verify files
ls -la ~/.config/systemd/user/redditview-*.service
```

**System-level:**
```bash
# List installed services
sudo systemctl list-unit-files | grep redditview

# Check service status
sudo systemctl status redditview-api
sudo systemctl status redditview-tui        # if using 'both' mode

# Verify files
sudo ls -la /etc/systemd/system/redditview-*.service
```

### Step 5: Enable on Boot (Optional)

**User-level - Enable user session to run at boot:**
```bash
sudo loginctl enable-linger $USER

# Enable services
systemctl --user enable redditview-api.service
systemctl --user enable redditview-tui.service  # if applicable
```

**System-level - Services already enabled:**
```bash
# Services are already system-wide on boot
# Just start them if not already running
sudo systemctl start redditview-api.service
sudo systemctl start redditview-tui.service  # if applicable
```

### Step 6: Verify Configuration

**User-level:**
```bash
cat ~/.config/systemd/user/redditview-api.service
cat ~/.config/systemd/user/redditview-tui.service

# Or view all RedditView services
ls -la ~/.config/systemd/user/redditview-*.service
```

**System-level:**
```bash
sudo cat /etc/systemd/system/redditview-api.service
sudo cat /etc/systemd/system/redditview-tui.service

# Or view all RedditView services
sudo ls -la /etc/systemd/system/redditview-*.service
```

---

## 🎮 Service Management

### Starting Services

**User-level - Start all services:**
```bash
systemctl --user start redditview-api
systemctl --user start redditview-tui      # if using 'both' mode
```

**System-level - Start all services:**
```bash
sudo systemctl start redditview-api
sudo systemctl start redditview-tui        # if using 'both' mode
```

**Start with auto-start on boot:**

User-level:
```bash
systemctl --user enable redditview-api
systemctl --user enable redditview-tui     # if using 'both' mode
systemctl --user start redditview-api
systemctl --user start redditview-tui
```

System-level:
```bash
sudo systemctl enable redditview-api
sudo systemctl enable redditview-tui       # if using 'both' mode
sudo systemctl start redditview-api
sudo systemctl start redditview-tui
```

### Stopping Services

**User-level - Stop services:**
```bash
systemctl --user stop redditview-api
systemctl --user stop redditview-tui       # if using 'both' mode
```

**System-level - Stop services:**
```bash
sudo systemctl stop redditview-api
sudo systemctl stop redditview-tui         # if using 'both' mode
```

**Disable auto-start:**

User-level:
```bash
systemctl --user disable redditview-api
systemctl --user disable redditview-tui    # if using 'both' mode
```

System-level:
```bash
sudo systemctl disable redditview-api
sudo systemctl disable redditview-tui      # if using 'both' mode
```

### Restarting Services

**User-level:**
```bash
# Restart API service
systemctl --user restart redditview-api

# Restart TUI service
systemctl --user restart redditview-tui

# Restart all services
systemctl --user restart redditview-api redditview-tui
```

**System-level:**
```bash
# Restart API service
sudo systemctl restart redditview-api

# Restart TUI service
sudo systemctl restart redditview-tui

# Restart all services
sudo systemctl restart redditview-api redditview-tui
```

### Checking Status

**User-level - View service status:**
```bash
systemctl --user status redditview-api
systemctl --user status redditview-tui     # if using 'both' mode
```

**System-level - View service status:**
```bash
sudo systemctl status redditview-api
sudo systemctl status redditview-tui       # if using 'both' mode
```

**Check if enabled:**

User-level:
```bash
systemctl --user is-enabled redditview-api
```

System-level:
```bash
sudo systemctl is-enabled redditview-api
```

# Enable
systemctl --user enable redditview-api

# Disable
systemctl --user disable redditview-api
```


### Viewing Logs

**User-level - Real-time logs:**
```bash
# API server logs
journalctl --user -u redditview-api -f

# TUI service logs
journalctl --user -u redditview-tui -f

# Both services
journalctl --user -u redditview-api -u redditview-tui -f
```

**System-level - Real-time logs:**
```bash
# API server logs
sudo journalctl -u redditview-api -f

# TUI service logs
sudo journalctl -u redditview-tui -f

# Both services
sudo journalctl -u redditview-api -u redditview-tui -f
```

**Historical logs (user-level):**
```bash
# Last 50 lines
journalctl --user -u redditview-api -n 50

# Last hour
journalctl --user -u redditview-api --since "1 hour ago"

# Since specific time
journalctl --user -u redditview-api --since "2024-02-22 14:00:00"

# Full log output
journalctl --user -u redditview-api --no-pager
```

**Historical logs (system-level):**
```bash
# Last 50 lines
sudo journalctl -u redditview-api -n 50

# Last hour
sudo journalctl -u redditview-api --since "1 hour ago"

# Since specific time
sudo journalctl -u redditview-api --since "2024-02-22 14:00:00"

# Full log output
sudo journalctl -u redditview-api --no-pager
```

### Accessing the TUI

**When running in tmux (both mode):**

```bash
# Attach to running session
tmux attach-session -t redditview

# Detach from session (Ctrl+B then D)
# Session continues running

# List all tmux sessions
tmux list-sessions

# Kill tmux session (stops TUI)
tmux kill-session -t redditview
```

**Tmux navigation:**
- `Ctrl+B` - Prefix key
- `Ctrl+B D` - Detach from session
- `Ctrl+B [` - Enter scroll mode
- `Ctrl+B :` - Enter command mode
- `Ctrl+B ?` - Show keybindings

### System Resource Monitoring

**Monitor resource usage:**
```bash
# View process details
ps aux | grep redditview

# Monitor in real-time
top -p $(pgrep -f redditview-api),$(pgrep -f "tmux.*redditview")

# Check open ports
netstat -tlnp | grep -E "3000|3002"
```

---

## 🔍 Troubleshooting

### Service Won't Start

**Check service status:**
```bash
systemctl --user status redditview-api
journalctl --user -u redditview-api --no-pager
```

**Common issues:**

1. **Port already in use:**
   ```bash
   # Check what's using port 3002
   lsof -i :3002
   
   # Stop the existing process
   kill -9 <PID>
   ```

2. **Node.js not found:**
   ```bash
   # Check Node.js installation
   which node
   
   # Update service file with correct path
   nano ~/.config/systemd/user/redditview-api.service
   
   # Find and update the ExecStart line
   # Reload services
   systemctl --user daemon-reload
   ```

3. **Installation path doesn't exist:**
   ```bash
   # Verify path
   ls -la /path/to/install
   
   # Reinstall with correct path
   ./setup.sh --path /correct/path --enable
   ```

### TUI Not Connecting to API

**Verify API is running:**
```bash
# Check API service
systemctl --user status redditview-api

# Test API endpoint
curl http://localhost:3002/api/r/sysadmin.json
```

**Check TUI logs:**
```bash
journalctl --user -u redditview-tui -f
```

**Manually start TUI to see errors:**
```bash
# Kill existing TUI
tmux kill-session -t redditview

# Stop TUI service
systemctl --user stop redditview-tui

# Run TUI manually
./apps/tui/redditview
```

### Tmux Session Already Exists

**When TUI service fails to start:**
```bash
# Check for existing session
tmux list-sessions

# Kill the orphaned session
tmux kill-session -t redditview

# Restart service
systemctl --user restart redditview-tui
```

### Permission Denied Errors

**If service runs as wrong user:**
```bash
# Check service user
grep "^User=" ~/.config/systemd/user/redditview-api.service

# Fix ownership if needed
chown -R $USER:$USER /path/to/redditiew-local

# Reload services
systemctl --user daemon-reload
systemctl --user restart redditview-api
```

### Logs Show "File not found"

**Check all paths in service file:**
```bash
# View service file
cat ~/.config/systemd/user/redditview-api.service

# Verify binary exists
ls -la /path/to/api-server.js
ls -la /path/to/apps/tui/redditview

# Reinstall if paths are wrong
./setup.sh --path /correct/path
```

### Service Keeps Restarting

**Check restart limits in service file:**
```bash
# View service configuration
cat ~/.config/systemd/user/redditview-api.service

# Look for RestartLimitInterval and RestartLimitBurst
# View logs for error pattern
journalctl --user -u redditview-api --no-pager
```

---

## ⚙️ Advanced Configuration

### Customizing Service Files

**Edit service file:**
```bash
nano ~/.config/systemd/user/redditview-api.service
```

**Common customizations:**

1. **Change restart behavior:**
   ```ini
   [Service]
   Restart=on-failure
   RestartSec=5           # Shorter wait time
   StartLimitInterval=120 # Longer interval
   StartLimitBurst=5      # More restart attempts
   ```

2. **Set environment variables:**
   ```ini
   [Service]
   Environment="PORT=3003"
   Environment="NODE_ENV=production"
   Environment="LOG_LEVEL=debug"
   ```

3. **Adjust tmux window size:**
   ```ini
   [Service]
   ExecStart=/usr/bin/tmux new-session -d -s redditview -x 200 -y 50 "..."
   ```

4. **Add pre-start checks:**
   ```ini
   [Service]
   ExecStartPre=/bin/bash -c "mkdir -p /tmp/redditview"
   ExecStartPre=/usr/bin/test -f /path/to/config.json
   ```

**After editing, reload:**
```bash
systemctl --user daemon-reload
systemctl --user restart redditview-api
```

### Running Multiple Instances

**Create second instance with different config:**

```bash
# Copy service file
cp ~/.config/systemd/user/redditview-api.service \
   ~/.config/systemd/user/redditview-api-2.service

# Edit second instance
nano ~/.config/systemd/user/redditview-api-2.service

# Change:
# - Description
# - Port (in Environment)
# - WorkingDirectory (if different)

# Reload and enable
systemctl --user daemon-reload
systemctl --user enable redditview-api-2.service
systemctl --user start redditview-api-2.service
```

### Using with Docker

**Run RedditView in Docker with systemd service:**

```dockerfile
# Dockerfile
FROM node:18-alpine
WORKDIR /redditiew
COPY . .
RUN npm install
EXPOSE 3000 3002
CMD ["npm", "start"]
```

**Create systemd service for Docker:**
```ini
[Unit]
Description=RedditView Docker Container
After=docker.service
Requires=docker.service

[Service]
Type=simple
User=redditview
WorkingDirectory=/path/to/redditiew
ExecStart=/usr/bin/docker run --rm \
  --name redditview \
  -p 3000:3000 \
  -p 3002:3002 \
  -v /path/to/redditiew:/redditiew \
  redditview:latest

[Install]
WantedBy=multi-user.target
```

### Systemd Timer for Scheduled Tasks

**Create periodic restart (e.g., daily):**

```bash
# Create timer unit
cat > ~/.config/systemd/user/redditview-restart.timer << EOF
[Unit]
Description=Daily RedditView Restart
Requires=redditview-restart.service

[Timer]
OnCalendar=daily
OnCalendar=02:00
Persistent=true

[Install]
WantedBy=timers.target
EOF

# Create service to run on timer
cat > ~/.config/systemd/user/redditview-restart.service << EOF
[Unit]
Description=RedditView Restart

[Service]
Type=oneshot
ExecStart=/usr/bin/systemctl --user restart redditview-api.service
EOF

# Enable timer
systemctl --user enable redditview-restart.timer
systemctl --user start redditview-restart.timer
```

---

## 🗑️ Uninstallation

### Remove Services

**Disable and remove:**

```bash
# Stop services
systemctl --user stop redditview-api redditview-tui

# Disable from boot
systemctl --user disable redditview-api redditview-tui

# Remove service files
rm ~/.config/systemd/user/redditview-api.service
rm ~/.config/systemd/user/redditview-tui.service
rm ~/.config/systemd/user/redditview-web.service  # if exists

# Reload systemd
systemctl --user daemon-reload

# Verify removal
systemctl --user list-unit-files | grep redditview
```

### Remove User Lingering (if enabled)

```bash
# Disable lingering for user
sudo loginctl disable-linger $USER
```

### Complete Cleanup

```bash
# Kill any remaining tmux sessions
tmux kill-session -t redditview

# Remove setup files
rm -f ~/setup.sh
rm -rf ~/systemd-templates/

# Remove application (optional)
rm -rf /path/to/redditiew-local
```

---

## 📊 Reference

### Service File Locations

| Type | Location | Scope |
|------|----------|-------|
| User services | `~/.config/systemd/user/` | Current user only |
| System services | `/etc/systemd/system/` | Entire system |
| System user services | `/usr/lib/systemd/user/` | Default user services |

### Port Reference

| Component | Port | Type |
|-----------|------|------|
| API Server | 3002 | REST API |
| Web UI | 3000 | HTTP |
| Tmux | N/A | Socket |

### Useful Systemctl Commands

```bash
# User service commands
systemctl --user list-units              # List all user units
systemctl --user list-unit-files         # List all unit files
systemctl --user status redditview-api   # Show status
systemctl --user start redditview-api    # Start service
systemctl --user stop redditview-api     # Stop service
systemctl --user restart redditview-api  # Restart service
systemctl --user reload redditview-api   # Reload config
systemctl --user enable redditview-api   # Enable at boot
systemctl --user disable redditview-api  # Disable at boot
systemctl --user is-enabled redditview-api  # Check if enabled

# Journal commands
journalctl --user-unit=redditview-api    # View logs
journalctl --user-unit=redditview-api -f # Tail logs
journalctl --user-unit=redditview-api -n 100 # Last 100 lines
journalctl --user-unit=redditview-api --since "1 hour ago" # Since time
```

---

## 🎓 Learning Resources

- [systemd User Documentation](https://www.freedesktop.org/software/systemd/man/systemd.unit.html)
- [systemctl Manual](https://www.freedesktop.org/software/systemd/man/systemctl.html)
- [journalctl Manual](https://www.freedesktop.org/software/systemd/man/journalctl.html)
- [Tmux Cheat Sheet](https://tmuxcheatsheet.com/)
- [RedditView Documentation](DOCS_INDEX.md)

---

## 📞 Support

For issues or questions:
- Check logs: `journalctl --user -u redditview-api -f`
- Review configuration: `cat ~/.config/systemd/user/redditview-api.service`
- See [Troubleshooting](#troubleshooting) section above
- Open an issue: https://github.com/nicthegarden/redditiew/issues

---

**Happy browsing! 🚀**
