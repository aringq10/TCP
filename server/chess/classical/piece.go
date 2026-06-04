package classical

type Piece interface {
    IsValidMove(b *Board, from, to Square) (execute func(), valid bool)
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
