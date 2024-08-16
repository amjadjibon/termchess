package game

type Player int

const (
	PlayerWhite Player = iota
	PlayerBlack
)

func (p Player) String() string {
	switch p {
	case PlayerWhite:
		return "white"
	case PlayerBlack:
		return "black"
	default:
		return "unknown"
	}
}

func (p Player) Switch() Player {
	switch p {
	case PlayerWhite:
		return PlayerBlack
	default:
		return PlayerWhite
	}
}
