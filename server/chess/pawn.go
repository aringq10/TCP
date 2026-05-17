package chess

type Pawn struct {
    basePiece
}

func (p *Pawn) String() string {
    if p.Color() == WHITE {
        return "P"
    } else {
        return "p"
    }
}

func (p *Pawn) IsValidMove(b *Board, fromCol, fromRow, toCol, toRow int) (valid bool, execute func()) {
    var dir, startRow int
    if p.Color() == WHITE {
        dir, startRow = 1, 1
    } else {
        dir, startRow = -1, 6
    }

    dCol := toCol - fromCol
    dRow := toRow - fromRow
    target := b[toRow][toCol]

    switch {
    case dCol == 0 && dRow == dir && target == nil:
        // single push
    case dCol == 0 && dRow == 2*dir && fromRow == startRow &&
        target == nil && b[fromRow+dir][fromCol] == nil:
        // double push from starting rank
    case (dCol == 1 || dCol == -1) && dRow == dir &&
        target != nil && target.Color() != p.Color():
        // diagonal capture
    case (dCol == 1 || dCol == -1) && dRow == dir &&
        target == nil && b[fromRow][toCol].MovesMade() == 1:
        // en pessant
        return true, func() {
            b[toRow][toCol] = b[fromRow][fromCol]
            b[fromRow][fromCol] = nil
            b[fromRow][toCol] = nil
        }
    default:
        return false, nil
    }

    return true, func() {
        b[toRow][toCol] = b[fromRow][fromCol]
        b[fromRow][fromCol] = nil
    }
}
