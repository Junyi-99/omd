package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var (
	accent             = lipgloss.NewStyle().Foreground(lipgloss.Color("#7AA2F7"))
	muted              = lipgloss.NewStyle().Foreground(lipgloss.Color("#8A93AE"))
	green              = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ECE6A"))
	red                = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7768E"))
	activeRow          = lipgloss.NewStyle().Foreground(lipgloss.Color("#C0CAF5")).Background(lipgloss.Color("#292E42"))
	installedActiveRow = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ECE6A")).Background(lipgloss.Color("#292E42"))
)

func (m model) View() tea.View {
	if m.width < 48 || m.height < 16 {
		return m.render([]string{"omd · Resize to at least 48×16", "q quit"})
	}

	width := min(m.width-4, 100)
	items := m.items()
	end := min(len(items), m.offset+m.visibleRows())
	lines := m.header(width)
	for i := m.offset; i < end; i++ {
		lines = append(lines, m.packageRow(items[i], i == m.cursor, width))
	}
	for i := end - m.offset; i < m.visibleRows(); i++ {
		lines = append(lines, "")
	}
	lines = append(lines, muted.Render(fmt.Sprintf(
		"%d–%d / %d   ·   m minimal   g agent   f full", m.offset+1, end, len(items),
	)))

	detail := "No packages in this group."
	actionHelp := ""
	if len(items) > 0 {
		pkg := items[m.cursor]
		detail = pkg.Description
		actionHelp = packageSelectionHelp(pkg, m.selected[pkg.ID])
	}
	notice := muted.Render(m.notice)
	if m.errorNotice {
		notice = red.Render(m.notice)
	}
	lines = append(lines,
		detail,
		actionHelp,
		notice,
		"",
		accent.Render("↑↓ move  Space enable/disable  Tab group  a all  r refresh  u update frequency"),
		accent.Render("Enter install selected tools & return to shell  s save only  q quit"),
	)
	return m.render(lines)
}

func packageSelectionHelp(p Package, selected bool) string {
	if p.Env != "" {
		if selected {
			return "Selected: OMD enables its shell integration; Enter installs it if missing."
		}
		return "Unselected: OMD omits this integration; installed files stay in place."
	}
	if selected {
		return "Selected: Enter installs it if missing; OMD adds no shell integration."
	}
	return "Unselected: OMD skips installation; existing files stay in place."
}

func (m model) header(width int) []string {
	tag := "saved"
	if m.dirty {
		tag = "unsaved"
	}
	title := accent.Bold(true).Render("OMD") + muted.Render("  /  OH MY DOTFILES  ·  "+tag)
	var tabs []string
	for i, name := range []string{"Basics", "CLI Agents"} {
		label := fmt.Sprintf(" %d %s ", i+1, name)
		style := muted
		if m.tab == i {
			style = accent.Bold(true).Underline(true)
		}
		tabs = append(tabs, style.Render(label))
	}

	selectedCount, enabledIntegrations, installedCount, missing, checkedCount := 0, 0, 0, 0, 0
	for _, p := range m.manager.packages {
		if m.selected[p.ID] {
			selectedCount++
			if p.Env != "" {
				enabledIntegrations++
			}
		}
		if installed, checked := m.installed[p.ID]; checked {
			checkedCount++
			if installed {
				installedCount++
			} else {
				missing++
			}
		}
	}
	summary := fmt.Sprintf("%d selected  ·  %d integrations enabled  ·  %d installed  ·  %d missing  ·  %d updates  ·  %d/%d checked",
		selectedCount, enabledIntegrations, installedCount, missing, len(m.updates), checkedCount, len(m.manager.packages))
	return []string{
		title,
		muted.Render("Your terminal, assembled.  ·  Auto-update: " + m.updateFrequency + " (u to change)"),
		"",
		strings.Join(tabs, "    "),
		muted.Render(summary),
		muted.Render(strings.Repeat("─", width)),
	}
}

func (m model) packageRow(p Package, focused bool, width int) string {
	pointer, box := " ", "[ ]"
	if focused {
		pointer = "›"
	}
	if m.selected[p.ID] {
		box = "[✓]"
	}
	selectionState := "unselected"
	if m.selected[p.ID] {
		selectionState = "selected"
	}
	if p.Env != "" {
		selectionState = "integration disabled"
		if m.selected[p.ID] {
			selectionState = "integration enabled"
		}
	}
	availability := "checking"
	installed, checked := m.installed[p.ID]
	if checked {
		availability = "missing"
		if installed {
			availability = "installed"
		}
		if latest, ok := m.updates[p.ID]; ok && installed {
			availability = "update " + latest
		}
	}
	status := availability + " · " + selectionState

	nameWidth := max(8, width-lipgloss.Width(status)-9)
	name := p.Name
	if width >= 76 {
		name += "  / " + p.Category
	}
	name = ansi.Truncate(name, nameWidth, "…")
	padding := strings.Repeat(" ", max(0, nameWidth-lipgloss.Width(name)))
	row := fmt.Sprintf("%s %s %s%s  %9s", pointer, box, name, padding, status)
	switch {
	case checked && installed && focused:
		return installedActiveRow.Width(width).Render(row)
	case checked && installed:
		return green.Render(row)
	case focused:
		return activeRow.Width(width).Render(row)
	default:
		return row
	}
}

func (m model) render(lines []string) tea.View {
	lines = lines[:min(len(lines), max(1, m.height))]
	for i, line := range lines {
		line = strings.ReplaceAll(line, "\n", " ")
		lines[i] = ansi.Truncate("  "+line, max(1, m.width), "…")
	}
	view := tea.NewView(strings.Join(lines, "\n"))
	view.AltScreen = true
	return view
}
