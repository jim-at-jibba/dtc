package tui

import "github.com/charmbracelet/lipgloss"

// 4 = blue
// 6 = magenta
// 7 = white

var (
	width          = 60
	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2, 1, 2).
			Margin(2).
			Width(width).
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(lipgloss.Color("4"))

	ErrorContainerStyle = lipgloss.NewStyle().
				Padding(1, 2, 1, 2).
				Margin(2).
				Width(width).
				Border(lipgloss.NormalBorder(), true).
				BorderForeground(lipgloss.Color("1"))

	ValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("4")).
			PaddingLeft(1)

	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("7")).
			PaddingLeft(1)

	Spacer = lipgloss.NewStyle().Height(1)

	Spinner = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	ContainerNoBorderStyle = lipgloss.NewStyle().Margin(1, 2)

	ListTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("4")).
			PaddingLeft(1)
)
