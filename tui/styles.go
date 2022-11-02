package tui

import "github.com/charmbracelet/lipgloss"

var ContainerStyle = lipgloss.NewStyle().
	Padding(1, 2, 1, 2).
	Margin(2).
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(lipgloss.Color("#7D56F4"))

var ValueStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4")).
	PaddingLeft(1)

var LabelStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFF")).
	PaddingTop(1).
	PaddingLeft(1)
