# TUI App Template

A **Go project template** for building terminal user interface (TUI) applications with a full-screen layout, top menubar, and navigable screens. Clone, customize, and ship.

---

## Features

- **Full terminal** — Uses the full terminal size (alt screen); content resizes with the window.
- **Top menubar** — Horizontal bar with app name and items (Home │ Dashboard │ Settings │ About │ Quit). Navigate with **←/→** or **h/l**, **Enter** to open.
- **Multi-screen** — Simple routing: add screen constants, menu items, and key handlers.
- **Stack** — [BubbleTea](https://github.com/charmbracelet/bubbletea) (Elm-style TUI) + [Lip Gloss](https://github.com/charmbracelet/lipgloss) (styling).

---

## Quick start

1. **Copy or clone** this repository.
2. **Rename the module** in `go.mod` (e.g. `module my-app`).
3. **Set your app name** in `main.go`: update the `appName` constant.
4. **Run:**

   ```bash
   go mod tidy
   go run .
   ```

**Requirements:** Go 1.18+ (1.21+ recommended).

---

## Navigation

| Key | Action |
|-----|--------|
| **← / h** | Move menubar selection left |
| **→ / l** | Move menubar selection right |
| **Enter** | Open selected screen |
| **Esc / b** | Back to main menu |
| **q** / **Ctrl+C** | Quit |

---

## Documentation

| Resource | Description |
|----------|-------------|
| **[GUIDE.md](GUIDE.md)** | How to build UI components: **menubar**, **sidebar**, **text editor**, and layout patterns with Lip Gloss. |
| **[LICENSE](LICENSE)** | MIT License. |

---

## Template structure

| Where | What to do |
|-------|------------|
| `appName` (main.go) | Your app’s display name. |
| Screen constants | Add `screenFoo = "foo"` for each new screen. |
| `menuItems` | Add the menu label (e.g. `"My Screen"`). |
| `handleMenuKeys()` / `handleScreenKeys()` | Add `case "My Screen": m.currentScreen = screenFoo`; update `screenForLabel()`. |
| `View()` | Add `case screenFoo: return m.renderScreen("My Screen", "content")`. |

Optional: add styles in the `var (...)` block and use them in `renderMenubar()` or `renderScreen()`.

---

## Default screens

| Screen | Purpose |
|--------|---------|
| **Home** | Placeholder welcome. |
| **Dashboard** | Placeholder list; replace with data or widgets. |
| **Settings** | Placeholder options; replace with toggles/inputs. |
| **About** | Uses `appName` and keybindings. |
| **Quit** | Exit the app. |

Remove or rename by editing constants, `menuItems`, and the switch cases in the key handlers and `View()`.

---

## Build

```bash
go build -o my-tui .
./my-tui
```

---

## License

This project is licensed under the **MIT License** — see [LICENSE](LICENSE) for details.
