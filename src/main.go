package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	darkWoodBg  = "\033[48;5;130m" // Dark wood background (Brownish)
	lightWoodBg = "\033[48;5;137m" // Light wood background (Yellowish-brown)
	resetBg     = "\033[0m"        // Reset background
)

type model struct {
	board                [8][8]string // Represents the chess board
	cursorX, cursorY     int          // Cursor position
	selectedX, selectedY int          // Position of the selected piece
	selected             bool         // Whether a piece is selected
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{
		board: [8][8]string{
			{"♜", "♞", "♝", "♛", "♚", "♝", "♞", "♜"},
			{"♟", "♟", "♟", "♟", "♟", "♟", "♟", "♟"},
			{" ", " ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " ", " "},
			{" ", " ", " ", " ", " ", " ", " ", " "},
			{"♙", "♙", "♙", "♙", "♙", "♙", "♙", "♙"},
			{"♖", "♘", "♗", "♕", "♔", "♗", "♘", "♖"},
		},
		cursorX:  0,
		cursorY:  0,
		selected: false,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.cursorX > 0 {
				m.cursorX--
			}
		case "right", "l":
			if m.cursorX < 7 {
				m.cursorX++
			}
		case "up", "k":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "down", "j":
			if m.cursorY < 7 {
				m.cursorY++
			}
		case "enter", " ":
			if m.selected {
				// Move the piece
				m.board[m.cursorY][m.cursorX] = m.board[m.selectedY][m.selectedX]
				m.board[m.selectedY][m.selectedX] = " "
				m.selected = false
			} else if m.board[m.cursorY][m.cursorX] != " " {
				// Select a piece
				m.selectedX = m.cursorX
				m.selectedY = m.cursorY
				m.selected = true
			}
		case "esc":
			// Deselect the piece
			m.selected = false
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Terminal Chess\n\n"
	s += "  A  B  C  D  E  F  G  H\n" // Column labels

	for y, row := range m.board {
		s += fmt.Sprintf("%d ", 8-y) // Row labels (in reverse)
		for x, col := range row {
			bg := darkWoodBg
			if (x+y)%2 == 0 {
				bg = lightWoodBg
			}

			if x == m.cursorX && y == m.cursorY {
				if m.selected {
					s += fmt.Sprintf("%s(%s)%s", bg, col, resetBg) // Highlight selected piece
				} else {
					s += fmt.Sprintf("%s[%s]%s", bg, col, resetBg) // Highlight cursor
				}
			} else {
				s += fmt.Sprintf("%s %s %s", bg, col, resetBg)
			}
		}
		s += fmt.Sprintf(" %d\n", 8-y) // Row labels on the right side (in reverse)
	}
	s += "  A  B  C  D  E  F  G  H\n" // Column labels

	if m.selected {
		s += fmt.Sprintf("\nSelected piece: %s\n", m.board[m.selectedY][m.selectedX])
	}
	s += "\nPress 'q' or 'Ctrl+C' to quit."
	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
