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
	"github.com/charmbracelet/bubbles/list"
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
	ID      string
	Author  string
	Body    string
	Score   int
	Created float64
	Depth   int
	Replies []*Comment
}

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

// ============= List Item Implementation =============

type PostItem struct {
	post RedditPostData
}

func (p PostItem) FilterValue() string {
	return strings.ToLower(p.post.Title + " " + p.post.Author)
}

func (p PostItem) Title() string {
	const maxLen = 65
	title := p.post.Title
	if len(title) > maxLen {
		title = title[:maxLen-3] + "..."
	}
	return title
}

func (p PostItem) Description() string {
	return fmt.Sprintf("u/%s • ⬆%d • 💬%d", p.post.Author, p.post.Score, p.post.Comments)
}

// ============= Screen Types =============

type Screen int

const (
	ScreenPostList Screen = iota
	ScreenPostDetail
	ScreenComments
	ScreenSubredditSelect
)

// ============= Main Model =============

type Model struct {
	// Navigation
	currentScreen Screen
	
	// Post Management
	posts         []RedditPostData
	filteredPosts []RedditPostData
	selectedPost  *RedditPostData
	selectedIndex int
	
	// Comments
	comments []Comment
	commentsLoading bool
	
	// UI Components
	list          list.Model
	searchInput   textinput.Model
	subredditInput textinput.Model
	spinner       spinner.Model
	
	// State
	subreddit     string
	loading       bool
	error         string
	searching     bool
	selectingSub  bool
	
	// Layout
	windowWidth  int
	windowHeight int
	
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
	
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = ""
	l.SetFilteringEnabled(false)
	
	m := Model{
		client:         NewAPIClient(),
		subreddit:      "golang",
		currentScreen:  ScreenPostList,
		spinner:        s,
		searchInput:    searchInput,
		subredditInput: subInput,
		list:           l,
		loading:        true,
		windowWidth:    200,
		windowHeight:   50,
	}
	
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

type windowSizeMsg struct {
	width  int
	height int
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
		
	case postsLoadedMsg:
		if msg.error != nil {
			m.error = msg.error.Error()
		} else {
			m.posts = msg.posts
			m.filteredPosts = msg.posts
			m.updateListItems()
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
		m.list.SetSize(m.windowWidth, m.windowHeight-8)
		return m, nil
		
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	
	// Handle screen-specific shortcuts
	switch m.currentScreen {
	case ScreenPostList:
		return m.handlePostListKeyPress(msg)
	case ScreenPostDetail:
		return m.handlePostDetailKeyPress(msg)
	case ScreenComments:
		return m.handleCommentsKeyPress(msg)
	}
	
	return m, nil
}

func (m Model) handlePostListKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
			if len(m.filteredPosts) > m.selectedIndex {
				m.selectedPost = &m.filteredPosts[m.selectedIndex]
			}
		}
	case "down", "j":
		if m.selectedIndex < len(m.filteredPosts)-1 {
			m.selectedIndex++
			if len(m.filteredPosts) > m.selectedIndex {
				m.selectedPost = &m.filteredPosts[m.selectedIndex]
			}
		}
	case "enter":
		if m.selectedPost != nil {
			m.currentScreen = ScreenPostDetail
		}
	case "c":
		if m.selectedPost != nil {
			m.currentScreen = ScreenComments
			m.commentsLoading = true
		}
	}
	
	return m, nil
}

func (m Model) handlePostDetailKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.currentScreen = ScreenPostList
	case "q", "ctrl+c":
		return m, tea.Quit
	case "c":
		m.currentScreen = ScreenComments
		m.commentsLoading = true
	}
	return m, nil
}

func (m Model) handleCommentsKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.currentScreen = ScreenPostDetail
	case "q", "ctrl+c":
		return m, tea.Quit
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
	m.updateListItems()
}

func (m *Model) updateListItems() {
	items := make([]list.Item, len(m.filteredPosts))
	for i, post := range m.filteredPosts {
		items[i] = PostItem{post}
	}
	m.list.SetItems(items)
}

// ============= Rendering =============

func (m Model) View() string {
	if m.error != "" {
		return m.renderError()
	}
	
	if m.loading {
		return m.renderLoading()
	}
	
	switch m.currentScreen {
	case ScreenPostList:
		return m.renderPostList()
	case ScreenPostDetail:
		return m.renderPostDetail()
	case ScreenComments:
		return m.renderComments()
	default:
		return m.renderPostList()
	}
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

func (m Model) renderPostList() string {
	header := headerStyle.Render(fmt.Sprintf("  🔥 r/%s  %d posts", m.subreddit, len(m.filteredPosts)))
	
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
			Render("j/k or ↓↑: navigate | /: search | s: subreddit | Enter: view | c: comments | q: quit")
	}
	
	listHeight := m.windowHeight - 8
	m.list.SetSize(m.windowWidth, listHeight)
	
	footer := footerStyle.Render("j/k: navigate | /: search | s: subreddit | Enter: view | c: comments | q: quit")
	
	return fmt.Sprintf("%s\n%s\n%s\n%s", header, searchBar, m.list.View(), footer)
}

func (m Model) renderPostDetail() string {
	if m.selectedPost == nil {
		return "No post selected"
	}
	
	post := m.selectedPost
	maxTitleLen := m.windowWidth - 6
	title := post.Title
	if len(title) > maxTitleLen {
		title = title[:maxTitleLen-3] + "..."
	}
	
	header := headerStyle.Render(fmt.Sprintf("  %s", title))
	
	meta := lipgloss.NewStyle().
		Foreground(colorGold).
		Padding(1, 1).
		Render(fmt.Sprintf("👤 u/%s  r/%s  ⬆ %s  💬 %s", 
			post.Author, post.SubName, m.formatNum(post.Score), m.formatNum(post.Comments)))
	
	contentWidth := m.windowWidth - 4
	content := lipgloss.NewStyle().
		Foreground(colorGray).
		Padding(1, 1).
		Render(m.wrapText(post.SelfText, contentWidth))
	
	var urlStr string
	if post.URL != "" && !strings.HasPrefix(post.URL, "https://www.reddit.com") {
		displayURL := post.URL
		if len(displayURL) > m.windowWidth-8 {
			displayURL = displayURL[:m.windowWidth-11] + "..."
		}
		urlStr = lipgloss.NewStyle().
			Foreground(colorBlue).
			Padding(1, 1).
			Render(fmt.Sprintf("🔗 %s", displayURL))
	}
	
	actions := lipgloss.NewStyle().
		Foreground(colorGold).
		Padding(1, 1).
		Render("⬆ Upvote  ⬇ Downvote  💾 Save  🔗 Open on Reddit")
	
	footer := footerStyle.Render("c: comments | b: back | q: quit")
	
	result := header + "\n" + meta + "\n" + content
	if urlStr != "" {
		result += "\n" + urlStr
	}
	result += "\n" + actions + "\n" + footer
	
	return result
}

func (m Model) renderComments() string {
	if m.selectedPost == nil {
		return "No post selected"
	}
	
	post := m.selectedPost
	header := headerStyle.Render(fmt.Sprintf("  Comments (%s)", m.formatNum(post.Comments)))
	
	var content string
	if m.commentsLoading {
		content = lipgloss.NewStyle().
			Foreground(colorGold).
			Padding(1, 1).
			Render(fmt.Sprintf("%s Loading comments...", m.spinner.View()))
	} else if len(m.comments) == 0 {
		content = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(1, 1).
			Render("No comments to display")
	} else {
		var sb strings.Builder
		for i, comment := range m.comments {
			if i >= m.windowHeight-8 {
				break
			}
			indent := strings.Repeat("  ", comment.Depth)
			sb.WriteString(fmt.Sprintf("%su/%s  ⬆ %d\n", indent, comment.Author, comment.Score))
			body := comment.Body
			if len(body) > m.windowWidth-int(comment.Depth)*2-10 {
				body = body[:m.windowWidth-int(comment.Depth)*2-13] + "..."
			}
			sb.WriteString(fmt.Sprintf("%s%s\n\n", indent, body))
		}
		content = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(1, 1).
			Render(sb.String())
	}
	
	footer := footerStyle.Render("b: back | q: quit")
	
	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
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

// ============= Main =============

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
