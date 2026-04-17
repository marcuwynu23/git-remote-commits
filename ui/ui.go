package ui

import (
	"fmt"
	"strings"
	"time"

	"git-remote-commits/git"

	"github.com/charmbracelet/glamour"
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
	ShowHelp        bool
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
	selected      = lipgloss.NewStyle().Background(lipgloss.Color("#312E81"))
	freshStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ADE80"))
	hashStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FACC15"))
	refStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#93C5FD"))
	authorMe      = lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")).Bold(true)
	authorElse    = lipgloss.NewStyle().Foreground(lipgloss.Color("#86EFAC"))
	msgStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	shortcutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 0)
	addFileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ADE80")).Bold(true)
	delFileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true)
	modFileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FCD34D")).Bold(true)
	renFileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#93C5FD")).Bold(true)
)

func Render(v ViewData) string {
	outerWidth := max(v.Width, 1)
	footerLines := renderShortcutFooterLines(outerWidth)
	footerReservedHeight := 1 + len(footerLines) // one spacer line + footer lines
	outerHeight := max(v.Height-footerReservedHeight, 1)
	frameW := frameStyle.GetHorizontalFrameSize()
	frameH := frameStyle.GetVerticalFrameSize()
	mainWidth := max(outerWidth-frameW, 1)
	mainHeight := max(outerHeight-frameH, 1)

	header := ""
	headerRendered := ""
	topLineRendered := ""
	if v.Loaded {
		header = titleStyle.Render(fmt.Sprintf("git remote-commits %s", emptyFallback(v.Version, "dev")))
		headerRendered = lipgloss.NewStyle().Width(mainWidth).Render(header)
		topLineRendered = lipgloss.NewStyle().Width(mainWidth).Render(renderTopStatusLine(v, mainWidth))
	}

	topFixedLines := 0
	bottomFixedLines := 0
	if v.Loaded {
		topFixedLines = lipgloss.Height(headerRendered) + lipgloss.Height(topLineRendered) + 1
	}
	availableBody := max(mainHeight-topFixedLines-bottomFixedLines, 1)
	topHeight := availableBody
	panelHeight := 0
	showPanel := !v.ShowHelp && v.ShowCommitPanel && availableBody > 1

	body := ""
	commitPanel := ""
	if v.ShowHelp {
		body = lipgloss.NewStyle().Height(availableBody).Render(renderHelpView(mainWidth, availableBody))
	} else {
		if showPanel {
			panelHeight = max(availableBody/4, 3)
			if panelHeight > availableBody-1 {
				panelHeight = availableBody - 1
			}
			topHeight = max(availableBody-panelHeight, 1)
		}

		body = renderCommitList(v, mainWidth, topHeight)
		if !v.Loaded {
			body = renderLoading(mainWidth, topHeight)
		}
		body = lipgloss.NewStyle().Height(topHeight).Render(body)

		if v.Snapshot.RepoError != "" {
			body = errStyle.Render("Error: " + v.Snapshot.RepoError + "\n\nOpen this app from a valid git repository.")
		}

		if showPanel {
			commitPanel = renderCommitPanel(v, mainWidth, panelHeight)
		}
	}

	sections := make([]string, 0, 8)
	if v.Loaded {
		sections = append(sections,
			headerRendered,
			topLineRendered,
			"",
		)
	}
	sections = append(sections, body)
	if showPanel && panelHeight > 0 {
		sections = append(sections, commitPanel)
	}
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	panel := frameStyle.Width(mainWidth).Height(mainHeight).Render(content)
	footerStyled := make([]string, 0, len(footerLines))
	for _, line := range footerLines {
		footerStyled = append(footerStyled, lipgloss.NewStyle().Width(outerWidth).Render(line))
	}
	return lipgloss.JoinVertical(lipgloss.Left, append([]string{panel, ""}, footerStyled...)...)
}

func renderTopStatusLine(v ViewData, width int) string {
	parts := []string{
		renderHeaderLine(v),
		renderFooterLine(v),
	}
	join := func(list []string) string {
		return lipgloss.JoinHorizontal(lipgloss.Left, strings.Join(list, "  "))
	}
	line := join(parts)
	if lipgloss.Width(line) <= width {
		return metaStyle.Render(line)
	}

	// Drop least-critical details as width shrinks.
	compact := []string{
		renderHeaderLine(v),
		renderCompactSyncLine(v),
	}
	line = join(compact)
	if lipgloss.Width(line) <= width {
		return metaStyle.Render(line)
	}

	// Keep only sync state in very narrow widths.
	syncOnly := labelStyle.Render("Sync ") + renderSyncChip(v)
	return metaStyle.Render(syncOnly)
}

func renderHelpView(width int, height int) string {
	lines := []string{
		titleStyle.Render("Help"),
		"",
		labelStyle.Render("Navigation"),
		msgStyle.Render("up/down, j/k         Move selection"),
		msgStyle.Render("g or Home            Jump to latest commit"),
		msgStyle.Render("G or End             Jump to initial commit"),
		"",
		labelStyle.Render("Commit Panel"),
		msgStyle.Render("p                    Toggle commit panel"),
		msgStyle.Render("[ / ]                Scroll panel by line"),
		msgStyle.Render("u / d                Scroll panel by chunk"),
		msgStyle.Render("PgUp / PgDn          Scroll panel by chunk"),
		msgStyle.Render("Ctrl+U / Ctrl+D      Scroll panel by chunk"),
		"",
		labelStyle.Render("General"),
		msgStyle.Render("r                    Refresh now"),
		msgStyle.Render("? or h               Toggle help"),
		msgStyle.Render("q or Ctrl+C          Quit"),
	}
	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().Width(width).Height(height).Render(content)
}

func renderCommitPanel(v ViewData, width int, height int) string {
	height = max(height, 1)
	contentHeight := max(height-1, 1) // account for the top border row
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
			Height(contentHeight).
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
	leftWidth := max((panelWidth*3)/5, 20)
	if leftWidth > panelWidth-8 {
		leftWidth = max(panelWidth-8, 10)
	}
	rightWidth := max(panelWidth-leftWidth-1, 10)
	viewportWidth := max(leftWidth-2, 8)

	bodyRaw := strings.TrimSpace(c.Body)
	hasBody := bodyRaw != ""
	titleText := strings.TrimSpace(c.Message)
	if titleText == "" {
		titleText = "-"
	}

	currentAuthor := strings.TrimSpace(v.Snapshot.CurrentAuthor)
	authorRaw := emptyFallback(authorDisplay(c.Author, c.AuthorEmail), "-")
	authorStyled := authorElse.Render(authorRaw)
	if currentAuthor != "" && strings.EqualFold(strings.TrimSpace(c.Author), currentAuthor) {
		authorStyled = authorMe.Render(authorRaw)
	}

	refsText := strings.TrimSpace(commitRefsLabel(c.Refs))

	infoLines := []string{
		labelStyle.Render("Commit Hash:") + " " + hashStyle.Render(emptyFallback(c.Hash, "-")),
		labelStyle.Render("Author:") + " " + authorStyled,
		labelStyle.Render("When:") + " " + metaStyle.Render(emptyFallback(when, "-")),
	}
	if refsText != "" {
		infoLines = append(infoLines, labelStyle.Render("Refs:")+" "+refStyle.Render(refsText))
	}
	infoLines = append(infoLines, labelStyle.Render("Title:")+" "+msgStyle.Render(titleText))

	if hasBody {
		mdLines := renderMarkdownLines(bodyRaw, viewportWidth)
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, mdLines...)
	}

	fileLines := buildFileLinesFromDiff(v.Snapshot.SelectedDiff, rightWidth)

	maxLines := max(len(infoLines), len(fileLines))
	maxScroll := max(maxLines-contentHeight, 0)
	scroll := v.PanelScroll
	if scroll < 0 {
		scroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}
	hasUp := scroll > 0
	hasDown := scroll < maxScroll

	if len(infoLines) > 0 {
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
			infoLines[0] = infoLines[0] + "  " + metaStyle.Render(indicator+" u/d")
		}
	}

	leftVisible := sliceLines(infoLines, scroll, contentHeight)
	rightVisible := sliceLines(fileLines, scroll, contentHeight)

	leftBlock := lipgloss.NewStyle().
		Width(leftWidth).
		Height(contentHeight).
		Render(strings.Join(leftVisible, "\n"))

	rightBlock := lipgloss.NewStyle().
		Width(rightWidth).
		Height(contentHeight).
		Render(strings.Join(rightVisible, "\n"))

	body := lipgloss.JoinHorizontal(lipgloss.Left, leftBlock, " ", rightBlock)
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderLeft(false).
		BorderRight(false).
		BorderBottom(false).
		BorderForeground(lipgloss.Color("#374151")).
		Padding(0, 1).
		Width(panelWidth).
		Height(contentHeight).
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

func renderShortcutFooterLines(width int) []string {
	width = max(width, 20)
	groups := []string{
		"[up/j | down/k]",
		"[PgUp/u | PgDn/d]",
		"[g/Home | G/End]",
		"[p]",
		"[r]",
		"[?/h]",
		"[q/Ctrl+C]",
		"[Ctrl+U | Ctrl+D]",
	}
	lines := make([]string, 0, 3)
	currentRaw := ""
	currentStyled := ""
	for _, group := range groups {
		styledGroup := shortcutStyle.Render(group)
		nextRaw := group
		nextStyled := styledGroup
		if currentRaw != "" {
			nextRaw = currentRaw + "  " + group
			nextStyled = currentStyled + "  " + styledGroup
		}
		if len([]rune(nextRaw)) > width && currentRaw != "" {
			lines = append(lines, currentStyled)
			currentRaw = group
			currentStyled = styledGroup
			continue
		}
		currentRaw = nextRaw
		currentStyled = nextStyled
	}
	if currentStyled != "" {
		lines = append(lines, currentStyled)
	}
	return lines
}

func renderMarkdownLines(md string, width int) []string {
	if width <= 0 {
		return []string{""}
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return wrapLines(md, width)
	}
	out, err := r.Render(md)
	if err != nil {
		return wrapLines(md, width)
	}
	out = strings.TrimRight(out, "\n")
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func authorDisplay(author, email string) string {
	author = strings.TrimSpace(author)
	email = strings.TrimSpace(email)
	if email != "" {
		if at := strings.Index(email, "@"); at > 0 {
			username := strings.TrimSpace(email[:at])
			if username != "" {
				return username
			}
		}
	}
	return author
}

func sliceLines(lines []string, start, height int) []string {
	if height <= 0 {
		return []string{}
	}
	if start < 0 {
		start = 0
	}
	if start > len(lines) {
		start = len(lines)
	}
	end := min(start+height, len(lines))
	out := make([]string, 0, height)
	for i := start; i < end; i++ {
		out = append(out, lines[i])
	}
	for len(out) < height {
		out = append(out, "")
	}
	return out
}

func buildFileLinesFromDiff(diff string, width int) []string {
	width = max(width-2, 8)
	raw := strings.TrimSpace(diff)
	if raw == "" || raw == "No commit selected." || raw == "Unable to load diff preview." {
		return []string{metaStyle.Render("No file changes.")}
	}

	type fileChange struct {
		path   string
		status string // added|modified|deleted|renamed
		from   string
		to     string
	}

	seen := make(map[string]struct{})
	changes := make([]fileChange, 0)

	var cur *fileChange
	flush := func() {
		if cur == nil {
			return
		}
		key := cur.path
		if cur.status == "renamed" && cur.from != "" && cur.to != "" {
			key = cur.from + "->" + cur.to
		}
		if key == "" {
			cur = nil
			return
		}
		if _, ok := seen[key]; ok {
			cur = nil
			return
		}
		seen[key] = struct{}{}
		changes = append(changes, *cur)
		cur = nil
	}

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "diff --git ") {
			flush()
			parts := strings.Fields(trim)
			if len(parts) >= 4 {
				path := strings.TrimPrefix(parts[2], "a/")
				path = strings.TrimSpace(path)
				if path != "" {
					tmp := fileChange{path: path, status: "modified"}
					cur = &tmp
				}
			}
			continue
		}
		if cur == nil {
			continue
		}
		switch {
		case strings.HasPrefix(trim, "new file mode "):
			cur.status = "added"
		case strings.HasPrefix(trim, "deleted file mode "):
			cur.status = "deleted"
		case strings.HasPrefix(trim, "rename from "):
			cur.status = "renamed"
			cur.from = strings.TrimSpace(strings.TrimPrefix(trim, "rename from "))
		case strings.HasPrefix(trim, "rename to "):
			cur.status = "renamed"
			cur.to = strings.TrimSpace(strings.TrimPrefix(trim, "rename to "))
		case strings.HasPrefix(trim, "--- ") && strings.Contains(trim, "/dev/null"):
			cur.status = "added"
		case strings.HasPrefix(trim, "+++ ") && strings.Contains(trim, "/dev/null"):
			cur.status = "deleted"
		}
	}
	flush()

	if len(changes) == 0 {
		return []string{metaStyle.Render("No file changes.")}
	}

	lines := []string{labelStyle.Render("File Changes")}
	for _, ch := range changes {
		symbol := "~"
		symbolStyle := modFileStyle
		pathStyle := modFileStyle
		display := ch.path
		switch ch.status {
		case "added":
			symbol = "+"
			symbolStyle = addFileStyle
			pathStyle = addFileStyle
		case "deleted":
			symbol = "-"
			symbolStyle = delFileStyle
			pathStyle = delFileStyle
		case "renamed":
			symbol = "→"
			symbolStyle = renFileStyle
			pathStyle = renFileStyle
			if ch.from != "" && ch.to != "" {
				display = ch.from + " -> " + ch.to
			}
		default:
			symbol = "~"
			symbolStyle = modFileStyle
			pathStyle = modFileStyle
		}

		wrapped := wrapLines(display, width)
		if len(wrapped) == 0 {
			continue
		}
		lines = append(lines, symbolStyle.Render(symbol)+" "+pathStyle.Render(wrapped[0]))
		for i := 1; i < len(wrapped); i++ {
			lines = append(lines, "  "+pathStyle.Render(wrapped[i]))
		}
	}
	return lines
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
		ts := c.When.Local().Format("2006-01-02")
		tm := c.When.Local().Format("3:04")
		ampm := c.When.Local().Format("PM")
		if len([]rune(tm)) < 5 {
			tm += strings.Repeat(" ", 5-len([]rune(tm)))
		}
		ts = ts + ": " + tm + " " + ampm
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
			refStyle.Render(padRight(trimToWidth(emptyFallback(r.commit.Graph, "•"), 4), 4)) + " " +
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
	syncChip := renderSyncChip(v)
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

func renderCompactSyncLine(v ViewData) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		labelStyle.Render("Sync "),
		renderSyncChip(v),
		"  ",
		labelStyle.Render("Refresh "),
		chipStyle.Render(v.Snapshot.LastRefresh.Format(time.Kitchen)),
	)
}

func renderSyncChip(v ViewData) string {
	rawRemoteStatus := strings.ToLower(strings.TrimSpace(v.Snapshot.RemoteStatus))
	syncText := "online"
	if v.Refreshing {
		syncText = "pending"
	} else if rawRemoteStatus == "" ||
		strings.Contains(rawRemoteStatus, "offline") ||
		strings.Contains(rawRemoteStatus, "unavailable") ||
		strings.Contains(rawRemoteStatus, "pull failed") {
		syncText = "offline"
	}
	syncChip := chipStyle.Render(syncText)
	if v.Refreshing {
		syncChip = chipInfoStyle.Render(syncText)
	} else if strings.EqualFold(syncText, "online") {
		syncChip = chipGoodStyle.Render(syncText)
	} else if strings.EqualFold(syncText, "offline") {
		syncChip = chipWarnStyle.Render(syncText)
	}
	return syncChip
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
	// Keep each rendered row within viewport width:
	// marker(2) + date + "  " + graph(4) + " " + hash + "  " + refs + "  " + author + "  " + message
	available := max(totalWidth-2, 1) // account for selected-row prefix "> "
	graphW := 4
	hashW = 8
	separators := 9 // 2 + 1 + 2 + 2 + 2 between columns including graph gap
	minDateW := 10
	minRefsW := 1
	minAuthorW := 4
	minMsgW := 6

	dateW = min(max(maxTimeW, minDateW), 22)
	refsW = min(max(maxRefsW, minRefsW), 24)
	authorW = min(max(maxAuthorW, minAuthorW), 22)

	fixed := 2 + dateW + graphW + hashW + refsW + authorW + separators
	msgW = available - fixed
	for msgW < minMsgW {
		switch {
		case authorW > minAuthorW:
			authorW--
		case refsW > minRefsW:
			refsW--
		case dateW > minDateW:
			dateW--
		default:
			// As a last resort, allow minimal message width.
			msgW = minMsgW
			return dateW, hashW, refsW, authorW, msgW
		}
		fixed = 2 + dateW + graphW + hashW + refsW + authorW + separators
		msgW = available - fixed
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
