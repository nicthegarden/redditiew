# System-Level Systemd Deployment Guide

This guide covers **system-level** installation of RedditView as systemd services in `/etc/systemd/system/`, suitable for production servers and always-on systems.

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [System-Level vs User-Level](#system-level-vs-user-level)
3. [Quick Start](#quick-start)
4. [Step-by-Step Installation](#step-by-step-installation)
5. [Verification](#verification)
6. [Service Management](#service-management)
7. [Boot Startup Configuration](#boot-startup-configuration)
8. [System-Wide Considerations](#system-wide-considerations)
9. [Monitoring & Troubleshooting](#monitoring--troubleshooting)
10. [Uninstallation](#uninstallation)

---

## Prerequisites

### System Requirements

- **Operating System**: Linux with systemd (Ubuntu 16.04+, Debian 8+, Fedora, Arch, etc.)
- **Root/Sudo Access**: Required for system-level installation
- **systemd version**: 230+ (check with `systemd --version`)
- **Disk Space**: ~500MB for application + dependencies

### Required Software

```bash
# Check systemd
systemctl --version

# Check Node.js
node --version  # v16 or higher recommended

# Check tmux (for TUI support)
tmux --version

# Check for required user creation capability
id  # Must have sudo access
```

### Network Ports

- **API Server**: 8765 (HTTP)
- **Web UI**: 5174 (HTTP, dev server)
- **Tmux**: Internal socket (no port needed)

**Note:** If ports are already in use, the setup will fail. Check with:
```bash
sudo netstat -tlnp | grep -E "8765|5174"
# or with newer systems:
sudo ss -tlnp | grep -E "8765|5174"
```

---

## System-Level vs User-Level

### Key Differences

| Aspect | User-Level | System-Level |
|--------|-----------|--------------|
| **Installation Path** | `~/.config/systemd/user/` | `/etc/systemd/system/` |
| **Root Required** | ❌ No | ✅ Yes (sudo) |
| **Runs At Boot** | ❌ No (requires login) | ✅ Yes (even headless) |
| **Service User** | Current user | Dedicated user (e.g., `redditview`) |
| **Best For** | Development, personal machines | Production servers, always-on systems |
| **Management** | `systemctl --user` | `sudo systemctl` |
| **Logs** | `journalctl --user` | `sudo journalctl` |

### When to Choose System-Level

✅ **Choose system-level if:**
- Running on a server (no GUI expected)
- Need services to start before user login
- Want always-on operation
- Running multiple RedditView instances
- Need production-grade reliability
- Services must survive user logout

❌ **Choose user-level if:**
- Running on desktop/laptop
- Development environment
- Want services isolated to your account
- Prefer no root/sudo involvement
- Testing or experimentation

---

## Quick Start

### Installation (3 steps)

```bash
# Step 1: Clone the repository (if not already done)
git clone git@github.com:nicthegarden/redditiew.git
cd redditiew

# Step 2: Run setup with sudo (system-level)
sudo ./setup-systemctl.sh \
  --scope system \
  --mode both \
  --user redditview \
  --enable \
  --start

# Step 3: Verify installation
sudo systemctl status redditview-api
sudo systemctl status redditview-tui
```

### Verify Services Are Running

```bash
# Check status of all RedditView services
sudo systemctl list-units --all | grep redditview

# Test API endpoint
curl http://localhost:8765/api/r/sysadmin.json

# View logs
sudo journalctl -u redditview-api -u redditview-tui -f
```

---

## Step-by-Step Installation

### Step 1: Create System User for Services

The system-level installation runs services as a dedicated user for security. You can either create this manually or let the setup script handle it:

**Option A: Let setup script create user (Recommended)**

```bash
# Script will create 'redditview' user automatically
sudo ./setup-systemctl.sh --scope system --user redditview
```

**Option B: Create user manually**

```bash
# Create user without login shell for security
sudo useradd -r -s /bin/false -d /var/lib/redditview redditview

# Create home directory with proper permissions
sudo mkdir -p /var/lib/redditview
sudo chown redditview:redditview /var/lib/redditview
sudo chmod 750 /var/lib/redditview
```

### Step 2: Prepare Installation Directory

Decide where to install RedditView:

**Option A: System-wide installation (/opt)**

```bash
# Copy or clone to /opt
sudo cp -r /path/to/redditiew /opt/redditiew
sudo chown -R redditview:redditview /opt/redditiew
sudo chmod 755 /opt/redditiew
```

**Option B: Keep in home directory**

```bash
# If in user home, ensure permissions are correct
sudo chown -R redditview:redditview /home/user/redditiew
```

### Step 3: Run Setup Script

Navigate to the RedditView directory and run the setup script with sudo:

```bash
cd /opt/redditiew  # or wherever you installed it

# Run setup as root
sudo ./setup-systemctl.sh \
  --scope system \
  --mode both \
  --path /opt/redditiew \
  --user redditview \
  --enable \
  --start
```

**Script options explained:**

- `--scope system`: Install in `/etc/systemd/system/`
- `--mode both`: Install both API and TUI services
- `--path /opt/redditiew`: Path to RedditView installation
- `--user redditview`: User to run services as
- `--enable`: Enable services to start at boot
- `--start`: Start services immediately after setup

### Step 4: Verify Service Files

Check that service files were created correctly:

```bash
# List all RedditView services
sudo ls -la /etc/systemd/system/redditview-*.service

# Check API service configuration
sudo cat /etc/systemd/system/redditview-api.service

# Check TUI service configuration
sudo cat /etc/systemd/system/redditview-tui.service
```

Expected output for API service:

```ini
[Unit]
Description=RedditView API Server
After=network.target

[Service]
Type=simple
User=redditview
WorkingDirectory=/opt/redditiew
ExecStart=/usr/bin/node /opt/redditiew/api-server.js
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=redditview-api

[Install]
WantedBy=multi-user.target
```

---

## Verification

### Step 1: Check Service Status

```bash
# Check all services
sudo systemctl status redditview-api
sudo systemctl status redditview-tui

# Or view in list format
sudo systemctl list-units --all | grep redditview
```

**Healthy status output:**

```
● redditview-api.service - RedditView API Server
     Loaded: loaded (/etc/systemd/system/redditview-api.service; enabled; vendor preset: disabled)
     Active: active (running) since Mon 2024-02-26 10:15:23 UTC; 2min ago
     ...
```

### Step 2: Test API Endpoint

```bash
# Test local API (from server)
curl http://localhost:8765/api/r/sysadmin.json | head -20

# Test from remote machine (replace with server IP)
curl http://192.168.1.104:8765/api/r/sysadmin.json | head -20
```

**Healthy response:**

```json
{
  "subreddit": "sysadmin",
  "displayName": "r/sysadmin",
  "subscribers": 123456,
  "posts": [
    {
      "id": "abc123",
      "title": "Post Title",
      ...
    }
  ]
}
```

### Step 3: Check Logs

```bash
# View recent logs (all services)
sudo journalctl -u redditview-api -u redditview-tui -n 50

# Follow logs in real-time
sudo journalctl -u redditview-api -u redditview-tui -f

# View logs since last boot
sudo journalctl -u redditview-api -b
```

### Step 4: Verify Boot Startup

Check if services are enabled for automatic startup:

```bash
# Check if enabled
sudo systemctl is-enabled redditview-api
sudo systemctl is-enabled redditview-tui

# Should output: "enabled"
```

To enable for boot (if not done during setup):

```bash
sudo systemctl enable redditview-api.service
sudo systemctl enable redditview-tui.service
```

### Step 5: Check System Resources

```bash
# View running RedditView processes
ps aux | grep -E "node.*api-server|redditview"

# Check memory usage
ps aux | grep redditview | awk '{print $6}' | awk '{sum+=$1} END {print "Total memory: " sum "KB"}'

# Check open ports
sudo netstat -tlnp | grep -E "8765|5174"

# Or with ss (newer systems)
sudo ss -tlnp | grep -E "8765|5174"
```

---

## Service Management

### Starting Services

```bash
# Start all RedditView services
sudo systemctl start redditview-api redditview-tui

# Or individually
sudo systemctl start redditview-api
sudo systemctl start redditview-tui
```

### Stopping Services

```bash
# Stop all services
sudo systemctl stop redditview-api redditview-tui

# Or individually
sudo systemctl stop redditview-api
sudo systemctl stop redditview-tui
```

### Restarting Services

```bash
# Restart all services
sudo systemctl restart redditview-api redditview-tui

# Or individually
sudo systemctl restart redditview-api
sudo systemctl restart redditview-tui

# Reload configuration (without restarting)
sudo systemctl reload redditview-api
```

### Checking Service Status

```bash
# Detailed status
sudo systemctl status redditview-api

# Quick status
sudo systemctl is-active redditview-api

# Check if enabled
sudo systemctl is-enabled redditview-api
```

---

## Boot Startup Configuration

### Enable Services to Start at Boot

The `--enable` flag during setup enables services automatically, but you can also enable manually:

```bash
# Enable services
sudo systemctl enable redditview-api.service
sudo systemctl enable redditview-tui.service

# Verify they're enabled
sudo systemctl is-enabled redditview-api
sudo systemctl is-enabled redditview-tui
```

### Disable Auto-Start

If you want services to NOT start at boot:

```bash
# Disable auto-start
sudo systemctl disable redditview-api.service
sudo systemctl disable redditview-tui.service

# Services will still run until you stop them
sudo systemctl stop redditview-api redditview-tui
```

### Test Boot Startup (Simulated)

To test without rebooting:

```bash
# Stop services
sudo systemctl stop redditview-api redditview-tui

# Verify they're stopped
sudo systemctl status redditview-api

# "Start" services (simulating boot)
sudo systemctl start redditview-api redditview-tui

# Verify they started
sudo systemctl status redditview-api
```

### View Boot Startup Sequence

```bash
# View systemd startup logs
sudo journalctl -b

# View RedditView startup during boot
sudo journalctl -u redditview-api -u redditview-tui -b
```

---

## System-Wide Considerations

### Security

#### File Permissions

```bash
# Verify correct ownership
ls -la /etc/systemd/system/redditview-*.service

# Should be owned by root
# -rw-r--r-- 1 root root

# Verify installation directory
ls -la /opt/redditiew | head -20
# Should be readable by redditview user
```

#### User Isolation

Services run as the `redditview` user, not root:

```bash
# Verify service runs as correct user
grep "^User=" /etc/systemd/system/redditview-api.service
# Should output: User=redditview
```

### Port Management

Ensure required ports are available:

```bash
# Check if ports are in use
sudo ss -tlnp | grep -E "8765|5174"

# If in use, identify the process
sudo lsof -i :8765

# Stop conflicting service and try again
sudo systemctl stop <service-name>
```

### Resource Limits

System-level services can enforce resource limits. Customize in service files:

```bash
# Edit service file
sudo nano /etc/systemd/system/redditview-api.service

# Add to [Service] section:
# MemoryLimit=512M
# CPUQuota=50%
# LimitNOFILE=65536

# Reload and restart
sudo systemctl daemon-reload
sudo systemctl restart redditview-api
```

### Firewall Configuration

If using firewall, allow API and Web UI ports:

```bash
# UFW (Ubuntu/Debian)
sudo ufw allow 8765/tcp
sudo ufw allow 5174/tcp

# firewalld (Fedora/CentOS)
sudo firewall-cmd --permanent --add-port=8765/tcp
sudo firewall-cmd --permanent --add-port=5174/tcp
sudo firewall-cmd --reload

# iptables
sudo iptables -A INPUT -p tcp --dport 8765 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 5174 -j ACCEPT
```

### Logging Configuration

System-level services log to journalctl:

```bash
# View all RedditView logs
sudo journalctl -u redditview-api -u redditview-tui -n 100

# Follow logs
sudo journalctl -u redditview-api -f

# Filter by priority
sudo journalctl -u redditview-api -p err  # Errors only

# Export logs
sudo journalctl -u redditview-api -o export > /tmp/redditview-logs.txt
```

---

## Monitoring & Troubleshooting

### Common Issues and Solutions

#### Services Won't Start

**Check status:**
```bash
sudo systemctl status redditview-api
sudo journalctl -u redditview-api -n 20
```

**Common causes:**

1. **Port already in use:**
   ```bash
   # Identify process using port 8765
   sudo lsof -i :8765
   
   # Kill conflicting process
   sudo kill -9 <PID>
   
   # Restart service
   sudo systemctl restart redditview-api
   ```

2. **Installation path doesn't exist:**
   ```bash
   # Check path in service file
   sudo cat /etc/systemd/system/redditview-api.service | grep WorkingDirectory
   
   # Verify path exists
   ls -la /opt/redditiew
   
   # If missing, reinstall:
   sudo cp -r /path/to/source /opt/redditiew
   sudo chown -R redditview:redditview /opt/redditiew
   ```

3. **Node.js not found:**
   ```bash
   # Verify Node.js location
   which node
   
   # Update service file with correct path
   sudo nano /etc/systemd/system/redditview-api.service
   # Change: ExecStart=/usr/bin/node ... to correct path
   
   # Reload
   sudo systemctl daemon-reload
   sudo systemctl restart redditview-api
   ```

#### TUI Not Connecting

**Check if API is running:**
```bash
sudo systemctl status redditview-api
curl http://localhost:8765/api/r/sysadmin.json
```

**Check TUI logs:**
```bash
sudo journalctl -u redditview-tui -f
```

**Manually test TUI:**
```bash
# Find tmux session
tmux list-sessions

# Attach if exists
tmux attach-session -t redditview

# Or kill and restart
sudo systemctl stop redditview-tui
sudo systemctl start redditview-tui
tmux attach-session -t redditview
```

#### Tmux Session Issues (TUI Service)

**The Robust 4-Phase Startup Strategy**

The TUI service uses a resilient approach to tmux session management:
- **Phase 1:** Create session (idempotent, won't fail if exists)
- **Phase 2:** Wait for session readiness (up to 3 seconds)
- **Phase 3:** Send TUI command to session
- **Phase 4:** Monitor session and keep service running

**Tmux session doesn't exist:**
```bash
# Check for existing session
sudo tmux list-sessions -t redditview 2>/dev/null || echo "No session found"

# Verify the service created it
sudo systemctl status redditview-tui

# Check if Phase 1-4 completed successfully
sudo journalctl -u redditview-tui -n 30 --no-pager

# Manually restart to trigger Phase 1-4
sudo systemctl restart redditview-tui

# Verify session now exists
sudo tmux list-sessions | grep redditview
```

**TUI process exists but not in tmux:**
```bash
# Kill orphaned processes
sudo pkill -f "apps/tui/redditview"
sudo pkill -f "tmux new-session"

# Kill the session
sudo tmux kill-session -t redditview 2>/dev/null || true

# Restart service (triggers Phase 1-4 again)
sudo systemctl restart redditview-tui

# Verify TUI is now in tmux
ps -ef | grep redditview | grep -v grep
# Parent process should be tmux, not bash
```

**Accessing the TUI in tmux:**
```bash
# Attach to the TUI session
sudo tmux attach-session -t redditview

# Detach without stopping (Ctrl+B then D from inside tmux)
# Or by sending keys:
sudo tmux send-keys -t redditview C-b
sudo tmux send-keys -t redditview D

# Service continues running in background
sudo systemctl status redditview-tui
```

#### Service Keeps Restarting

**Check restart limits:**
```bash
sudo cat /etc/systemd/system/redditview-api.service | grep Restart

# Should show something like:
# Restart=on-failure
# RestartSec=5
```

**View crash logs:**
```bash
sudo journalctl -u redditview-api --no-pager
```

**Increase restart interval temporarily:**
```bash
sudo nano /etc/systemd/system/redditview-api.service

# Change:
# RestartSec=30  # Longer wait between restarts
# StartLimitInterval=300
# StartLimitBurst=3

sudo systemctl daemon-reload
sudo systemctl restart redditview-api
```

#### Permission Issues

**Verify file permissions:**
```bash
# Check ownership of installation
ls -la /opt/redditiew | head -10

# Should be: redditview:redditview with read/execute

# Fix permissions if needed
sudo chown -R redditview:redditview /opt/redditiew
sudo chmod -R 755 /opt/redditiew
```

### Monitoring Dashboard

Create a monitoring script to check status:

```bash
#!/bin/bash
# Save as: /usr/local/bin/redditview-status

echo "=== RedditView System Services ==="
echo ""
echo "Service Status:"
sudo systemctl status redditview-api --no-pager | head -10
echo ""
sudo systemctl status redditview-tui --no-pager | head -10
echo ""
echo "Recent Logs:"
sudo journalctl -u redditview-api -n 5 --no-pager
echo ""
echo "Open Ports:"
sudo ss -tlnp | grep -E "8765|5174"
```

Make executable and run:
```bash
sudo chmod +x /usr/local/bin/redditview-status
redditview-status
```

---

## Uninstallation

### Complete Removal

```bash
# Step 1: Stop services
sudo systemctl stop redditview-api redditview-tui

# Step 2: Disable auto-start
sudo systemctl disable redditview-api redditview-tui

# Step 3: Remove service files
sudo rm /etc/systemd/system/redditview-*.service

# Step 4: Reload systemd
sudo systemctl daemon-reload

# Step 5: Remove installation directory (optional)
sudo rm -rf /opt/redditiew

# Step 6: Remove system user (optional)
sudo userdel -r redditview

# Step 7: Verify removal
sudo systemctl list-units --all | grep redditview
```

### Partial Removal (Keep Application)

If you want to keep the application but remove systemd services:

```bash
# Stop services
sudo systemctl stop redditview-api redditview-tui

# Disable from boot
sudo systemctl disable redditview-api redditview-tui

# Remove only service files
sudo rm /etc/systemd/system/redditview-*.service

# Reload systemd
sudo systemctl daemon-reload
```

---

## Advanced Topics

### Running Multiple Instances

Create additional services for different configurations:

```bash
# Copy service file for second instance
sudo cp /etc/systemd/system/redditview-api.service \
        /etc/systemd/system/redditview-api-2.service

# Edit new service
sudo nano /etc/systemd/system/redditview-api-2.service

# Change:
# Description=RedditView API Server (Instance 2)
# Add to [Service] section:
# Environment="PORT=8766"

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable redditview-api-2.service
sudo systemctl start redditview-api-2.service
```

### Integration with Monitoring Tools

#### With Prometheus

Add health check endpoint to service file:

```bash
sudo nano /etc/systemd/system/redditview-api.service

# Add to [Service] section:
# ExecStartPost=/bin/bash -c 'sleep 2 && curl -f http://localhost:8765/health || exit 1'
```

#### With Nagios

Create monitoring script:

```bash
#!/bin/bash
# Save as: /usr/local/lib/nagios/plugins/check_redditview

STATUS=$(sudo systemctl is-active redditview-api)
if [ "$STATUS" = "active" ]; then
  echo "OK: RedditView API is running"
  exit 0
else
  echo "CRITICAL: RedditView API is not running"
  exit 2
fi
```

### Backup and Restore

Backup service configuration:

```bash
# Backup all RedditView services
sudo cp -r /etc/systemd/system/redditview-*.service /tmp/redditview-backup/

# Backup application
sudo tar czf /tmp/redditview-app-backup.tar.gz /opt/redditiew
```

Restore:

```bash
# Restore services
sudo cp /tmp/redditview-backup/redditview-*.service /etc/systemd/system/
sudo systemctl daemon-reload

# Restore application
sudo tar xzf /tmp/redditview-app-backup.tar.gz -C /
sudo chown -R redditview:redditview /opt/redditiew
```

---

## Reference

### Service File Locations

- **System-level**: `/etc/systemd/system/`
- **User-level**: `~/.config/systemd/user/`
- **System defaults**: `/usr/lib/systemd/system/`

### Useful Commands

```bash
# Reload systemd after changes
sudo systemctl daemon-reload

# View all units
sudo systemctl list-units --all

# View specific service details
sudo systemctl show redditview-api

# View service unit file
sudo systemctl cat redditview-api.service

# Edit service (with proper reloading)
sudo systemctl edit redditview-api.service

# View dependencies
sudo systemctl list-dependencies redditview-api

# Journal commands
sudo journalctl -u redditview-api          # View logs
sudo journalctl -u redditview-api -f       # Follow logs
sudo journalctl -u redditview-api -n 50    # Last 50 lines
sudo journalctl -u redditview-api -p err   # Errors only
sudo journalctl -u redditview-api -b       # Since last boot
```

---

**Last Updated**: 2024-02-26
**Version**: 1.0
**Compatible**: Linux systems with systemd 230+

