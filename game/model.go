package game

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	boardSize = 8
)

var (
	// Cursor on white square style
	whiteCursorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#e3d5ca")).
				Foreground(lipgloss.Color("#000000")).
				Align(lipgloss.Center).
				Padding(1, 3)

	// Cursor on black square style
	blackCursorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#81b583")).
				Foreground(lipgloss.Color("#ffffff")).
				Align(lipgloss.Center).
				Padding(1, 3)

	// Selected square style
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("21")).
			Foreground(lipgloss.Color("21")).
			Align(lipgloss.Center).
			Padding(1, 3)

	// Black square style
	blackSquare = lipgloss.NewStyle().
			Background(lipgloss.Color("#4e7837")).
			Foreground(lipgloss.Color("#ffffff")).
			Align(lipgloss.Center).
			Padding(1, 3)

	// White square style
	whiteSquare = lipgloss.NewStyle().
			Background(lipgloss.Color("#ffffff")).
			Foreground(lipgloss.Color("#4e7837")).
			Align(lipgloss.Center).
			Padding(1, 3)
)

type Model struct {
	board                Board // Represents the chess board
	cursorX, cursorY     int   // Cursor position on the board
	selectedX, selectedY int   // Position of the selected piece
	selectedPiece        Piece // Piece that is selected
	selected             bool  // Whether a piece is selected
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel() Model {
	return Model{
		board:    NewBoard(),
		cursorX:  0,
		cursorY:  0,
		selected: false,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
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
				m.board[m.cursorY][m.cursorX] = m.selectedPiece
				m.board[m.selectedY][m.selectedX] = Empty
				m.selected = false
			} else if m.board[m.cursorY][m.cursorX] != Empty {
				// Select a piece
				m.selectedX = m.cursorX
				m.selectedY = m.cursorY
				m.selectedPiece = m.board[m.selectedY][m.selectedX]
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

func (m Model) View() string {
	re := lipgloss.NewRenderer(os.Stdout)

	// Create the table with alternating black and white squares
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderRow(false).
		BorderColumn(false).
		Rows(m.board.Display()...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if m.cursorX == col && m.cursorY == row-1 && m.selected {
				return selectedStyle
			} else if (row+col)%2 == 0 {
				if m.cursorX == col && m.cursorY == row-1 {
					return blackCursorStyle
				}
				return blackSquare
			} else {
				if m.cursorX == col && m.cursorY == row-1 {
					return whiteCursorStyle
				}
				return whiteSquare
			}
		})

	// Labels for ranks (1-8) and files (A-H)
	labelStyle := re.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Center)
	ranks := labelStyle.Render(strings.Join([]string{"\n    A", "B", "C", "D", "E", "F", "G", "H"}, "      "))
	files := strings.Join([]string{
		labelStyle.Render("\n8"),
		labelStyle.Render("\n\n7"),
		labelStyle.Render("\n\n6"),
		labelStyle.Render("\n\n5"),
		labelStyle.Render("\n\n4"),
		labelStyle.Render("\n\n3"),
		labelStyle.Render("\n\n2"),
		labelStyle.Render("\n\n1"),
	}, "\n")

	header := labelStyle.Render("                      Terminal Chess\n")

	footer := ranks
	footerSelectedPiece := lipgloss.NewStyle().
		Background(lipgloss.Color("23")).
		Foreground(lipgloss.Color("23")).
		Align(lipgloss.Center)
	if m.selected {
		footer += fmt.Sprintf("\nSelected piece: %s\n", footerSelectedPiece.Render(m.selectedPiece.String()))
	}
	footer += "\n\n\nPress 'q' or 'Ctrl+C' to quit.\n"

	return header + lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left, "", files),
			t.Render(),
		),
	) + footer
}
