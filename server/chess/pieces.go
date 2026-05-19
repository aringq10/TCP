package chess

const (
    WHITE Color = iota
    BLACK
)

var colorNames = map[Color]string {
    WHITE: "WHTE",
    BLACK: "BLCK",
}

type Color int

func (c Color) String() string {
    return colorNames[c]
}

type Piece interface {
    IsValidMove(b *Board, m Move) (execute func(), valid bool)
    String() string
    Color() Color
    SetColor(c Color)
    MovesMade() int
    IncMovesMade()
}

type BasePiece struct {
    color Color
    movesMade int
}

func NewBasePiece(c Color) BasePiece {
    return BasePiece{color: c}
}

func (p *BasePiece) Color() Color {
    return p.color
}

func (p *BasePiece) SetColor(c Color) {
    p.color = c
}

func (p *BasePiece) MovesMade() int {
    return p.movesMade
}

func (p *BasePiece) IncMovesMade() {
    p.movesMade++
}
