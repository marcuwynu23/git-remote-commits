package ui

import (
	"fmt"
	"strings"
	"time"

	"git-remote-commits/git"

	"github.com/charmbracelet/lipgloss"
)

type ViewData struct {
	Width         int
	Height        int
	RepoName      string
	Selected      int
	Loaded        bool
	NewCommitHash map[string]struct{}
	Snapshot      git.Snapshot
}

var (
	frameStyle = lipgloss.NewStyle().
			Padding(0, 1)
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#22C55E"))
	metaStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true)
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	chipStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Background(lipgloss.Color("#1F2937")).
			Padding(0, 1)
	chipGoodStyle = chipStyle.Copy().
			Foreground(lipgloss.Color("#052E16")).
			Background(lipgloss.Color("#86EFAC"))
	chipWarnStyle = chipStyle.Copy().
			Foreground(lipgloss.Color("#3F2A00")).
			Background(lipgloss.Color("#FCD34D"))
	chipInfoStyle = chipStyle.Copy().
			Foreground(lipgloss.Color("#082F49")).
			Background(lipgloss.Color("#93C5FD"))
	selected   = lipgloss.NewStyle().Background(lipgloss.Color("#312E81"))
	freshStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ADE80"))
	hashStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FACC15"))
	refStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#93C5FD"))
	authorMe   = lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Bold(true)
	authorElse = lipgloss.NewStyle().Foreground(lipgloss.Color("#86EFAC"))
	msgStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
)

func Render(v ViewData) string {
	outerWidth := max(v.Width, 1)
	outerHeight := max(v.Height-1, 1)
	frameW := frameStyle.GetHorizontalFrameSize()
	frameH := frameStyle.GetVerticalFrameSize()
	mainWidth := max(outerWidth-frameW, 1)
	mainHeight := max(outerHeight-frameH, 1)

	header := titleStyle.Render("git remote-commits")
	statusLine := renderHeaderLine(v)
	footer := metaStyle.Render("Loading repository data...")
	if v.Loaded {
		footer = renderFooterLine(v)
	}

	// Keep the panel height stable by giving commits a fixed viewport,
	// so header/footer positions don't jump when commit count is small.
	commitHeight := max(mainHeight-5, 1)
	body := renderCommitList(v, mainWidth, commitHeight)
	if !v.Loaded {
		body = renderLoading(mainWidth, commitHeight)
	}
	body = lipgloss.NewStyle().Height(commitHeight).Render(body)

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

func renderCommitList(v ViewData, width int, height int) string {
	if len(v.Snapshot.Commits) == 0 {
		return metaStyle.Render("No commits found.")
	}

	total := len(v.Snapshot.Commits)
	height = max(height, 1)

	start := 0
	if total > height && v.Selected >= height {
		start = v.Selected - height + 1
	}
	if start+height > total {
		start = max(total-height, 0)
	}
	end := min(start+height, total)

	lines := make([]string, 0, end-start+2)
	currentAuthor := strings.TrimSpace(v.Snapshot.CurrentAuthor)
	for i := start; i < end; i++ {
		c := v.Snapshot.Commits[i]
		prefix := "  "
		ts := c.When.Local().Format("January 2, 2006: 03:04 PM")
		if c.When.IsZero() {
			ts = c.Time
		}
		rowPrefix := fmt.Sprintf("● %s | ", ts)
		authorLabel := formatAuthorLabel(c.Author, c.AuthorEmail)
		refsLabel := commitRefsLabel(c.Refs)
		isNew := false
		if _, ok := v.NewCommitHash[c.Hash]; ok {
			isNew = true
		}

		plainPrefix := prefix + rowPrefix + c.Hash + " | "
		if refsLabel != "" {
			plainPrefix += refsLabel + " | "
		}
		plainPrefix += authorLabel + " | "
		if isNew {
			plainPrefix = "NEW " + plainPrefix
		}
		msgWidth := max(width-8-len([]rune(plainPrefix)), 0)
		msg := trimToWidth(c.Message, msgWidth)

		authorText := authorElse.Render(authorLabel)
		if currentAuthor != "" && strings.EqualFold(strings.TrimSpace(c.Author), currentAuthor) {
			authorText = authorMe.Render(authorLabel)
		}

		row := rowPrefix + hashStyle.Render(c.Hash) + " | "
		if refsLabel != "" {
			row += refStyle.Render(refsLabel) + " | "
		}
		row += authorText + " | " + msgStyle.Render(msg)
		if isNew {
			row = freshStyle.Render("NEW ") + row
		}

		if i == v.Selected {
			lines = append(lines, selected.Render("> "+row))
			continue
		}
		lines = append(lines, prefix+row)
	}
	return strings.Join(lines, "\n")
}

func renderLoading(width int, height int) string {
	barWidth := min(max(width/3, 18), 36)
	progressSteps := barWidth
	pos := int(time.Now().UnixNano()/1e8) % progressSteps
	filled := strings.Repeat("█", pos+1)
	empty := strings.Repeat("░", barWidth-pos-1)
	bar := "[" + filled + empty + "]"
	loading := lipgloss.JoinVertical(
		lipgloss.Center,
		metaStyle.Render("Loading commits..."),
		freshStyle.Render(bar),
	)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, loading)
}

func renderHeaderLine(v ViewData) string {
	repo := chipStyle.Render(emptyFallback(v.RepoName, "-"))
	branch := chipInfoStyle.Render(emptyFallback(v.Snapshot.Branch, "-"))
	statusText := emptyFallback(v.Snapshot.Status, "-")
	statusChip := chipStyle.Render(statusText)
	if strings.EqualFold(statusText, "clean") {
		statusChip = chipGoodStyle.Render(statusText)
	} else if strings.EqualFold(statusText, "dirty") {
		statusChip = chipWarnStyle.Render(statusText)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("Repository "),
		repo,
		"  ",
		labelStyle.Render("Branch "),
		branch,
		"  ",
		labelStyle.Render("Status "),
		statusChip,
	)
}

func renderFooterLine(v ViewData) string {
	remote := chipInfoStyle.Render(emptyFallback(v.Snapshot.RemoteTrackName, "none"))
	syncText := emptyFallback(v.Snapshot.RemoteStatus, "checking remote...")
	syncChip := chipStyle.Render(syncText)
	if strings.EqualFold(syncText, "up to date") {
		syncChip = chipGoodStyle.Render(syncText)
	} else if strings.Contains(strings.ToLower(syncText), "behind") || strings.Contains(strings.ToLower(syncText), "diverged") {
		syncChip = chipWarnStyle.Render(syncText)
	}
	refresh := chipStyle.Render(v.Snapshot.LastRefresh.Format(time.Kitchen))
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("Remote "),
		remote,
		"  ",
		labelStyle.Render("Sync "),
		syncChip,
		"  ",
		labelStyle.Render("Refresh "),
		refresh,
	)
}

func commitRefsLabel(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	if trimmed == "" {
		return ""
	}

	parts := strings.Split(trimmed, ",")
	labels := make([]string, 0, len(parts))
	for _, p := range parts {
		item := strings.TrimSpace(p)
		if strings.Contains(item, "HEAD") || strings.HasPrefix(item, "tag: ") {
			labels = append(labels, item)
		}
	}
	return strings.Join(labels, " ")
}

func formatAuthorLabel(author, email string) string {
	author = strings.TrimSpace(author)
	email = strings.TrimSpace(email)
	if author == "" || email == "" {
		return author
	}
	at := strings.Index(email, "@")
	if at <= 0 {
		return author
	}
	username := strings.TrimSpace(email[:at])
	if username == "" || strings.EqualFold(author, username) {
		return author
	}
	return fmt.Sprintf("%s (%s)", author, username)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
