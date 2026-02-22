# RedditView

A **keyboard-friendly Reddit browser** with a split-pane interface, built for efficient browsing.

![](https://img.shields.io/badge/React-19.2.0-blue) ![](https://img.shields.io/badge/TypeScript-Latest-blue) ![](https://img.shields.io/badge/Vite-7.3.1-green) ![](https://img.shields.io/badge/Node-18+-green)

## Features

✨ **Split-pane interface** — Posts on left, content on right (email-client style)  
⌨️ **Keyboard-first** — Navigate with arrow keys, Tab, Enter  
🔍 **Smart search** — Filter posts locally or search all of Reddit  
💾 **Intelligent caching** — 1-hour per-subreddit cache + 60s proxy cache  
🛡️ **Rate-limit recovery** — Auto-retry with exponential backoff  
📱 **Responsive design** — Works on desktop, tablet, and mobile  
🌙 **Dark theme ready** — Dark mode by default, light theme CSS included  
⚡ **Fast** — Vite hot reload, lazy-loaded thumbnails, optimized rendering  

## Quick Start

```bash
# Install
npm install

# Run (starts proxy + Vite together)
npm run dev

# Opens: http://localhost:5173
```

Then:
1. Enter a subreddit name (e.g., `linux`, `sysadmin`)
2. Use ↑↓ keys to navigate posts
3. Press Enter to view
4. Type to filter posts, press Enter to search Reddit

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Cycle focus: search → filter → post list |
| `Shift+Tab` | Cycle backwards |
| `↑` / `↓` | Navigate posts |
| `Enter` | Open selected post / search Reddit |
| `Ctrl+F` | Focus iframe for page search |

## Configuration

RedditView supports a `config.json` file in the root directory for easy customization of both the TUI and Web interfaces.

### Setup

Create a `config.json` file in the project root:

```json
{
  "tui": {
    "default_subreddit": "golang",
    "posts_per_page": 50,
    "list_height": 12,
    "max_title_length": 80
  },
  "web": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 20,
    "theme": "dark"
  },
  "api": {
    "base_url": "http://localhost:3002/api",
    "timeout_seconds": 10
  }
}
```

### Configuration Options

#### TUI Settings (`tui`)
- **`default_subreddit`** (string): Subreddit loaded on startup (default: `"golang"`)
- **`posts_per_page`** (number): Number of posts to fetch (default: `50`)
- **`list_height`** (number): Maximum lines shown in list view (default: `12`)
- **`max_title_length`** (number): Truncate titles longer than this (default: `80`)

#### Web Settings (`web`)
- **`default_subreddit`** (string): Initial subreddit in web browser (default: `"sysadmin"`)
- **`posts_per_page`** (number): Posts per page in web interface (default: `20`)
- **`theme`** (string): Default theme: `"dark"` or `"light"` (default: `"dark"`)

#### API Settings (`api`)
- **`base_url`** (string): Backend API endpoint (default: `"http://localhost:3002/api"`)
- **`timeout_seconds`** (number): Request timeout in seconds (default: `10`)

### Usage

- **TUI**: Loads config from `../../config.json` relative to the TUI binary directory
- **Web**: Loads config from `/config.json` (served by HTTP)
- If config file is missing, sensible defaults are used

## Architecture

**Frontend:** React 19 + TypeScript + Vite  
**Backend:** Node.js proxy server (CORS bypass + caching)  
**Storage:** localStorage (1h per subreddit) + in-memory proxy cache
**TUI:** Go + Bubble Tea (Terminal UI)

## Recent Updates

✅ **TypeScript migration** — Full type safety  
✅ **Rate-limit handling** — Automatic retry logic  
✅ **Reddit-wide search** — Search all subreddits from filter box  
✅ **Better pagination** — Manual "Load More" button  
✅ **Mobile responsive** — Optimized for phones & tablets  
✅ **Theme support** — Light theme CSS ready  

## Documentation

- **Setup & development:** See [DEVELOPMENT.md](./DEVELOPMENT.md)
- **Feature specs:** See [SPEC.md](./SPEC.md)

## Project Structure

```
src/
├── App.tsx          # Main component (React + TS)
├── main.tsx         # Entry point
└── index.css        # Dark/light themes

proxy.ts            # CORS proxy (Node.js)
vite.config.ts      # Build config
tsconfig.json       # Type checking
```

## Commands

```bash
npm run dev         # Start dev (proxy + Vite)
npm run dev:vite    # Vite only
npm run dev:proxy   # Proxy only
npm run build       # Production build
npm run preview     # Preview built version
npm run lint        # ESLint check
```

## Troubleshooting

**"Proxy error: connection refused"**
- Make sure proxy is running on port 3001
- Check: `lsof -i :3001`

**"Rate limited" error**
- Wait 60 seconds (proxy retries automatically)
- Or clear cache: `localStorage.clear()`

**Posts not loading**
- Check Network tab (F12) for errors
- Verify subreddit exists

**TypeScript errors**
- Run: `npx tsc --noEmit`

## Future Ideas

- [ ] Theme toggle button in UI
- [ ] User settings (cache TTL, pagination size)
- [ ] Keyboard shortcut customization
- [ ] Post filtering by score/date/time
- [ ] Multi-subreddit view
- [ ] PWA for offline reading
- [ ] Auto dark-mode detection

## License

Built with ❤️ for efficient Reddit browsing.

---

**Need help?** Check [DEVELOPMENT.md](./DEVELOPMENT.md) for detailed setup and architecture info.
