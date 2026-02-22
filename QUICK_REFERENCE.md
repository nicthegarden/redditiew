# 🚀 RedditView Quick Reference & Checklists

## Installation Checklist

- [ ] Clone the repository
- [ ] Install Node.js 16+ and Go 1.21+
- [ ] Run `npm install` in project root
- [ ] Run `npm run build` to build core package

## Running the Web App

```bash
./launch.sh web
```

Then open http://localhost:5173 in browser.

**What it does**:
- Starts API server on port 3002
- Starts React web app on port 5173
- Vite proxy routes `/api` requests to API server

**Port**: 5173  
**API Port**: 3002

## Running the Go TUI

### Single Command (Recommended)
```bash
./launch.sh tui
```

This starts:
1. API Server on port 3002
2. Go TUI app in terminal

### Manual Setup (Two Terminals)

**Terminal 1: Start API Server**
```bash
node api-server.js
```
Verify: `curl http://localhost:3002/health`

**Terminal 2: Run TUI App**
```bash
cd apps/tui
go run main.go
```

**Controls**:
- `↑`/`↓` or `k`/`j` - Navigate posts
- `q` or `Ctrl+C` - Quit

## Monorepo Commands Reference

### Building
```bash
npm run build              # Build everything
npm run dev:core          # Watch and rebuild core
npm run build --workspace=@redditview/core  # Build just core
```

### Development
```bash
./launch.sh web                    # Start web app (port 5173)
./launch.sh tui                    # Start TUI app (uses API on 3002)
./launch.sh all                    # Start web + TUI + API
./launch.sh api                    # Start API server only (port 3002)

# Alternative (manual control):
npm run dev                        # Start web app (Vite, port 5173)
node api-server.js                # Start API server (port 3002)
cd apps/tui && go run main.go     # Start TUI app
```

### Maintenance
```bash
npm install              # Install all dependencies
npm run lint             # Lint code
npm run preview          # Preview production build
```

## Core Package Usage

### Importing from Core
```typescript
// Models
import type { 
  RedditPost, 
  Comment, 
  FetchPostsResult 
} from '@redditview/core'

// API Client
import { RedditApiClient, apiClient } from '@redditview/core'

// Cache
import { PostCache, LocalStorageCache } from '@redditview/core'

// Utils
import { 
  formatTime, 
  formatNum, 
  getPostThumbnail 
} from '@redditview/core'
```

### API Client Examples
```typescript
// Create client
const client = new RedditApiClient({
  baseUrl: '/api',      // For web
  timeout: 15000
})

// Fetch posts
const result = await client.fetchPosts('javascript', {
  limit: 50,
  useCache: true
})

// Fetch comments
const comments = await client.fetchComments('/r/javascript/comments/abc123/')

// Search
const search = await client.search('TypeScript', { limit: 20 })
```

## Ports Cheat Sheet

| Service | Port | Command |
|---------|------|---------|
| React Web | 5173 | `npm run dev` |
| API Server | 3002 | `npm run dev:api` |
| Old Proxy | 3001 | `npm run dev:proxy` |
| Reddit | 443 | (external) |

## File Locations

### Core Package
- Models: `packages/core/src/models/index.ts`
- API Client: `packages/core/src/api/index.ts`
- Caching: `packages/core/src/cache/index.ts`
- Utilities: `packages/core/src/utils/index.ts`
- Exports: `packages/core/src/index.ts`

### Web App
- Main: `packages/web/src/App.tsx`
- Components: `packages/web/src/components/`
- Styles: `packages/web/src/index.css`
- Config: `packages/web/vite.config.ts`

### TUI App
- Main: `apps/tui/main.go`
- Module: `apps/tui/go.mod`

### Backend
- API Server: `api-server.js` (root, main app starter)
- Old Proxy: `proxy.ts` (deprecated)
- Launcher: `launch.sh` (multi-platform starter)

## Common Tasks

### Add a New Feature to Core

1. **Define types in core/src/models/**
```typescript
export interface NewType {
  field: string
}
```

2. **Add logic in core/src/api/ or other modules**
```typescript
async newMethod() {
  // implementation
}
```

3. **Export from core/src/index.ts**
```typescript
export { newMethod, NewType }
```

4. **Build**
```bash
npm run build
```

5. **Use in web**
```typescript
import { NewType, newMethod } from '@redditview/core'
```

6. **Use in TUI**
```bash
# If it requires a new API endpoint, add it to api-server.ts
# Then call from Go
resp, _ := http.Get("http://localhost:3002/api/new-endpoint")
```

### Update React App

1. Edit files in `packages/web/src/`
2. Web app hot-reloads automatically
3. If you changed core, rebuild: `npm run build`

### Update TUI App

1. Edit `apps/tui/main.go`
2. Changes take effect on next `go run`
3. For hot reload, install and use `air`: `go install github.com/cosmtrek/air@latest && air`

### Build for Production

```bash
# Build core and web
npm run build

# Output in packages/web/dist/

# Build TUI executable
cd apps/tui
go build -o redditview main.go
```

## Debugging

### Check Web App
```bash
# Open browser DevTools (F12)
# Check Application > Cache Storage
# Check Network tab for /api calls
```

### Check TUI App
```bash
# Check API server logs
npm run dev:api

# Check if API server is responding
curl http://localhost:3002/health

# Check if Go TUI can reach server
curl http://localhost:3002/api/r/golang.json?limit=5
```

### Check Core Package
```bash
# Rebuild to catch TypeScript errors
npm run build

# Check packages/core/dist/ for compiled code
ls packages/core/dist/

# Check exports
node -e "import('packages/core/dist/index.js').then(m => console.log(Object.keys(m)))"
```

## Troubleshooting by Symptom

### "Cannot find module '@redditview/core'"
```bash
npm install
npm run build
```

### Web app shows empty posts
```bash
# Check browser console for errors
# Verify Vite proxy is working
# Check Network tab - look for /api/r/... requests
# Make sure comments load (separate endpoint)
```

### TUI shows "Loading posts..." forever
```bash
# Check if API server is running:
curl http://localhost:3002/health

# If not, start it:
node api-server.js

# Then start TUI:
cd apps/tui && go run main.go
```

### Port X already in use
```bash
# Find process using port
lsof -i :PORT_NUMBER

# Kill it
kill -9 PID

# Or use different port by editing config
```

### TypeScript errors when building
```bash
npm run build  # Shows all errors
npm run dev:core  # Watch mode to fix as you go
```

### Go compilation errors
```bash
cd apps/tui
go mod tidy    # Fix dependencies
go run main.go # Run and see errors
```

## Performance Tips

1. **Enable caching** - Always use `useCache: true` in API calls
2. **Lazy load** - Load posts/comments only when needed
3. **Pagination** - Use `limit` parameter to reduce data transfer
4. **Batch requests** - Combine multiple API calls when possible
5. **Monitor cache** - Check `npm run dev:api` logs for cache hits

## Deployment Ideas

### Web App
```bash
npm run build
# Deploy packages/web/dist/ to:
# - Vercel
# - Netlify
# - GitHub Pages
# - AWS S3
```

### TUI App
```bash
cd apps/tui
go build -o redditview main.go
# Distribute redditview binary
# Can be run on any OS with Go runtime
```

### API Server
```bash
# Run on a server (Heroku, DigitalOcean, AWS, etc.)
# Environment: Node.js 16+
# Port: 3002
# Or containerize with Docker
```

## Version Control

```bash
# Commit all monorepo changes together
git add .
git commit -m "feat: add cool feature to core and use in web"

# View changes across workspaces
git log --oneline

# Check what changed in core
git log packages/core/
```

## Getting Help

1. **Check documentation files**
   - SETUP_COMPLETE.md - You are here!
   - MONOREPO_ARCHITECTURE.md - Deep dive
   - TUI_SETUP_GUIDE.md - TUI specifics
   - ARCHITECTURE.md - Diagrams and flows

2. **Check relevant README files**
   - packages/core/ - Core package docs
   - packages/web/ - Web app docs
   - apps/tui/ - TUI app docs

3. **Debug step by step**
   - Read error messages carefully
   - Use curl to test API endpoints
   - Check browser DevTools for web issues
   - Check terminal output for build errors

4. **Try rebuilding**
   ```bash
   npm run build
   npm install
   ```

---

**Last Updated**: Feb 22, 2026  
**Monorepo Version**: 1.0  
**Status**: ✅ Ready for development
