package game

import (
	"fmt"
)

type Board struct {
	grid [8][8]Piece
}

func (b *Board) Get(x, y int) Piece {
	return b.grid[x][y]
}

func (b *Board) Set(x, y int, p Piece) {
	b.grid[x][y] = p
}

func (b *Board) Replace(fromX, fromY int, toX, toY int) {
	b.grid[toX][toY] = b.grid[fromX][fromY]
	b.grid[fromX][fromY] = Empty
}

// Display method to convert the board into a 2D slice of strings
func (b *Board) Display() [][]string {
	display := make([][]string, len(b.grid))

	for i := range b.grid {
		display[i] = make([]string, len(b.grid[i]))
		for j := range b.grid[i] {
			display[i][j] = b.grid[i][j].Render()
		}
	}

	return display
}

func NewBoard() *Board {
	grid := [8][8]Piece{
		{BlackRook, BlackKnight, BlackBishop, BlackQueen, BlackKing, BlackBishop, BlackKnight, BlackRook},
		{BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty},
		{WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn},
		{WhiteRook, WhiteKnight, WhiteBishop, WhiteQueen, WhiteKing, WhiteBishop, WhiteKnight, WhiteRook},
	}

	return &Board{grid: grid}
}

// Position converts board coordinates to chess notation (e.g., (6, 4) -> "e2")
func Position(x, y int) string {
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		return ""
	}

	columns := "abcdefgh" // Column letters for chess notation
	row := 8 - x          // Row numbers in chess notation
	column := columns[y]  // Get the corresponding column letter

	return fmt.Sprintf("%c%d", column, row) // Format as chess position
}

// Coordinates converts chess notation (e.g., "e2") to board coordinates (e.g., "e2" -> (6, 4)).
func coordinates(pos string) (x, y int) {
	if len(pos) != 2 {
		panic("invalid position format")
	}

	// Convert column letter (a-h) to index (0-7)
	y = int(pos[0] - 'a')

	// Convert row number (1-8) to index (0-7)
	x = 8 - int(pos[1]-'0')

	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		panic("position out of bounds")
	}

	return x, y
}
