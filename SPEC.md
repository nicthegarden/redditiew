# Reddit Viewer - Split Pane Webapp

## Project Overview
- **Type**: Single Page Web Application (React + Vite)
- **Core functionality**: Reddit browser with dual-pane layout - posts list on left, content/comments on right
- **Target users**: Reddit users who want an email-client-like browsing experience

## UI/UX Specification

### Layout Structure
- **Split view**: Fixed 40% left pane / 60% right pane (resizable divider)
- **Left pane**: Subreddit selector + posts list
- **Right pane**: Post content + comments thread
- **Responsive**: On mobile (<768px), show single pane with navigation back

### Visual Design

**Color Palette (Dark Theme - Reddit-inspired)**
- Background primary: `#0e1113` (dark)
- Background secondary: `#1a1d21` (card/list bg)
- Background tertiary: `#242729` (hover states)
- Accent orange: `#ff4500` (Reddit orange)
- Accent blue: `#7193ff` (links)
- Text primary: `#d7dadc`
- Text secondary: `#818384`
- Border color: `#343536`

**Typography**
- Font family: `"IBM Plex Sans", -apple-system, sans-serif`
- Post title: 15px, font-weight 500
- Body text: 14px
- Metadata: 12px, text-secondary
- Line height: 1.5

**Spacing**
- Base unit: 8px
- Card padding: 12px
- List item gap: 2px
- Section gap: 16px

### Components

**Header Bar (left pane top)**
- Subreddit input with r/ prefix
- Subreddit search button
- Quick links: r/popular, r/all

**Post List Item**
- Thumbnail (40x40, rounded) if available
- Title (truncate 2 lines)
- Metadata: subreddit, author, time ago, score
- Comment count indicator
- Hover: background shift to tertiary

**Right Pane - Post View**
- Post title (large)
- Post body (selftext) if text post
- Media/Link preview if link post
- Action bar: upvote, downvote, share, open in Reddit
- Comments section header with count

**Comment Component**
- Author + time ago
- Comment body (markdown rendered)
- Score
- Collapsible children (nested replies)
- Indentation for nesting levels (max 8 levels visual)

**Empty States**
- Left pane initial: "Select a subreddit to start"
- Right pane initial: "Click a post to view"

## Functionality Specification

### Core Features
1. **Subreddit browsing**: Enter subreddit name, load posts
2. **Post list**: Display 25 posts per page, infinite scroll
3. **Post detail**: Click post to load in right pane
4. **Comments**: Load and display nested comments
5. **Navigation**: Browser URL updates with subreddit (e.g., /r/funny)
6. **Deep linking**: Direct URL to specific post

### User Interactions
- Click subreddit → loads posts in left pane
- Click post → loads in right pane
- Click comment → collapses/expands
- Click "load more" → pagination
- Keyboard: Arrow keys to navigate posts

### Data Handling
- Reddit JSON API: `https://www.reddit.com/r/{subreddit}.json`
- Comments: `https://www.reddit.com/r/{subreddit}/{postId}.json`
- No authentication required (public API)
- Cache responses in memory

### Edge Cases
- Subreddit not found → show error message
- Post removed/deleted → show placeholder
- API rate limit → show retry button
- Empty subreddit → show "no posts"

## Acceptance Criteria
1. App loads with r/popular by default
2. Can search any public subreddit
3. Posts display with title, metadata, thumbnail
4. Clicking post shows content + comments in right pane
5. Comments are nested and collapsible
6. URL reflects current subreddit
7. Works on desktop and mobile
8. Dark theme matches spec colors
