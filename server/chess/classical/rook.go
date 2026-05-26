package classical

type Rook struct {
    BasePiece
}

func NewRook(c Color) *Rook {
    return &Rook{NewBasePiece(c)}
}

func (r *Rook) String() string {
    if r.Color() == White {
        return "R"
    } else {
        return "r"
    }
}

func (r *Rook) IsValidMove(b *Board, m Move) (execute func(), valid bool) {
    s := &b.squares
    fromRow, fromCol, toRow, toCol := m.From.Row, m.From.Col, m.To.Row, m.To.Col

    dRow := toRow - fromRow
    dCol := toCol - fromCol
    if (dRow == 0) == (dCol == 0) {
        return nil, false
    }

    stepRow, stepCol := sign(dRow), sign(dCol)
    for row, col := fromRow+stepRow, fromCol+stepCol; row != toRow || col != toCol; row, col = row+stepRow, col+stepCol {
        if s[row][col] != nil {
            return nil, false
        }
    }

    if target := s[toRow][toCol]; target != nil && target.Color() == r.Color() {
        return nil, false
    }

    return func() {
        s[toRow][toCol] = s[fromRow][fromCol]
        s[fromRow][fromCol] = nil
    }, true
}

func sign(n int) int {
    switch {
    case n > 0:
        return 1
    case n < 0:
        return -1
    }
    return 0
}
