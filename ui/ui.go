package ui

import (
	"fmt"
	"strings"
	"time"

	"tuiapp/git"

	"github.com/charmbracelet/lipgloss"
)

type ViewData struct {
	Width         int
	Height        int
	Selected      int
	NewCommitHash map[string]struct{}
	Snapshot      git.Snapshot
}

var (
	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A78BFA"))
	metaStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true)
	selected   = lipgloss.NewStyle().Foreground(lipgloss.Color("#111827")).Background(lipgloss.Color("#C4B5FD"))
	freshStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ADE80"))
	oldStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
)

func Render(v ViewData) string {
	outerWidth := max(v.Width, 1)
	outerHeight := max(v.Height-1, 1)
	frameW := frameStyle.GetHorizontalFrameSize()
	frameH := frameStyle.GetVerticalFrameSize()
	mainWidth := max(outerWidth-frameW, 1)
	mainHeight := max(outerHeight-frameH, 1)

	header := titleStyle.Render("Git Live Monitor")
	statusLine := fmt.Sprintf("Branch: %s | Status: %s", emptyFallback(v.Snapshot.Branch, "-"), emptyFallback(v.Snapshot.Status, "-"))

	body := renderCommitList(v, mainWidth)
	remote := fmt.Sprintf("Remote: %s\nStatus: %s", v.Snapshot.RemoteTrackName, emptyFallback(v.Snapshot.RemoteStatus, "checking remote..."))
	refresh := fmt.Sprintf("Last refresh: %s", v.Snapshot.LastRefresh.Format(time.Kitchen))
	footer := metaStyle.Render(remote + "\n" + refresh)

	if v.Snapshot.RepoError != "" {
		body = errStyle.Render("Error: " + v.Snapshot.RepoError + "\n\nOpen this app from a valid git repository.")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		metaStyle.Render(statusLine),
		"",
		body,
		"",
		footer,
	)

	panel := frameStyle.Width(mainWidth).Height(mainHeight).Render(content)
	help := helpStyle.Width(outerWidth).Render("up/down or j/k: select • r: refresh • q: quit")
	return lipgloss.JoinVertical(lipgloss.Left, panel, help)
}

func renderCommitList(v ViewData, width int) string {
	if len(v.Snapshot.Commits) == 0 {
		return metaStyle.Render("No commits found.")
	}
	lines := make([]string, 0, len(v.Snapshot.Commits))
	for i, c := range v.Snapshot.Commits {
		prefix := "  "
		ts := c.When.Local().Format("January 2, 2006: 03:04 PM")
		if c.When.IsZero() {
			ts = c.Time
		}
		row := fmt.Sprintf("● %s | %s | %s", ts, c.Author, c.Message)
		if _, ok := v.NewCommitHash[c.Hash]; ok {
			row = "NEW " + row
		}
		row = trimToWidth(row, max(width-8, 8))

		if i == v.Selected {
			lines = append(lines, selected.Render("> "+row))
			continue
		}
		if _, ok := v.NewCommitHash[c.Hash]; ok {
			lines = append(lines, freshStyle.Render(prefix+row))
		} else {
			lines = append(lines, oldStyle.Render(prefix+row))
		}
	}
	return strings.Join(lines, "\n")
}

func trimToWidth(s string, w int) string {
	rs := []rune(s)
	if len(rs) <= w {
		return s
	}
	if w < 4 {
		return string(rs[:w])
	}
	return string(rs[:w-3]) + "..."
}

func emptyFallback(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
