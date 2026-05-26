package classical

type Bishop struct {
    BasePiece
}

func NewBishop(c Color) *Bishop {
    return &Bishop{NewBasePiece(c)}
}

func (b *Bishop) String() string {
    if b.Color() == White {
        return "B"
    } else {
        return "b"
    }
}

func (bp *Bishop) IsValidMove(b *Board, m Move) (execute func(), valid bool) {
    s := &b.squares
    fromRow, fromCol, toRow, toCol := m.From.Row, m.From.Col, m.To.Row, m.To.Col

    dRow := toRow - fromRow
    dCol := toCol - fromCol
    if dRow == 0 || abs(dRow) != abs(dCol) {
        return nil, false
    }

    stepRow, stepCol := sign(dRow), sign(dCol)
    for row, col := fromRow+stepRow, fromCol+stepCol; row != toRow; row, col = row+stepRow, col+stepCol {
        if s[row][col] != nil {
            return nil, false
        }
    }

    if target := s[toRow][toCol]; target != nil && target.Color() == bp.Color() {
        return nil, false
    }

    return func() {
        s[toRow][toCol] = s[fromRow][fromCol]
        s[fromRow][fromCol] = nil
    }, true
}

func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}
