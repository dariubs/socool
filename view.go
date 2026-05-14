package main

import (
	"fmt"
	"strings"

	"github.com/dariubs/socool/scanner"
)

func (m model) View() string {
	switch m.state {
	case stateRunning:
		subject := map[mode]string{
			modeFiles: "files",
			modeDirs:  "directories",
			modeDups:  "duplicate files",
		}[m.mode]
		return renderHeader() + "\n\n" +
			titleStyle.Render("  Scanning your system for "+subject+"...") + "\n" +
			dimStyle.Render("  This may take a moment.\n")
	case stateResult:
		if m.mode == modeDups {
			return renderDupResult(m)
		}
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
		header := "  Top 20 Largest Files"
		if m.mode == modeDirs {
			header = "  Top 20 Largest Directories"
		}
		b.WriteString(resultHeaderStyle.Render(header) + "\n\n")
		for i, f := range m.result {
			idxStr := fmt.Sprintf("  %2d.", i+1)
			sizeStr := fmt.Sprintf("%-10s", scanner.FormatSize(f.Size))
			if i == m.listCursor {
				b.WriteString(selectedItemStyle.Render(idxStr+"  "+sizeStr+"  "+f.Path) + "\n")
			} else {
				b.WriteString(resultIndexStyle.Render(idxStr) + "  " +
					sizeStyle.Render(sizeStr) + "  " +
					resultRowStyle.Render(f.Path) + "\n")
			}
		}
	}

	if m.statusMsg != "" {
		b.WriteString("\n" + errorStyle.Render("  "+m.statusMsg))
	}

	hint := "  ↑/↓ navigate • enter delete • q back"
	if m.mode == modeDirs {
		hint = "  ↑/↓ navigate • q back"
	}
	b.WriteString("\n" + dimStyle.Render(hint))
	return b.String()
}

func renderConfirmDelete(m model) string {
	f := m.result[m.listCursor]
	var b strings.Builder
	b.WriteString(renderHeader() + "\n\n")
	b.WriteString(resultHeaderStyle.Render("  Delete File") + "\n\n")
	b.WriteString(resultRowStyle.Render("  "+f.Path) + "\n")
	b.WriteString(sizeStyle.Render(fmt.Sprintf("  Size: %s", scanner.FormatSize(f.Size))) + "\n\n")
	b.WriteString(warnStyle.Render("  ⚠  This cannot be undone.") + "\n\n")
	b.WriteString(dimStyle.Render("  y confirm • n/esc cancel"))
	return b.String()
}

func renderDupResult(m model) string {
	var b strings.Builder
	b.WriteString(renderHeader() + "\n")

	if m.err != nil {
		b.WriteString("\n" + errorStyle.Render("  Error: "+m.err.Error()) + "\n")
	} else if len(m.dupResult) == 0 {
		b.WriteString(resultHeaderStyle.Render("  Duplicate Files") + "\n\n")
		b.WriteString(resultRowStyle.Render("  No duplicate files found.") + "\n")
	} else {
		b.WriteString(resultHeaderStyle.Render(fmt.Sprintf("  Top %d Duplicate Groups", len(m.dupResult))) + "\n\n")
		for i, g := range m.dupResult {
			wasted := g.Size * int64(len(g.Paths)-1)
			row := fmt.Sprintf("  %d copies  ×  %-10s  =  %s wasted",
				len(g.Paths), scanner.FormatSize(g.Size), scanner.FormatSize(wasted))
			idxStr := fmt.Sprintf("  %2d.", i+1)
			if i == m.listCursor {
				b.WriteString(selectedItemStyle.Render(idxStr+row) + "\n")
			} else {
				b.WriteString(resultIndexStyle.Render(idxStr) + resultRowStyle.Render(row) + "\n")
			}
			for _, p := range g.Paths {
				b.WriteString(dimStyle.Render("       "+p) + "\n")
			}
		}
	}

	b.WriteString("\n" + dimStyle.Render("  ↑/↓ navigate • q back"))
	return b.String()
}
