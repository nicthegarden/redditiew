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
	colorBlack = lipgloss.Color("#000000")
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
		Foreground(colorWhite).
		Bold(true)
	
	footerStyle = lipgloss.NewStyle().
		Background(colorDarkGray).
		Foreground(colorWhite).
		Padding(0, 1)
	
	focusedStyle = lipgloss.NewStyle().
		Foreground(colorOrange).
		Bold(true)
	
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

// ============= List Item Implementation =============

type PostItem struct {
	post RedditPostData
}

func (p PostItem) FilterValue() string {
	return strings.ToLower(p.post.Title + " " + p.post.Author)
}

func (p PostItem) Title() string {
	return p.post.Title
}

func (p PostItem) Description() string {
	return fmt.Sprintf("u/%s  •  ⬆ %s  •  💬 %s", p.post.Author, formatNum(p.post.Score), formatNum(p.post.Comments))
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
	
	// List component
	list list.Model
	
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
	showDetails   bool
	
	// Detail view scroll
	detailScrollY int
	detailMaxScroll int
	
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
	
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = ""
	l.SetFilteringEnabled(false)
	l.SetShowFilter(false)
	l.DisableQuitKeybindings()
	
	m := Model{
		client:         NewAPIClient(),
		subreddit:      "golang",
		spinner:        s,
		searchInput:    searchInput,
		subredditInput: subInput,
		list:           l,
		loading:        true,
		windowWidth:    120,
		windowHeight:   40,
		showDetails:    false,
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
			m.updateListItems()
		}
		m.loading = false
		m.showDetails = false
		m.detailScrollY = 0
		return m, nil
		
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.updateListSize()
		return m, nil
		
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	
	// List update
	if !m.showDetails {
		m.list, cmd = m.list.Update(msg)
	}
	
	return m, cmd
}

func (m *Model) updateListSize() {
	if m.showDetails {
		// Show both list and details
		listHeight := (m.windowHeight - 8) / 2
		m.list.SetSize(m.windowWidth-2, listHeight)
	} else {
		// Show full list
		m.list.SetSize(m.windowWidth-2, m.windowHeight-8)
	}
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
				m.showDetails = false
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
	
	// Detail view navigation
	if m.showDetails {
		switch msg.String() {
		case "esc", "tab":
			m.showDetails = false
			m.detailScrollY = 0
			return m, nil
		case "up", "k":
			if m.detailScrollY > 0 {
				m.detailScrollY--
			}
			return m, nil
		case "down", "j":
			if m.detailScrollY < m.detailMaxScroll {
				m.detailScrollY++
			}
			return m, nil
		case "home":
			m.detailScrollY = 0
			return m, nil
		case "end":
			m.detailScrollY = m.detailMaxScroll
			return m, nil
		}
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
		m.showDetails = false
		return m, m.loadPosts(m.subreddit)
	case "enter":
		if len(m.filteredPosts) > 0 && m.list.Index() < len(m.filteredPosts) {
			m.showDetails = true
			m.detailScrollY = 0
		}
		return m, nil
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
	header := headerStyle.Render(fmt.Sprintf("  🔥 r/%s  %d posts", m.subreddit, len(m.filteredPosts)))
	
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
		infoBar = m.renderInfoBar()
	}
	
	// Content
	var content string
	if m.showDetails && len(m.filteredPosts) > 0 {
		content = m.renderWithDetails()
	} else {
		content = m.renderListOnly()
	}
	
	// Footer
	footer := m.renderFooter()
	
	return fmt.Sprintf("%s\n%s\n%s\n%s", header, infoBar, content, footer)
}

func (m Model) renderInfoBar() string {
	if m.showDetails {
		return lipgloss.NewStyle().
			Foreground(colorGreen).
			Padding(0, 1).
			Render("▲/▼ (k/j): scroll  Home/End: jump  Esc/Tab: back  Ctrl+F: search  F5: refresh  q: quit")
	}
	return lipgloss.NewStyle().
		Foreground(colorGreen).
		Padding(0, 1).
		Render("▲/▼ (k/j): navigate  Enter: view  Ctrl+F: search  Ctrl+R: subreddit  F5: refresh  q: quit")
}

func (m Model) renderListOnly() string {
	m.updateListSize()
	return m.list.View()
}

func (m Model) renderWithDetails() string {
	// Split view: list on top, details on bottom
	listHeight := (m.windowHeight - 8) / 2
	detailsHeight := m.windowHeight - 8 - listHeight - 1
	
	m.list.SetSize(m.windowWidth-2, listHeight)
	listView := m.list.View()
	
	// Details section
	detailsView := m.renderDetailsSection(detailsHeight)
	
	separator := lipgloss.NewStyle().
		Foreground(colorOrange).
		Render(strings.Repeat("─", m.windowWidth-2))
	
	return fmt.Sprintf("%s\n%s\n%s", listView, separator, detailsView)
}

func (m Model) renderDetailsSection(height int) string {
	if m.list.Index() >= len(m.filteredPosts) {
		return ""
	}
	
	post := m.filteredPosts[m.list.Index()]
	
	var sb strings.Builder
	
	// Title
	sb.WriteString(focusedStyle.Render(fmt.Sprintf("📄 %s\n", post.Title)))
	
	// Meta
	sb.WriteString(lipgloss.NewStyle().Foreground(colorGold).Render(
		fmt.Sprintf("👤 u/%s  •  r/%s  •  ⬆ %s  •  💬 %s\n\n",
			post.Author, post.SubName, formatNum(post.Score), formatNum(post.Comments))))
	
	// Content
	if post.SelfText != "" {
		content := wrapText(post.SelfText, m.windowWidth-4)
		lines := strings.Split(content, "\n")
		
		// Apply scrolling
		startLine := m.detailScrollY
		endLine := startLine + height - 3
		if endLine > len(lines) {
			endLine = len(lines)
		}
		
		// Calculate max scroll
		m.detailMaxScroll = max(0, len(lines)-height+3)
		
		if startLine < len(lines) {
			visibleLines := lines[startLine:endLine]
			sb.WriteString(strings.Join(visibleLines, "\n"))
		}
	}
	
	// URL
	if post.URL != "" && !strings.HasPrefix(post.URL, "https://www.reddit.com") {
		displayURL := post.URL
		if len(displayURL) > m.windowWidth-10 {
			displayURL = displayURL[:m.windowWidth-13] + "..."
		}
		sb.WriteString(lipgloss.NewStyle().Foreground(colorBlue).Render("\n🔗 " + displayURL))
	}
	
	return lipgloss.NewStyle().
		Foreground(colorGray).
		Padding(0, 1).
		Height(height).
		Render(sb.String())
}

func (m Model) renderFooter() string {
	if m.showDetails {
		return footerStyle.Render("Esc/Tab: back to list  •  Ctrl+F: search  •  F5: refresh  •  q: quit")
	}
	
	status := "no posts"
	if len(m.filteredPosts) > 0 {
		status = fmt.Sprintf("%d/%d", m.list.Index()+1, len(m.filteredPosts))
	}
	
	return footerStyle.Render(fmt.Sprintf("Post %s  •  Enter: view details  •  Ctrl+F: search  •  F5: refresh  •  q: quit", status))
}

// ============= Utilities =============

func wrapText(text string, width int) string {
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

func formatNum(n int) string {
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
