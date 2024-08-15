package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	boardSize = 8
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
			if m.cursorX < boardSize-1 {
				m.cursorX++
			}
		case "up", "k":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "down", "j":
			if m.cursorY < boardSize-1 {
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
	var b strings.Builder

	// Create styles
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("23"))

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("34"))

	lightSquareStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#ffffff"))

	darkSquareStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#4e7837"))

	// labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	// Build the board with labels
	for y := 0; y < boardSize; y++ {
		if y == 0 {
			// Add the file labels (A-H) above the board
			b.WriteString("  A B C D E F G H\n")
		}
		for x := 0; x < boardSize; x++ {
			var style lipgloss.Style
			if (x+y)%2 == 0 {
				style = lightSquareStyle
			} else {
				style = darkSquareStyle
			}

			if x == m.cursorX && y == m.cursorY {
				if m.selected {
					style = selectedStyle
				} else {
					style = cursorStyle
				}
			}

			// Apply the style to the piece
			b.WriteString(style.Render(" " + m.board[y][x] + " "))
		}
		// Add the rank label (1-8) at the end of each row
		b.WriteString(fmt.Sprintf(" %d\n", y+1))
	}

	if m.selected {
		b.WriteString(fmt.Sprintf("\nSelected piece: %s\n", m.board[m.selectedY][m.selectedX]))
	}
	b.WriteString("\nPress 'q' or 'Ctrl+C' to quit.")

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
