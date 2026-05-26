package classical

type Knight struct {
    BasePiece
}

func NewKnight(c Color) *Knight {
    return &Knight{NewBasePiece(c)}
}

func (k *Knight) String() string {
    if k.Color() == White {
        return "N"
    } else {
        return "n"
    }
}

func (k *Knight) IsValidMove(b *Board, m Move) (execute func(), valid bool) {
    s := &b.squares
    fromRow, fromCol, toRow, toCol := m.From.Row, m.From.Col, m.To.Row, m.To.Col

    dRow := abs(toRow - fromRow)
    dCol := abs(toCol - fromCol)
    if !((dRow == 2 && dCol == 1) || (dRow == 1 && dCol == 2)) {
        return nil, false
    }

    if target := s[toRow][toCol]; target != nil && target.Color() == k.Color() {
        return nil, false
    }

    return func() {
        s[toRow][toCol] = s[fromRow][fromCol]
        s[fromRow][fromCol] = nil
    }, true
}
