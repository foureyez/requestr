package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"

	"github.com/foureyez/requestr/internal/home"
)

func main() {
	home := home.NewHome()
	p := tea.NewProgram(home)
	if _, err := p.Run(); err != nil {
		log.Fatal().Msgf("Unable to initialize app: %s", err.Error())
	}
}
