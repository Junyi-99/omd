package main

import (
	"context"
	"maps"
	"os"
	"os/exec"
	"slices"

	tea "charm.land/bubbletea/v2"
)

type checkedMsg struct {
	id        string
	installed bool
}

// releaseMsg carries the latest release of a tool omd installed from GitHub.
type releaseMsg struct {
	id, installed, latest string
}
type savedMsg struct {
	err     error
	install bool
}
type installedMsg struct{ err error }

type model struct {
	manager         *manager
	selected        map[string]bool
	installed       map[string]bool
	updates         map[string]string // id → newer release tag
	updateFrequency string

	tab    int
	cursor int
	offset int
	width  int
	height int

	busy        bool
	dirty       bool
	notice      string
	errorNotice bool
}

func newModel(m *manager, selected map[string]bool, frequencies ...string) model {
	frequency := defaultUpdateFrequency
	if len(frequencies) > 0 {
		frequency = frequencies[0]
	}
	return model{
		manager:         m,
		selected:        selected,
		installed:       make(map[string]bool),
		updates:         make(map[string]string),
		updateFrequency: frequency,
		width:           80,
		height:          24,
		notice:          "Space selects tools; for shell plugins it also toggles integration. Enter installs selected tools.",
	}
}

func (m model) Init() tea.Cmd {
	return m.checks()
}

func (m model) checks() tea.Cmd {
	var cmds []tea.Cmd
	for _, p := range m.manager.packages {
		cmds = append(cmds, func() tea.Msg {
			return checkedMsg{id: p.ID, installed: m.manager.check(p.Check)}
		})
		if m.manager.releaseFor(p) != nil {
			cmds = append(cmds, func() tea.Msg {
				// Errors only mean no update is shown; omd outdated reports them.
				installed, latest, _ := m.manager.checkRelease(context.Background(), p)
				return releaseMsg{id: p.ID, installed: installed, latest: latest}
			})
		}
	}
	return tea.Batch(cmds...)
}

func (m model) items() []Package {
	var result []Package
	for _, p := range m.manager.packages {
		if p.Tab == []string{"Basics", "CLI Agents"}[m.tab] {
			result = append(result, p)
		}
	}
	return result
}

func (m model) visibleRows() int { return max(1, m.height-13) }

func (m *model) scroll() {
	n := len(m.items())
	m.cursor = max(0, min(m.cursor, n-1))
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.visibleRows() {
		m.offset = m.cursor - m.visibleRows() + 1
	}
	m.offset = max(0, min(m.offset, n-m.visibleRows()))
}

func (m model) save(install bool) tea.Cmd {
	selected := maps.Clone(m.selected)
	frequency := m.updateFrequency
	return func() tea.Msg {
		err := m.manager.saveConfig(selected)
		if err == nil {
			err = m.manager.saveUpdateFrequency(frequency)
		}
		if err == nil && !install {
			err = m.manager.generateEnv(selected)
		}
		return savedMsg{err: err, install: install}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.scroll()
	case checkedMsg:
		m.installed[msg.id] = msg.installed
	case releaseMsg:
		if msg.latest != "" && msg.latest != msg.installed {
			m.updates[msg.id] = msg.latest
		}
	case savedMsg:
		m.busy = false
		if msg.err != nil {
			m.notice, m.errorNotice = msg.err.Error(), true
			return m, nil
		}
		m.dirty = false
		if !msg.install {
			m.notice, m.errorNotice = "Saved. Exit omd to load the environment in this shell.", false
			return m, nil
		}
		exe, err := os.Executable()
		if err != nil {
			m.notice, m.errorNotice = err.Error(), true
			return m, nil
		}
		m.busy = true
		// Hand the terminal to the installer, including sudo/password prompts.
		// Bubble Tea restores the UI when the child exits.
		cmd := exec.Command(exe, "install")
		cmd.Env = m.manager.environment()
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return installedMsg{err: err} })
	case installedMsg:
		m.busy, m.errorNotice = false, msg.err != nil
		m.notice = "Installation complete. Returning to your shell to load the environment."
		if msg.err != nil {
			m.notice = "Installation failed. See cache/omd/install.log; Enter retries."
			m.installed, m.updates = make(map[string]bool), make(map[string]string)
			return m, m.checks()
		}
		return m, tea.Quit
	case tea.KeyPressMsg:
		return m.updateKey(msg.String())
	}
	return m, nil
}

func (m model) updateKey(key string) (tea.Model, tea.Cmd) {
	if m.busy {
		return m, nil
	}
	items := m.items()
	switch key {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit
	case "tab", "right", "l", "shift+tab", "left", "h":
		m.tab = (m.tab + 1) % 2
		m.cursor, m.offset = 0, 0
	case "1", "2":
		m.tab = int(key[0] - '1')
		m.cursor, m.offset = 0, 0
	case "j", "down":
		m.cursor++
	case "k", "up":
		m.cursor--
	case "pgdown":
		m.cursor += m.visibleRows()
	case "pgup":
		m.cursor -= m.visibleRows()
	case "home":
		m.cursor = 0
	case "end":
		m.cursor = len(items) - 1
	case "space":
		if len(items) > 0 {
			p := items[m.cursor]
			m.selected[p.ID] = !m.selected[p.ID]
			m.dirty = true
		}
	case "a":
		all := true
		for _, p := range items {
			all = all && m.selected[p.ID]
		}
		for _, p := range items {
			m.selected[p.ID] = !all
		}
		m.dirty = true
	case "u":
		for i, frequency := range updateFrequencies {
			if frequency == m.updateFrequency {
				m.updateFrequency = updateFrequencies[(i+1)%len(updateFrequencies)]
				break
			}
		}
		m.dirty = true
		m.notice, m.errorNotice = "Auto-update frequency: "+m.updateFrequency+". Press s to save.", false
	case "m", "g", "f":
		preset := map[string]string{"m": "minimal", "g": "agent", "f": "full"}[key]
		for _, p := range m.manager.packages {
			m.selected[p.ID] = slices.Contains(p.Presets, preset)
		}
		m.dirty = true
		m.notice, m.errorNotice = "Applied "+preset+" preset. Enter installs; s saves only.", false
	case "s", "enter":
		m.busy = true
		m.notice, m.errorNotice = "Saving selections…", false
		return m, m.save(key == "enter")
	case "r":
		m.installed, m.updates = make(map[string]bool), make(map[string]string)
		return m, m.checks()
	}
	m.scroll()
	return m, nil
}
