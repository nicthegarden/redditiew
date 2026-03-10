# RedditView SystemD Quick Start Guide

Get RedditView running as systemd services in just a few seconds!

## Prerequisites

- Linux with systemd
- Node.js installed
- tmux (for TUI mode)
- Terminal access

## Installation

### Option 1: User-Level (Recommended - No Sudo)

Perfect for development machines, desktops, and personal use. Services run as your user.

```bash
cd /path/to/redditiew
./setup-systemctl.sh --scope user --mode both --enable --start
```

**This command:**
- Creates services in `~/.config/systemd/user/`
- Enables them to run at startup
- Starts them immediately

### Option 2: System-Level (Production - Requires Sudo)

Perfect for servers and production deployments. Services run at system boot.

```bash
cd /path/to/redditiew
sudo ./setup-systemctl.sh --scope system --mode both --user redditview --enable --start
```

**This command:**
- Creates services in `/etc/systemd/system/`
- Creates a dedicated `redditview` user
- Enables automatic startup on system boot
- Starts services immediately

### Option 3: API-Only (Headless Servers)

Use this for servers without a display, or if you only need the API.

```bash
./setup-systemctl.sh --scope user --mode api-only --enable --start
```

## What You Get

### Running Services

After installation, three services are created and running:

1. **redditview-api** (Port 8765)
   - REST API for Reddit data
   - Caches responses for 60 seconds
   - Always required for other components

2. **redditview-tui** (Terminal UI)
   - Interactive terminal interface
   - Runs in tmux session named "redditview"
   - Requires display/terminal access
   - Only created in "both" mode

3. **redditview-web** (Port 5174)
   - React-based web interface
   - Browser-based access
   - Can be accessed from any computer on the network

## Accessing RedditView

### Web Interface

Open your browser to one of:
- **Local:** http://localhost:5174
- **Remote:** http://your-ip-address:5174

### Terminal Interface (TUI)

If running in "both" mode with TUI:

```bash
# Attach to the tmux session
tmux attach-session -t redditview

# Detach (keep running): Press Ctrl+B then D
# Kill session: tmux kill-session -t redditview
```

### API Server

Test the API directly:

```bash
# Health check
curl http://localhost:8765/health

# Get config
curl http://localhost:8765/config

# Get subreddit posts
curl http://localhost:8765/api/r/sysadmin

# Get top posts
curl http://localhost:8765/api/r/sysadmin/top
```

## Managing Services

### Check Status

```bash
# User-level services
systemctl --user status redditview-api
systemctl --user status redditview-tui

# System-level services (add sudo)
sudo systemctl status redditview-api
```

### Start/Stop Services

```bash
# User-level
systemctl --user start redditview-api
systemctl --user stop redditview-api
systemctl --user restart redditview-api

# System-level
sudo systemctl start redditview-api
sudo systemctl stop redditview-api
sudo systemctl restart redditview-api
```

### View Logs

```bash
# User-level - real-time logs
journalctl --user -u redditview-api -f

# User-level - last 50 lines
journalctl --user -u redditview-api -n 50

# System-level (add sudo)
sudo journalctl -u redditview-api -f
```

### Enable/Disable Auto-Start

```bash
# User-level
systemctl --user enable redditview-api
systemctl --user disable redditview-api

# System-level
sudo systemctl enable redditview-api
sudo systemctl disable redditview-api
```

## Verification

### Quick Test

```bash
# Test API
curl http://localhost:8765/health

# Test Web (should return HTML)
curl -I http://localhost:5174

# Check service status
systemctl --user status redditview-api

# View recent logs
journalctl --user -u redditview-api -n 10
```

### Port Check

Verify services are listening on correct ports:

```bash
netstat -tlnp | grep -E "8765|5174|tmux"
```

## Troubleshooting

### Service Won't Start

```bash
# Check full error
systemctl --user status redditview-api

# See detailed logs
journalctl --user -u redditview-api --no-pager

# Try starting manually for error messages
cd /path/to/redditiew && /usr/bin/node api-server.js
```

### TUI Not Connecting

```bash
# Verify API is running
systemctl --user status redditview-api

# Test API endpoint
curl http://localhost:8765/health

# Check TUI logs
journalctl --user -u redditview-tui -f

# Kill and restart TUI
systemctl --user restart redditview-tui
```

### Port Already in Use

```bash
# Find what's using the port
lsof -i :8765

# Kill the process
kill -9 <PID>

# Or change the port in the service file
nano ~/.config/systemd/user/redditview-api.service
# Change Environment="PORT=8765" to a different port
# Then reload and restart
systemctl --user daemon-reload
systemctl --user restart redditview-api
```

### Can't Connect Externally

Ensure services are listening on all interfaces:
- Check service file has `--host 0.0.0.0` (web)
- Check no firewall is blocking ports 8765 and 5174
- Use your machine's IP address instead of localhost

```bash
# Find your IP
hostname -I

# Test from another machine
curl http://<your-ip>:8765/health
```

## Configuration

### Change Ports

Edit the service file:

```bash
nano ~/.config/systemd/user/redditview-api.service
```

Find the line `Environment="PORT=8765"` and change to your desired port.

Then reload:

```bash
systemctl --user daemon-reload
systemctl --user restart redditview-api
```

### Environment Variables

Set additional variables in the service file `[Service]` section:

```ini
[Service]
Environment="NODE_ENV=production"
Environment="PORT=8765"
Environment="LOG_LEVEL=debug"
```

## Common Tasks

### Reboot and Verify Auto-Start

```bash
# Enable auto-start (if not already done)
systemctl --user enable redditview-api

# After reboot, verify
systemctl --user status redditview-api
```

### Backup Service Configuration

```bash
cp ~/.config/systemd/user/redditview-*.service ~/backup/
```

### Uninstall Services

```bash
# Stop services
systemctl --user stop redditview-api redditview-tui

# Disable from boot
systemctl --user disable redditview-api redditview-tui

# Remove service files
rm ~/.config/systemd/user/redditview-*.service

# Reload systemd
systemctl --user daemon-reload
```

## Advanced Usage

### Multiple Instances

Run multiple independent RedditView instances:

```bash
# Copy service file
cp ~/.config/systemd/user/redditview-api.service \
   ~/.config/systemd/user/redditview-api-2.service

# Edit second instance - change port and working directory
nano ~/.config/systemd/user/redditview-api-2.service

# Reload and start
systemctl --user daemon-reload
systemctl --user start redditview-api-2
```

### View Service File

```bash
# User-level
cat ~/.config/systemd/user/redditview-api.service

# System-level
sudo cat /etc/systemd/system/redditview-api.service
```

### Customize Service

Edit the service file directly:

```bash
nano ~/.config/systemd/user/redditview-api.service
```

After changes:

```bash
systemctl --user daemon-reload
systemctl --user restart redditview-api
```

## Getting Help

### Full Documentation

See the complete guide with all options and advanced topics:

```bash
cat SYSTEMD_SETUP.md

# Or view the setup script help
./setup-systemctl.sh --help
```

### Support

- **GitHub Issues:** https://github.com/nicthegarden/redditiew/issues
- **Logs:** `journalctl --user -u redditview-api -f`
- **Status:** `systemctl --user status redditview-api`

## Next Steps

1. ✅ Install services: `./setup-systemctl.sh --scope user --mode both --enable --start`
2. ✅ Verify: `systemctl --user status redditview-api`
3. ✅ Access Web: http://localhost:5174
4. ✅ Explore: Check the main README.md for features and usage
5. ✅ Customize: Edit service files or configuration as needed

## Port Reference

| Service | Port | Protocol | Purpose |
|---------|------|----------|---------|
| API | 8765 | HTTP/REST | Backend API for all clients |
| Web UI | 5174 | HTTP | Browser-based interface |
| TUI | tmux | Terminal | Terminal user interface |

---

**Happy browsing with RedditView! 🚀**
