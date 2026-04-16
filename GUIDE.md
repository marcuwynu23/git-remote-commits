# TUI Components Guide

This guide shows how to build UI components in your BubbleTea + Lip Gloss app: **menubar**, **sidebar**, **text editor**, and layout patterns you can reuse.

---

## 1. How the UI is built

- **BubbleTea**: Your app has a single `model` and one `View()` that returns a **string**. Every frame, `View()` is called and the whole screen is that string.
- **Lip Gloss**: You build that string from smaller pieces using **styles** (colors, padding, width, height) and **layout** (placing blocks next to or below each other).
- **Full screen**: Store `width` and `height` from `tea.WindowSizeMsg` and use them when rendering so the UI fills the terminal and resizes correctly.

```go
// In Update:
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    return m, nil

// In View: give your root layout a fixed size so it fills the terminal
rootStyle := lipgloss.NewStyle().Width(m.width).Height(m.height)
return rootStyle.Render(menubar + "\n" + content)
```

---

## 2. Menubar (top bar)

The template already has a horizontal menubar. Here’s what it does and how to change it.

### What it is

- One row at the top: **App name │ Item1 │ Item2 │ …**
- State: `menuIndex` (which item is focused).
- Keys: **←/→** or **h/l** to move, **Enter** to activate.

### Customizing the menubar

**Add items**  
Append to `menuItems` in `initialModel()`, then add a `case "NewItem":` in `handleMenuKeys` (and in `handleScreenKeys` for Enter) and in `screenForLabel()`.

**Change look**  
Edit the `var` block at the top of `main.go`:

```go
menubarStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("#2D2D2D")).  // bar background
    Foreground(lipgloss.Color("#E0E0E0")).  // default text
    Padding(0, 1)

menubarItemSelectedStyle = lipgloss.NewStyle().
    Padding(0, 2).
    Bold(true).
    Foreground(lipgloss.Color("#7D56F4")).   // selected item color
    Background(lipgloss.Color("#3D3D3D"))     // selected item background
```

**Menubar with dropdowns (concept)**  
Keep a `menubarOpen bool` and `menubarDropdownIndex int`. When `menubarOpen` is true, draw a small box under the selected item (using `lipgloss.Place` or a padded block) with sub-items, and handle up/down/Enter/Esc in `Update`. Same idea as the main menu: index + key handling + render a list under the item.

---

## 3. Sidebar (left panel)

A sidebar is a **left column** (fixed or proportional width) and a **main content** area to its right. Use Lip Gloss to build two blocks and join them horizontally.

### Layout

- Row 1: menubar (full width).
- Row 2+: `[ Sidebar | Main content ]` — same for all remaining height.

### Example: fixed-width sidebar

```go
const sidebarWidth = 24

func (m model) renderWithSidebar(mainContent string) string {
    // Content area height (below menubar)
    contentHeight := m.height - 1

    sidebarContent := "  Files\n  • doc.txt\n  • notes.md\n  Settings\n  About"
    sidebarStyle := lipgloss.NewStyle().
        Width(sidebarWidth).
        Height(contentHeight).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#444")).
        Padding(0, 1)

    mainStyle := lipgloss.NewStyle().
        Width(m.width - sidebarWidth).
        Height(contentHeight).
        Padding(0, 2)

    sidebar := sidebarStyle.Render(sidebarContent)
    main := mainStyle.Render(mainContent)
    return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
}
```

### Sidebar with selection

Add state for the selected sidebar item and handle keys when that screen is active:

```go
type model struct {
    // ... existing fields ...
    sidebarIndex int
    sidebarItems  []string
}

// In Update (when current screen uses sidebar):
case "up", "k":
    if m.sidebarIndex > 0 { m.sidebarIndex-- }
case "down", "j":
    if m.sidebarIndex < len(m.sidebarItems)-1 { m.sidebarIndex++ }
case "enter":
    // Navigate or open selected sidebar item
```

Render the sidebar by iterating `sidebarItems` and highlighting the one at `sidebarIndex` (e.g. with a different style and a `▸` prefix).

### Proportional sidebar

Use a fraction of `m.width` instead of a constant:

```go
sidebarWidth := m.width / 4  // 25% of terminal width
```

---

## 4. Text editor (editable text area)

You can use the **Bubbles textarea** component or a **simple custom editor**. Both need: a string or lines, cursor position, and key handling (typing, backspace, enter, arrows).

### Option A: Bubbles textarea (recommended)

Add the dependency:

```bash
go get github.com/charmbracelet/bubbles/textarea
```

In your model, embed the component and delegate size/keys:

```go
import "github.com/charmbracelet/bubbles/textarea"

type model struct {
    // ... existing ...
    editor textarea.Model
}

func initialModel() model {
    m := model{ /* ... */ }
    ti := textarea.New()
    ti.Placeholder = "Type here..."
    ti.SetWidth(60)
    ti.SetHeight(10)
    m.editor = ti
    return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // When this screen is active and you want the editor focused:
    var cmd tea.Cmd
    m.editor, cmd = m.editor.Update(msg)
    return m, cmd
}

func (m model) View() string {
    // ... menubar ...
    return menubar + "\n" + m.editor.View()
}
```

Resize the textarea on `tea.WindowSizeMsg` (e.g. `m.editor.SetWidth(m.width - 4)` and `SetHeight(m.height - 2)`).

### Option B: Simple custom editor (no deps)

Store lines and cursor; handle keys in `Update`; render in `View`.

```go
type model struct {
    // ...
    editorLines []string  // lines of text
    editorRow   int       // cursor row (0-based)
    editorCol   int       // cursor column (0-based)
}

// In Update when editor is focused:
case tea.KeyMsg:
    switch msg.String() {
    case "backspace":
        // Remove rune before cursor or merge lines
    case "enter":
        // Split line at cursor, insert new line
    case "up", "down", "left", "right":
        // Move editorRow / editorCol (clamp to line length and line count)
    default:
        // Insert rune at cursor (msg.Type == tea.KeyRune)
        if msg.Type == tea.KeyRune {
            // Insert msg.Rune into editorLines[editorRow] at editorCol
        }
    }
```

Rendering: loop over `editorLines`, add a visible cursor (e.g. block or underscore) at `(editorRow, editorCol)` using a helper that inserts the cursor into the current line string.

```go
func (m model) renderEditor() string {
    var b strings.Builder
    for i, line := range m.editorLines {
        if i == m.editorRow {
            // Insert cursor at editorCol (e.g. show a block character)
            line = line[:m.editorCol] + "█" + line[m.editorCol:]
        }
        b.WriteString(line + "\n")
    }
    style := lipgloss.NewStyle().Width(m.width - 4).Height(m.height - 2)
    return style.Render(b.String())
}
```

Use a proper rune-aware cursor (e.g. `rune` slice or `utf8` functions) so cursor position and backspace work with multi-byte characters.

---

## 5. Layout patterns with Lip Gloss

### Width and height

- **Fixed size**: `lipgloss.NewStyle().Width(40).Height(10).Render(s)` — content is truncated or padded to that size.
- **Full width**: `Width(m.width)` so the block uses the whole terminal width.
- **Remaining height**: e.g. `Height(m.height - 1)` for content below the menubar.

### Joining blocks

- **Side by side**: `lipgloss.JoinHorizontal(alignment, block1, block2, ...)` — alignment is `lipgloss.Top`, `Center`, or `Bottom`.
- **Stacked**: `lipgloss.JoinVertical(alignment, block1, block2, ...)` — use `lipgloss.Left`, `Center`, or `Right`.

Example: sidebar + main:

```go
row := lipgloss.JoinHorizontal(lipgloss.Top, sidebarBlock, mainBlock)
```

### Placing a block

- **Place** a small block inside a larger area (e.g. a dropdown under a menubar item):

```go
placed := lipgloss.Place(m.width, m.height-1, lipgloss.Center, lipgloss.Center, smallBlock, lipgloss.WithWhitespaceChars(" "))
```

### Borders

- Add a border to any block: `.Border(lipgloss.RoundedBorder())` or `lipgloss.NormalBorder()`, and `.BorderForeground(color)` to style it.

### Padding and margin

- **Padding**: space inside the style (around the text).
- **Margin**: space outside (e.g. between menubar and content).

Use them to separate sidebar from main area or to indent content.

---

## 6. Putting it together: one possible layout

A typical “app” layout could look like this:

```
┌─────────────────────────────────────────────────────────┐
│ AppName  │  Home  │  Dashboard  │  Settings  │  Quit    │  ← menubar (full width)
├──────────┬──────────────────────────────────────────────┤
│  Sidebar │                                              │
│  • Item1 │         Main content (e.g. text editor        │
│  • Item2 │         or dashboard widgets)                 │
│          │                                              │
└──────────┴──────────────────────────────────────────────┘
```

In code:

1. **Menubar**: already in the template; keep it as the first line.
2. **Second row**: `sidebarBlock := sidebarStyle.Render(sidebarContent)` and `mainBlock := mainStyle.Render(mainContent)`; then `contentRow := lipgloss.JoinHorizontal(lipgloss.Top, sidebarBlock, mainBlock)`.
3. **Sizing**: `sidebarStyle.Width(sidebarWidth).Height(m.height-1)`, `mainStyle.Width(m.width-sidebarWidth).Height(m.height-1)`.
4. **View**: `return topBarStyle.Render(menubar) + "\n" + contentRow`.

Then plug in “main content” per screen: static text, a Bubbles textarea, or your custom editor. Use the same `Update` branching you already have (e.g. by `currentScreen`) to route keys to the menubar, sidebar, or editor.

---

## 7. Quick reference

| Goal              | Approach |
|-------------------|----------|
| **Menubar**       | One row of items; `menuIndex` + ←/→ and Enter; style with `menubarStyle`, `menubarItemSelectedStyle`. |
| **Sidebar**       | Fixed- or %-width block + main block; `lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)`; optional `sidebarIndex` + up/down/Enter. |
| **Text editor**   | Use `github.com/charmbracelet/bubbles/textarea` and wire it in `Update`/`View`, or build a simple line+cursor model and handle KeyRune/Backspace/Enter/arrows. |
| **Full screen**   | Store `tea.WindowSizeMsg` in `m.width`/`m.height`; apply `.Width(m.width).Height(m.height)` (or minus menubar height) to the root or content area. |
| **Borders**      | `lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(...)`. |
| **Focus**        | Track `focus` (e.g. "menubar" / "sidebar" / "editor") and in `Update` only handle keys for the focused part; in `View` highlight only that part. |

For more components (tables, lists, inputs), see [Charm Bubbles](https://github.com/charmbracelet/bubbles) and [BubbleTea examples](https://github.com/charmbracelet/bubbletea/tree/master/examples).
