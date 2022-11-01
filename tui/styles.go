package tui

import "github.com/charmbracelet/lipgloss"

var DocStyle = lipgloss.NewStyle().
	Padding(1, 2, 1, 2).
	Margin(2).
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(lipgloss.Color("#7D56F4"))
