# Contributing to RedditView

Thank you for your interest in contributing to RedditView! This document provides guidelines and instructions for contributing to the project.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)
- [Feature Requests](#feature-requests)

---

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. We are committed to providing a harassment-free experience for everyone.

### Expected Behavior
- Be respectful and inclusive
- Welcome diverse perspectives
- Provide constructive feedback
- Focus on code, not the person
- Help others learn and grow

### Unacceptable Behavior
- Harassment or discrimination
- Offensive language or behavior
- Trolling or spam
- Personal attacks

---

## Getting Started

### 1. Fork the Repository
```bash
# Go to GitHub and fork: https://github.com/yourusername/redditiew-local
# Clone your fork
git clone https://github.com/YOUR-USERNAME/redditiew-local.git
cd redditiew-local
```

### 2. Add Upstream Remote
```bash
# Add the original repository as upstream
git remote add upstream https://github.com/yourusername/redditiew-local.git
git fetch upstream
```

### 3. Create a Feature Branch
```bash
# Always create a new branch for your work
git checkout -b feature/your-feature-name
```

---

## Development Setup

### Prerequisites
- Go 1.19+ (1.21+ recommended)
- Node.js 16+ (18 LTS recommended)
- npm 7+

### Setup Steps

```bash
# Install dependencies
npm install

# Build TUI
cd apps/tui
go build -o redditview .
cd ../..

# Start API server
npm start

# In another terminal, run TUI
./apps/tui/redditview
```

For detailed setup instructions, see [INSTALLATION.md](../INSTALLATION.md).

---

## Making Changes

### Before You Start
1. Check existing [Issues](https://github.com/yourusername/redditiew-local/issues)
2. Check existing [Pull Requests](https://github.com/yourusername/redditiew-local/pulls)
3. Open an issue for major changes (get approval first)
4. For small fixes/improvements, you can go straight to PR

### File Organization

```
redditiew-local/
├── apps/tui/main.go          # TUI source code
├── api-server.js             # API server
├── config.json               # Configuration
├── package.json              # Dependencies
└── README.md                 # Documentation
```

### Making Code Changes

1. **Edit files** in `apps/tui/main.go` or `api-server.js`
2. **Test locally** before committing
3. **Follow code style** (see below)
4. **Update documentation** if needed
5. **Test again** to ensure no regressions

### Key Source Files

**TUI Application** (`apps/tui/main.go`)
- Model struct (state management)
- Update() function (message handling)
- View() function (rendering)
- Key handlers (keyboard input)

**API Server** (`api-server.js`)
- Route handlers
- Reddit API integration
- Error handling

**Configuration** (`config.json`)
- Settings for TUI and Web UI
- API configuration

---

## Testing

### Run the Application
```bash
# Start API server
npm start

# In another terminal, run TUI
./apps/tui/redditview
```

### Manual Testing Checklist
- [ ] Posts load from API
- [ ] Navigation works (j/k, arrows)
- [ ] Details view opens (Enter key)
- [ ] Comments load (c key)
- [ ] Scrolling works (arrows, Page Up/Down)
- [ ] Search works (Ctrl+F)
- [ ] Subreddit switching works (s key)
- [ ] URL opening works (w key)
- [ ] Terminal resize handled properly
- [ ] Error messages display correctly

### Build Testing
```bash
# Rebuild TUI
cd apps/tui
go build -o redditview .

# Check for errors
go vet ./...

# Run tests (if any)
go test ./...
```

---

## Submitting Changes

### Step 1: Commit Your Changes

```bash
# Stage changes
git add .

# Create meaningful commit
git commit -m "feat: Add new feature

Description of what changed and why.
Reference issue #123 if applicable."
```

See [Commit Messages](#commit-messages) below for format.

### Step 2: Push to Your Fork
```bash
git push origin feature/your-feature-name
```

### Step 3: Create Pull Request
1. Go to GitHub
2. Click "Compare & pull request"
3. Fill out PR template (see below)
4. Submit PR

---

## Code Style

### Go Code (TUI)

**Format:** Use `gofmt`
```bash
gofmt -w apps/tui/main.go
```

**Style Guidelines:**
- Use camelCase for variables and functions
- Use PascalCase for types and exported functions
- Keep functions under 50 lines when possible
- Add comments for exported functions
- Handle errors explicitly (don't ignore)

**Example:**
```go
// fetchPosts retrieves posts from the API
func (m *Model) fetchPosts(subreddit string) error {
    // Implementation
}
```

### JavaScript Code (API)

**Format:** Use 2-space indentation
```javascript
// Bad
const fetchPosts = (subreddit) => {
    return api.get(`/r/${subreddit}`)
}

// Good
const fetchPosts = (subreddit) => {
  return api.get(`/r/${subreddit}`)
}
```

### Documentation

**Markdown Format:**
- Use clear, concise language
- Code blocks with language specified
- Proper heading hierarchy (# ## ###)
- Links to related documentation
- Examples where helpful

---

## Commit Messages

### Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- **feat:** New feature
- **fix:** Bug fix
- **docs:** Documentation changes
- **style:** Code style (formatting, semicolons, etc.)
- **refactor:** Code refactoring
- **test:** Adding or updating tests
- **chore:** Build, dependencies, etc.

### Examples

**Good:**
```
feat: Add comment scrolling with arrow keys

- Implement scroll state management
- Add scroll boundary calculations
- Handle window resize properly
- Tested with multiple window sizes

Fixes #123
```

**Bad:**
```
fixed stuff
changed code
updated files
```

### Guidelines
- Use present tense ("add feature" not "added feature")
- Be specific and descriptive
- Reference issues/PRs when applicable
- Keep subject under 50 characters
- Explain "why" in body, not "what"

---

## Pull Request Process

### PR Title
- Use same format as commit messages
- Clear, descriptive title

### PR Description Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring

## Testing
- [ ] Tested locally
- [ ] No breaking changes
- [ ] Related tests pass

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] No new warnings/errors
- [ ] Tested on Windows/Linux
- [ ] Commit messages are clear

## Related Issues
Closes #123
```

### PR Review Process

1. **Automated Checks**
   - Code builds without errors
   - No obvious style issues

2. **Code Review**
   - Maintainers review code
   - May request changes
   - Constructive feedback provided

3. **Approval**
   - After review, changes may be requested
   - Make updates in same PR
   - Request review again

4. **Merge**
   - Maintainer merges when approved
   - Branch deleted automatically

---

## Reporting Issues

### Bug Reports

**Include:**
- Clear title describing the bug
- Step-by-step reproduction steps
- Expected behavior
- Actual behavior
- OS and version (Windows 10, Ubuntu 20.04, etc.)
- Go version (`go version`)
- Node version (`node --version`)
- Screenshots if applicable

**Example:**
```
Title: Comment scrolling stops at post 5

Steps:
1. Load post list
2. Select post 1
3. Press 'c' to open comments
4. Scroll down with arrow keys
5. Switch to post 5 with 'l' key
6. Try to scroll comments - stops scrolling

Expected: Comments scroll up/down
Actual: Pressing arrows does nothing
```

### Security Issues
- **Do NOT** open a public issue for security vulnerabilities
- Email security details to: security@example.com
- Include: Description, impact, steps to reproduce

---

## Feature Requests

### Propose Features

**Include:**
- Clear title and description
- Motivation (why is this needed?)
- Proposed implementation (if you have ideas)
- Examples of similar features in other apps

**Example:**
```
Title: Add support for multireddits

Description:
Would like to browse multiple subreddits at once

Motivation:
Some users want combined feed from several related subreddits

Proposed Solution:
- Add subreddit list (e.g., "programming+compsci")
- Show posts from both subreddits mixed
- Sort by newest first

References:
- Reddit native subreddit syntax
- Reddit Enhancement Suite does this
```

---

## Workflow Summary

1. **Fork** the repository
2. **Clone** your fork locally
3. **Create** a feature branch
4. **Make** your changes
5. **Test** locally
6. **Commit** with clear messages
7. **Push** to your fork
8. **Create** a Pull Request
9. **Respond** to review feedback
10. **Done!** PR is merged

---

## Development Tips

### Useful Commands

```bash
# Update your fork from upstream
git fetch upstream
git rebase upstream/main

# See what changed
git diff main

# See commit history
git log --oneline

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Clean up local branches
git branch -D feature-branch
```

### Debugging

**TUI Issues:**
```bash
# Rebuild with debug symbols
cd apps/tui
go build -gcflags="all=-N -l" -o redditview .

# Run with verbose output
GODEBUG=gctrace=1 ./redditview
```

**API Issues:**
```bash
# Start with logging
NODE_ENV=development npm start

# Check what requests are made
curl -v http://localhost:3002/api/r/sysadmin.json
```

---

## Getting Help

- **Questions?** Open a GitHub Discussion
- **Bug?** Open an Issue
- **Feature idea?** Open a Feature Request
- **Code help?** Ask in PR comments
- **General question?** Check [README.md](../README.md) or [DOCS_INDEX.md](../DOCS_INDEX.md)

---

## Recognition

Contributors are:
- Mentioned in PR comments
- Listed in commit history
- Credited in future releases (if applicable)

---

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project (MIT License).

---

## Questions?

If anything is unclear:
1. Check existing [Issues](https://github.com/yourusername/redditiew-local/issues)
2. Open a [Discussion](https://github.com/yourusername/redditiew-local/discussions)
3. Email: support@example.com

---

**Thank you for contributing to RedditView! 🎉**

We appreciate all types of contributions:
- Code improvements
- Bug reports
- Documentation fixes
- Feature suggestions
- Testing and feedback

Every contribution helps make RedditView better for everyone! 🚀

---

Happy contributing! 🎓
