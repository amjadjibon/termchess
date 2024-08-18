package game

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/notnil/chess"
)

const (
	boardSize = 8
)

type Model struct {
	board                *Board // Represents the chess board
	cursorX, cursorY     int    // Cursor position on the board
	selectedX, selectedY int    // Position of the selected piece
	selectedPiece        Piece  // Piece that is selected
	selected             bool   // Whether a piece is selected
	currentPlayer        Player
	gameEngine           *chess.Game
	mousePosition        string
	enPassantTarget      string
}

func InitialModel() *Model {
	return &Model{
		board:         NewBoard(),
		cursorX:       4, // Column 'e'
		cursorY:       6, // Row '2' (reversed ranks: 6 for row 2)
		selected:      false,
		currentPlayer: PlayerWhite,
		gameEngine:    chess.NewGame(chess.UseNotation(chess.UCINotation{})),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	re := lipgloss.NewRenderer(os.Stdout)

	// create the table with alternating black and white squares
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

	// Labels for ranks (1-8) and files (a-h)
	labelStyle := re.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Center)
	ranks := labelStyle.Render(strings.Join([]string{"\n      a", "b", "c", "d", "e", "f", "g", "h"}, "      "))
	files := strings.Join([]string{
		labelStyle.Render("\n 8"),
		labelStyle.Render("\n\n 7"),
		labelStyle.Render("\n\n 6"),
		labelStyle.Render("\n\n 5"),
		labelStyle.Render("\n\n 4"),
		labelStyle.Render("\n\n 3"),
		labelStyle.Render("\n\n 2"),
		labelStyle.Render("\n\n 1"),
	}, "\n")

	header := labelStyle.Render("                      Terminal Chess\n")

	// Render the PGN on the right side of the board
	pgnSplit := splitPGN(m.gameEngine.String())
	pgnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	pgnMoves := "\n" + pgnStyle.Render(strings.Join(pgnSplit, "\n"))

	footer := ranks
	footerSelectedPiece := lipgloss.NewStyle().
		Background(lipgloss.Color("#ffffff")).
		Foreground(lipgloss.Color("#ffffff"))
	if m.selected {
		footer += fmt.Sprintf("\nSelected piece: %s\n",
			footerSelectedPiece.Render(m.selectedPiece.Render()))
	}

	footer += "\n\n\nCurrent player: " + m.currentPlayer.String()
	footer += "\nPress 'q' or 'Ctrl+C' to quit.\n"

	footer += "\n" + m.mousePosition

	return header + lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left, "", files),
			t.Render(),
			pgnMoves,
		),
	) + footer
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "left", "h":
			m.moveCursorLeft()
		case "right", "l":
			m.moveCursorRight()
		case "up", "k":
			m.moveCursorUp()
		case "down", "j":
			m.moveCursorDown()
		case "enter", " ":
			m.handleSelectOrMove()
		case "esc":
			m.deselectPiece()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.MouseMsg:
		switch msgType.Action {
		case tea.MouseActionPress:
			m.handleMouseClick(msgType.X, msgType.Y)
		default:

		}

	}

	return m, nil
}

func (m *Model) moveCursorLeft() {
	if m.cursorX > 0 {
		m.cursorX--
	}
}

func (m *Model) moveCursorRight() {
	if m.cursorX < boardSize-1 {
		m.cursorX++
	}
}

func (m *Model) moveCursorUp() {
	if m.cursorY > 0 {
		m.cursorY--
	}
}

func (m *Model) moveCursorDown() {
	if m.cursorY < boardSize-1 {
		m.cursorY++
	}
}

func (m *Model) deselectPiece() {
	m.selected = false
}

func (m *Model) handleSelectOrMove() {
	if m.selected {
		m.applyMove()
	} else {
		m.selectPiece()
	}
}

func (m *Model) canApplyMove() bool {
	// prevent moving a piece to the same place
	if m.cursorX == m.selectedX && m.cursorY == m.selectedY {
		return false
	}

	// prevent moving to a square occupied by a piece of the same color
	if m.selectedPiece.IsWhite() && m.board.Get(m.cursorY, m.cursorX).IsWhite() ||
		m.selectedPiece.IsBlack() && m.board.Get(m.cursorY, m.cursorX).IsBlack() {
		return false
	}

	return true
}

func (m *Model) applyMove() {
	from := coordsToUCI(m.selectedX, m.selectedY)
	to := coordsToUCI(m.cursorX, m.cursorY)
	move := from + to

	// Handle en passant
	if m.selectedPiece.IsPawn() && m.enPassantTarget == to {
		if m.selectedPiece.IsWhite() && m.cursorY == 2 {
			m.board.grid[3][m.cursorX] = Empty // Remove captured black pawn
			m.enPassantTarget = ""
		} else if m.selectedPiece.IsBlack() && m.cursorY == 5 {
			m.board.grid[4][m.cursorX] = Empty // Remove captured white pawn
			m.enPassantTarget = ""
		}
	}

	// Handle pawn promotion
	if canPiecePromote(m.selectedPiece, m.cursorY) {
		move += m.handlePromotion()
	}

	if err := m.gameEngine.MoveStr(move); err != nil {
		return
	}

	defer func() {
		m.currentPlayer = m.currentPlayer.Switch()
		m.selected = false
	}()

	if m.selectedPiece.IsPawn() && abs(m.selectedY-m.cursorY) == 2 {
		m.enPassantTarget = coordsToUCI(m.cursorX, (m.selectedY+m.cursorY)/2)
	}

	switch move {
	case "e1g1": // White king-side castling
		m.board.Replace(7, 4, 7, 6) // Move the king
		m.board.Replace(7, 7, 7, 5) // Move the rook
		return
	case "e1c1": // White queen-side castling
		m.board.Replace(7, 4, 7, 2) // Move the king
		m.board.Replace(7, 0, 7, 3) // Move the rook
		return
	case "e8g8": // Black king-side castling
		m.board.Replace(0, 4, 0, 6) // Move the king
		m.board.Replace(0, 7, 0, 5) // Move the rook
		return
	case "e8c8": // Black queen-side castling
		m.board.Replace(0, 4, 0, 2) // Move the king
		m.board.Replace(0, 0, 0, 3) // Move the rook
		return
	}

	// Regular move
	m.board.Set(m.cursorY, m.cursorX, m.selectedPiece)
	m.board.grid[m.selectedY][m.selectedX] = Empty
	m.selected = false
}

func (m *Model) canSelect() bool {
	// no player can select an empty space
	if m.board.Get(m.cursorY, m.cursorX) == Empty {
		return false
	}

	// ensure the current player can only select their own pieces
	if (m.currentPlayer == PlayerWhite && m.board.Get(m.cursorY, m.cursorX).IsBlack()) ||
		(m.currentPlayer == PlayerBlack && m.board.Get(m.cursorY, m.cursorX).IsWhite()) {
		return false
	}

	return true
}

func (m *Model) selectPiece() {
	if !m.canSelect() {
		return
	}

	m.selectedX = m.cursorX
	m.selectedY = m.cursorY
	m.selectedPiece = m.board.Get(m.selectedY, m.selectedX)
	m.selected = true
}

func coordsToUCI(x, y int) string {
	files := "abcdefgh"
	ranks := "87654321" // reversed for standard board setup
	return string(files[x]) + string(ranks[y])
}

// Check if the selected piece is a pawn that reaches the promotion rank
func canPiecePromote(piece Piece, targetY int) bool {
	if !piece.IsPawn() {
		return false
	}

	if (piece.IsWhite() && targetY == 0) ||
		(piece.IsBlack() && targetY == 7) {
		return true
	}

	return false
}

func promotionForm() *huh.Select[string] {
	return huh.NewSelect[string]().
		Title("choose a piece").
		Options(
			huh.NewOption("Queen", "q"),
			huh.NewOption("Rook", "r"),
			huh.NewOption("Bishop", "b"),
			huh.NewOption("Knight", "n"),
		)
}

func (m *Model) handlePromotion() string {
	form := promotionForm()
	if err := form.Run(); err != nil {
		return ""
	}

	var piece string
	form.Value(&piece)

	m.updateBoardForPromotion(piece)
	return piece
}

func (m *Model) updateBoardForPromotion(piece string) {
	promotedPiece := m.getPromotionPiece(piece)
	m.selectedPiece = promotedPiece
}

func (m *Model) getPromotionPiece(piece string) Piece {
	switch piece {
	case "q":
		if m.selectedPiece.IsWhite() {
			return WhiteQueen
		}
		return BlackQueen
	case "r":
		if m.selectedPiece.IsWhite() {
			return WhiteRook
		}
		return BlackRook
	case "b":
		if m.selectedPiece.IsWhite() {
			return WhiteBishop
		}
		return BlackBishop
	case "n":
		if m.selectedPiece.IsWhite() {
			return WhiteKnight
		}
		return BlackKnight
	default:
		// If for some reason an invalid piece is provided, return Empty
		return Empty
	}
}

func splitPGN(pgn string) []string {
	// Regular expression to match move numbers (e.g., 1., 2., 3.)
	re := regexp.MustCompile(`\d+\.`)

	// Find all matches in the PGN string
	matches := re.FindAllStringSubmatchIndex(pgn, -1)
	if len(matches) == 0 {
		return nil
	}

	var segments []string

	// Initialize the start index for the first segment
	start := matches[0][0]

	// Extract and append each segment including the move numbers
	for i := 1; i < len(matches); i++ {
		end := matches[i][0]
		segment := strings.TrimSpace(pgn[start:end])
		if segment != "" {
			segments = append(segments, segment)
		}
		start = end
	}

	// Append the last segment from the last move number to the end of the string
	segment := strings.TrimSpace(pgn[start:])
	if segment != "" {
		segments = append(segments, segment)
	}

	return segments
}

func (m *Model) handleMouseClick(x, y int) {
	boardOffsetX := 2
	boardOffsetY := 2
	cellWidth := 7
	cellHeight := 3

	col := (x - boardOffsetX) / cellWidth
	row := (y - boardOffsetY) / cellHeight

	if col >= 0 && col < boardSize && row >= 0 && row < boardSize {
		m.cursorX = col
		m.cursorY = row
		m.handleSelectOrMove()
	}
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}
