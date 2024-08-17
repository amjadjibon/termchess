package game

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	whitePieceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000"))

	blackPieceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000"))
)

type Piece int

const (
	Empty Piece = iota
	WhitePawn
	WhiteRook
	WhiteKnight
	WhiteBishop
	WhiteQueen
	WhiteKing
	BlackPawn
	BlackRook
	BlackKnight
	BlackBishop
	BlackQueen
	BlackKing
)

var PieceMap = map[Piece]string{
	WhitePawn:   "♙",
	WhiteRook:   "♖",
	WhiteKnight: "♘",
	WhiteBishop: "♗",
	WhiteQueen:  "♕",
	WhiteKing:   "♔",
	BlackPawn:   "♟",
	BlackRook:   "♜",
	BlackKnight: "♞",
	BlackBishop: "♝",
	BlackQueen:  "♛",
	BlackKing:   "♚",
	Empty:       " ", // Represents an empty square
}

func (p Piece) String() string {
	if emoji, exists := PieceMap[p]; exists {
		return emoji
	}
	return " "
}

func (p Piece) Name() string {
	switch p {
	case WhiteRook, BlackRook:
		return "R"
	case WhiteKnight, BlackKnight:
		return "N"
	case WhiteBishop, BlackBishop:
		return "B"
	case WhiteQueen, BlackQueen:
		return "Q"
	case WhiteKing, BlackKing:
		return "K"
	default:
		return ""
	}
}

func (p Piece) IsWhite() bool {
	return p >= 1 && p <= 6
}

func (p Piece) IsBlack() bool {
	return p >= 7 && p <= 12
}

func (p Piece) IsPawn() bool {
	return p == 1 || p == 7
}

func (p Piece) Render() string {
	switch p {
	case WhitePawn, WhiteRook, WhiteKnight, WhiteBishop, WhiteQueen, WhiteKing:
		return whitePieceStyle.Render(p.String())
	case BlackPawn, BlackRook, BlackKnight, BlackBishop, BlackQueen, BlackKing:
		return blackPieceStyle.Render(p.String())
	case Empty:
		return " "
	default:
		return " "
	}
}
