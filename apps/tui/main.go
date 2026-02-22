package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/spinner"
)

// ============= Constants =============

const apiBaseURL = "http://localhost:3002/api"

var (
	// Colors
	colorOrange = lipgloss.Color("#FF4500")
	colorOrangeDark = lipgloss.Color("#FF6B35")
	colorGold = lipgloss.Color("#FFD700")
	colorGreen = lipgloss.Color("#90EE90")
	colorWhite = lipgloss.Color("#FFFFFF")
	colorGray = lipgloss.Color("#CCCCCC")
	colorDarkGray = lipgloss.Color("#333333")
	colorBlue = lipgloss.Color("#87CEEB")
	colorRed = lipgloss.Color("#FF0000")
	
	// Styles
	headerStyle = lipgloss.NewStyle().
		Background(colorOrange).
		Foreground(colorWhite).
		Bold(true).
		Padding(0, 1)
	
	selectedStyle = lipgloss.NewStyle().
		Background(colorOrangeDark).
		Foreground(colorWhite)
	
	focusedBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorOrange)
	
	unfocusedBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorDarkGray)
	
	footerStyle = lipgloss.NewStyle().
		Background(colorDarkGray).
		Foreground(colorWhite).
		Padding(0, 1)
	
	errorStyle = lipgloss.NewStyle().
		Foreground(colorRed).
		Bold(true)
)

// ============= Data Models =============

type RedditPostData struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Score     int     `json:"score"`
	Created   float64 `json:"created_utc"`
	Comments  int     `json:"num_comments"`
	SelfText  string  `json:"selftext"`
	URL       string  `json:"url"`
	SubName   string  `json:"subreddit"`
	Permalink string  `json:"permalink"`
}

type RedditPost struct {
	Kind string           `json:"kind"`
	Data RedditPostData `json:"data"`
}

type RedditResponse struct {
	Data struct {
		Children []RedditPost `json:"children"`
		After    string       `json:"after"`
	} `json:"data"`
}

type Comment struct {
	ID        string
	Author    string
	Body      string
	Score     int
	Created   float64
	Depth     int
	Replies   []*Comment
	Collapsed bool
}

// ============= Pane Types =============

type Pane int

const (
	PanePostList Pane = iota
	PanePostDetail
	PaneComments
)

// ============= API Client =============

type APIClient struct {
	baseURL string
}

func NewAPIClient() *APIClient {
	return &APIClient{baseURL: apiBaseURL}
}

func (c *APIClient) FetchPosts(subreddit string) ([]RedditPostData, error) {
	resp, err := http.Get(fmt.Sprintf("%s/r/%s.json?limit=50", c.baseURL, subreddit))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var result RedditResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Reddit API response: %w", err)
	}

	posts := make([]RedditPostData, 0, len(result.Data.Children))
	for _, post := range result.Data.Children {
		if post.Kind == "t3" {
			posts = append(posts, post.Data)
		}
	}

	return posts, nil
}

// ============= Main Model =============

type Model struct {
	// Navigation
	focusedPane Pane
	
	// Post Management
	posts         []RedditPostData
	filteredPosts []RedditPostData
	selectedPost  *RedditPostData
	selectedIndex int
	
	// Comments
	comments      []*Comment
	commentsLoaded bool
	
	// UI Components
	searchInput     textinput.Model
	subredditInput  textinput.Model
	spinner         spinner.Model
	
	// State
	subreddit     string
	loading       bool
	error         string
	searching     bool
	selectingSub  bool
	
	// Scroll positions
	detailScrollY  int
	commentsScrollY int
	
	// Layout
	windowWidth  int
	windowHeight int
	paneWidth    int
	paneHeight   int
	
	// API
	client *APIClient
}

func initialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	
	searchInput := textinput.New()
	searchInput.Placeholder = "Search posts..."
	searchInput.CharLimit = 100
	
	subInput := textinput.New()
	subInput.Placeholder = "Enter subreddit..."
	subInput.CharLimit = 50
	
	m := Model{
		client:        NewAPIClient(),
		subreddit:     "golang",
		focusedPane:   PanePostList,
		spinner:       s,
		searchInput:   searchInput,
		subredditInput: subInput,
		loading:       true,
		windowWidth:   120,
		windowHeight:  40,
	}
	
	m.calculatePaneDimensions()
	return m
}

// ============= Message Types =============

type postsLoadedMsg struct {
	posts []RedditPostData
	error error
}

type commentsLoadedMsg struct {
	comments []*Comment
	error    error
}

type errMsg struct {
	err error
}

// ============= Commands =============

func (m Model) loadPosts(subreddit string) tea.Cmd {
	return func() tea.Msg {
		posts, err := m.client.FetchPosts(subreddit)
		if err != nil {
			return postsLoadedMsg{nil, err}
		}
		return postsLoadedMsg{posts, nil}
	}
}

// ============= Update Logic =============

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadPosts(m.subreddit),
		m.spinner.Tick,
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd = m.handleKeyPress(msg)
		return m, cmd
		
	case postsLoadedMsg:
		if msg.error != nil {
			m.error = msg.error.Error()
		} else {
			m.posts = msg.posts
			m.filteredPosts = msg.posts
			if len(msg.posts) > 0 {
				m.selectedPost = &msg.posts[0]
				m.selectedIndex = 0
			}
		}
		m.loading = false
		return m, nil
		
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.calculatePaneDimensions()
		return m, nil
		
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	
	return m, nil
}

func (m *Model) calculatePaneDimensions() {
	// Divide equally: Total width minus borders/margins
	// Each pane takes (width - 4) / 3 (4 is minimum for separators)
	availableWidth := m.windowWidth - 4
	if availableWidth < 60 {
		availableWidth = 60 // Minimum for 3 panes
	}
	m.paneWidth = availableWidth / 3
	m.paneHeight = m.windowHeight - 4 // Header + footer
	if m.paneHeight < 10 {
		m.paneHeight = 10
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle subreddit selection mode
	if m.selectingSub {
		switch msg.String() {
		case "esc":
			m.selectingSub = false
			m.subredditInput.Reset()
			return m, nil
		case "enter":
			newSub := strings.TrimSpace(m.subredditInput.Value())
			if newSub != "" {
				m.subreddit = newSub
				m.subredditInput.Reset()
				m.selectingSub = false
				m.loading = true
				return m, m.loadPosts(newSub)
			}
		}
		var cmd tea.Cmd
		m.subredditInput, cmd = m.subredditInput.Update(msg)
		return m, cmd
	}
	
	// Handle search mode
	if m.searching {
		switch msg.String() {
		case "esc":
			m.searching = false
			m.searchInput.Reset()
			m.filterPosts("")
			return m, nil
		case "enter":
			m.searching = false
			query := m.searchInput.Value()
			m.filterPosts(query)
			return m, nil
		}
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.filterPosts(m.searchInput.Value())
		return m, cmd
	}
	
	// Global shortcuts
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "/":
		m.searching = true
		m.searchInput.Focus()
		return m, nil
	case "s":
		m.selectingSub = true
		m.subredditInput.Focus()
		return m, nil
	case "tab":
		// Cycle through panes
		m.focusedPane = (m.focusedPane + 1) % 3
		return m, nil
	}
	
	// Pane-specific navigation
	switch m.focusedPane {
	case PanePostList:
		return m.handlePostListKeyPress(msg)
	case PanePostDetail:
		return m.handlePostDetailKeyPress(msg)
	case PaneComments:
		return m.handleCommentsKeyPress(msg)
	}
	
	return m, nil
}

func (m Model) handlePostListKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
			if len(m.filteredPosts) > m.selectedIndex {
				m.selectedPost = &m.filteredPosts[m.selectedIndex]
				m.detailScrollY = 0
				m.commentsScrollY = 0
				m.commentsLoaded = false
			}
		}
	case "down", "j":
		if m.selectedIndex < len(m.filteredPosts)-1 {
			m.selectedIndex++
			if len(m.filteredPosts) > m.selectedIndex {
				m.selectedPost = &m.filteredPosts[m.selectedIndex]
				m.detailScrollY = 0
				m.commentsScrollY = 0
				m.commentsLoaded = false
			}
		}
	}
	return m, nil
}

func (m Model) handlePostDetailKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.detailScrollY > 0 {
			m.detailScrollY--
		}
	case "down", "j":
		m.detailScrollY++
	}
	return m, nil
}

func (m Model) handleCommentsKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.commentsScrollY > 0 {
			m.commentsScrollY--
		}
	case "down", "j":
		m.commentsScrollY++
	case "c":
		// Toggle collapse (placeholder for now)
	}
	return m, nil
}

func (m *Model) filterPosts(query string) {
	query = strings.ToLower(query)
	if query == "" {
		m.filteredPosts = m.posts
	} else {
		m.filteredPosts = []RedditPostData{}
		for _, post := range m.posts {
			if strings.Contains(strings.ToLower(post.Title), query) ||
				strings.Contains(strings.ToLower(post.Author), query) {
				m.filteredPosts = append(m.filteredPosts, post)
			}
		}
	}
	m.selectedIndex = 0
	if len(m.filteredPosts) > 0 {
		m.selectedPost = &m.filteredPosts[0]
	}
}

// ============= Rendering =============

func (m Model) View() string {
	if m.error != "" {
		return m.renderError()
	}
	
	if m.loading {
		return m.renderLoading()
	}
	
	return m.render3PaneLayout()
}

func (m Model) renderError() string {
	return errorStyle.Render(fmt.Sprintf("❌ Error: %s", m.error))
}

func (m Model) renderLoading() string {
	return lipgloss.NewStyle().
		Foreground(colorGold).
		Padding(1, 2).
		Render(fmt.Sprintf("%s Loading posts from r/%s...", m.spinner.View(), m.subreddit))
}

func (m Model) render3PaneLayout() string {
	// Header
	header := headerStyle.Render(fmt.Sprintf("  🔥 r/%s  %d posts", m.subreddit, len(m.filteredPosts)))
	
	// Search/subreddit bar
	var searchBar string
	if m.searching {
		searchBar = lipgloss.NewStyle().
			Foreground(colorGold).
			Render(fmt.Sprintf("🔍 Search: %s", m.searchInput.View()))
	} else if m.selectingSub {
		searchBar = lipgloss.NewStyle().
			Foreground(colorGold).
			Render(fmt.Sprintf("📍 Go to: %s", m.subredditInput.View()))
	} else {
		searchBar = lipgloss.NewStyle().
			Foreground(colorGreen).
			Render("Tab: focus | j/k: navigate | /: search | s: subreddit | q: quit")
	}
	
	// Render three panes side by side
	leftPane := m.renderPostListPane()
	middlePane := m.renderPostDetailPane()
	rightPane := m.renderCommentsPane()
	
	// Join panes horizontally
	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		middlePane,
		rightPane,
	)
	
	// Footer
	footer := footerStyle.Render("Tab: focus | j/k: navigate | /: search | s: subreddit | q: quit")
	
	return fmt.Sprintf("%s\n%s\n%s\n%s", header, searchBar, panes, footer)
}

func (m Model) renderPostListPane() string {
	title := "📬 Posts"
	if m.focusedPane == PanePostList {
		title = " 📬 Posts (focused) "
	}
	
	var content string
	if len(m.filteredPosts) == 0 {
		content = "No posts loaded"
	} else {
		var sb strings.Builder
		for i, post := range m.filteredPosts {
			if i >= m.paneHeight-3 {
				break
			}
			
			// Post title (truncated)
			title := post.Title
			if len(title) > m.paneWidth-2 {
				title = title[:m.paneWidth-5] + "..."
			}
			
			// Highlight selected
			if i == m.selectedIndex {
				title = selectedStyle.Render(title)
			} else {
				title = lipgloss.NewStyle().Foreground(colorGray).Render(title)
			}
			sb.WriteString(title)
			sb.WriteString("\n")
			
			// Meta line
			meta := fmt.Sprintf("u/%s ⬆%d 💬%d", post.Author, post.Score, post.Comments)
			if len(meta) > m.paneWidth-2 {
				meta = meta[:m.paneWidth-5] + "..."
			}
			sb.WriteString(lipgloss.NewStyle().Foreground(colorGold).Render(meta))
			sb.WriteString("\n")
		}
		content = sb.String()
	}
	
	// Apply border based on focus
	var borderStyle lipgloss.Style
	if m.focusedPane == PanePostList {
		borderStyle = focusedBorderStyle
	} else {
		borderStyle = unfocusedBorderStyle
	}
	
	pane := borderStyle.
		Width(m.paneWidth).
		Height(m.paneHeight).
		Padding(0, 1).
		Render(title + "\n" + content)
	
	return pane
}

func (m Model) renderPostDetailPane() string {
	title := "📄 Details"
	if m.focusedPane == PanePostDetail {
		title = " 📄 Details (focused) "
	}
	
	var content string
	if m.selectedPost == nil {
		content = "Select a post..."
	} else {
		post := m.selectedPost
		var sb strings.Builder
		
		// Title
		postTitle := post.Title
		if len(postTitle) > m.paneWidth-2 {
			postTitle = m.wrapText(postTitle, m.paneWidth-2)
		}
		sb.WriteString(postTitle)
		sb.WriteString("\n\n")
		
		// Meta
		sb.WriteString(fmt.Sprintf("u/%s\n", post.Author))
		sb.WriteString(fmt.Sprintf("r/%s\n", post.SubName))
		sb.WriteString(fmt.Sprintf("⬆ %s | 💬 %s\n\n", m.formatNum(post.Score), m.formatNum(post.Comments)))
		
		// Content
		if post.SelfText != "" {
			contentText := m.wrapText(post.SelfText, m.paneWidth-2)
			sb.WriteString(contentText)
			sb.WriteString("\n\n")
		}
		
		// URL
		if post.URL != "" && !strings.HasPrefix(post.URL, "https://www.reddit.com") {
			displayURL := post.URL
			if len(displayURL) > m.paneWidth-4 {
				displayURL = displayURL[:m.paneWidth-7] + "..."
			}
			sb.WriteString("🔗 " + displayURL)
		}
		
		content = sb.String()
	}
	
	// Apply border based on focus
	var borderStyle lipgloss.Style
	if m.focusedPane == PanePostDetail {
		borderStyle = focusedBorderStyle
	} else {
		borderStyle = unfocusedBorderStyle
	}
	
	pane := borderStyle.
		Width(m.paneWidth).
		Height(m.paneHeight).
		Padding(0, 1).
		Render(title + "\n" + content)
	
	return pane
}

func (m Model) renderCommentsPane() string {
	title := "💬 Comments"
	if m.focusedPane == PaneComments {
		title = " 💬 Comments (focused) "
	}
	
	var content string
	if m.selectedPost == nil {
		content = "Select a post..."
	} else if !m.commentsLoaded {
		content = "No comments loaded"
	} else if len(m.comments) == 0 {
		content = "No comments"
	} else {
		var sb strings.Builder
		m.renderCommentsTree(&sb, m.comments, 0, m.paneWidth)
		content = sb.String()
	}
	
	// Apply border based on focus
	var borderStyle lipgloss.Style
	if m.focusedPane == PaneComments {
		borderStyle = focusedBorderStyle
	} else {
		borderStyle = unfocusedBorderStyle
	}
	
	pane := borderStyle.
		Width(m.paneWidth).
		Height(m.paneHeight).
		Padding(0, 1).
		Render(title + "\n" + content)
	
	return pane
}

func (m Model) renderCommentsTree(sb *strings.Builder, comments []*Comment, depth int, maxWidth int) {
	for _, comment := range comments {
		if comment == nil {
			continue
		}
		
		indent := strings.Repeat("  ", depth)
		lineWidth := maxWidth - (depth * 2) - 2
		if lineWidth < 10 {
			lineWidth = 10
		}
		
		// Author and score
		author := fmt.Sprintf("%su/%s ⬆%d\n", indent, comment.Author, comment.Score)
		sb.WriteString(author)
		
		// Body with wrapping
		body := comment.Body
		if len(body) > lineWidth {
			body = body[:lineWidth-3] + "..."
		}
		sb.WriteString(fmt.Sprintf("%s%s\n\n", indent, body))
		
		// Recursively render replies
		if !comment.Collapsed && len(comment.Replies) > 0 {
			m.renderCommentsTree(sb, comment.Replies, depth+1, maxWidth)
		}
	}
}

// ============= Utilities =============

func (m Model) wrapText(text string, width int) string {
	if width <= 0 || text == "" {
		return ""
	}
	
	words := strings.Fields(text)
	var lines []string
	var currentLine string
	
	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return strings.Join(lines, "\n")
}

func (m Model) formatNum(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ============= Main =============

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
