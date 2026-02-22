# RedditView TUI v4 - Quick Reference

## 🚀 Launch

```bash
./launch.sh tui
```

## ⌨️ Keyboard Shortcuts (Quick Lookup)

### Navigation
```
j or ↓    Navigate down (next post)
k or ↑    Navigate up (previous post)
Home      Jump to first post
End       Jump to last post
```

### Feature Shortcuts (Ctrl+)
```
Ctrl+F    Search posts (Find)
Ctrl+R    Edit subreddit (Reddit name)
F5        Refresh current subreddit
```

### Exit
```
q         Quit
Ctrl+C    Quit
```

## 📖 How to Use

### Browse Posts
1. Posts load automatically when you start
2. Use `j`/`k` or arrow keys to navigate
3. Selected post shows full details below
4. Scroll down to see more posts

### Search for Topic
1. Press `Ctrl+F`
2. Type your search (title or author)
3. Press `Enter` to apply search
4. Posts filter in real-time
5. Press `Esc` to cancel search

### Change Subreddit
1. Press `Ctrl+R`
2. Type subreddit name (e.g., "rust", "python")
3. Press `Enter` to load
4. Posts load from new subreddit
5. Press `Esc` to cancel

### Refresh Posts
1. Press `F5`
2. Current subreddit reloads
3. Returns to first post

## 🎨 Display

```
Header       : 🔥 r/golang  Posts: 50
Info Bar     : Keyboard shortcuts or search field
Post List    : Collapsed posts (▶ = selected)
Post Details : Full title, author, content, URL
Comments     : Placeholder (coming soon)
Footer       : Position + available shortcuts
```

## 🔧 Settings

Current subreddit: shown in header
Default subreddit on startup: `golang`

To change default, edit main.go:
```go
subreddit: "golang",  // Change to your preferred subreddit
```

## 🆘 Troubleshooting

| Problem | Solution |
|---------|----------|
| Posts not loading | Check API server: `node api-server.js` |
| Can't search | Make sure you press `Ctrl+F` first |
| Shortcuts not working | Press `Esc` first to exit search mode |
| Text running off screen | Make terminal wider or use different font size |
| Unicode symbols broken | Check terminal UTF-8 support |

## 📋 Feature Checklist

- ✅ Load posts from any subreddit
- ✅ Display post list
- ✅ Expand post to see full details
- ✅ Search by title and author
- ✅ Change subreddit on the fly
- ✅ Keyboard navigation
- ✅ Refresh posts
- ⏳ Comment tree display (v4.1)
- ⏳ Post voting (v4.2)
- ⏳ Post sorting (v4.2)

## 💡 Tips

1. **Use Ctrl+ shortcuts** instead of traditional shortcuts
2. **Home/End keys** jump to first/last post quickly
3. **j/k faster than arrows** for fast navigation
4. **F5 to refresh** if posts seem stale
5. **Search persists** - F5 refreshes search results

## 🎯 Common Tasks

**View top posts from r/golang**
```
./launch.sh tui
# Posts load automatically
# Navigate with j/k
```

**Search for "concurrency" posts**
```
./launch.sh tui
Ctrl+F
type: concurrency
Enter
# See filtered results
```

**Switch to r/rust**
```
Ctrl+R
type: rust
Enter
# New posts load
```

**See latest posts**
```
F5
# Refreshes current subreddit
```

**Go to first post**
```
Home
# Jumps to post 1
```

**Exit program**
```
q
# or Ctrl+C
```

## 📞 Help

Press `Ctrl+F` to see available features in header
Check info bar (line 2) for current mode shortcuts

---

**Version**: 4.0.0  
**Updated**: February 22, 2026
