package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	re := lipgloss.NewRenderer(os.Stdout)
	labelStyle := re.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Center)

	// Define piece styles
	whitePieceLightSquareStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")) // White pieces on light squares
	whitePieceDarkSquareStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))  // White pieces on dark squares 	// White pieces
	blackPieceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))            // Black pieces

	// Chessboard setup with alternative larger symbols
	board := [][]string{
		{blackPieceStyle.Render("♜"), blackPieceStyle.Render("♞"), blackPieceStyle.Render("♝"), blackPieceStyle.Render("♛"), blackPieceStyle.Render("♚"), blackPieceStyle.Render("♝"), blackPieceStyle.Render("♞"), blackPieceStyle.Render("♜")}, // Black pieces
		{blackPieceStyle.Render("♟"), blackPieceStyle.Render("♟"), blackPieceStyle.Render("♟"), blackPieceStyle.Render("♟"), blackPieceStyle.Render("♟"), blackPieceStyle.Render("♟"), blackPieceStyle.Render("♟"), blackPieceStyle.Render("♟")}, // Black pawns
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{whitePieceLightSquareStyle.Render("♙"), whitePieceDarkSquareStyle.Render("♙"), whitePieceLightSquareStyle.Render("♙"), whitePieceDarkSquareStyle.Render("♙"), whitePieceLightSquareStyle.Render("♙"), whitePieceDarkSquareStyle.Render("♙"), whitePieceLightSquareStyle.Render("♙"), whitePieceDarkSquareStyle.Render("♙")}, // White pawns
		{whitePieceDarkSquareStyle.Render("♖"), whitePieceLightSquareStyle.Render("♘"), whitePieceDarkSquareStyle.Render("♗"), whitePieceLightSquareStyle.Render("♕"), whitePieceDarkSquareStyle.Render("♔"), whitePieceLightSquareStyle.Render("♗"), whitePieceDarkSquareStyle.Render("♘"), whitePieceLightSquareStyle.Render("♖")}, // White pieces

	}

	// Generate styles for black and white squares with extra padding
	blackSquare := lipgloss.NewStyle().
		Background(lipgloss.Color("#4e7837")).
		Foreground(lipgloss.Color("#ffffff")).
		Align(lipgloss.Center).
		Padding(1, 3) // Adjusted padding for visual space

	whiteSquare := lipgloss.NewStyle().
		Background(lipgloss.Color("#ffffff")).
		Foreground(lipgloss.Color("#4e7837")).
		Align(lipgloss.Center).
		Padding(1, 3) // Adjusted padding for visual space

	// Create the table with alternating black and white squares
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderRow(false).
		BorderColumn(false).
		Rows(board...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if (row+col)%2 == 0 {
				return blackSquare
			}
			return whiteSquare
		})

	// Labels for ranks (1-8) and files (A-H)
	ranks := labelStyle.Render(strings.Join([]string{"  A", "B", "C", "D", "E", "F", "G", "H"}, "      "))
	files := strings.Join([]string{
		labelStyle.Render("\n1"),
		labelStyle.Render("\n\n2"),
		labelStyle.Render("\n\n3"),
		labelStyle.Render("\n\n4"),
		labelStyle.Render("\n\n5"),
		labelStyle.Render("\n\n6"),
		labelStyle.Render("\n\n7"),
		labelStyle.Render("\n\n8"),
	}, "\n")

	// Display the board with properly aligned labels
	fmt.Println(lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left, "", files),
			t.Render(),
		),
	) + "\n  " + ranks)
}
