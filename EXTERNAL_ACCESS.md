# RedditView - External Access Configuration

## 🚀 Quick Start

Your RedditView application is now **running and accessible** from any computer on your network.

### Immediate Access URLs

**API Server:**
```
http://192.168.1.104:8765
```

**Web Interface:**
```
http://192.168.1.104:5174
```

---

## 📊 System Status

| Component | Status | Local | External |
|-----------|--------|-------|----------|
| API Server | ✓ Running | `localhost:8765` | `192.168.1.104:8765` |
| Web UI | ✓ Running | `localhost:5174` | `192.168.1.104:5174` |
| TUI | Available | `go run main.go` | N/A (local only) |

---

## 🌐 Network Access

### From Your Computer
```bash
# Test API
curl http://localhost:8765/health

# Open Web UI
# Browser: http://localhost:5174
```

### From Another Computer on Your Network
```bash
# Test API
curl http://192.168.1.104:8765/health

# Open Web UI
# Browser: http://192.168.1.104:5174

# Get Reddit posts
curl http://192.168.1.104:8765/api/r/sysadmin

# Search Reddit
curl "http://192.168.1.104:8765/api/search.json?q=kubernetes"
```

---

## 🔧 Configuration Details

### What Changed

1. **api-server.js**
   - Server now listens on `0.0.0.0` instead of localhost
   - Displays external IP address on startup
   - CORS enabled for cross-origin requests

2. **vite.config.js**
   - Server now listens on `0.0.0.0` instead of localhost
   - Maintains port 5173 (auto-selected 5174 due to existing process)

### How It Works

When a server listens on `0.0.0.0`, it accepts connections from:
- `localhost` (127.0.0.1) - local access
- `192.168.1.104` - your computer's IP
- Any other computer on your LAN

---

## 📡 API Endpoints

All endpoints are accessible from external computers:

```
GET  /health                        - Health check
GET  /api/config                    - Server configuration
GET  /api/stats                     - Server statistics
GET  /api/r/:subreddit              - Get posts from subreddit
GET  /api/r/:subreddit/:sort        - Get sorted posts (hot/new/top)
GET  /api/search.json?q=:query      - Search Reddit
GET  /api/r/:subreddit/comments/:id - Get comments on post
```

### Example Requests

```bash
# From external computer
curl http://192.168.1.104:8765/api/r/golang
curl http://192.168.1.104:8765/api/r/programming/hot
curl "http://192.168.1.104:8765/api/search.json?q=docker"
```

---

## 🛑 Stopping/Starting Services

### Stop API Server
```bash
kill $(cat /tmp/api_server.pid)
# or
pkill -f "node api-server.js"
```

### Stop Web UI
```bash
kill $(cat /tmp/web_server.pid)
# or
pkill -f "vite"
```

### Restart Both
```bash
cd /home/edve/2/redditiew

# Terminal 1: Start API
npm run dev:api

# Terminal 2: Start Web UI
npm run dev
```

### View Logs
```bash
# API Server logs
tail -f /tmp/api_server.log

# Web UI logs
tail -f /tmp/web_server.log
```

---

## 🔒 Security Notes

### Current Configuration
- ✓ Accessible from all network interfaces
- ✓ No authentication (internal network)
- ✓ CORS enabled
- ✓ Good for LAN/development

### For Production
Consider implementing:
- Authentication/authorization
- HTTPS instead of HTTP
- Rate limiting
- API key validation
- Reverse proxy (nginx)
- IP whitelist/blacklist

---

## 🐛 Troubleshooting

### Cannot connect from external computer
1. Verify both computers are on the same network
2. Check firewall isn't blocking ports 8765 or 5174
3. Use the correct IP address (192.168.1.104)
4. Verify servers are running: `ps aux | grep node`

### Port already in use
```bash
# Find which process is using the port
lsof -i :8765

# Kill it
kill -9 <PID>
```

### Slow response
- Check network connection
- Look at logs: `cat /tmp/api_server.log`
- Monitor cache: `curl http://192.168.1.104:8765/api/stats`
- Check system resources

---

## 💡 Usage Examples

### JavaScript/Fetch API
```javascript
// Get posts from r/golang
fetch('http://192.168.1.104:8765/api/r/golang')
  .then(r => r.json())
  .then(data => console.log(data))

// Search Reddit
fetch('http://192.168.1.104:8765/api/search.json?q=kubernetes')
  .then(r => r.json())
  .then(data => console.log(data.data.children))
```

### Python
```python
import requests

# Get posts
response = requests.get('http://192.168.1.104:8765/api/r/sysadmin')
posts = response.json()

# Search
search = requests.get('http://192.168.1.104:8765/api/search.json', 
                      params={'q': 'docker'})
results = search.json()
```

### cURL
```bash
# Get subreddit posts
curl http://192.168.1.104:8765/api/r/golang

# Get hot posts
curl http://192.168.1.104:8765/api/r/programming/hot

# Get new posts
curl http://192.168.1.104:8765/api/r/linux/new

# Search
curl "http://192.168.1.104:8765/api/search.json?q=kubernetes"

# Save to file
curl http://192.168.1.104:8765/api/r/sysadmin > posts.json
```

---

## 📋 System Information

- **System IP:** 192.168.1.104
- **OS:** Arch Linux
- **Node.js:** v25.7.0
- **npm:** 11.11.0
- **Go:** 1.26.1
- **API Port:** 8765
- **Web Port:** 5174

---

## ✅ Verification

Last verified: March 9, 2025

- ✓ API Server: Running (PID 50065)
- ✓ Web UI: Running (PID 50224)
- ✓ External Access: Configured
- ✓ All Endpoints: Responding
- ✓ Network: Accessible

---

## 📞 Support

For issues or questions:
1. Check the logs: `/tmp/api_server.log`, `/tmp/web_server.log`
2. Review the troubleshooting section above
3. Check the main README.md for documentation
4. Verify network connectivity between computers

---

**Your RedditView application is ready for use! 🚀**
