package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

// ============= Configuration =============

type AppConfig struct {
	TUI struct {
		DefaultSubreddit string `json:"default_subreddit"`
		PostsPerPage     int    `json:"posts_per_page"`
		ListHeight       int    `json:"list_height"`
		MaxTitleLength   int    `json:"max_title_length"`
	} `json:"tui"`
	Web struct {
		DefaultSubreddit string `json:"default_subreddit"`
		PostsPerPage     int    `json:"posts_per_page"`
		Theme            string `json:"theme"`
	} `json:"web"`
	API struct {
		BaseURL        string `json:"base_url"`
		TimeoutSeconds int    `json:"timeout_seconds"`
	} `json:"api"`
}

var appConfig AppConfig

func loadConfig() error {
	// Try to load from config.json in parent directory
	configPath := "../../config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		// If not found, use defaults
		appConfig = AppConfig{}
		appConfig.TUI.DefaultSubreddit = "sysadmin"
		appConfig.TUI.PostsPerPage = 50
		appConfig.TUI.ListHeight = 10
		appConfig.TUI.MaxTitleLength = 80
		appConfig.API.BaseURL = "http://localhost:3002/api"
		appConfig.API.TimeoutSeconds = 10
		return nil
	}
	
	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config.json: %w", err)
	}
	
	// Set defaults for any missing values
	if appConfig.TUI.DefaultSubreddit == "" {
		appConfig.TUI.DefaultSubreddit = "sysadmin"
	}
	if appConfig.TUI.PostsPerPage == 0 {
		appConfig.TUI.PostsPerPage = 50
	}
	if appConfig.TUI.ListHeight == 0 {
		appConfig.TUI.ListHeight = 10
	}
	if appConfig.TUI.MaxTitleLength == 0 {
		appConfig.TUI.MaxTitleLength = 80
	}
	if appConfig.API.BaseURL == "" {
		appConfig.API.BaseURL = "http://localhost:3002/api"
	}
	if appConfig.API.TimeoutSeconds == 0 {
		appConfig.API.TimeoutSeconds = 10
	}
	
	return nil
}

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
	return &APIClient{baseURL: appConfig.API.BaseURL}
}

func (c *APIClient) FetchPosts(subreddit string) ([]RedditPostData, error) {
	resp, err := http.Get(fmt.Sprintf("%s/r/%s.json?limit=%d", c.baseURL, subreddit, appConfig.TUI.PostsPerPage))
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

// SearchPosts performs a Reddit-wide search
func (c *APIClient) SearchPosts(query string) ([]RedditPostData, error) {
	if query == "" {
		return []RedditPostData{}, nil
	}
	
	searchURL := fmt.Sprintf("%s/search.json?q=%s&type=link&limit=%d", 
		c.baseURL, 
		url.QueryEscape(query),
		appConfig.TUI.PostsPerPage)
	
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var result RedditResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	posts := make([]RedditPostData, 0, len(result.Data.Children))
	for _, post := range result.Data.Children {
		if post.Kind == "t3" {
			posts = append(posts, post.Data)
		}
	}

	return posts, nil
}

// FetchComments fetches top-level comments for a post
func (c *APIClient) FetchComments(subreddit, postID string) ([]*Comment, error) {
	commentsURL := fmt.Sprintf("%s/r/%s/comments/%s/", c.baseURL, subreddit, postID)
	
	resp, err := http.Get(commentsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	// Try to parse as array first (Reddit's native format)
	var resultArray []map[string]interface{}
	var resultSingle map[string]interface{}
	
	// Try array format first
	if err := json.Unmarshal(data, &resultArray); err == nil && len(resultArray) >= 2 {
		// Reddit returns [posts, comments] format
		commentsData := resultArray[1]
		if dataField, ok := commentsData["data"]; ok {
			dataMap := dataField.(map[string]interface{})
			return parseComments(dataMap)
		}
	}
	
	// Fall back to single object format
	if err := json.Unmarshal(data, &resultSingle); err == nil {
		if dataField, ok := resultSingle["data"]; ok {
			dataMap := dataField.(map[string]interface{})
			return parseComments(dataMap)
		}
	}

	return nil, nil
}

func parseComments(dataMap map[string]interface{}) ([]*Comment, error) {
	childrenInterface, ok := dataMap["children"].([]interface{})
	if !ok {
		return nil, nil
	}

	comments := make([]*Comment, 0)
	for _, childInterface := range childrenInterface {
		childMap, ok := childInterface.(map[string]interface{})
		if !ok {
			continue
		}
		
		kind, ok := childMap["kind"].(string)
		if !ok || kind != "t1" {
			continue // Skip non-comments
		}

		data, ok := childMap["data"].(map[string]interface{})
		if !ok {
			continue
		}

		comment := &Comment{
			Author: toString(data["author"]),
			Body:   toString(data["body"]),
			Score:  toInt(data["score"]),
		}

		if idVal, ok := data["id"].(string); ok {
			comment.ID = idVal
		}

		comments = append(comments, comment)
		if len(comments) >= 5 {
			break // Limit to top 5 comments
		}
	}

	return comments, nil
}

// Helper functions for type conversions
func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	}
	return 0
}

// openURL opens a URL in the default browser
func openURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("empty URL")
	}
	
	// Determine the command based on OS
	var cmd *exec.Cmd
	switch os.Getenv("GOOS") {
	case "darwin":
		cmd = exec.Command("open", urlStr)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", urlStr)
	default:
		// Linux and other Unix-like systems
		cmd = exec.Command("xdg-open", urlStr)
	}
	
	return cmd.Run()
}

// ============= Main Model =============

type Model struct {
	// Post Management
	posts         []RedditPostData
	filteredPosts []RedditPostData
	
	// Comments
	comments         []*Comment
	showComments     bool
	commentsScrollY  int
	commentsMaxScroll int
	commentsLoading  bool
	
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
		subreddit:      appConfig.TUI.DefaultSubreddit,
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

type searchResultsMsg struct {
	posts  []RedditPostData
	query  string
	error  error
}

type commentsLoadedMsg struct {
	comments []*Comment
	error    error
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

func (m Model) searchReddit(query string) tea.Cmd {
	return func() tea.Msg {
		if query == "" {
			return searchResultsMsg{[]RedditPostData{}, "", nil}
		}
		posts, err := m.client.SearchPosts(query)
		if err != nil {
			return searchResultsMsg{nil, query, err}
		}
		return searchResultsMsg{posts, query, nil}
	}
}

func (m Model) loadComments(subreddit, postID string) tea.Cmd {
	return func() tea.Msg {
		comments, err := m.client.FetchComments(subreddit, postID)
		return commentsLoadedMsg{comments, err}
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
	var handled bool
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd, handled = m.handleKeyPress(msg)
		if handled {
			return m, cmd
		}
		// If not handled, fall through to list update
		
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
	
	case searchResultsMsg:
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
	
	case commentsLoadedMsg:
		if msg.error != nil {
			m.error = msg.error.Error()
		} else {
			m.comments = msg.comments
			m.commentsLoading = false
			// Calculate max scroll for comments using actual details height
			listHeight := (m.windowHeight - 8) / 2
			detailsHeight := m.windowHeight - 8 - listHeight - 1
			m = m.calculateCommentsMaxScroll(detailsHeight)
		}
		return m, nil
		
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.updateListSize()
		// Recalculate max scroll for comments on window resize
		if m.showComments {
			detailsHeight := m.windowHeight - 8 - (m.windowHeight-8)/2 - 1
			m = m.calculateCommentsMaxScroll(detailsHeight)
		}
		return m, nil
		
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	
	// List update for navigation keys
	if !m.showDetails && !m.searching && !m.selectingSub {
		m.list, cmd = m.list.Update(msg)
	}
	
	return m, cmd
}

func (m *Model) updateListSize() {
	if m.showDetails {
		// Show both list and details: list gets 6 items max
		listHeight := min(6, m.windowHeight/6)
		if listHeight < 3 {
			listHeight = 3
		}
		m.list.SetSize(m.windowWidth-2, listHeight)
	} else {
		// Show full list: limit to 10 items for better readability
		listHeight := min(10, m.windowHeight-8)
		m.list.SetSize(m.windowWidth-2, listHeight)
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd, bool) {
	// Handle subreddit selection
	if m.selectingSub {
		switch msg.String() {
		case "esc":
			m.selectingSub = false
			m.subredditInput.Reset()
			return m, nil, true
		case "enter":
			newSub := strings.TrimSpace(m.subredditInput.Value())
			if newSub != "" {
				m.subreddit = newSub
				m.subredditInput.Reset()
				m.selectingSub = false
				m.loading = true
				m.showDetails = false
				return m, m.loadPosts(newSub), true
			}
		}
		var cmd tea.Cmd
		m.subredditInput, cmd = m.subredditInput.Update(msg)
		return m, cmd, true
	}
	
	// Handle search
	if m.searching {
		switch msg.String() {
		case "esc":
			m.searching = false
			m.searchInput.Reset()
			m.filterPosts("")
			return m, nil, true
		case "enter":
			m.searching = false
			query := m.searchInput.Value()
			if query != "" {
				// Perform Reddit-wide search
				m.loading = true
				return m, m.searchReddit(query), true
			}
			return m, nil, true
		}
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		// Live filter as user types
		m.filterPosts(m.searchInput.Value())
		return m, cmd, true
	}
	
	// Detail view navigation
	if m.showDetails {
		switch msg.String() {
		case "esc", "tab":
			if m.showComments {
				m.showComments = false
				m.commentsScrollY = 0
				return m, nil, true
			}
			m.showDetails = false
			m.detailScrollY = 0
			return m, nil, true
		case "up":
			// Up arrow: scroll comments if open, else scroll details
			if m.showComments {
				if m.commentsScrollY > 0 {
					m.commentsScrollY--
				}
			} else {
				if m.detailScrollY > 0 {
					m.detailScrollY--
				}
			}
			return m, nil, true
		case "down":
			// Down arrow: scroll comments if open, else scroll details
			if m.showComments {
				if m.commentsScrollY < m.commentsMaxScroll {
					m.commentsScrollY++
				}
			} else {
				if m.detailScrollY < m.detailMaxScroll {
					m.detailScrollY++
				}
			}
			return m, nil, true
		case "k":
			// k: scroll details only (not comments)
			if !m.showComments {
				if m.detailScrollY > 0 {
					m.detailScrollY--
				}
			}
			return m, nil, true
		case "j":
			// j: scroll details only (not comments)
			if !m.showComments {
				if m.detailScrollY < m.detailMaxScroll {
					m.detailScrollY++
				}
			}
			return m, nil, true
		case "pgup":
			// Page up: scroll content
			if m.showComments {
				if m.commentsScrollY > 5 {
					m.commentsScrollY -= 5
				} else {
					m.commentsScrollY = 0
				}
			} else {
				if m.detailScrollY > 10 {
					m.detailScrollY -= 10
				} else {
					m.detailScrollY = 0
				}
			}
			return m, nil, true
		case "pgdn":
			// Page down: scroll content
			if m.showComments {
				if m.commentsScrollY+5 < m.commentsMaxScroll {
					m.commentsScrollY += 5
				} else {
					m.commentsScrollY = m.commentsMaxScroll
				}
			} else {
				if m.detailScrollY+10 < m.detailMaxScroll {
					m.detailScrollY += 10
				} else {
					m.detailScrollY = m.detailMaxScroll
				}
			}
			return m, nil, true
		case "home":
			if m.showComments {
				m.commentsScrollY = 0
			} else {
				m.detailScrollY = 0
			}
			return m, nil, true
		case "end":
			if m.showComments {
				m.commentsScrollY = m.commentsMaxScroll
			} else {
				m.detailScrollY = m.detailMaxScroll
			}
			return m, nil, true
		case "left", "h":
			// Left/h: previous post (works even when comments are open)
			if len(m.filteredPosts) > 0 {
				idx := m.list.Index()
				if idx > 0 {
					m.list.CursorUp()
					m.detailScrollY = 0
					m.commentsScrollY = 0
				}
			}
			return m, nil, true
		case "right", "l":
			// Right/l: next post (works even when comments are open)
			if len(m.filteredPosts) > 0 {
				idx := m.list.Index()
				if idx < len(m.filteredPosts)-1 {
					m.list.CursorDown()
					m.detailScrollY = 0
					m.commentsScrollY = 0
				}
			}
			return m, nil, true
		}
	}
	
	// Global shortcuts (also support left/right to navigate in list view)
	switch msg.String() {
	case "left", "h":
		if !m.showDetails && len(m.filteredPosts) > 0 {
			idx := m.list.Index()
			if idx > 0 {
				m.list.CursorUp()
			}
			return m, nil, true
		}
	case "right", "l":
		if !m.showDetails && len(m.filteredPosts) > 0 {
			idx := m.list.Index()
			if idx < len(m.filteredPosts)-1 {
				m.list.CursorDown()
			}
			return m, nil, true
		}
	}
	
	// Global shortcuts (remaining)
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "ctrl+f":
		m.searching = true
		m.searchInput.Focus()
		return m, nil, true
	case "ctrl+r":
		m.selectingSub = true
		m.subredditInput.Focus()
		m.subredditInput.SetValue(m.subreddit)
		return m, nil, true
	case "f5":
		m.loading = true
		m.showDetails = false
		return m, m.loadPosts(m.subreddit), true
	case "c":
		if m.showComments {
			// Close comments panel
			m.showComments = false
			m.commentsScrollY = 0
		} else if m.showDetails && len(m.filteredPosts) > 0 {
			// Open comments panel
			m.showComments = true
			m.commentsScrollY = 0
			m.commentsLoading = true
			post := m.filteredPosts[m.list.Index()]
			return m, m.loadComments(m.subreddit, post.ID), true
		}
		return m, nil, true
	case "w":
		// Open current post URL in browser
		if len(m.filteredPosts) > 0 {
			post := m.filteredPosts[m.list.Index()]
			// Use permalink if available, fallback to URL
			postURL := post.Permalink
			if postURL == "" {
				postURL = post.URL
			}
			if postURL != "" {
				// Prepend reddit.com domain if it's just a permalink
				if strings.HasPrefix(postURL, "/") {
					postURL = "https://reddit.com" + postURL
				}
				if err := openURL(postURL); err != nil {
					m.error = fmt.Sprintf("Failed to open URL: %v", err)
				}
			}
		}
		return m, nil, true
	case "enter":
		if len(m.filteredPosts) > 0 && m.list.Index() < len(m.filteredPosts) {
			m.showDetails = true
			m.detailScrollY = 0
			m.showComments = false
		}
		return m, nil, true
	}
	
	// Key not handled - let list component handle it
	return m, nil, false
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

// ============= Helpers =============

func (m Model) calculateCommentsMaxScroll(height int) Model {
	if len(m.comments) == 0 {
		m.commentsMaxScroll = 0
		return m
	}
	
	// Build comment lines to calculate total
	var commentLines []string
	for _, comment := range m.comments {
		commentLines = append(commentLines, "author line")
		if comment.Body != "" {
			wrapped := wrapText(comment.Body, m.windowWidth-6)
			wrappedLines := strings.Split(wrapped, "\n")
			commentLines = append(commentLines, wrappedLines...)
		}
		commentLines = append(commentLines, "") // Blank line between comments
	}
	
	m.commentsMaxScroll = max(0, len(commentLines)-height+4)
	return m
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

func (m *Model) renderMain() string {
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

func (m *Model) renderListOnly() string {
	m.updateListSize()
	return m.list.View()
}

func (m *Model) renderWithDetails() string {
	// Split view: list on top, details on bottom
	listHeight := (m.windowHeight - 8) / 2
	detailsHeight := m.windowHeight - 8 - listHeight - 1
	
	m.list.SetSize(m.windowWidth-2, listHeight)
	listView := m.list.View()
	
	// Details section or comments
	var contentView string
	if m.showComments {
		contentView = m.renderCommentsPanel(detailsHeight)
	} else {
		contentView = m.renderDetailsSection(detailsHeight)
	}
	
	separator := lipgloss.NewStyle().
		Foreground(colorOrange).
		Render(strings.Repeat("─", m.windowWidth-2))
	
	return fmt.Sprintf("%s\n%s\n%s", listView, separator, contentView)
}

func (m Model) renderCommentsPanel(height int) string {
	if len(m.comments) == 0 {
		if m.commentsLoading {
			return lipgloss.NewStyle().
				Foreground(colorGray).
				Padding(1, 1).
				Height(height).
				Render("💬 Loading comments...")
		}
		return lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(1, 1).
			Height(height).
			Render("💬 No comments found")
	}
	
	var sb strings.Builder
	sb.WriteString(focusedStyle.Render("💬 Comments\n\n"))
	
	// Build comment lines
	var commentLines []string
	for _, comment := range m.comments {
		// Author and score
		authorLine := lipgloss.NewStyle().Foreground(colorGold).Render(
			fmt.Sprintf("👤 u/%s  •  ⬆ %s", comment.Author, formatNum(comment.Score)))
		commentLines = append(commentLines, authorLine)
		
		// Comment body with wrapping
		if comment.Body != "" {
			wrapped := wrapText(comment.Body, m.windowWidth-6)
			wrappedLines := strings.Split(wrapped, "\n")
			for _, line := range wrappedLines {
				commentLines = append(commentLines, "  "+line)
			}
		}
		commentLines = append(commentLines, "") // Blank line between comments
	}
	
	// Apply scrolling
	startLine := m.commentsScrollY
	endLine := startLine + height - 4
	if endLine > len(commentLines) {
		endLine = len(commentLines)
	}
	
	if startLine < len(commentLines) {
		visibleLines := commentLines[startLine:endLine]
		sb.WriteString(strings.Join(visibleLines, "\n"))
	}
	
	return lipgloss.NewStyle().
		Foreground(colorGray).
		Padding(0, 1).
		Height(height).
		Render(sb.String())
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
	var contentLines []string
	if post.SelfText != "" {
		content := wrapText(post.SelfText, m.windowWidth-4)
		contentLines = strings.Split(content, "\n")
	}
	
	// Add URL if present
	if post.URL != "" && !strings.HasPrefix(post.URL, "https://www.reddit.com") {
		contentLines = append(contentLines, "")
		displayURL := post.URL
		if len(displayURL) > m.windowWidth-10 {
			displayURL = displayURL[:m.windowWidth-13] + "..."
		}
		contentLines = append(contentLines, "🔗 "+displayURL)
	}
	
	// Note: Comments section disabled until API endpoint is fixed
	// Comments fetching and display will be re-enabled in a future update
	
	// Apply scrolling
	startLine := m.detailScrollY
	endLine := startLine + height - 3
	if endLine > len(contentLines) {
		endLine = len(contentLines)
	}
	
	// Calculate max scroll
	m.detailMaxScroll = max(0, len(contentLines)-height+3)
	
	if startLine < len(contentLines) {
		visibleLines := contentLines[startLine:endLine]
		sb.WriteString(strings.Join(visibleLines, "\n"))
	}
	
	return lipgloss.NewStyle().
		Foreground(colorGray).
		Padding(0, 1).
		Height(height).
		Render(sb.String())
}

func (m Model) renderFooter() string {
	if m.showDetails {
		if m.showComments {
			return footerStyle.Render("↑↓: scroll comments  •  h/l: switch posts  •  w: open URL  •  Esc: close comments  •  Ctrl+F: search  •  q: quit")
		}
		return footerStyle.Render("↑↓: scroll details  •  h/l: switch posts  •  w: open URL  •  Esc/Tab: back to list  •  c: view comments  •  q: quit")
	}
	
	status := "no posts"
	if len(m.filteredPosts) > 0 {
		status = fmt.Sprintf("%d/%d", m.list.Index()+1, len(m.filteredPosts))
	}
	
	return footerStyle.Render(fmt.Sprintf("Post %s  •  Enter: view details  •  w: open URL  •  Ctrl+F: search  •  F5: refresh  •  q: quit", status))
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
	// Load configuration
	if err := loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config.json: %v\n", err)
	}
	
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
