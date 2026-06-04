package classical

type King struct {
    BasePiece
}

func NewKing(c Color) *King {
    return &King{NewBasePiece(c)}
}

func (p *King) String() string {
    if p.Color() == White {
        return "K"
    } else {
        return "k"
    }
}

func (p *King) IsValidMove(b *Board, from, to Square) (execute func(), valid bool) {
    s := &b.squares
    fromRow, fromCol, toRow, toCol := from.Row, from.Col, to.Row, to.Col

    dRow := toRow - fromRow
    dCol := toCol - fromCol
    if abs(dRow) > 1 || abs(dCol) > 1 || abs(dRow) + abs(dCol) == 0 {
        return nil, false
    }

    target := s[toRow][toCol]

    if target != nil && target.Color() == p.Color() {
        return nil, false
    }

    return func() {
        s[toRow][toCol] = s[fromRow][fromCol]
        s[fromRow][fromCol] = nil
    }, true
}
