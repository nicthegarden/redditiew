# RedditView - Development Guide

## Overview

RedditView is a keyboard-friendly Reddit browser built with React + TypeScript and Vite, featuring:
- Split-pane interface (posts list on left, content on right)
- Email-client-like browsing experience
- Subreddit navigation with smart caching
- Post filtering and Reddit-wide search
- Rate-limit recovery with exponential backoff
- Mobile-responsive design
- Dark/light theme support

## Architecture

### Frontend Stack
- **React 19.2.0** - UI framework
- **TypeScript** - Type safety
- **Vite 7.3.1** - Build tool & dev server
- **React Router 7.13.0** - Routing ready

### Backend Stack
- **Node.js** - Runtime
- **Express 5.2.1** - HTTP server (optional, proxy uses native http/https)
- **Concurrently** - Run proxy + Vite in parallel

### Project Structure

```
redditview/
├── src/
│   ├── App.tsx              # Main React component
│   ├── main.tsx             # Entry point
│   ├── index.css            # Styling (dark/light themes)
│   └── assets/              # Images
├── proxy.ts                 # CORS proxy server (port 3001)
├── index.html               # HTML template
├── vite.config.ts           # Vite configuration
├── tsconfig.json            # TypeScript config
├── package.json             # Dependencies
└── dist/                    # Production build
```

## Installation & Setup

### Prerequisites
- Node.js 18+
- npm 9+

### Install Dependencies

```bash
cd redditview
npm install
npm install --save-dev typescript @types/node @types/express
npm install concurrently
```

### Start Development

```bash
npm run dev
```

This runs both:
- **Vite dev server** on http://localhost:5173 (auto-reload)
- **Proxy server** on http://localhost:3001 (CORS bypass)

**Or run separately:**
```bash
npm run dev:vite    # Terminal 1
npm run dev:proxy   # Terminal 2
```

### Build for Production

```bash
npm run build
npm run preview    # Preview built version
```

## Features

### 1. Subreddit Browsing
- Enter subreddit name (with/without r/ prefix)
- Quick-link buttons for favorites
- Smart caching (1-hour TTL per subreddit)
- Auto-loads on startup (defaults to r/sysadmin)

### 2. Post List & Selection
- Shows: thumbnail, title, subreddit, time ago, score, comments
- Keyboard navigation: ↑↓ arrow keys, Tab to focus, Enter to select
- Hover effects & visual feedback

### 3. Post Content
- Displays post body/link in iframe
- Comments loaded from Reddit
- Full Reddit interface available

### 4. Local Filtering
- Type in filter box to search post titles
- Real-time results as you type
- Separate from Reddit-wide search

### 5. Reddit-Wide Search
- Press Enter in filter box to search all of Reddit
- Returns results across all subreddits
- Shows search query in subreddit name

### 6. Rate Limit Handling
- Detects 429/503 responses
- Automatic retry with exponential backoff (max 3 attempts)
- User-friendly error messages with retry times
- Cache hits reduce API calls

### 7. Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Cycle through: search → filter → post list |
| `Shift+Tab` | Cycle backwards |
| `↑` / `↓` | Navigate posts (when list focused) |
| `Enter` | Open selected post |
| `Ctrl+F` | Focus iframe for page search |
| `PageUp` / `PageDown` | Scroll post content |

### 8. Mobile Support
- Responsive breakpoints: 768px (tablet), 480px (phone)
- Touch-friendly buttons & spacing
- Stacked layout on small screens

## API Integration

### Reddit API (Public, No Auth)

```
GET https://www.reddit.com/r/{subreddit}.json?limit=50
GET https://www.reddit.com/search.json?q={query}&type=link&limit=50
GET https://www.reddit.com{permalink}.json
```

### Proxy Server

Routes requests through `http://localhost:3001` to:
- Bypass CORS restrictions
- Add caching layer
- Handle rate limits gracefully
- Strip security headers that block iframes

**Endpoints:**
```
/api/*           → proxies to www.reddit.com (cached)
/search/*        → proxies to search API (not cached)
/*               → proxies to old.reddit.com (cached)
```

## Caching Strategy

### Frontend (localStorage)
- Stores posts per subreddit
- TTL: 1 hour per subreddit
- Used on app load if available

### Backend (in-memory)
- Caches raw API responses
- TTL: 60 seconds
- Cleared on new requests (per unique URL)
- Bypassed for search queries

## Error Handling

### Rate Limiting (429)
```
→ Shows: "Rate limited. Please try again in 60 seconds."
→ Proxy retries automatically (exponential backoff, max 3 attempts)
→ Frontend displays user-friendly error with retry time
```

### Not Found (404)
```
→ Shows: "Subreddit r/{name} not found"
→ User can try a different subreddit
```

### Server Error (5xx)
```
→ Shows: "Reddit server error. Please try again later."
→ Retries after user intervention
```

### Network Error
```
→ Shows: "Proxy error: {message}"
→ Check proxy server is running on port 3001
```

## Performance Optimizations

1. **Lazy Loading**
   - Post thumbnails loaded on-demand (`loading="lazy"`)
   
2. **Caching**
   - Frontend: localStorage (1h per subreddit)
   - Backend: in-memory (60s global)
   
3. **Pagination**
   - Loads 50 posts per page
   - Manual "Load More" button
   - Optional auto-scroll trigger (can be enabled)
   
4. **Rendering**
   - React hooks for efficient state updates
   - Filtered lists re-compute only when needed
   - Memoized callbacks prevent unnecessary renders

## TypeScript Types

```typescript
interface RedditPost {
  kind: string
  data: {
    id: string
    title: string
    subreddit: string
    author: string
    created_utc: number
    score: number
    num_comments: number
    thumbnail?: string
    preview?: {
      images: Array<{ source: { url: string } }>
    }
    permalink: string
  }
}

interface CacheEntry {
  posts: RedditPost[]
  time: number
}
```

## Theming

### Dark Theme (Default)
```css
--bg-primary: #0e1113
--bg-secondary: #1a1d21
--accent-orange: #ff4500 (Reddit orange)
--accent-blue: #7193ff
--text-primary: #d7dadc
```

### Light Theme (CSS Variables Ready)
```css
--bg-primary: #ffffff
--bg-secondary: #f5f5f5
--text-primary: #1a1a1a
```

To enable: `document.body.classList.add('light')`

## Troubleshooting

### "Proxy error: connection refused"
- Check proxy running: `lsof -i :3001`
- Restart: `npm run dev:proxy`

### Posts not loading
- Check Network tab (F12) for 429 responses
- Wait for retry or clear cache: `localStorage.clear()`

### TypeScript errors
- Run: `npx tsc --noEmit` to check for type errors
- Import types are in App.tsx interfaces

### Port conflicts
- Vite: `npm run dev:vite -- --port 5174`
- Proxy: Edit `proxy.ts` line 2: `const proxyPort = 3002`

## Future Enhancements

1. Theme toggle button in UI
2. User settings (cache TTL, pagination size)
3. Keyboard shortcut customization
4. Comment tree navigation
5. Post filtering by score/date
6. Multi-subreddit browsing in one view
7. PWA support for offline reading
8. Dark mode auto-detection

## Commits & Version History

All changes tracked in git. View history:
```bash
git log --oneline
git show <commit>
```

Recent changes:
1. ✅ Dev setup with concurrently
2. ✅ Full TypeScript migration
3. ✅ Rate-limit retry logic
4. ✅ Reddit-wide search
5. ✅ Improved pagination UX
6. ✅ Mobile responsiveness
7. ✅ Light theme CSS variables

## Contributing

When adding features:
1. Keep types strict (no `any`)
2. Test keyboard navigation
3. Add responsive CSS (test at 768px, 480px)
4. Update this documentation
5. Commit with clear message

## License

Built with ❤️ for efficient Reddit browsing.
