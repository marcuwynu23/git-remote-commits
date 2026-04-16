package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git-remote-commits/git"
	"git-remote-commits/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultRefresh = 3 * time.Second
	commitLimit    = 0
)

type tickMsg time.Time
type snapshotMsg git.Snapshot

type Model struct {
	RepoPath        string
	RemoteName      string
	RefreshInterval time.Duration
	Width           int
	Height          int
	Selected        int
	Loaded          bool
	Refreshing      bool
	PendingRefresh  bool
	Quitting        bool

	Snapshot      git.Snapshot
	KnownHashes   map[string]struct{}
	NewCommitHash map[string]struct{}
}

func Initial(remoteName string) Model {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	if remoteName == "" {
		remoteName = "origin"
	}

	return Model{
		RepoPath:        wd,
		RemoteName:      remoteName,
		RefreshInterval: defaultRefresh,
		Refreshing:      true,
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
		m.Refreshing = false
		shouldRing := m.applySnapshot(git.Snapshot(v))
		cmds := make([]tea.Cmd, 0, 3)
		if shouldRing {
			cmds = append(cmds, bellCmd())
		}
		if m.PendingRefresh {
			m.PendingRefresh = false
			m.Refreshing = true
			cmds = append(cmds, m.pollCmd())
		}
		return m, tea.Batch(cmds...)
	case tickMsg:
		if m.Refreshing {
			return m, m.tickCmd()
		}
		m.Refreshing = true
		return m, tea.Batch(m.tickCmd(), m.pollCmd())
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
		if m.Refreshing {
			m.PendingRefresh = true
			return m, nil
		}
		m.Refreshing = true
		return m, m.pollCmd()
	}
	return m, nil
}

func (m *Model) applySnapshot(s git.Snapshot) bool {
	firstLoad := !m.Loaded
	previousTopHash := ""
	if len(m.Snapshot.Commits) > 0 {
		previousTopHash = m.Snapshot.Commits[0].Hash
	}
	m.Snapshot = s
	m.Loaded = true
	m.NewCommitHash = make(map[string]struct{})

	for _, c := range s.Commits {
		if _, ok := m.KnownHashes[c.Hash]; !ok {
			m.NewCommitHash[c.Hash] = struct{}{}
		}
	}
	newCount := len(m.NewCommitHash)

	m.KnownHashes = make(map[string]struct{}, len(s.Commits))
	for _, c := range s.Commits {
		m.KnownHashes[c.Hash] = struct{}{}
	}

	if len(s.Commits) == 0 {
		m.Selected = 0
		return !firstLoad && newCount > 0
	}

	// If the newest commit changed after refresh/pull, jump to top so updated content is visible.
	if previousTopHash != "" && s.Commits[0].Hash != previousTopHash {
		m.Selected = 0
	}

	if m.Selected >= len(s.Commits) {
		m.Selected = len(s.Commits) - 1
	}
	if m.Selected < 0 {
		m.Selected = 0
	}
	return !firstLoad && newCount > 0
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(m.RefreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) pollCmd() tea.Cmd {
	return func() tea.Msg {
		s := git.CollectSnapshot(m.RepoPath, m.RemoteName, commitLimit)
		return snapshotMsg(s)
	}
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	repoName := strings.TrimSpace(m.Snapshot.RepoName)
	if repoName == "" {
		repoName = filepath.Base(m.RepoPath)
	}
	return ui.Render(ui.ViewData{
		Width:         m.Width,
		Height:        m.Height,
		RepoName:      repoName,
		Selected:      m.Selected,
		Loaded:        m.Loaded,
		Refreshing:    m.Refreshing,
		NewCommitHash: m.NewCommitHash,
		Snapshot:      m.Snapshot,
	})
}

func bellCmd() tea.Cmd {
	return func() tea.Msg {
		fmt.Print("\a")
		return nil
	}
}
