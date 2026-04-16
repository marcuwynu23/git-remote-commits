// TUI App Template — A Go terminal UI with navigable screens.
//
// Customize: set appName below, then add screens (constants + menu item + handler + view).
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Template: set your app name (used in menu title and About) ---
const appName = "My TUI App"

// --- Template: add a new screen constant when you add a screen ---
const (
	screenMenu      = "menu"
	screenHome      = "home"
	screenDashboard = "dashboard"
	screenSettings  = "settings"
	screenAbout     = "about"
)

var (
	// Menubar (top bar)
	menubarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2D2D2D")).
			Foreground(lipgloss.Color("#E0E0E0")).
			Padding(0, 1)

	menubarItemStyle = lipgloss.NewStyle().
				Padding(0, 2)

	menubarItemSelectedStyle = lipgloss.NewStyle().
					Padding(0, 2).
					Bold(true).
					Foreground(lipgloss.Color("#7D56F4")).
					Background(lipgloss.Color("#3D3D3D"))

	menubarSeparator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555")).
			Render("│")

	// Content area
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			MarginTop(1)

	screenTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				MarginBottom(2)

	bodyStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2)
)

type model struct {
	currentScreen string
	menuIndex     int
	menuItems     []string
	quitting      bool
	width        int
	height       int
}

func initialModel() model {
	return model{
		currentScreen: screenMenu,
		menuIndex:     0,
		// --- Template: add menu items here; keep "Quit" last if you use it ---
		menuItems: []string{
			"Home",
			"Dashboard",
			"Settings",
			"About",
			"Quit",
		},
		width:  60,
		height: 20,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.currentScreen {
		case screenMenu:
			return m.handleMenuKeys(msg)
		default:
			return m.handleScreenKeys(msg)
		}
	}
	return m, nil
}

func (m model) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "left", "h":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
		return m, nil

	case "right", "l":
		if m.menuIndex < len(m.menuItems)-1 {
			m.menuIndex++
		}
		return m, nil

	case "enter":
		// --- Template: add a case for each new menu item (match menuItems labels) ---
		choice := m.menuItems[m.menuIndex]
		switch choice {
		case "Quit":
			m.quitting = true
			return m, tea.Quit
		case "Home":
			m.currentScreen = screenHome
		case "Dashboard":
			m.currentScreen = screenDashboard
		case "Settings":
			m.currentScreen = screenSettings
		case "About":
			m.currentScreen = screenAbout
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleScreenKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc", "b":
		m.currentScreen = screenMenu
		return m, nil
	case "left", "h":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
		return m, nil
	case "right", "l":
		if m.menuIndex < len(m.menuItems)-1 {
			m.menuIndex++
		}
		return m, nil
	case "enter":
		// From any screen, Enter activates the currently focused menubar item
		choice := m.menuItems[m.menuIndex]
		if choice == "Quit" {
			m.quitting = true
			return m, tea.Quit
		}
		m.currentScreen = m.screenForLabel(choice)
		return m, nil
	}
	return m, nil
}

func (m model) screenForLabel(label string) string {
	switch label {
	case "Home":
		return screenHome
	case "Dashboard":
		return screenDashboard
	case "Settings":
		return screenSettings
	case "About":
		return screenAbout
	default:
		return screenMenu
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	menubar := m.renderMenubar()
	var content string
	// --- Template: add a case for each screen; use renderScreen(title, body) or custom view ---
	switch m.currentScreen {
	case screenMenu:
		content = m.renderMenuContent()
	case screenHome:
		content = m.renderScreen("Home", "Welcome home.\n\nThis is the Home screen. Replace this text with your own content.")
	case screenDashboard:
		content = m.renderScreen("Dashboard", "Dashboard\n\n• Item A: 42\n• Item B: 17\n• Item C: 99\n\nReplace with your widgets or data.")
	case screenSettings:
		content = m.renderScreen("Settings", "Settings\n\n[ ] Option 1\n[ ] Option 2\n\nUse Esc to go back.")
	case screenAbout:
		content = m.renderScreen("About", fmt.Sprintf("About %s\n\nA Terminal User Interface built with Go and BubbleTea.\n\n←/→ or h/l: menu • Enter: select • Esc: back • q: quit", appName))
	default:
		content = m.renderMenuContent()
	}

	// Full terminal: menubar on first line, content fills the rest
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 1)
	topBarStyle := lipgloss.NewStyle().Width(m.width)
	full := topBarStyle.Render(menubar) + "\n" + contentStyle.Render(content)
	return full
}

func (m model) renderMenubar() string {
	var parts []string
	for i, item := range m.menuItems {
		if i > 0 {
			parts = append(parts, menubarSeparator)
		}
		// Highlight the focused menu item (cursor position)
		selected := i == m.menuIndex
		if selected {
			parts = append(parts, menubarItemSelectedStyle.Render(item))
		} else {
			parts = append(parts, menubarItemStyle.Render(item))
		}
	}
	return menubarStyle.Render(appName + "  " + menubarSeparator + "  " + strings.Join(parts, " "))
}

func (m model) renderMenuContent() string {
	return bodyStyle.Render("Select an item from the menubar above.\n\n←/→ or h/l: move • Enter: open • q: quit")
}

func (m model) renderScreen(title, body string) string {
	contentHeight := m.height - 1
	helpLine := helpStyle.Render("Esc or b: back to menu • q: quit")
	mainContent := screenTitleStyle.Render(title) + bodyStyle.Render(body)
	mainLines := strings.Count(mainContent, "\n") + 1
	padding := contentHeight - 1 - mainLines
	if padding < 0 {
		padding = 0
	}
	var b strings.Builder
	b.WriteString(mainContent)
	for i := 0; i < padding; i++ {
		b.WriteString("\n")
	}
	b.WriteString(helpLine)
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
