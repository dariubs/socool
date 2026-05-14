package main

import "github.com/charmbracelet/lipgloss"

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
