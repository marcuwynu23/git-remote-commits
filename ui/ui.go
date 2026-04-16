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
	ShowDiff      bool
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
	diffTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FDE68A"))
)

func Render(v ViewData) string {
	mainWidth := max(v.Width-2, 70)
	mainHeight := max(v.Height-2, 18)

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
		"",
		helpStyle.Render("up/down or j/k: select • d: toggle diff • r: refresh • q: quit"),
	)

	mainPanel := frameStyle.Width(mainWidth).Height(mainHeight).Render(content)

	if !v.ShowDiff {
		return mainPanel
	}

	leftW := max((v.Width/2)-2, 40)
	rightW := max(v.Width-leftW-6, 30)
	left := frameStyle.Width(leftW).Height(mainHeight).Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		metaStyle.Render(statusLine),
		"",
		renderCommitList(v, leftW),
		"",
		footer,
	))
	diffBody := clampLines(v.Snapshot.SelectedDiff, max(mainHeight-4, 8))
	right := frameStyle.Width(rightW).Height(mainHeight).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			diffTitle.Render("Diff Preview"),
			metaStyle.Render("git show "+emptyFallback(v.Snapshot.SelectedHash, "<none>")),
			"",
			diffBody,
		),
	)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func renderCommitList(v ViewData, width int) string {
	if len(v.Snapshot.Commits) == 0 {
		return metaStyle.Render("No commits found.")
	}
	lines := make([]string, 0, len(v.Snapshot.Commits))
	for i, c := range v.Snapshot.Commits {
		prefix := "  "
		row := fmt.Sprintf("● %s %s (%s, %s)", c.Hash, c.Message, c.Author, c.Time)
		if _, ok := v.NewCommitHash[c.Hash]; ok {
			row = "NEW " + row
		}
		row = trimToWidth(row, max(width-8, 24))

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

func clampLines(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	return strings.Join(lines[:maxLines], "\n") + "\n... (truncated)"
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
