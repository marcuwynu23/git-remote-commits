package model

import (
	"os"
	"time"

	"tuiapp/git"
	"tuiapp/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultRefresh = 3 * time.Second
	commitLimit    = 20
)

type tickMsg time.Time
type snapshotMsg git.Snapshot

type Model struct {
	RepoPath        string
	RefreshInterval time.Duration
	Width           int
	Height          int
	Selected        int
	Quitting        bool

	Snapshot      git.Snapshot
	KnownHashes   map[string]struct{}
	NewCommitHash map[string]struct{}
}

func Initial() Model {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	return Model{
		RepoPath:        wd,
		RefreshInterval: defaultRefresh,
		KnownHashes:     make(map[string]struct{}),
		NewCommitHash:   make(map[string]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.pollCmd(), m.tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = v.Width
		m.Height = v.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(v)
	case snapshotMsg:
		m.applySnapshot(git.Snapshot(v))
		return m, m.tickCmd()
	case tickMsg:
		return m, m.pollCmd()
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		}
		return m, nil
	case "down", "j":
		if m.Selected < len(m.Snapshot.Commits)-1 {
			m.Selected++
		}
		return m, nil
	case "r":
		return m, m.pollCmd()
	}
	return m, nil
}

func (m *Model) applySnapshot(s git.Snapshot) {
	m.Snapshot = s
	m.NewCommitHash = make(map[string]struct{})

	for _, c := range s.Commits {
		if _, ok := m.KnownHashes[c.Hash]; !ok {
			m.NewCommitHash[c.Hash] = struct{}{}
		}
	}

	m.KnownHashes = make(map[string]struct{}, len(s.Commits))
	for _, c := range s.Commits {
		m.KnownHashes[c.Hash] = struct{}{}
	}

	if len(s.Commits) == 0 {
		m.Selected = 0
		return
	}

	if m.Selected >= len(s.Commits) {
		m.Selected = len(s.Commits) - 1
	}
	if m.Selected < 0 {
		m.Selected = 0
	}
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(m.RefreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) pollCmd() tea.Cmd {
	return func() tea.Msg {
		s := git.CollectSnapshot(m.RepoPath, commitLimit)
		return snapshotMsg(s)
	}
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	return ui.Render(ui.ViewData{
		Width:         m.Width,
		Height:        m.Height,
		Selected:      m.Selected,
		NewCommitHash: m.NewCommitHash,
		Snapshot:      m.Snapshot,
	})
}
