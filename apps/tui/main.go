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
	
	footerStyle = lipgloss.NewStyle().
		Background(colorDarkGray).
		Foreground(colorWhite).
		Padding(0, 1)
	
	errorStyle = lipgloss.NewStyle().
		Foreground(colorRed).
		Bold(true)
	
	sectionTitleStyle = lipgloss.NewStyle().
		Foreground(colorOrange).
		Bold(true).
		Padding(1, 0)
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
	// Post Management
	posts         []RedditPostData
	filteredPosts []RedditPostData
	selectedIndex int
	
	// UI Components
	searchInput    textinput.Model
	subredditInput textinput.Model
	spinner        spinner.Model
	
	// State
	subreddit     string
	loading       bool
	error         string
	searching     bool
	selectingSub  bool
	
	// Scroll
	scrollY      int
	scrollHeight int
	
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
	subInput.Placeholder = "Enter subreddit (e.g., golang, rust)..."
	subInput.CharLimit = 50
	
	m := Model{
		client:         NewAPIClient(),
		subreddit:      "golang",
		spinner:        s,
		searchInput:    searchInput,
		subredditInput: subInput,
		loading:        true,
		windowWidth:    120,
		windowHeight:   40,
	}
	
	return m
}

// ============= Message Types =============

type postsLoadedMsg struct {
	posts []RedditPostData
	error error
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
			m.selectedIndex = 0
		}
		m.loading = false
		return m, nil
		
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.scrollHeight = m.windowHeight - 8
		return m, nil
		
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle subreddit selection
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
				m.scrollY = 0
				return m, m.loadPosts(newSub)
			}
		}
		var cmd tea.Cmd
		m.subredditInput, cmd = m.subredditInput.Update(msg)
		return m, cmd
	}
	
	// Handle search
	if m.searching {
		switch msg.String() {
		case "esc":
			m.searching = false
			m.searchInput.Reset()
			m.filterPosts("")
			m.scrollY = 0
			return m, nil
		case "enter":
			m.searching = false
			query := m.searchInput.Value()
			m.filterPosts(query)
			m.scrollY = 0
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
	case "ctrl+f":
		m.searching = true
		m.searchInput.Focus()
		return m, nil
	case "ctrl+r":
		m.selectingSub = true
		m.subredditInput.Focus()
		m.subredditInput.SetValue(m.subreddit)
		return m, nil
	case "f5":
		m.loading = true
		m.scrollY = 0
		return m, m.loadPosts(m.subreddit)
	}
	
	// Navigation
	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case "down", "j":
		if m.selectedIndex < len(m.filteredPosts)-1 {
			m.selectedIndex++
		}
	case "home":
		m.selectedIndex = 0
		m.scrollY = 0
	case "end":
		m.selectedIndex = len(m.filteredPosts) - 1
		m.scrollY = 0
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
}

// ============= Rendering =============

func (m Model) View() string {
	if m.error != "" {
		return m.renderError()
	}
	
	if m.loading {
		return m.renderLoading()
	}
	
	return m.renderMain()
}

func (m Model) renderError() string {
	return errorStyle.Render(fmt.Sprintf("❌ Error: %s\n\nPress q to quit", m.error))
}

func (m Model) renderLoading() string {
	return lipgloss.NewStyle().
		Foreground(colorGold).
		Padding(2, 4).
		Render(fmt.Sprintf("%s Loading r/%s...", m.spinner.View(), m.subreddit))
}

func (m Model) renderMain() string {
	// Header
	header := headerStyle.Render(fmt.Sprintf("  🔥 r/%s  Posts: %d", m.subreddit, len(m.filteredPosts)))
	
	// Info bar
	var infoBar string
	if m.searching {
		infoBar = lipgloss.NewStyle().
			Foreground(colorGold).
			Padding(0, 1).
			Render(fmt.Sprintf("🔍 Search: %s", m.searchInput.View()))
	} else if m.selectingSub {
		infoBar = lipgloss.NewStyle().
			Foreground(colorGold).
			Padding(0, 1).
			Render(fmt.Sprintf("📍 Subreddit: %s", m.subredditInput.View()))
	} else {
		infoBar = lipgloss.NewStyle().
			Foreground(colorGreen).
			Padding(0, 1).
			Render(fmt.Sprintf("▼/▲ (j/k): navigate  Ctrl+F: search  Ctrl+R: subreddit  F5: refresh  q: quit"))
	}
	
	// Content area
	contentHeight := m.windowHeight - 4
	content := m.renderContent(contentHeight)
	
	// Footer
	footer := m.renderFooter()
	
	return fmt.Sprintf("%s\n%s\n%s\n%s", header, infoBar, content, footer)
}

func (m Model) renderContent(height int) string {
	if len(m.filteredPosts) == 0 {
		return lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(2, 2).
			Render("No posts found. Try a different search or subreddit.")
	}
	
	var sb strings.Builder
	
	// Render all posts with selected one expanded
	for i, post := range m.filteredPosts {
		if m.selectedIndex == i {
			sb.WriteString(m.renderSelectedPost(post, m.windowWidth-4))
		} else {
			sb.WriteString(m.renderPostListItem(post, i == m.selectedIndex))
		}
		sb.WriteString("\n")
	}
	
	content := sb.String()
	
	// Apply styling
	styled := lipgloss.NewStyle().
		Foreground(colorGray).
		Padding(0, 1).
		Render(content)
	
	return styled
}

func (m Model) renderPostListItem(post RedditPostData, selected bool) string {
	var sb strings.Builder
	
	// Indicator
	indicator := "  "
	if selected {
		indicator = "▶ "
	}
	
	// Title (truncated)
	title := post.Title
	if len(title) > m.windowWidth-10 {
		title = title[:m.windowWidth-13] + "..."
	}
	
	if selected {
		title = selectedStyle.Render(title)
	}
	
	sb.WriteString(indicator)
	sb.WriteString(title)
	sb.WriteString("\n")
	
	// Meta line
	meta := fmt.Sprintf("  u/%s  •  ⬆ %s  •  💬 %s", 
		post.Author, m.formatNum(post.Score), m.formatNum(post.Comments))
	sb.WriteString(lipgloss.NewStyle().Foreground(colorGold).Render(meta))
	sb.WriteString("\n")
	
	return sb.String()
}

func (m Model) renderSelectedPost(post RedditPostData, width int) string {
	var sb strings.Builder
	
	// Indicator
	sb.WriteString("▼ ")
	
	// Title
	sb.WriteString(selectedStyle.Render(post.Title))
	sb.WriteString("\n\n")
	
	// Meta information
	metaStr := fmt.Sprintf("👤 u/%s  •  r/%s  •  ⬆ %s  •  💬 %s",
		post.Author, post.SubName, m.formatNum(post.Score), m.formatNum(post.Comments))
	sb.WriteString(lipgloss.NewStyle().Foreground(colorGold).Render(metaStr))
	sb.WriteString("\n")
	
	// Separator
	sb.WriteString(strings.Repeat("─", min(width, 80)))
	sb.WriteString("\n\n")
	
	// Post content
	if post.SelfText != "" {
		content := m.wrapText(post.SelfText, width-2)
		sb.WriteString(content)
		sb.WriteString("\n\n")
	}
	
	// URL
	if post.URL != "" && !strings.HasPrefix(post.URL, "https://www.reddit.com") {
		displayURL := post.URL
		if len(displayURL) > width-4 {
			displayURL = displayURL[:width-7] + "..."
		}
		sb.WriteString(lipgloss.NewStyle().Foreground(colorBlue).Render("🔗 " + displayURL))
		sb.WriteString("\n\n")
	}
	
	// Separator before comments
	sb.WriteString(strings.Repeat("─", min(width, 80)))
	sb.WriteString("\n")
	sb.WriteString(sectionTitleStyle.Render("💬 Top Comments"))
	sb.WriteString("\n")
	
	// Placeholder for comments (TODO: load actual comments)
	sb.WriteString(lipgloss.NewStyle().Foreground(colorGray).Italic(true).Render(
		"(Comments loading would go here)"))
	sb.WriteString("\n")
	
	return sb.String()
}

func (m Model) renderFooter() string {
	selectedNum := m.selectedIndex + 1
	totalNum := len(m.filteredPosts)
	
	statusStr := fmt.Sprintf("%d/%d", selectedNum, totalNum)
	if len(m.filteredPosts) == 0 {
		statusStr = "0/0"
	}
	
	footerText := fmt.Sprintf("Post %s  •  Ctrl+F: search  •  Ctrl+R: subreddit  •  F5: refresh  •  q: quit", statusStr)
	return footerStyle.Render(footerText)
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
