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
    IsValidMove(b *Board, fromCol, fromRow, toCol, toRow int) (valid bool, execute func())
    String() string
    Color() Color
    SetColor(c Color)
    MovesMade() int
    IncMovesMade()
}

type basePiece struct {
    color Color
    movesMade int
}

func (p *basePiece) Color() Color {
    return p.color
}

func (p *basePiece) SetColor(c Color) {
    p.color = c
}

func (p *basePiece) MovesMade() int {
    return p.movesMade
}

func (p *basePiece) IncMovesMade() {
    p.movesMade++
}
