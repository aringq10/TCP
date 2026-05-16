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
    Color() Color
}

type basePiece struct {
    Clr Color
}

func (p *basePiece) Color() Color {
    return p.Clr
}
