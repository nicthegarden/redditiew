package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	resp, err := http.Get(fmt.Sprintf("%s/r/%s.json?limit=20", c.baseURL, subreddit))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var result RedditResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// Log the error with context
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
	posts    []RedditPostData
	selected int
	loading  bool
	error    string
	sub      string
	client   *APIClient
}

func initialModel() Model {
	return Model{
		client:   NewAPIClient(),
		selected: 0,
		sub:      "golang",
		loading:  true,  // Set to true initially so we show "Loading..." message
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

func (m Model) Init() tea.Cmd {
	return m.loadPosts("golang")
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.posts)-1 {
				m.selected++
			}
		}
	case postsLoadedMsg:
		m.posts = msg.posts
		m.loading = false
	case errMsg:
		m.error = msg.Error()
	}
	return m, nil
}

func (m Model) View() string {
	if m.error != "" {
		return errorStyle.Render("❌ Error: " + m.error)
	}

	if m.loading || len(m.posts) == 0 {
		return loadingStyle.Render("\n  ⏳ Loading posts from r/" + m.sub + "...\n")
	}

	var s string
	
	// Header
	s += headerBg.Render("  🔥 RedditView TUI  ") + "\n"
	s += subredditStyle.Render("  r/" + m.sub) + "\n"
	s += dividerStyle.Render(getDivider(80)) + "\n\n"

	// Posts
	for i, post := range m.posts {
		var postStr string
		
		if i == m.selected {
			// Selected post - highlighted
			postStr += selectedIndicator + " "
			postStr += selectedTitleStyle.Render(truncateTitle(post.Title, 70))
			postStr += "\n"
			postStr += selectedAuthorStyle.Render("     👤 u/" + post.Author)
			postStr += selectedStatsStyle.Render(fmt.Sprintf("  ⬆ %d  💬 %d", post.Score, post.Comments))
			postStr += "\n"
			postStr += selectedBgStyle.Render("  " + postStr)
			s += selectedBgStyle.Render(postStr + "\n")
		} else {
			// Normal post
			postStr += "  " + normalIndicator + " "
			postStr += normalTitleStyle.Render(truncateTitle(post.Title, 70))
			postStr += "\n"
			postStr += normalAuthorStyle.Render("     👤 u/" + post.Author)
			postStr += normalStatsStyle.Render(fmt.Sprintf("  ⬆ %d  💬 %d", post.Score, post.Comments))
			postStr += "\n"
			s += postStr + "\n"
		}
	}

	// Footer
	s += dividerStyle.Render(getDivider(80)) + "\n"
	s += footerStyle.Render("  ⬆⬇ navigate · j/k vim · q quit · ? help  ")

	return s
}

func truncateTitle(title string, maxLen int) string {
	if len(title) > maxLen {
		return title[:maxLen-3] + "..."
	}
	return title
}

func getDivider(length int) string {
	divider := ""
	for i := 0; i < length; i++ {
		divider += "─"
	}
	return divider
}

var (
	// Constants
	selectedIndicator = "❯"
	normalIndicator   = "┃"

	// Header styles
	headerBg = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF4500")).
		Padding(0, 1).
		MarginBottom(1)

	subredditStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B35")).
		MarginBottom(1)

	// Title styles
	selectedTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF6B35"))

	normalTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	// Author and stats styles
	selectedAuthorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Background(lipgloss.Color("#FF6B35"))

	selectedStatsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90EE90")).
		Background(lipgloss.Color("#FF6B35"))

	normalAuthorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700"))

	normalStatsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90EE90"))

	// Background for selected post
	selectedBgStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#FF6B35")).
		Padding(0, 0)

	// Divider and utilities
	dividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4500"))

	// Footer
	footerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		MarginTop(1)

	// Error style
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF0000")).
		Padding(1, 2)

	// Loading style
	loadingStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(1, 2)
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
