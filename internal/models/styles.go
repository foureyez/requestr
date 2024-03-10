package models

import "github.com/charmbracelet/lipgloss"

var appStyle = lipgloss.NewStyle().
	Margin(1, 1).Border(lipgloss.DoubleBorder(), true).Width(150)

var (
	inputStyle      = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Width(100)
	headerStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true)
	httpMethodStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Bold(true)
)
