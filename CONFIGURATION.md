# Configuration Guide

Complete reference for configuring RedditView to your preferences.

## Table of Contents
- [Configuration File](#configuration-file)
- [TUI Settings](#tui-settings)
- [Web Settings](#web-settings)
- [API Settings](#api-settings)
- [Advanced Configuration](#advanced-configuration)
- [Configuration Examples](#configuration-examples)
- [Troubleshooting](#troubleshooting)

---

## Configuration File

Configuration is managed through `config.json` in the project root directory.

### Default Configuration
```json
{
  "tui": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 200,
    "list_height": 10,
    "max_title_length": 80,
    "default_sort": "popular",
    "subreddit_shortcuts": {
      "1": "sysadmin",
      "2": "golang",
      "3": "programming",
      "4": "linux",
      "5": "devops",
      "6": "webdev",
      "7": "learnprogramming",
      "8": "100DaysOfCode",
      "9": "codereview"
    }
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

### Creating a Custom Configuration

1. Copy the default configuration (above)
2. Edit values as needed (see sections below)
3. Save as `config.json` in project root
4. Restart the application

⚠️ **Note:** Configuration is loaded when the application starts. Changes to `config.json` require a restart.

---

## TUI Settings

### default_subreddit
**Type:** `string`  
**Default:** `"sysadmin"`  
**Description:** Subreddit to load when the TUI starts

**Valid Values:**
- Any subreddit name (without "r/" prefix)
- Examples: `"programming"`, `"AskReddit"`, `"pics"`, `"news"`

**Example:**
```json
"default_subreddit": "programming"
```

**Notes:**
- Subreddit names are case-insensitive
- You can change subreddits at runtime with `s` key
- Invalid subreddit names will show an error

---

### posts_per_page
**Type:** `integer`  
**Default:** `200`  
**Valid Range:** `10` - `500`  
**Description:** Number of posts to load per request

| Value | Effect | Use Case |
|-------|--------|----------|
| 50 | Few posts, fast loading | Slow connection, low memory |
| 100 | Balanced | Most users |
| 200 | Many posts, more to browse | Fast connection, powerful computer |
| 300+ | Lots of posts, slower loading | Powerful computer, patient user |

**Example:**
```json
"posts_per_page": 150
```

**Performance Impact:**
- **Higher values:** More posts visible, slower initial load, more memory usage
- **Lower values:** Fewer posts, faster load time, better for slow devices

**Recommendation:**
- Modern computers: `200` (current default)
- Older computers or slow connection: `100`
- Very old devices: `50`

---

### list_height
**Type:** `integer`  
**Default:** `10`  
**Valid Range:** `5` - `30`  
**Description:** Height of the post list panel (in lines)

**Example:**
```json
"list_height": 15
```

**Effect on Layout:**
- Increasing list height → Details section gets smaller
- Decreasing list height → More space for details

**Typical Values:**
| Terminal Size | Recommended |
|---|---|
| Small (24x80) | 8-10 |
| Medium (30x100) | 12-15 |
| Large (40x140) | 18-20 |

---

### max_title_length
**Type:** `integer`  
**Default:** `80`  
**Valid Range:** `40` - `200`  
**Description:** Maximum characters to display for post titles

**Example:**
```json
"max_title_length": 100
```

**Effect:**
- Titles longer than this value will be truncated with "..."
- Helps prevent line wrapping and keeps layout clean
- Set higher if you have wide terminal, lower for narrow terminals

---

### default_sort
**Type:** `string`  
**Default:** `"popular"`  
**Valid Values:** `"popular"`, `"new"`, `"top"`, `"controversial"`, `"rising"`  
**Description:** Default sorting for posts when loading a subreddit

**Example:**
```json
"default_sort": "new"
```

**Options:**
- `"popular"` (Hot) - Most upvoted posts
- `"new"` - Newest posts
- `"top"` - Top posts (all time)
- `"controversial"` - Most controversial posts
- `"rising"` - Rising posts

**Notes:**
- Can be toggled at runtime with `t` key
- Current sort preference persists when switching subreddits
- Maps to Reddit's `/hot`, `/new`, `/top`, etc. endpoints

---

### subreddit_shortcuts
**Type:** `object` (map of keys to subreddit names)  
**Default:** See default config above  
**Description:** Quick keyboard shortcuts to jump to frequently used subreddits

**Example:**
```json
"subreddit_shortcuts": {
  "1": "programming",
  "2": "golang",
  "3": "rust",
  "4": "python",
  "5": "learnprogramming",
  "6": "webdev",
  "7": "devops",
  "8": "sysadmin",
  "9": "linux"
}
```

**Usage:**
- Press keys **1-9** to instantly jump to configured subreddit
- Current sort preference is maintained when switching
- Search/filter is reset when switching subreddits

**Customization Tips:**
- Configure shortcuts to your most-used subreddits
- Remove unused shortcuts by setting to empty object `{}`
- Keys must be strings "1" through "9"
- Values must be valid subreddit names (without "r/" prefix)

---

### default_subreddit (Web)
**Type:** `string`  
**Default:** `"sysadmin"`  
**Description:** Subreddit to load when the web UI starts

Same format as TUI setting (see [TUI Settings](#tui-settings)).

---

### posts_per_page (Web)
**Type:** `integer`  
**Default:** `20`  
**Valid Range:** `5` - `100`  
**Description:** Posts to display per page in web UI

**Typical Values:**
| Value | Effect |
|-------|--------|
| 10 | Few posts per page, more pagination |
| 20 | Balanced (current default) |
| 50 | Many posts per page, less scrolling |

**Note:** Web UI typically loads fewer posts than TUI for better performance.

---

### theme
**Type:** `string`  
**Default:** `"dark"`  
**Valid Values:** `"dark"`, `"light"`  
**Description:** Visual theme for the web UI

**Example:**
```json
"theme": "light"
```

---

## API Settings

### base_url
**Type:** `string`  
**Default:** `"http://localhost:3002/api"`  
**Description:** URL of the API server

**Default Format:**
```
http://localhost:PORT/api
```

**Important Notes:**
- This should point to your API server
- Default assumes API server runs on same machine on port 3002
- Do NOT include trailing slash

**Use Cases:**
- **Local development:** `"http://localhost:3002/api"` (default)
- **Remote server:** `"http://api.example.com:3002/api"`
- **Docker:** `"http://api-container:3002/api"`

**Example (Remote Server):**
```json
"base_url": "http://reddit-api.myserver.com:3002/api"
```

---

### timeout_seconds
**Type:** `integer`  
**Default:** `10`  
**Valid Range:** `5` - `60`  
**Description:** API request timeout in seconds

**Example:**
```json
"timeout_seconds": 15
```

**When to Adjust:**
| Situation | Value |
|-----------|-------|
| Fast, local network | 5-7 seconds |
| Normal internet | 10 seconds (default) |
| Slow connection | 15-20 seconds |
| Very slow connection | 30+ seconds |

**Note:** Higher values mean longer waits if server is unresponsive.

---

## Advanced Configuration

### Environment Variables

You can override configuration values with environment variables:

```bash
# Linux/macOS
export REDDIT_SUBREDDIT=programming
export REDDIT_POSTS_PER_PAGE=100
export API_BASE_URL=http://api.example.com:3002/api

# Windows PowerShell
$env:REDDIT_SUBREDDIT = "programming"
$env:REDDIT_POSTS_PER_PAGE = "100"
$env:API_BASE_URL = "http://api.example.com:3002/api"

# Then start the application
./apps/tui/redditview
```

**Available Environment Variables:**
- `REDDIT_SUBREDDIT` - Overrides `tui.default_subreddit`
- `REDDIT_POSTS_PER_PAGE` - Overrides `tui.posts_per_page`
- `API_BASE_URL` - Overrides `api.base_url`
- `API_TIMEOUT` - Overrides `api.timeout_seconds`

---

## Configuration Examples

### Example 1: Low-Bandwidth Setup
For users with slow internet or old computers:

```json
{
  "tui": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 50,
    "list_height": 8,
    "max_title_length": 60
  },
  "api": {
    "base_url": "http://localhost:3002/api",
    "timeout_seconds": 20
  }
}
```

**Changes:**
- Reduced posts per page (50 → faster loading)
- Smaller list height (8 → more room for details)
- Increased timeout (20 → tolerates slow network)

---

### Example 2: Power User Setup
For users with high-speed internet and large displays:

```json
{
  "tui": {
    "default_subreddit": "programming",
    "posts_per_page": 300,
    "list_height": 18,
    "max_title_length": 120
  },
  "api": {
    "base_url": "http://localhost:3002/api",
    "timeout_seconds": 8
  }
}
```

**Changes:**
- More posts per page (300 → more content to browse)
- Larger list height (18 → bigger post list)
- Longer titles (120 → see full titles)
- Lower timeout (8 → expects fast network)

---

### Example 3: Multiple Subreddits
If you want to quickly switch between subreddits, use the `s` key in TUI rather than changing config. However, you can set your most-used subreddit as default:

```json
{
  "tui": {
    "default_subreddit": "AskReddit",
    "posts_per_page": 150,
    "list_height": 12,
    "max_title_length": 80
  }
}
```

Then use keyboard shortcut `s` to switch to other subreddits at runtime.

---

### Example 4: Remote API Server
If running API server on a different machine:

```json
{
  "tui": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 200,
    "list_height": 10,
    "max_title_length": 80
  },
  "api": {
    "base_url": "http://192.168.1.100:3002/api",
    "timeout_seconds": 15
  }
}
```

**Changes:**
- Updated API URL to remote server IP
- Increased timeout (15 → network latency)

---

## Troubleshooting

### "Configuration file not found"

**Solution:** Create `config.json` in project root with default content:

```bash
cd redditiew-local
cat > config.json << 'EOF'
{
  "tui": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 200,
    "list_height": 10,
    "max_title_length": 80,
    "default_sort": "popular",
    "subreddit_shortcuts": {
      "1": "sysadmin",
      "2": "golang",
      "3": "programming",
      "4": "linux",
      "5": "devops",
      "6": "webdev",
      "7": "learnprogramming",
      "8": "100DaysOfCode",
      "9": "codereview"
    }
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
EOF
```

---

### "Invalid JSON in config.json"

**Symptoms:** Application won't start, shows JSON parsing error

**Solution:**
1. Open `config.json` in a text editor
2. Check for:
   - Missing commas between properties
   - Trailing commas (not allowed in JSON)
   - Unmatched quotes
   - Missing braces

**Validation:** Copy config to https://jsonlint.com to check syntax

---

### "Unknown configuration field"

RedditView ignores unknown fields. This is safe for forward/backward compatibility.

---

### "Configuration changes not applied"

Changes to `config.json` require a restart:
1. Stop the running application (press `q` or `Ctrl+C`)
2. Edit `config.json`
3. Start the application again

---

### "Very slow loading"

If posts load slowly:
1. Reduce `posts_per_page` (start with 100)
2. Check internet connection
3. Increase `timeout_seconds` to 20-30
4. Check if API server is running: `curl http://localhost:3002/api/r/sysadmin.json`

---

### "Subreddit not found"

Verify the subreddit name:
- Check spelling (case-insensitive)
- Don't include "r/" prefix
- Try a popular subreddit: `"programming"`, `"AskReddit"`, `"pics"`
- Some subreddits may be private or deleted

---

## Configuration Best Practices

1. **Start with defaults** - The default configuration is optimized for most users
2. **Make small changes** - Change one setting at a time to see the effect
3. **Test after changes** - Always restart and verify changes work as expected
4. **Keep backups** - Save a copy of your working config before major changes
5. **Match hardware** - Adjust `posts_per_page` based on your computer's power

---

## Default Configuration Quick Reference

| Setting | Default | Notes |
|---------|---------|-------|
| default_subreddit | sysadmin | - |
| posts_per_page | 200 | Range: 10-500 |
| list_height | 10 | Range: 5-30 |
| max_title_length | 80 | Range: 40-200 |
| default_sort | popular | Options: popular, new, top, controversial, rising |
| subreddit_shortcuts | (see default config) | Keys 1-9 for quick access |
| timeout_seconds | 10 | Range: 5-60 |

---

## See Also

- 📖 [README.md](README.md) - Project overview
- 🚀 [QUICKSTART.md](QUICKSTART.md) - Get started quickly
- ⌨️ [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md) - Keyboard shortcuts
- 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md) - Technical details

---

Happy configuring! 🎉
