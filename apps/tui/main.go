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
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Score    int    `json:"score"`
	Created  int64  `json:"created_utc"`
	Comments int    `json:"num_comments"`
	SelfText string `json:"selftext"`
	URL      string `json:"url"`
	SubName  string `json:"subreddit"`
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
		return nil, err
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
		return errorStyle.Render("Error: " + m.error)
	}

	if m.loading || len(m.posts) == 0 {
		return "Loading posts from r/" + m.sub + "..."
	}

	var s string
	s += titleStyle.Render("RedditView TUI") + "\n\n"
	s += subtitleStyle.Render("r/"+m.sub) + "\n"
	s += dividerStyle.Render("─────────────────────────────────────────────────────────────") + "\n\n"

	for i, post := range m.posts {
		selected := " "
		if i == m.selected {
			selected = "▶"
		}

		title := post.Title
		if len(title) > 55 {
			title = title[:52] + "..."
		}

		score := fmt.Sprintf("↑ %d", post.Score)
		comments := fmt.Sprintf("💬 %d", post.Comments)

		var style lipgloss.Style
		if i == m.selected {
			style = selectedStyle
		} else {
			style = normalStyle
		}

		s += style.Render(fmt.Sprintf("%s %s\n", selected, title))
		s += dimStyle.Render(fmt.Sprintf("  u/%s · %s · %s\n\n", post.Author, score, comments))
	}

	s += dividerStyle.Render("─────────────────────────────────────────────────────────────") + "\n"
	s += helpStyle.Render("↑↓/jk: navigate · q: quit")

	return s
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("238")).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle()

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Margin(1, 0, 0, 0)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
