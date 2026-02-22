# ✅ RedditView Monorepo - Setup Complete!

Congratulations! Your RedditView project has been restructured as a professional monorepo with a shared core that both your React web app and Go TUI app can use.

## 📦 What Was Created

### 1. **Shared Core Package** (`packages/core/`)
A reusable TypeScript library with:
- **Models**: Post, Comment, CacheEntry, API response types
- **API Client**: RedditApiClient with fetchPosts(), fetchComments(), search()
- **Cache**: In-memory PostCache and browser LocalStorageCache
- **Utils**: formatTime(), formatNum(), getPostThumbnail(), etc.

**Published as**: `@redditview/core` (npm package)

### 2. **React Web App** (`packages/web/`)
Your existing React app, now refactored to:
- Use `@redditview/core` for all business logic
- Leverage Vite's built-in proxy (no need for separate Node.js proxy)
- Share the same data models as TUI app

### 3. **Go TUI App** (`apps/tui/`)
A new Terminal UI application built with Bubble Tea:
- Fetches data from Node.js API server (port 3002)
- Uses same data models as React app
- Full keyboard navigation support
- Displays Reddit posts in a beautiful terminal interface

### 4. **API Server** (`api-server.ts`)
A Node.js backend server (port 3002) that:
- Serves API endpoints for the Go TUI
- Caches responses to reduce Reddit API calls
- Proxies requests to Reddit API with proper headers
- Can be reused by other clients (mobile, desktop, etc.)

### 5. **Documentation**
Complete guides for:
- `MONOREPO_ARCHITECTURE.md` - Technical architecture overview
- `TUI_SETUP_GUIDE.md` - Step-by-step setup and usage
- `ARCHITECTURE.md` - System diagrams and data flow

## 🚀 Quick Start

### 🌐 Run React Web App (Easiest)
```bash
npm install
npm run build
./launch.sh web
# Visit http://localhost:5173
```

### 💻 Run Go TUI App
```bash
npm install
npm run build
./launch.sh tui
# Navigate with ↑↓/jk, press q to quit
```

### 🔄 Run Both Together
```bash
npm install
npm run build
./launch.sh all
```

### 🛠️ Run API Server Only
```bash
npm install
npm run build
./launch.sh api
```

## 📁 File Structure

```
redditiew-monorepo/
├── packages/
│   ├── core/
│   │   ├── src/
│   │   │   ├── api/index.ts        → RedditApiClient
│   │   │   ├── models/index.ts     → Data types
│   │   │   ├── cache/index.ts      → Cache implementations
│   │   │   ├── utils/index.ts      → Helper functions
│   │   │   └── index.ts            → Main exports
│   │   └── package.json
│   └── web/
│       ├── src/App.tsx             → Uses @redditview/core
│       ├── src/components/
│       └── vite.config.ts
├── apps/
│   └── tui/
│       ├── main.go                 → Bubble Tea TUI
│       └── go.mod
├── api-server.js                   → Node.js backend (port 3002)
├── api-server.ts                   → TypeScript version (deprecated)
├── launch.sh                        → Multi-platform launcher (recommended)
├── proxy.ts                        → Old proxy (deprecated)
├── package.json                    → Monorepo with workspaces
├── MONOREPO_ARCHITECTURE.md
├── TUI_SETUP_GUIDE.md
└── ARCHITECTURE.md
```

## 🔄 How to Use the Shared Core

### In React Web App
```typescript
import { 
  RedditApiClient, 
  formatTime, 
  formatNum,
  type RedditPost 
} from '@redditview/core'

const client = new RedditApiClient()
const result = await client.fetchPosts('javascript')
console.log(formatNum(result.posts[0].data.score))
```

### In Go TUI App
```go
// Models defined in Go matching TypeScript definitions
type RedditPostData struct {
    ID      string
    Title   string
    Score   int
    // ...
}

// Call API server endpoints
resp, _ := http.Get("http://localhost:3002/api/r/golang.json")
```

## 📊 Architecture Overview

```
React Web (localhost:5173)
    ↓
Vite Proxy (/api → old.reddit.com)
    ↓
Reddit API

Go TUI (Terminal)
    ↓
API Server (localhost:3002)
    ↓
Reddit API

Both use:
    ↓
@redditview/core
(shared models, API client, utilities)
```

## ✨ Key Features

✅ **Shared Core Logic**
- One implementation used by multiple apps
- Consistent data models everywhere
- Easy to maintain and update

✅ **Type Safety**
- TypeScript interfaces guide API contracts
- Compile-time type checking
- Better IDE support and documentation

✅ **Flexible Architecture**
- Web uses Vite proxy (simple, no backend needed)
- TUI uses Node.js backend (can be deployed separately)
- Easy to add more clients (mobile, desktop, CLI)

✅ **Professional Structure**
- npm workspaces for dependency management
- Monorepo best practices
- Clear separation of concerns
- Comprehensive documentation

## 📝 Next Steps

### Immediate
1. Test the web app: `npm run dev`
2. Test the TUI: `npm run dev:api` + `go run main.go` in apps/tui/
3. Review the documentation files

### Short Term
- [ ] Add more features to TUI (post details, comments)
- [ ] Implement search functionality in TUI
- [ ] Add pagination support
- [ ] Add unit tests for core package

### Medium Term
- [ ] Setup CI/CD pipeline
- [ ] Deploy React app to Vercel/Netlify
- [ ] Build TUI executable with `go build`
- [ ] Create Docker containers for distribution

### Long Term
- [ ] Add mobile app using React Native
- [ ] Add desktop app using Electron
- [ ] Expand API server with user accounts
- [ ] Add real-time updates with WebSockets

## 🆘 Troubleshooting

**Web app won't load:**
```bash
npm install
npm run build
npm run dev
```

**TUI says "Error: Network error":**
```bash
# Make sure API server is running
npm run dev:api
# Check it's responding
curl http://localhost:3002/health
```

**Port already in use:**
```bash
# Find what's using the port
lsof -i :3002  # or :5173 for web

# Kill it if needed
kill -9 <PID>
```

**Go dependencies not found:**
```bash
cd apps/tui
go mod download
go run main.go
```

## 📚 Documentation

- **MONOREPO_ARCHITECTURE.md** - Deep dive into architecture
- **TUI_SETUP_GUIDE.md** - Complete TUI setup and usage
- **ARCHITECTURE.md** - System diagrams and data flow
- **packages/core/README.md** - Core package documentation
- **apps/tui/README.md** - TUI app documentation

## 💡 Tips

1. **Develop core features first** - Add types/logic to `packages/core`
2. **Use in both apps** - React and Go both benefit from shared code
3. **Keep models in sync** - TypeScript types drive API contracts
4. **Test thoroughly** - Core logic should have unit tests
5. **Deploy separately** - Web on Vercel, TUI as executable, API on a server

## 🎯 Benefits Summary

| Aspect | Before | After |
|--------|--------|-------|
| Code Duplication | High | Low |
| Code Reuse | Limited | Extensive |
| Maintenance | Multiple places | Single source of truth |
| Type Safety | Partial | Full |
| Time to Add Features | Slow | Fast |
| Number of Clients | 1 (web) | Many (web, TUI, future) |

## 🤝 Contributing

When adding features:
1. **Identify if it's core logic** - Does both web and TUI need it?
2. **Add to core first** - `packages/core/src/`
3. **Update exports** - Add to `packages/core/src/index.ts`
4. **Use in web** - Import from `@redditview/core`
5. **Use in TUI** - Call via HTTP API server
6. **Test both** - Ensure works in all clients

## 📞 Support

- Check the relevant documentation file
- Review error messages carefully
- Try rebuilding: `npm run build`
- Verify ports are available (5173, 3002)
- Check internet connection (need Reddit API access)

---

**You now have a professional, scalable monorepo architecture ready for multiple clients!** 🎉

Good luck with your RedditView project! 🚀
