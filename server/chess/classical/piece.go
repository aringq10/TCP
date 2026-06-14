package classical

type Piece interface {
    IsValidMove(b *Board, from, to Square) (execute func(), valid bool)
    String() string
    Color() Color
    MoveCount() int
    IncMoveCount()
}

type BasePiece struct {
    color Color
    moveCount int
}

func NewBasePiece(c Color) BasePiece {
    return BasePiece{color: c}
}

func (p *BasePiece) Color() Color {
    return p.color
}

func (p *BasePiece) MoveCount() int{
    return p.moveCount
}

func (p *BasePiece) IncMoveCount() {
    p.moveCount++
}
