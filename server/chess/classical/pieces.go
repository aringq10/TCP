package classical

type Color int

const (
    White Color = iota
    Black
)

var colorNames = map[Color]string {
    White: "WHTE",
    Black: "BLCK",
}

func (c Color) String() string {
    return colorNames[c]
}

type Piece interface {
    IsValidMove(b *Board, m Move) (execute func(), valid bool)
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
