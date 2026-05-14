package main

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dariubs/socool/scanner"
)

type appState int

const (
	stateMenu appState = iota
	stateRunning
	stateResult
	stateConfirmDelete
)

type mode int

const (
	modeFiles mode = iota
	modeDirs
	modeDups
)

type (
	tickMsg        time.Time
	filesResultMsg struct {
		files []scanner.FileEntry
		err   error
	}
	dirsResultMsg struct {
		dirs []scanner.FileEntry
		err  error
	}
	dupResultMsg struct {
		groups []scanner.DupGroup
		err    error
	}
)

type menuItem struct {
	label       string
	description string
}

var menuItems = []menuItem{
	{"Biggest Files", "Find the largest files on your system"},
	{"Largest Dirs", "Find the largest directories on your system"},
	{"Duplicate Files", "Find duplicate files and wasted space"},
}

type model struct {
	state      appState
	mode       mode
	cursor     int
	listCursor int
	shineTick  int
	result     []scanner.FileEntry
	dupResult  []scanner.DupGroup
	err        error
	statusMsg  string
}

func initialModel() model { return model{} }

func tickCmd() tea.Cmd {
	return tea.Tick(110*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
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
				m.dupResult = nil
				m.err = nil
				m.statusMsg = ""
				m.listCursor = 0
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
				if m.listCursor > 0 {
					m.listCursor--
				}
			}

		case "down", "j":
			switch m.state {
			case stateMenu:
				if m.cursor < len(menuItems)-1 {
					m.cursor++
				}
			case stateResult:
				limit := len(m.result) - 1
				if m.mode == modeDups {
					limit = len(m.dupResult) - 1
				}
				if m.listCursor < limit {
					m.listCursor++
				}
			}

		case "enter":
			switch m.state {
			case stateMenu:
				return m.startScan()
			case stateResult:
				if m.mode == modeFiles && len(m.result) > 0 {
					m.statusMsg = ""
					m.state = stateConfirmDelete
				}
			}

		case " ":
			if m.state == stateMenu {
				return m.startScan()
			}

		case "y", "Y":
			if m.state == stateConfirmDelete {
				path := m.result[m.listCursor].Path
				if err := os.Remove(path); err != nil {
					m.statusMsg = "Delete failed: " + err.Error()
				} else {
					m.result = append(m.result[:m.listCursor], m.result[m.listCursor+1:]...)
					if m.listCursor >= len(m.result) && m.listCursor > 0 {
						m.listCursor--
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
		m.state = stateResult
		m.result = msg.files
		m.err = msg.err

	case dirsResultMsg:
		m.state = stateResult
		m.result = msg.dirs
		m.err = msg.err

	case dupResultMsg:
		m.state = stateResult
		m.dupResult = msg.groups
		m.err = msg.err
	}

	return m, nil
}

func (m model) startScan() (model, tea.Cmd) {
	m.state = stateRunning
	m.listCursor = 0
	switch m.cursor {
	case 0:
		m.mode = modeFiles
		return m, scanFiles()
	case 1:
		m.mode = modeDirs
		return m, scanDirs()
	case 2:
		m.mode = modeDups
		return m, scanDups()
	}
	return m, nil
}

func scanFiles() tea.Cmd {
	return func() tea.Msg {
		files, err := scanner.FindLargestFiles("/", scanner.TopN)
		return filesResultMsg{files: files, err: err}
	}
}

func scanDirs() tea.Cmd {
	return func() tea.Msg {
		dirs, err := scanner.FindLargestDirs("/", scanner.TopN)
		return dirsResultMsg{dirs: dirs, err: err}
	}
}

func scanDups() tea.Cmd {
	return func() tea.Msg {
		groups, err := scanner.FindDuplicateFiles("/", scanner.TopN)
		return dupResultMsg{groups: groups, err: err}
	}
}
