package common

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type (
	TickMsg  struct{}
	frameMsg struct{}
)

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return TickMsg{}
	})
}

func Frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

// Color a string's foreground with the given value.
func ColorFg(val, color string) string {
	term := termenv.EnvColorProfile()
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func MakeFgStyle(color string) func(string) string {
	term := termenv.EnvColorProfile()
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}
