# RedditView Application - System Requirements & Setup Guide

**Version:** 2.0.0  
**Last Updated:** March 11, 2026

## 📋 System Requirements

### Minimum Requirements

#### Operating System
- **Linux** (Ubuntu 20.04+, Debian 11+, CentOS 8+, etc.)
- **macOS** (10.15+)
- **Windows** (10 or 11 with WSL2 recommended)

#### Node.js & npm
- **Node.js:** v18.0.0 or higher (v25.7.0 tested and recommended)
- **npm:** v8.0.0 or higher (v11.11.0 tested and recommended)
- Check versions: `node --version` and `npm --version`

#### System Resources
- **CPU:** 2+ cores (for development)
- **RAM:** 4GB minimum (8GB+ recommended)
- **Disk Space:** 2GB for node_modules + build artifacts
- **Network:** Internet connection required (for fetching Reddit data)

---

## 🏗️ Application Architecture

RedditView is a **monorepo** with the following structure:

```
redditiew/
├── packages/
│   ├── core/          # Shared business logic & data models (@redditview/core)
│   └── web/           # React Web UI (@redditview/web)
├── api-server.js      # HTTP API Server (Port 8765)
├── web-server.js      # Express Web Server (Port 5174)
├── proxy.js           # Development proxy (optional)
├── vite.config.js     # Vite build configuration
├── config.json        # Application configuration
└── dist/              # Production build output
```

---

## 🗂️ Project Structure

### Core Package (`packages/core/`)
- **Purpose:** Shared TypeScript business logic
- **Main Exports:**
  - `./dist/api` - Reddit API client
  - `./dist/models` - Data type definitions
  - `./dist/utils` - Helper utilities
  - `./dist/cache` - Caching system
- **Build Output:** TypeScript compiled to JavaScript in `dist/`

### Web Package (`packages/web/`)
- **Purpose:** React-based Web UI
- **Stack:** React 19 + TypeScript
- **Build Tool:** Vite
- **Build Output:** Optimized SPA in `dist/`

---

## 📦 Dependencies

### Runtime Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `react` | ^19.2.0 | Frontend UI framework |
| `react-dom` | ^19.2.0 | React DOM rendering |
| `react-router-dom` | ^7.13.0 | Client-side routing |
| `express` | ^5.2.1 | Web server framework |
| `node-fetch` | ^3.3.2 | HTTP client for API calls |

### Development Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `typescript` | ^5.0.0 | Type checking |
| `vite` | ^7.3.1 | Build tool & dev server |
| `@vitejs/plugin-react` | ^5.1.1 | Vite React plugin |
| `eslint` | ^9.39.1 | Code linting |
| `@types/react` | ^19.2.7 | React TypeScript definitions |
| `@types/react-dom` | ^19.2.3 | React DOM TypeScript definitions |

---

## 🚀 Installation & Setup

### Step 1: Clone Repository
```bash
git clone https://github.com/nicthegarden/redditiew.git
cd redditiew
```

### Step 2: Install Dependencies
```bash
npm install
```

This will:
- Install root dependencies
- Install workspace dependencies (`packages/core` and `packages/web`)
- Set up monorepo linking

### Step 3: Build Core Package
```bash
npm run build --workspace=@redditview/core
```

Or use the root build command:
```bash
npm run build
```

---

## 🛠️ Available npm Scripts

### Development
```bash
npm run dev           # Start Vite dev server (port 5173)
npm run dev:api      # Start API server (port 8765)
npm run dev:proxy    # Start development proxy
npm run dev:vite     # Start Vite dev server explicitly
npm run dev:core     # Watch TypeScript compilation for core package
```

### Production
```bash
npm run build        # Build for production (compiles TS + bundles with Vite)
npm run preview      # Preview production build locally (port 4173)
```

### Development Tools
```bash
npm run lint         # Run ESLint on all files
```

---

## 📝 Configuration

### `config.json`
Centralized configuration file for the application:

```json
{
  "tui": {
    "default_subreddit": "sysadmin",     // Default subreddit to load
    "posts_per_page": 200,
    "default_sort": "popular"            // hot or popular
  },
  "web": {
    "default_subreddit": "sysadmin",     // Web UI default subreddit
    "posts_per_page": 20,
    "theme": "dark"                      // dark or light
  },
  "api": {
    "base_url": "http://localhost:8765/api",
    "timeout_seconds": 10
  }
}
```

**Environment Variables (optional):**
```bash
VITE_API_BASE_URL=http://localhost:8765/api
PORT=5174  # For web-server.js
```

---

## 🌐 Server Architecture

RedditView uses a **local proxy server architecture** to safely fetch Reddit data:

### API Server (Port 8765)
- **File:** `api-server.js`
- **Purpose:** Provides HTTP endpoints for Reddit API access
- **Features:**
  - Server-side Reddit data fetching
  - 60-second response caching
  - CORS support
  - Public endpoints only (no authentication required)

**Key Endpoints:**
```
GET /api/r/:subreddit.json          # Get posts from subreddit
GET /api/r/:subreddit/search.json   # Search within subreddit
GET /api/:permalink/comments.json   # Get post comments
```

### Web Server (Port 5174)
- **File:** `web-server.js`
- **Purpose:** Serves static React build + proxies API requests
- **Features:**
  - Serves `dist/` directory
  - Proxies `/api` requests to API server
  - Production deployment ready

### Vite Dev Server (Port 5173)
- **File:** Configured in `vite.config.js`
- **Purpose:** Development server with hot module replacement
- **Proxy:** Configured to route `/api` to Reddit API during development

---

## 🔌 Running the Application

### Development Mode (Recommended)

**Terminal 1: API Server**
```bash
npm run dev:api
# Starts on http://localhost:8765
```

**Terminal 2: Vite Dev Server**
```bash
npm run dev
# Starts on http://localhost:5173
# Access: http://localhost:5173
```

**Features in Dev Mode:**
- Hot Module Replacement (HMR)
- Source maps for debugging
- Live reload on file changes

### Production Mode

**Step 1: Build**
```bash
npm run build
# Creates optimized build in dist/
```

**Step 2: Serve**
```bash
npm run preview
# Serves production build on http://localhost:4173
```

Or use the web server:
```bash
npm run dev:api &           # Start API server in background
node web-server.js          # Start web server
# Serves on http://localhost:5174
```

---

## 🌍 Network & API Access

### Reddit API Access
- **Method:** Public endpoints (no authentication required)
- **Base URL:** `https://old.reddit.com`
- **Rate Limiting:** Managed by proxy server (60-second caching)
- **Data Format:** JSON

### Network Requirements
- Outbound HTTPS access to `old.reddit.com` required
- No proxy authentication needed
- Can work behind corporate firewalls if HTTPS outbound is allowed

---

## 🔒 Security Considerations

### What RedditView Does NOT Do
- ❌ Does NOT store Reddit credentials
- ❌ Does NOT require OAuth authentication
- ❌ Does NOT upload data anywhere
- ❌ Does NOT track user behavior

### What RedditView DOES Do
- ✅ Fetches public Reddit data
- ✅ Caches responses locally (60-second TTL)
- ✅ Runs entirely on your machine
- ✅ Uses HTTPS for API calls

---

## 📊 Build Output

### Production Build (`npm run build`)
Creates optimized files in `dist/`:
- `index.html` - Entry HTML file
- `assets/index-*.css` - Minified CSS (13KB gzipped)
- `assets/index-*.js` - Minified JavaScript (208KB, 65KB gzipped)

**File Sizes:**
- Total: ~220KB uncompressed
- Gzipped: ~68KB
- Load time: <1 second on typical internet

---

## 🧪 Testing the Installation

### Verify Installation
```bash
# Check Node/npm
node --version  # Should be v18+
npm --version   # Should be v8+

# Check build works
npm run build

# Test dev server
npm run dev
# Open http://localhost:5173 in browser
```

### Test API Connectivity
```bash
# In a new terminal
curl http://localhost:8765/api/r/sysadmin.json | head -20
```

---

## 🆘 Troubleshooting

### Port Already in Use
```bash
# Find process using port
lsof -i :5173  # For Vite dev server
lsof -i :8765  # For API server
lsof -i :5174  # For web server

# Kill process
kill -9 <PID>
```

### Module Not Found
```bash
# Clean install
rm -rf node_modules package-lock.json
npm install
npm run build
```

### Build Fails
```bash
# Ensure TypeScript compilation works
npm run build --workspace=@redditview/core

# Check for TypeScript errors
npx tsc --noEmit
```

### API Not Responding
```bash
# Check API server is running
curl http://localhost:8765/api/r/test.json

# Check network connectivity
curl https://old.reddit.com/r/test.json
```

---

## 📈 Performance Characteristics

### Build Time
- Development: ~500ms (Vite)
- Production: ~550ms (with TypeScript compilation)

### Runtime Performance
- Initial load: <1 second
- Post navigation: <100ms
- Comment loading: 1-2 seconds (API dependent)
- Memory usage: ~100MB for web UI

### Network Usage
- Per subreddit load: ~50KB (gzipped)
- Per post comments: ~20-100KB (variable)
- Cached responses: No network usage for 60 seconds

---

## 🚢 Deployment Options

### Option 1: Local Development
```bash
npm run dev:api &
npm run dev
```

### Option 2: Production Build + Preview
```bash
npm run build
npm run preview
```

### Option 3: Using Web Server
```bash
npm run build
npm run dev:api &
node web-server.js
```

### Option 4: Docker (Requires Dockerfile)
Not currently provided but can be created based on Node.js image.

---

## 📝 Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `VITE_API_BASE_URL` | `http://localhost:8765/api` | API server base URL |
| `PORT` | `5174` | Web server port |
| `NODE_ENV` | `development` | Environment mode |

---

## 🔄 Workflow Summary

```
1. Install: npm install
2. Build:   npm run build
3. Develop: npm run dev:api & npm run dev
4. Test:    Open http://localhost:5173
5. Deploy:  npm run build && npm run preview
```

---

## 📚 Additional Resources

- **GitHub:** https://github.com/nicthegarden/redditiew
- **React Docs:** https://react.dev
- **Vite Docs:** https://vitejs.dev
- **TypeScript Docs:** https://www.typescriptlang.org/docs
- **Node.js Docs:** https://nodejs.org/docs

---

## ✅ Verification Checklist

- [ ] Node.js v18+ installed
- [ ] npm v8+ installed
- [ ] Git repository cloned
- [ ] `npm install` completed successfully
- [ ] `npm run build` completes without errors
- [ ] `npm run dev:api` starts on port 8765
- [ ] `npm run dev` starts on port 5173
- [ ] Can access http://localhost:5173 in browser
- [ ] Can load subreddit posts
- [ ] Can select and view post details
- [ ] Can view post comments
- [ ] Version shows as v2.0.0

---

**Document Version:** 2.0.0  
**Generated:** March 11, 2026
