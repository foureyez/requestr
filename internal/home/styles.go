package home

import "github.com/charmbracelet/lipgloss"

var appStyle = lipgloss.NewStyle().
	Margin(1, 1).Border(lipgloss.DoubleBorder(), true).Width(100)

var inputStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true)
