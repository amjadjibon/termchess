package main

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess/uci"

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

	// set up engine to use stockfish exe
	eng, err := uci.New("stockfish/stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()

	// initialize uci with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}

	// Start the TUI program
	p := tea.NewProgram(game.InitialModel(eng), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err = p.Run(); err != nil {
		panic(err)
	}
}
