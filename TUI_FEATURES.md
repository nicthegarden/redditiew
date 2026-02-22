# RedditView TUI - Split-View Features

## Overview

The RedditView Terminal User Interface now includes a **professional split-view layout** for browsing Reddit posts with real-time search, filtering, and detailed post viewing.

## Layout

```
┌─────────────────────────────────────────────────────────────┐
│ 🔥 r/golang                                                 │
├─────────────────────────────────────────────────────────────┤
│ 🔍 Search: rust  ▌                                           │
├──────────────────────────────┬──────────────────────────────┤
│ Left Sidebar                 │ Right Sidebar                │
│ (Post List)                  │ (Post Details)               │
│                              │                              │
│ ❯ Rust vs Go comparison      │ ╭────────────────────────╮  │
│   u/john_dev                 │ │ Rust vs Go comparison  │  │
│   ⬆ 3240  💬 156             │ │                        │  │
│                              │ │ 👤 u/john_dev         │  │
│ • Memory safety in Go        │ │ ⬆ 3240 upvotes        │  │
│   u/alice_rust               │ │ 💬 156 comments        │  │
│   ⬆ 2891  💬 203             │ │                        │  │
│                              │ │ Content:               │  │
│ • Concurrency patterns       │ │ This post compares... │  │
│   u/bob_gopher               │ │                        │  │
│   ⬆ 2445  💬 89              │ │ Link: (external URL)   │  │
│                              │ │ ╰────────────────────────╯  │
│ (more posts...)              │                              │
│                              │                              │
├──────────────────────────────┴──────────────────────────────┤
│ ↑↓/jk: navigate · /: search · ENTER: view · q: quit         │
└──────────────────────────────────────────────────────────────┘
```

## Features

### 1. Split-View Display

**Left Sidebar - Post List**
- Shows subreddit posts in a scrollable list
- Each post displays:
  - Title (truncated to fit width)
  - Author (username)
  - Upvote count (⬆)
  - Comment count (💬)
- Selected post highlighted with orange background (#FF6B35)
- Reddit-inspired color scheme (orange, gold, green)

**Right Sidebar - Post Details**
- Shows full details of selected post:
  - Complete title (wrapped text)
  - Author with 👤 icon
  - Full statistics (upvotes and comment count)
  - Post content/selftext (wrapped and formatted)
  - External link if available (blue, underlined)
- Graceful "Select a post to view" message when no post selected

### 2. Search and Filtering

**Activate Search Mode**
- Press `/` key to enter search mode
- Search bar appears at top with 🔍 icon and text cursor (▌)
- Transitions from "🔍 Press '/' to search · Found X posts" to "🔍 Search: (your query) ▌"

**Type to Filter**
- Posts are filtered in real-time by:
  - Title (case-insensitive)
  - Author username (case-insensitive)
- Post count updates as you type: "Found X posts"
- Filtered post list updates immediately in left sidebar
- Right sidebar updates to show selected filtered post

**Exit Search**
- `ESC` key: Cancel search and restore full list
- `ENTER` key: Apply search and continue searching
- `Backspace` key: Delete last character from search query

### 3. Navigation

**Movement Controls**
- `↑` or `k`: Move up in post list
- `↓` or `j`: Move down in post list
- Navigation wraps at boundaries (cannot scroll past first/last post)
- Selected post automatically shows in right pane

**Post Selection**
- `ENTER` key: Explicitly select current post for detail view
- Automatic selection when navigating with arrow keys

**Quit Application**
- `q` key: Close TUI cleanly
- `Ctrl+C`: Force quit (emergency escape)

### 4. Responsive Design

- Window size detection on startup
- Dynamic calculation of sidebar widths (50/50 split)
- Text wrapping for titles and content
- Proper padding and spacing
- Automatic adjustment to terminal window changes

### 5. Visual Elements

**Color Scheme**
- Header background: Reddit orange (#FF4500)
- Selected item: Lighter orange (#FF6B35)
- Text: White (#FFFFFF)
- Author/username: Gold (#FFD700)
- Stats: Green (#90EE90)
- Post content: Light gray (#CCCCCC)
- Links: Sky blue (#87CEEB)
- Footer: Dark gray (#333333)
- Dividers: Orange (#FF4500)

**Icons**
- 🔥 - Subreddit header
- 🔍 - Search
- 👤 - Author
- ⬆ - Upvotes
- 💬 - Comments
- ❯ - Selected item indicator
- ▌ - Search cursor

## Usage Examples

### Basic Browsing
```
1. Launch: ./launch.sh tui
2. See first 30 posts from r/golang
3. Use ↑↓ or jk to browse
4. Press ENTER to view selected post details
5. Press q to quit
```

### Search by Keyword
```
1. Press / to enter search mode
2. Type: "rust" (filters to posts with "rust" in title/author)
3. Results update in real-time
4. Use ↑↓ to navigate filtered results
5. Press ESC to cancel search and restore full list
```

### View Post Details
```
1. Navigate to a post with ↑↓/jk
2. Right pane auto-updates with full post details
3. View title, author, content, and external link
4. Text wraps automatically for readability
5. Use ↑↓/jk to move to next post or ESC to search again
```

## Technical Details

### Implementation

**File:** `/home/nd/GIT/redditiew-local/apps/tui/main.go`

**Key Functions:**
- `renderSplitView()` - Main rendering function combining both sidebars
- `renderPostList()` - Left sidebar with post list
- `renderPostDetails()` - Right sidebar with post details
- `renderSearchBar()` - Search input or instruction bar
- `renderFooter()` - Context-aware control hints
- `filterPosts()` - Real-time post filtering logic
- `wrapText()` - Text wrapping for long content
- `truncateTitle()` - Truncate titles to fit width

**Model Fields:**
```go
type Model struct {
    posts         []RedditPostData    // All loaded posts
    filteredPosts []RedditPostData    // Filtered by search
    selected      int                 // Index of selected post
    loading       bool                // Loading state
    error         string              // Error message
    sub           string              // Subreddit name
    client        *APIClient          // Reddit API client
    searchQuery   string              // Current search text
    inputMode     bool                // In search mode?
    selectedPost  *RedditPostData     // Detailed view post
    windowWidth   int                 // Terminal width
    windowHeight  int                 // Terminal height
}
```

**API Integration:**
- Fetches 30 posts from specified subreddit
- Uses Go's Bubble Tea framework for TUI
- Lipgloss for styling and colors
- HTTP client for Reddit API communication

## Performance Notes

- Posts loaded on startup (30 posts from r/golang)
- Search filters in-memory (no new API calls)
- Real-time filtering with each keystroke
- Responsive layout updates
- No lag with up to 30+ posts
- Text wrapping computed dynamically

## Requirements

**Runtime:**
- Go 1.16+
- Terminal with ANSI color support
- POSIX-compatible shell (for launch.sh)
- Node.js 16+ (for API server)

**Dependencies (Go):**
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling

**API Server:**
- Node.js API server running on port 3002
- Proxies to Reddit API (www.reddit.com)
- Caches responses (1 minute TTL)

## Troubleshooting

### TUI Won't Start
```
Error: "connect: connection refused"
→ Make sure API server is running: ./launch.sh api
```

### Layout Broken on Small Terminal
```
Error: Text overlapping or missing
→ Expand terminal window or use: stty rows 40 cols 120
```

### Search Not Working
```
Error: Filtering not updating
→ Verify you're in input mode (typing in search bar)
→ Check that posts loaded successfully
```

### Colors Look Wrong
```
Error: Color display issue
→ Terminal might not support 24-bit color
→ Try: export COLORTERM=truecolor
```

## Future Enhancements

Potential improvements for future versions:
- [ ] Sort posts by upvotes, comments, or date
- [ ] Load more posts (pagination)
- [ ] Subscribe to multiple subreddits (tabs)
- [ ] Comment viewing and navigation
- [ ] Save/bookmark posts
- [ ] User profile viewing
- [ ] Advanced search filters
- [ ] Customizable color schemes
- [ ] Mouse support
- [ ] Vim keybindings customization

---

**Status:** ✅ Fully Implemented and Tested
**Last Updated:** 2026-02-22
