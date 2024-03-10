package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"

	"github.com/foureyez/requestr/internal/models"
)

func main() {
	requestModel := models.NewRequestModal()
	p := tea.NewProgram(&requestModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal().Msgf("Unable to initialize app: %s", err.Error())
	}
}
