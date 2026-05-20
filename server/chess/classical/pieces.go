package classical

import (
	"github.com/aringq10/TCP/server/chess"
)

type Color int

const (
    White Color = iota
    Black
)

var colorNames = map[chess.Color]string {
    White: "WHTE",
    Black: "BLCK",
}

func (c Color) String() string {
    return colorNames[c]
}

type Piece interface {
    IsValidMove(b *Board, m chess.Move) (execute func(), valid bool)
    String() string
    Color() Color
}

type BasePiece struct {
    color Color
}

func NewBasePiece(c Color) BasePiece {
    return BasePiece{color: c}
}

func (p *BasePiece) Color() Color {
    return p.color
}
