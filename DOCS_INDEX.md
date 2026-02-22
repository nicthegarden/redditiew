# Documentation Index

Welcome to RedditView documentation! This comprehensive guide covers everything you need to know about installing, configuring, and using the RedditView Reddit browser.

## 🚀 Getting Started (5 minutes)

**New to RedditView?** Start here:

### [QUICKSTART.md](QUICKSTART.md) - Get Up and Running Fast
- ⏱️ 5-minute setup for Windows
- ⏱️ 5-minute setup for Linux
- Basic usage instructions
- Troubleshooting quick fixes

👉 **Choose your OS and follow the step-by-step guide**

---

## 📖 Documentation Roadmap

### 1. Installation & Setup
- **[QUICKSTART.md](QUICKSTART.md)** - Fast setup (5 min)
  - Windows installation
  - Linux installation
  - Basic first run
  
- **[INSTALLATION.md](INSTALLATION.md)** - Detailed technical setup
  - Prerequisite verification
  - OS-specific instructions
  - Docker setup
  - Troubleshooting common issues
  - Build from source

### 2. Configuration & Customization
- **[CONFIGURATION.md](CONFIGURATION.md)** - Customize your setup
  - Configuration file reference
  - TUI settings (posts per page, subreddit, etc.)
  - Web settings
  - API settings
  - Configuration examples
  - Environment variables

### 3. Using the Application
- **[README.md](README.md)** - Project overview
  - Feature highlights
  - Screenshots with PNG examples
  - Basic usage guide
  - System requirements
  - Troubleshooting
  
- **[TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md)** - Complete keyboard shortcuts
  - Navigation keys
  - Application controls
  - View modes
  - Tips & tricks
  - Accessibility features

### 4. Technical Deep Dive
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Technical architecture
  - System design
  - Component overview
  - Data flow
  - API endpoints
  - Development guide

### 5. Reference & Support
- **[DEVELOPMENT.md](DEVELOPMENT.md)** - Contributing & development
  - Development setup
  - Code structure
  - Testing
  - Contributing guidelines

---

## 📸 Visual Guide

### Screenshots Included

The documentation includes actual screenshots of the application:

**1. TUI Post List View** ([TUI.png](TUI.png))
```
Shows the main post browsing interface with:
- Post list on the left
- Post details on the right
- Keybinding footer
- Post scores and metadata
```

**2. TUI Comments View** ([TUI-Comment.png](TUI-Comment.png))
```
Shows the comments panel with:
- Comment list with author, score, and content
- Scrollable comment body
- Navigation indicators
- Context-aware keybindings
```

**3. Web UI** ([WebUI.png](WebUI.png))
```
Shows the web interface with:
- Modern responsive design
- Post browser
- Comment viewer
- Theme support
```

---

## 🎯 Quick Navigation by Task

### "I want to..."

**...install and run it quickly**
→ [QUICKSTART.md](QUICKSTART.md)

**...understand all keyboard shortcuts**
→ [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md)

**...customize the application**
→ [CONFIGURATION.md](CONFIGURATION.md)

**...set it up for development**
→ [INSTALLATION.md](INSTALLATION.md) → [DEVELOPMENT.md](DEVELOPMENT.md)

**...understand the architecture**
→ [ARCHITECTURE.md](ARCHITECTURE.md)

**...fix a problem**
→ [QUICKSTART.md#troubleshooting](QUICKSTART.md#troubleshooting) → [INSTALLATION.md#troubleshooting](INSTALLATION.md#troubleshooting)

**...use advanced features**
→ [TUI_KEYBINDINGS.md#advanced-tips--tricks](TUI_KEYBINDINGS.md#advanced-tips--tricks)

**...report a bug or contribute**
→ [DEVELOPMENT.md](DEVELOPMENT.md)

---

## 📊 Documentation Overview

| Document | Purpose | Length | Technical Depth |
|----------|---------|--------|-----------------|
| [QUICKSTART.md](QUICKSTART.md) | Fast setup guide | 10 min read | Beginner |
| [README.md](README.md) | Project overview | 5 min read | Beginner |
| [CONFIGURATION.md](CONFIGURATION.md) | Config reference | 15 min read | Intermediate |
| [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md) | Keyboard shortcuts | 10 min read | Beginner |
| [INSTALLATION.md](INSTALLATION.md) | Technical setup | 20 min read | Advanced |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Technical design | 15 min read | Advanced |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Contributing guide | 10 min read | Advanced |

---

## 🎓 Learning Path

### For First-Time Users
1. Read [README.md](README.md) - Get overview
2. Follow [QUICKSTART.md](QUICKSTART.md) - Install in 5 minutes
3. Try basic features
4. Reference [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md) - Learn shortcuts
5. Explore [CONFIGURATION.md](CONFIGURATION.md) - Customize setup

### For Advanced Users
1. Skim [README.md](README.md) - Get overview
2. Read [INSTALLATION.md](INSTALLATION.md) - Custom build
3. Study [CONFIGURATION.md](CONFIGURATION.md) - Advanced settings
4. Read [ARCHITECTURE.md](ARCHITECTURE.md) - Understanding internals
5. Read [DEVELOPMENT.md](DEVELOPMENT.md) - Contribute code

### For Operators/DevOps
1. [INSTALLATION.md](INSTALLATION.md) - Deployment options
2. [CONFIGURATION.md](CONFIGURATION.md) - Environment variables
3. [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
4. Production deployment section in [INSTALLATION.md](INSTALLATION.md)

---

## 📋 Quick Reference Checklists

### Installation Checklist
- [ ] Install Go 1.19+ 
- [ ] Install Node.js 16+
- [ ] Clone repository
- [ ] Run `npm install`
- [ ] Build TUI with `go build`
- [ ] Start API server with `npm start`
- [ ] Run TUI
- [ ] Verify posts load

### First-Time Usage Checklist
- [ ] Browse posts with j/k or arrow keys
- [ ] View post with Enter
- [ ] Read comments with c
- [ ] Scroll with arrow keys and Page Up/Down
- [ ] Search with Ctrl+F
- [ ] Change subreddit with s
- [ ] Open in browser with w
- [ ] Quit with q

### Configuration Checklist
- [ ] Edit config.json if desired
- [ ] Set preferred default subreddit
- [ ] Adjust posts_per_page for your hardware
- [ ] Verify API server address if remote
- [ ] Test with `curl` before running TUI

---

## 🔧 Configuration Files

### config.json - Application Configuration
```json
{
  "tui": {
    "default_subreddit": "sysadmin",
    "posts_per_page": 200,
    "list_height": 10,
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

**See [CONFIGURATION.md](CONFIGURATION.md) for complete reference**

---

## 🐛 Troubleshooting by Symptom

| Symptom | Likely Cause | Solution |
|---------|------------|----------|
| "command not found: go" | Go not installed | [INSTALLATION.md](INSTALLATION.md) Step 1 |
| "Cannot connect to API" | API server not running | Run `npm start` |
| TUI won't start | API not responding | Check `curl http://localhost:3002/api/r/sysadmin.json` |
| Keys not working | Terminal size too small | Resize to minimum 80×24 |
| Posts won't load | Network/API issue | Check internet connection, API logs |
| Comments not showing | API issue | Refresh with F5, check API logs |
| Layout broken/garbled | Terminal too small | Increase to minimum 80×24 |

**For more detailed troubleshooting:**
- [QUICKSTART.md#troubleshooting](QUICKSTART.md#troubleshooting)
- [INSTALLATION.md#troubleshooting](INSTALLATION.md#troubleshooting)

---

## 📞 Support & Contact

### Getting Help

1. **Check Documentation** - Most answers are here
2. **Check FAQ** - See QUICKSTART.md troubleshooting
3. **Search Issues** - Check GitHub Issues
4. **Create Issue** - Report bugs on GitHub
5. **Email Support** - support@example.com (if applicable)

### Reporting Bugs

Include in bug report:
- OS and version (Windows 10, Ubuntu 20.04, macOS 12)
- Go version (`go version`)
- Node.js version (`node --version`)
- What you were doing when bug occurred
- Full error message
- Steps to reproduce

### Feature Requests

Create GitHub issue with:
- Clear title
- Detailed description
- Why it would be useful
- Any relevant screenshots

---

## 🎯 Key Features Overview

### TUI Features
✨ Keyboard-driven navigation  
✨ Split-view layout  
✨ Smooth scrolling  
✨ Search functionality  
✨ Comment viewing  
✨ Browser integration  
✨ Responsive to terminal size  

### Web UI Features
✨ Modern responsive design  
✨ Mouse support  
✨ Theme switching  
✨ Real-time updates  

### Cross-Platform
✅ Linux  
✅ Windows  
✅ macOS  

---

## 📊 By The Numbers

- **5** minutes to get started
- **50+** keyboard shortcuts
- **200** posts per page (default)
- **5** comments per post
- **3** main view modes
- **11** MB TUI binary
- **4** main configuration files

---

## 🔄 Recent Updates

### Latest Features (v0.2.0)
- ✨ Enhanced comment scrolling with proper height calculation
- ✨ Open posts directly in browser with 'w' key
- ✨ Increased page scroll distance for faster navigation
- ✨ 200 posts per page by default (up from 50)

### Latest Fixes
- 🐛 Fixed comment scrolling state propagation
- 🐛 Fixed list display bug affecting post visibility
- 🐛 Fixed hardcoded height in comment scroll calculation

**See git log for complete history**

---

## 📚 Document Index

### User Guides
- [README.md](README.md) - Project overview and introduction
- [QUICKSTART.md](QUICKSTART.md) - Fast 5-minute setup
- [TUI_KEYBINDINGS.md](TUI_KEYBINDINGS.md) - Keyboard shortcut reference

### Configuration
- [CONFIGURATION.md](CONFIGURATION.md) - All configuration options
- [config.json](config.json) - Configuration file

### Technical
- [INSTALLATION.md](INSTALLATION.md) - Detailed installation guide
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development guide

### This Document
- [DOCS_INDEX.md](DOCS_INDEX.md) - You are here!

---

## 🎉 Getting Started Now

### Choose Your Path

**Fast Track (5 min):**
```
→ Read: QUICKSTART.md (your OS section)
→ Run: npm start && ./apps/tui/redditview
→ Done!
```

**Thorough Setup (20 min):**
```
→ Read: README.md
→ Read: INSTALLATION.md (full)
→ Run: Custom installation
→ Configure: CONFIGURATION.md
→ Learn: TUI_KEYBINDINGS.md
```

**Development Setup (30 min):**
```
→ Read: INSTALLATION.md
→ Read: ARCHITECTURE.md
→ Read: DEVELOPMENT.md
→ Setup: Development environment
```

---

## 📖 Complete Documentation Map

```
RedditView Documentation/
├── README.md                    ← Start here: Project overview
├── QUICKSTART.md               ← Fast setup (5 min)
├── INSTALLATION.md             ← Detailed installation
├── CONFIGURATION.md            ← All configuration options
├── TUI_KEYBINDINGS.md         ← Keyboard shortcut reference
├── ARCHITECTURE.md             ← Technical architecture
├── DEVELOPMENT.md              ← Contributing guide
├── DOCS_INDEX.md              ← You are here!
├── config.json                ← Configuration file
└── Screenshots/
    ├── TUI.png                ← TUI interface screenshot
    ├── TUI-Comment.png        ← Comments view screenshot
    └── WebUI.png              ← Web interface screenshot
```

---

## ✨ Pro Tips

1. **Bookmark QUICKSTART.md** - Fast reference for setup
2. **Print TUI_KEYBINDINGS.md** - Keep by your desk
3. **Keep config.json handy** - Easy to customize
4. **Check ARCHITECTURE.md** - Understand how it works
5. **Set terminal size minimum 120x40** - Better experience

---

**Ready to get started? [Go to QUICKSTART.md →](QUICKSTART.md)**

Happy browsing! 🚀

---

*Last updated: February 22, 2026*  
*RedditView v0.2.0*
