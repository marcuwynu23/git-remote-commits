package ui

import (
	"fmt"
	"strings"
	"time"

	"git-remote-commits/git"

	"github.com/charmbracelet/lipgloss"
)

type ViewData struct {
	Width           int
	Height          int
	RepoName        string
	Version         string
	Selected        int
	Loaded          bool
	Refreshing      bool
	ShowCommitPanel bool
	PanelScroll     int
	NewCommitHash   map[string]struct{}
	Snapshot        git.Snapshot
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

	header := ""
	statusLine := ""
	footer := ""
	if v.Loaded {
		header = titleStyle.Render(fmt.Sprintf("git remote-commits %s", emptyFallback(v.Version, "dev")))
		statusLine = renderHeaderLine(v)
		footer = renderFooterLine(v)
	}

	staticLines := 0
	if v.Loaded {
		// header + status + spacer before body + spacer before footer + footer
		staticLines = 5
	}
	availableBody := max(mainHeight-staticLines, 1)
	topHeight := availableBody
	panelHeight := 0
	showPanel := v.ShowCommitPanel && availableBody > 1
	if showPanel {
		panelHeight = max(availableBody/4, 3)
		if panelHeight > availableBody-1 {
			panelHeight = availableBody - 1
		}
		topHeight = max(availableBody-panelHeight, 1)
	}

	body := renderCommitList(v, mainWidth, topHeight)
	if !v.Loaded {
		body = renderLoading(mainWidth, topHeight)
	}
	body = lipgloss.NewStyle().Height(topHeight).Render(body)

	if v.Snapshot.RepoError != "" {
		body = errStyle.Render("Error: " + v.Snapshot.RepoError + "\n\nOpen this app from a valid git repository.")
	}

	commitPanel := ""
	if showPanel {
		commitPanel = renderCommitPanel(v, mainWidth, panelHeight)
	}

	sections := make([]string, 0, 8)
	if v.Loaded {
		sections = append(sections,
			header,
			metaStyle.Render(statusLine),
			"",
		)
	}
	sections = append(sections, body)
	if showPanel && panelHeight > 0 {
		sections = append(sections, commitPanel)
	}
	if v.Loaded {
		sections = append(sections, "", footer)
	}
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	panel := frameStyle.Width(mainWidth).Height(mainHeight).Render(content)
	help := helpStyle.Width(outerWidth).Render("up/down or j/k: select • [/], u/d, pgup/pgdown: panel scroll • p: toggle commit panel • r: refresh • q: quit")
	return lipgloss.JoinVertical(lipgloss.Left, panel, help)
}

func renderCommitPanel(v ViewData, width int, height int) string {
	height = max(height, 1)
	panelWidth := max(width-2, 10)
	if len(v.Snapshot.Commits) == 0 {
		empty := metaStyle.Render("No commit selected.")
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderLeft(false).
			BorderRight(false).
			BorderBottom(false).
			BorderForeground(lipgloss.Color("#374151")).
			Padding(0, 1).
			Width(panelWidth).
			Height(height).
			Render(empty)
	}

	selectedIndex := v.Selected
	if selectedIndex < 0 {
		selectedIndex = 0
	}
	if selectedIndex >= len(v.Snapshot.Commits) {
		selectedIndex = len(v.Snapshot.Commits) - 1
	}
	c := v.Snapshot.Commits[selectedIndex]
	when := c.When.Local().Format("January 2, 2006: 03:04 PM")
	if c.When.IsZero() {
		when = c.Time
	}
	viewportWidth := max(panelWidth-2, 8)

	titleText := emptyFallback(c.Message, "-")
	titleLines := wrapLines(titleText, viewportWidth)

	bodyRaw := strings.TrimSpace(c.Body)
	hasBody := bodyRaw != ""
	var bodyLines []string
	if hasBody {
		bodyLines = wrapLines(bodyRaw, viewportWidth)
	}

	lines := []string{
		labelStyle.Render("Hash  ") + hashStyle.Render(emptyFallback(c.Hash, "-")),
		labelStyle.Render("Author") + " " + msgStyle.Render(emptyFallback(formatAuthorLabel(c.Author, c.AuthorEmail), "-")),
		labelStyle.Render("When  ") + " " + metaStyle.Render(emptyFallback(when, "-")),
		labelStyle.Render("Refs  ") + " " + refStyle.Render(emptyFallback(commitRefsLabel(c.Refs), "-")),
		"",
		labelStyle.Render("Title"),
	}
	for _, tl := range titleLines {
		lines = append(lines, msgStyle.Render(tl))
	}
	if hasBody {
		lines = append(lines, "", labelStyle.Render("Body"))
		for _, bl := range bodyLines {
			lines = append(lines, msgStyle.Render(bl))
		}
	}

	maxScroll := max(len(lines)-height, 0)
	scroll := v.PanelScroll
	if scroll < 0 {
		scroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}
	hasUp := scroll > 0
	hasDown := scroll < maxScroll

	if len(lines) > 0 {
		indicator := ""
		if hasUp {
			indicator += "↑"
		}
		if hasDown {
			if indicator != "" {
				indicator += " "
			}
			indicator += "↓"
		}
		if indicator != "" {
			lines[0] = lines[0] + "  " + metaStyle.Render(indicator+" u/d")
		}
	}

	end := min(scroll+height, len(lines))
	visible := lines[scroll:end]
	body := lipgloss.NewStyle().Height(height).Render(strings.Join(visible, "\n"))
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderLeft(false).
		BorderRight(false).
		BorderBottom(false).
		BorderForeground(lipgloss.Color("#374151")).
		Padding(0, 1).
		Width(panelWidth).
		Height(height).
		Render(body)
}

func wrapLines(text string, width int) []string {
	if width <= 0 {
		return []string{""}
	}
	srcLines := strings.Split(text, "\n")
	out := make([]string, 0, len(srcLines))
	for _, src := range srcLines {
		runes := []rune(src)
		if len(runes) == 0 {
			out = append(out, "")
			continue
		}
		for len(runes) > width {
			out = append(out, string(runes[:width]))
			runes = runes[width:]
		}
		out = append(out, string(runes))
	}
	return out
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

	type rowData struct {
		commit      git.Commit
		timeLabel   string
		authorLabel string
		refsLabel   string
		isNew       bool
	}
	rows := make([]rowData, 0, end-start)
	maxTimeW := 0
	maxRefsW := 0
	maxAuthorW := 0
	for i := start; i < end; i++ {
		c := v.Snapshot.Commits[i]
		ts := c.When.Local().Format("January 2, 2006: 03:04 PM")
		if c.When.IsZero() {
			ts = c.Time
		}
		authorLabel := formatAuthorLabel(c.Author, c.AuthorEmail)
		refsLabel := emptyFallback(commitRefsLabel(c.Refs), "-")
		_, isNew := v.NewCommitHash[c.Hash]

		rows = append(rows, rowData{
			commit:      c,
			timeLabel:   ts,
			authorLabel: authorLabel,
			refsLabel:   refsLabel,
			isNew:       isNew,
		})
		maxTimeW = max(maxTimeW, len([]rune(ts)))
		maxRefsW = max(maxRefsW, len([]rune(refsLabel)))
		maxAuthorW = max(maxAuthorW, len([]rune(authorLabel)))
	}

	dateW, hashW, refsW, authorW, msgW := commitColumnWidths(width, maxTimeW, maxRefsW, maxAuthorW)
	lines := make([]string, 0, len(rows)+2)
	currentAuthor := strings.TrimSpace(v.Snapshot.CurrentAuthor)
	for idx, r := range rows {
		i := start + idx
		dateText := padRight(r.timeLabel, dateW)
		hashText := padRight(trimToWidth(r.commit.Hash, hashW), hashW)
		refsText := padRight(r.refsLabel, refsW)
		authorTextRaw := padRight(r.authorLabel, authorW)
		msgText := trimToWidth(r.commit.Message, msgW)

		authorText := authorElse.Render(authorTextRaw)
		if currentAuthor != "" && strings.EqualFold(strings.TrimSpace(r.commit.Author), currentAuthor) {
			authorText = authorMe.Render(authorTextRaw)
		}

		marker := "  "
		if r.isNew {
			marker = freshStyle.Render("● ")
		}

		row := marker +
			metaStyle.Render(dateText) + "  " +
			hashStyle.Render(hashText) + "  " +
			refStyle.Render(refsText) + "  " +
			authorText + "  " +
			msgStyle.Render(msgText)

		if i == v.Selected {
			lines = append(lines, selected.Render("> "+row))
			continue
		}
		lines = append(lines, row)
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
	if v.Refreshing {
		syncText = "refreshing..."
	}
	syncChip := chipStyle.Render(syncText)
	if v.Refreshing {
		syncChip = chipInfoStyle.Render(syncText)
	} else if strings.EqualFold(syncText, "up to date") {
		syncChip = chipGoodStyle.Render(syncText)
	} else if strings.Contains(strings.ToLower(syncText), "behind") || strings.Contains(strings.ToLower(syncText), "diverged") {
		syncChip = chipWarnStyle.Render(syncText)
	}
	refresh := chipStyle.Render(v.Snapshot.LastRefresh.Format(time.Kitchen))
	parts := []string{
		labelStyle.Render("Remote "),
		remote,
		"  ",
		labelStyle.Render("Sync "),
		syncChip,
		"  ",
		labelStyle.Render("Refresh "),
		refresh,
	}
	if v.Refreshing {
		parts = append(parts,
			"  ",
			labelStyle.Render("Loading "),
			freshStyle.Render(renderMiniLoadingBar(12)),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		parts...,
	)
}

func renderMiniLoadingBar(width int) string {
	if width <= 0 {
		return ""
	}
	pos := int(time.Now().UnixNano()/1e8) % width
	filled := strings.Repeat("█", pos+1)
	empty := strings.Repeat("░", width-pos-1)
	return "[" + filled + empty + "]"
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
	if w <= 0 {
		return ""
	}
	rs := []rune(s)
	if len(rs) <= w {
		return s
	}
	if w < 4 {
		return string(rs[:w])
	}
	return string(rs[:w-3]) + "..."
}

func padRight(s string, w int) string {
	if w <= 0 {
		return ""
	}
	rs := []rune(s)
	if len(rs) >= w {
		return string(rs[:w])
	}
	return s + strings.Repeat(" ", w-len(rs))
}

func commitColumnWidths(totalWidth int, maxTimeW int, maxRefsW int, maxAuthorW int) (dateW, hashW, refsW, authorW, msgW int) {
	usable := max(totalWidth-8, 20)
	hashW = 8
	dateW = max(16, maxTimeW)
	refsW = max(1, maxRefsW)
	authorW = max(1, maxAuthorW)
	separators := 8 // spaces between columns

	msgW = usable - (dateW + hashW + refsW + authorW + separators)
	if msgW < 0 {
		msgW = 0
	}
	return dateW, hashW, refsW, authorW, msgW
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
