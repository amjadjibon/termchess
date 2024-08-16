package main

import (
	tea "github.com/charmbracelet/bubbletea"

	"termchess/game"
)

func main() {
	// Start the TUI program
	p := tea.NewProgram(game.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
