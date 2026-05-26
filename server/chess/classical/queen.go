package classical

type Queen struct {
    BasePiece
}

func NewQueen(c Color) *Queen {
    return &Queen{NewBasePiece(c)}
}

func (q *Queen) String() string {
    if q.Color() == White {
        return "Q"
    } else {
        return "q"
    }
}

func (q *Queen) IsValidMove(b *Board, m Move) (execute func(), valid bool) {
    s := &b.squares
    fromRow, fromCol, toRow, toCol := m.From.Row, m.From.Col, m.To.Row, m.To.Col

    dRow := toRow - fromRow
    dCol := toCol - fromCol
    straight := (dRow == 0) != (dCol == 0)
    diagonal := dRow != 0 && abs(dRow) == abs(dCol)
    if !straight && !diagonal {
        return nil, false
    }

    stepRow, stepCol := sign(dRow), sign(dCol)
    for row, col := fromRow+stepRow, fromCol+stepCol; row != toRow || col != toCol; row, col = row+stepRow, col+stepCol {
        if s[row][col] != nil {
            return nil, false
        }
    }

    if target := s[toRow][toCol]; target != nil && target.Color() == q.Color() {
        return nil, false
    }

    return func() {
        s[toRow][toCol] = s[fromRow][fromCol]
        s[fromRow][fromCol] = nil
    }, true
}
