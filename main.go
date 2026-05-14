package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(110*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

var (
	orange  = lipgloss.Color("#FF8C00")
	white   = lipgloss.Color("#FFFFFF")
	dimGray = lipgloss.Color("#555555")
	darkBg  = lipgloss.Color("#1A1A1A")
	red     = lipgloss.Color("#FF4444")
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(orange).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(white)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(white).
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(darkBg).
				Background(orange).
				Bold(true).
				Padding(0, 2)

	dimStyle = lipgloss.NewStyle().
			Foreground(dimGray)

	resultHeaderStyle = lipgloss.NewStyle().
				Foreground(orange).
				Bold(true).
				MarginTop(1)

	resultRowStyle = lipgloss.NewStyle().
			Foreground(white)

	resultIndexStyle = lipgloss.NewStyle().
				Foreground(orange).
				Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(red)

	warnStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	sizeStyle = lipgloss.NewStyle().
			Foreground(orange)
)

type state int

const (
	stateMenu state = iota
	stateRunning
	stateResult
	stateConfirmDelete
)

type menuItem struct {
	label       string
	description string
}

var menuItems = []menuItem{
	{"Biggest Files", "Find the largest files on your system"},
}

type model struct {
	state      state
	cursor     int
	fileCursor int
	shineTick  int
	result     []fileEntry
	err        error
	scanning   bool
	statusMsg  string
}

func initialModel() model {
	return model{state: stateMenu}
}

func (m model) Init() tea.Cmd { return tickCmd() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q", "esc":
			switch m.state {
			case stateMenu:
				return m, tea.Quit
			case stateResult:
				m.state = stateMenu
				m.result = nil
				m.err = nil
				m.statusMsg = ""
				m.fileCursor = 0
			case stateConfirmDelete:
				m.state = stateResult
			}

		case "up", "k":
			switch m.state {
			case stateMenu:
				if m.cursor > 0 {
					m.cursor--
				}
			case stateResult:
				if m.fileCursor > 0 {
					m.fileCursor--
				}
			}

		case "down", "j":
			switch m.state {
			case stateMenu:
				if m.cursor < len(menuItems)-1 {
					m.cursor++
				}
			case stateResult:
				if m.fileCursor < len(m.result)-1 {
					m.fileCursor++
				}
			}

		case "enter":
			switch m.state {
			case stateMenu:
				m.state = stateRunning
				m.scanning = true
				return m, scanBiggestFiles()
			case stateResult:
				if len(m.result) > 0 {
					m.statusMsg = ""
					m.state = stateConfirmDelete
				}
			}

		case " ":
			if m.state == stateMenu {
				m.state = stateRunning
				m.scanning = true
				return m, scanBiggestFiles()
			}

		case "y", "Y":
			if m.state == stateConfirmDelete {
				path := m.result[m.fileCursor].path
				if err := os.Remove(path); err != nil {
					m.statusMsg = "Delete failed: " + err.Error()
				} else {
					m.result = append(m.result[:m.fileCursor], m.result[m.fileCursor+1:]...)
					if m.fileCursor >= len(m.result) && m.fileCursor > 0 {
						m.fileCursor--
					}
					m.statusMsg = ""
				}
				m.state = stateResult
			}

		case "n", "N":
			if m.state == stateConfirmDelete {
				m.state = stateResult
			}
		}

	case tickMsg:
		m.shineTick++
		return m, tickCmd()

	case filesResultMsg:
		m.scanning = false
		m.state = stateResult
		m.result = msg.files
		m.err = msg.err
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateRunning:
		return renderHeader() + "\n\n" +
			titleStyle.Render("  Scanning your system for largest files...") + "\n" +
			dimStyle.Render("  This may take a moment.\n")
	case stateResult:
		return renderResult(m)
	case stateConfirmDelete:
		return renderConfirmDelete(m)
	default:
		return renderMenu(m)
	}
}

func renderHeader() string {
	return titleStyle.Render("  🦆 socool") + subtitleStyle.Render(" — your cool system toolkit")
}

func renderMenu(m model) string {
	var b strings.Builder
	b.WriteString(buildLogo(m.shineTick))
	b.WriteString(titleStyle.Render("  socool") + subtitleStyle.Render(" — your cool system toolkit") + "\n\n")

	for i, item := range menuItems {
		label := fmt.Sprintf("%-20s %s", item.label, dimStyle.Render(item.description))
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▶ "+label) + "\n")
		} else {
			b.WriteString(menuItemStyle.Render("  "+label) + "\n")
		}
	}

	b.WriteString("\n" + dimStyle.Render("  ↑/↓ navigate • enter select • q quit"))
	return b.String()
}

func renderResult(m model) string {
	var b strings.Builder
	b.WriteString(renderHeader() + "\n")

	if m.err != nil {
		b.WriteString("\n" + errorStyle.Render("  Error: "+m.err.Error()) + "\n")
	} else {
		b.WriteString(resultHeaderStyle.Render("  Top 20 Largest Files") + "\n\n")
		for i, f := range m.result {
			idxStr := fmt.Sprintf("  %2d.", i+1)
			sizeStr := fmt.Sprintf("%-10s", formatSize(f.size))
			if i == m.fileCursor {
				line := idxStr + "  " + sizeStr + "  " + f.path
				b.WriteString(selectedItemStyle.Render(line) + "\n")
			} else {
				idx := resultIndexStyle.Render(idxStr)
				size := sizeStyle.Render(sizeStr)
				path := resultRowStyle.Render(f.path)
				b.WriteString(idx + "  " + size + "  " + path + "\n")
			}
		}
	}

	if m.statusMsg != "" {
		b.WriteString("\n" + errorStyle.Render("  "+m.statusMsg))
	}

	b.WriteString("\n" + dimStyle.Render("  ↑/↓ navigate • enter delete • q back"))
	return b.String()
}

func renderConfirmDelete(m model) string {
	f := m.result[m.fileCursor]

	var b strings.Builder
	b.WriteString(renderHeader() + "\n\n")
	b.WriteString(resultHeaderStyle.Render("  Delete File") + "\n\n")
	b.WriteString(resultRowStyle.Render("  "+f.path) + "\n")
	b.WriteString(sizeStyle.Render(fmt.Sprintf("  Size: %s", formatSize(f.size))) + "\n\n")
	b.WriteString(warnStyle.Render("  ⚠  This cannot be undone.") + "\n\n")
	b.WriteString(dimStyle.Render("  y confirm • n/esc cancel"))
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
