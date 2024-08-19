package main

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"termchess/game"
)

func main() {
	file, err := os.OpenFile(".log/chess.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = file.Close()
	}()

	logger := slog.New(slog.NewTextHandler(file, nil))
	slog.SetLogLoggerLevel(slog.LevelInfo)
	slog.SetDefault(logger)

	// Start the TUI program
	p := tea.NewProgram(game.InitialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err = p.Run(); err != nil {
		panic(err)
	}
}
