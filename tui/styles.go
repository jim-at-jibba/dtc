package tui

import "github.com/charmbracelet/lipgloss"

// 4 = blue
// 6 = magenta
// 7 = white

var width = 60

var ContainerStyle = lipgloss.NewStyle().
	Padding(1, 2, 1, 2).
	Margin(2).
	Width(width).
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(lipgloss.Color("4"))

var ErrorContainerStyle = lipgloss.NewStyle().
	Padding(1, 2, 1, 2).
	Margin(2).
	Width(width).
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(lipgloss.Color("1"))

var ValueStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("4")).
	PaddingLeft(1)

var LabelStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("7")).
	PaddingLeft(1)

var Spacer = lipgloss.NewStyle().Height(1)

var Spinner = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

var ContainerNoBorderStyle = lipgloss.NewStyle().Margin(1, 2)

var ListTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("4")).
	PaddingLeft(1)
