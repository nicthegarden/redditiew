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
)

// Models matching @redditview/core

type RedditPostData struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Author   string      `json:"author"`
	Score    int         `json:"score"`
	Created  float64     `json:"created_utc"`
	Comments int         `json:"num_comments"`
	SelfText string      `json:"selftext"`
	URL      string      `json:"url"`
	SubName  string      `json:"subreddit"`
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

// API Client

const apiBaseURL = "http://localhost:3002/api"

type APIClient struct {
	baseURL string
}

func NewAPIClient() *APIClient {
	return &APIClient{baseURL: apiBaseURL}
}

func (c *APIClient) FetchPosts(subreddit string) ([]RedditPostData, error) {
	resp, err := http.Get(fmt.Sprintf("%s/r/%s.json?limit=30", c.baseURL, subreddit))
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
		if post.Kind == "t3" { // Post kind
			posts = append(posts, post.Data)
		}
	}

	return posts, nil
}

// TUI Model

type Model struct {
	posts        []RedditPostData
	filteredPosts []RedditPostData
	selected     int
	loading      bool
	error        string
	sub          string
	client       *APIClient
	
	// New fields for split view
	searchQuery  string
	inputMode    bool
	selectedPost *RedditPostData
	windowWidth  int
	windowHeight int
}

func initialModel() Model {
	return Model{
		client:        NewAPIClient(),
		selected:      0,
		sub:           "golang",
		loading:       true,
		inputMode:     false,
		selectedPost:  nil,
		windowWidth:   200,
		windowHeight:  50,
	}
}

type postsLoadedMsg struct {
	posts []RedditPostData
}

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

type windowSizeMsg struct {
	width  int
	height int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadPosts("golang"),
		tea.EnterAltScreen,
	)
}

func (m Model) loadPosts(subreddit string) tea.Cmd {
	return func() tea.Msg {
		posts, err := m.client.FetchPosts(subreddit)
		if err != nil {
			return errMsg{err}
		}
		return postsLoadedMsg{posts}
	}
}

func (m Model) filterPosts() {
	if m.searchQuery == "" {
		m.filteredPosts = m.posts
		return
	}

	filtered := []RedditPostData{}
	query := strings.ToLower(m.searchQuery)
	for _, post := range m.posts {
		if strings.Contains(strings.ToLower(post.Title), query) ||
			strings.Contains(strings.ToLower(post.Author), query) {
			filtered = append(filtered, post)
		}
	}
	m.filteredPosts = filtered
	m.selected = 0
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input mode
		if m.inputMode {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.inputMode = false
				return m, nil
			case "enter":
				m.inputMode = false
				m.filterPosts()
				return m, nil
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.filterPosts()
				}
				return m, nil
			default:
				if len(msg.String()) == 1 && msg.String() >= " " && msg.String() <= "~" {
					m.searchQuery += msg.String()
					m.filterPosts()
				}
				return m, nil
			}
		}

		// Handle navigation mode
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
				if len(m.filteredPosts) > m.selected {
					m.selectedPost = &m.filteredPosts[m.selected]
				}
			}
		case "down", "j":
			if m.selected < len(m.filteredPosts)-1 {
				m.selected++
				if len(m.filteredPosts) > m.selected {
					m.selectedPost = &m.filteredPosts[m.selected]
				}
			}
		case "enter":
			// View detailed post
			if len(m.filteredPosts) > m.selected {
				m.selectedPost = &m.filteredPosts[m.selected]
			}
		case "/":
			// Enter search mode
			m.inputMode = true
			m.searchQuery = ""
			m.filterPosts()
		}

	case postsLoadedMsg:
		m.posts = msg.posts
		m.filteredPosts = msg.posts
		m.loading = false
		if len(msg.posts) > 0 {
			m.selectedPost = &msg.posts[0]
		}

	case errMsg:
		m.error = msg.Error()

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	if m.error != "" {
		return errorStyle.Render("❌ Error: " + m.error)
	}

	if m.loading {
		return loadingStyle.Render("\n  ⏳ Loading posts from r/" + m.sub + "...\n")
	}

	return m.renderSplitView()
}

func (m Model) renderSplitView() string {
	// Calculate dimensions
	totalWidth := m.windowWidth
	totalHeight := m.windowHeight
	leftWidth := (totalWidth / 2) - 1
	rightWidth := totalWidth / 2

	// Header
	header := headerBg.Render("  🔥 r/" + m.sub + "  ")
	
	// Search bar
	searchBar := m.renderSearchBar()

	// Left sidebar - Post list
	leftContent := m.renderPostList(leftWidth, totalHeight-6)
	
	// Right sidebar - Post details
	rightContent := m.renderPostDetails(rightWidth, totalHeight-6)

	// Combine
	var result strings.Builder
	result.WriteString(header + "\n")
	result.WriteString(searchBar + "\n")
	result.WriteString(dividerStyle.Render(strings.Repeat("─", totalWidth)) + "\n")

	// Render both sidebars side by side
	leftLines := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for i := 0; i < maxLines; i++ {
		left := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		left = padRight(left, leftWidth)

		right := ""
		if i < len(rightLines) {
			right = rightLines[i]
		}

		result.WriteString(left + "│" + right + "\n")
	}

	result.WriteString(dividerStyle.Render(strings.Repeat("─", totalWidth)) + "\n")
	result.WriteString(m.renderFooter())

	return result.String()
}

func (m Model) renderSearchBar() string {
	if m.inputMode {
		cursor := "▌"
		return searchInputStyle.Render(fmt.Sprintf("  🔍 Search: %s%s", m.searchQuery, cursor))
	}
	return normalStatsStyle.Render(fmt.Sprintf("  🔍 Press '/' to search · Found %d posts", len(m.filteredPosts)))
}

func (m Model) renderPostList(width int, height int) string {
	var result strings.Builder

	for i, post := range m.filteredPosts {
		if i >= height {
			break
		}

		title := truncateTitle(post.Title, width-3)

		if i == m.selected {
			result.WriteString(selectedTitleStyle.Render("❯ " + title + "\n"))
			result.WriteString(selectedAuthorStyle.Render(padRight(fmt.Sprintf("  u/%s", post.Author), width) + "\n"))
			result.WriteString(selectedStatsStyle.Render(padRight(fmt.Sprintf("  ⬆ %d  💬 %d", post.Score, post.Comments), width) + "\n"))
		} else {
			result.WriteString(normalTitleStyle.Render("  " + title + "\n"))
			result.WriteString(normalAuthorStyle.Render(padRight(fmt.Sprintf("  u/%s", post.Author), width) + "\n"))
			result.WriteString(normalStatsStyle.Render(padRight(fmt.Sprintf("  ⬆ %d  💬 %d", post.Score, post.Comments), width) + "\n"))
		}
		result.WriteString("\n")
	}

	return result.String()
}

func (m Model) renderPostDetails(width int, height int) string {
	if m.selectedPost == nil {
		return detailsEmptyStyle.Render(padRight("  Select a post to view", width))
	}

	var result strings.Builder

	post := m.selectedPost

	// Title
	titleLines := wrapText(post.Title, width-4)
	result.WriteString(detailsTitleStyle.Render(padRight("  "+titleLines[0], width) + "\n"))
	for i := 1; i < len(titleLines); i++ {
		result.WriteString(detailsTitleStyle.Render(padRight("  "+titleLines[i], width) + "\n"))
	}
	result.WriteString("\n")

	// Author and stats
	result.WriteString(detailsAuthorStyle.Render(padRight("  👤 u/"+post.Author, width) + "\n"))
	result.WriteString(detailsStatsStyle.Render(padRight(fmt.Sprintf("  ⬆ %d upvotes  💬 %d comments", post.Score, post.Comments), width) + "\n"))
	result.WriteString("\n")

	// Content
	result.WriteString(detailsLabelStyle.Render(padRight("  Content:", width) + "\n"))
	contentLines := wrapText(post.SelfText, width-4)
	for _, line := range contentLines {
		if len(line) == 0 {
			result.WriteString("\n")
		} else {
			result.WriteString(detailsContentStyle.Render(padRight("  "+line, width) + "\n"))
		}
	}

	// URL
	if post.URL != "" && !strings.HasPrefix(post.URL, "https://www.reddit.com") {
		result.WriteString("\n")
		result.WriteString(detailsLabelStyle.Render(padRight("  Link:", width) + "\n"))
		result.WriteString(detailsLinkStyle.Render(padRight("  "+truncateTitle(post.URL, width-4), width) + "\n"))
	}

	return result.String()
}

func (m Model) renderFooter() string {
	var controls string
	if m.inputMode {
		controls = "  ESC: cancel search · ENTER: apply"
	} else {
		controls = "  ↑↓/jk: navigate · /: search · ENTER: view · q: quit"
	}
	return footerStyle.Render(controls)
}

// Utility functions

func truncateTitle(title string, maxLen int) string {
	if len(title) > maxLen {
		return title[:maxLen-3] + "..."
	}
	return title
}

func wrapText(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
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

	return lines
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

// Styles

var (
	// Header
	headerBg = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF4500")).
		Padding(0, 1).
		MarginBottom(0)

	// Search
	searchInputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Background(lipgloss.Color("#1a1a1a"))

	// Post list - left sidebar
	selectedTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF6B35"))

	normalTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	selectedAuthorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Background(lipgloss.Color("#FF6B35"))

	normalAuthorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700"))

	selectedStatsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90EE90")).
		Background(lipgloss.Color("#FF6B35"))

	normalStatsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90EE90"))

	// Post details - right sidebar
	detailsTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B35"))

	detailsAuthorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700"))

	detailsStatsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90EE90"))

	detailsLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF4500"))

	detailsContentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC"))

	detailsLinkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB")).
		Underline(true)

	detailsEmptyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	// Dividers and utilities
	dividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4500"))

	// Footer
	footerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333"))

	// Error and loading
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF0000")).
		Padding(1, 2)

	loadingStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(1, 2)
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
