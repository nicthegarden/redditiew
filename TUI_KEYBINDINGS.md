# TUI Keyboard Shortcuts Reference

Complete guide to all keyboard shortcuts in the RedditView Terminal User Interface (TUI).

## Quick Reference

| Mode | Action | Keys |
|------|--------|------|
| **List** | Move up | `↑` / `k` |
| **List** | Move down | `↓` / `j` |
| **List** | View post | `Enter` |
| **List** | Search | `Ctrl+F` |
| **List** | Change subreddit | `s` |
| **List** | Refresh | `F5` |
| **List** | Quit | `q` |
| **Details** | Scroll up | `↑` / `k` |
| **Details** | Scroll down | `↓` / `j` |
| **Details** | Page up | `Page Up` / `b` |
| **Details** | Page down | `Page Down` / `f` |
| **Details** | Go to top | `Home` / `g` |
| **Details** | Go to bottom | `End` / `G` |
| **Details** | Previous post | `h` |
| **Details** | Next post | `l` |
| **Details** | View comments | `c` |
| **Details** | Open in browser | `w` |
| **Details** | Back to list | `Esc` / `Tab` |
| **Comments** | Scroll up | `↑` |
| **Comments** | Scroll down | `↓` |
| **Comments** | Page up | `Page Up` |
| **Comments** | Page down | `Page Down` |
| **Comments** | Go to top | `Home` |
| **Comments** | Go to bottom | `End` |
| **Comments** | Previous post | `h` |
| **Comments** | Next post | `l` |
| **Comments** | Open in browser | `w` |
| **Comments** | Close comments | `Esc` |

---

## Navigation Keys

### Post List Navigation

**Move between posts:**
```
↑ / k      Move up one post
↓ / j      Move down one post
```

**View post details:**
```
Enter      Open selected post in detail view
```

### Detail View Navigation

**Scroll within details:**
```
↑ / k      Scroll up one line
↓ / j      Scroll down one line
```

**Scroll by page:**
```
Page Up / b     Scroll up 10 lines
Page Down / f   Scroll down 10 lines
```

**Jump to position:**
```
Home / g   Jump to top of details
End / G    Jump to bottom of details
```

**Switch posts:**
```
h          View previous post in list
l          View next post in list
```

---

## Comment Navigation

When viewing comments (press `c` in detail view):

**Scroll through comments:**
```
↑          Scroll up one line
↓          Scroll down one line
```

**Scroll by page:**
```
Page Up    Scroll up 10 lines
Page Down  Scroll down 10 lines
```

**Jump to position:**
```
Home       Jump to top of comments
End        Jump to bottom of comments
```

**Switch posts while viewing comments:**
```
h          View comments for previous post
l          View comments for next post
```

**Close comments:**
```
Esc        Close comments, return to details
```

---

## Application Controls

### Search

**Open search:**
```
Ctrl+F     Open post search (list view only)
```

**In search mode:**
```
Type       Enter search query
Enter      Filter posts by search
Esc        Close search and cancel
```

**Search tips:**
- Searches post titles and authors
- Case-insensitive
- Partial matches supported

### Subreddit Selection

**Change subreddit:**
```
s          Open subreddit selector
```

**In subreddit mode:**
```
Type       Enter subreddit name (without r/)
Enter      Load selected subreddit
Esc        Cancel and keep current subreddit
```

### Browser Integration

**Open URL:**
```
w          Open current post in default browser
```

This works in:
- List view (opens selected post)
- Detail view (opens current post)
- Comment view (opens current post)

**Supported URL types:**
- Reddit permalinks (e.g., /r/sysadmin/comments/xyz)
- External URLs (e.g., https://example.com)
- Automatically prepends reddit.com to permalinks

### Refresh

**Reload posts:**
```
F5         Refresh post list from server
```

This is useful when:
- Switching subreddits
- Getting fresh data
- Recovering from errors

### Quit

**Exit application:**
```
q          Quit RedditView
```

Pressing `q`:
- Closes the TUI
- Returns to terminal prompt
- Does NOT save any data (stateless)

---

## View Modes

### List View
Shows all posts from current subreddit.

**Active keys:**
- Arrow keys / j/k - Move between posts
- Enter - View post details
- Ctrl+F - Search
- s - Change subreddit
- F5 - Refresh
- q - Quit

**Footer shows:**
```
Post 5/200 • Enter: view details • w: open URL • Ctrl+F: search • F5: refresh • q: quit
```

### Detail View
Shows full post text with details.

**Active keys:**
- Arrow keys / j/k - Scroll details
- h/l - Switch posts
- c - View comments
- w - Open in browser
- Page Up/Down - Scroll by page
- Home/End - Jump to position
- Esc/Tab - Back to list

**Footer shows:**
```
↑↓: scroll details • h/l: switch posts • w: open URL • Esc/Tab: back to list • c: view comments • q: quit
```

### Comments View
Shows comments for current post.

**Active keys:**
- Arrow keys - Scroll comments
- h/l - Switch posts
- w - Open post URL
- Page Up/Down - Scroll by page
- Home/End - Jump to position
- Esc - Close comments

**Footer shows:**
```
↑↓: scroll comments • h/l: switch posts • w: open URL • Esc: close comments • Ctrl+F: search • q: quit
```

---

## Advanced Tips & Tricks

### Quick Navigation
- Use `j`/`k` for fast single-line scrolling
- Use `Page Up`/`Page Down` for faster scrolling
- Use `Home`/`End` to jump to extremes

### Efficient Browsing
```
1. j/k to browse posts in list
2. Enter to view interesting post
3. c to read comments
4. h/l to quickly switch between nearby posts
5. Esc to return to detail
6. j/k to continue browsing
```

### Switching Subreddits
```
s                  # Open subreddit selector
programming        # Type subreddit (without r/)
Enter              # Load new subreddit
```

### Opening Multiple Posts
```
# In detail view:
Enter              # View post 1
...read...
h                  # Jump to previous post
# or
l                  # Jump to next post
# and continue reading without returning to list
```

### Power User Flow
```
Ctrl+F             # Search for topic
# Results appear in list
Enter              # View interesting result
c                  # Open comments
↑/↓                # Read through comments
w                  # Want to discuss? Open in browser
l                  # Check next result
```

---

## Keyboard Layouts

### QWERTY (English)
All examples above assume QWERTY layout.

### DVORAK / AZERTY / Other Layouts
The special keys work as expected:
- Arrow keys - Always work
- Page Up/Down - Always work
- Home/End - Always work
- Ctrl+F - Always work
- Function keys (F5) - Always work

The letter keys (j/k, h/l, etc.) depend on your keyboard layout.

---

## Accessibility Features

### For Users with Limited Mobility
- All functions accessible via arrow keys + modifiers
- No complex key combinations required
- No repetitive strain (single keypresses work)

### Terminal Size Requirements
- **Minimum:** 80 columns × 24 rows
- **Recommended:** 120 columns × 40 rows
- Smaller terminals may have text truncation

### Terminal Compatibility
Tested and working on:
- Linux (xterm, GNOME Terminal, Konsole, etc.)
- Windows (PowerShell, Windows Terminal, ConEmu)
- macOS (Terminal.app, iTerm2)

---

## Troubleshooting Keybindings

### "Keys aren't working"

**Solution 1: Check terminal focus**
- Click on the terminal window to ensure it has focus
- Some terminals require explicit focus

**Solution 2: Check for key conflicts**
- Some terminals intercept Ctrl+F, Ctrl+L, etc.
- Try using alternatives (e.g., use `h/l` instead of arrows if arrows don't work)

**Solution 3: Check keyboard layout**
- If using non-QWERTY layout, letter keys may differ
- Arrow keys and modifier keys always work regardless of layout

### "Ctrl+F not working"

**Solutions:**
1. Some terminals capture Ctrl+F for their own find feature
2. Try pressing it twice (once for terminal, once for app)
3. Workaround: Use `s` to change subreddit instead of searching

### "Page Up/Down not working"

**Solutions:**
1. Check that terminal has focus
2. Try `b` (page up) and `f` (page down) instead
3. Use Home/End keys instead

### "Terminal looks broken"

**Solutions:**
1. Resize terminal to at least 80x24
2. Try pressing Ctrl+L to redraw
3. Quit (q) and restart the application

---

## Default Bindings Summary

### Motion
| Action | Primary | Alternative |
|--------|---------|-------------|
| Up | `↑` | `k` |
| Down | `↓` | `j` |
| Page Up | `Page Up` | `b` |
| Page Down | `Page Down` | `f` |
| Top | `Home` | `g` |
| Bottom | `End` | `G` |
| Previous | `h` | - |
| Next | `l` | - |

### Features
| Action | Key |
|--------|-----|
| Open post | `Enter` |
| Comments | `c` |
| Browser | `w` |
| Search | `Ctrl+F` |
| Subreddit | `s` |
| Refresh | `F5` |
| Back | `Esc` / `Tab` |
| Quit | `q` |

---

## See Also

- 📖 [README.md](README.md) - Project overview
- 🚀 [QUICKSTART.md](QUICKSTART.md) - Getting started
- ⚙️ [CONFIGURATION.md](CONFIGURATION.md) - Configuration options
- 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md) - Technical details

---

Happy browsing! 🚀

Need help? Check [QUICKSTART.md](QUICKSTART.md) for troubleshooting.
