package classical

type Pawn struct {
    BasePiece
}

func NewPawn(c Color) *Pawn {
    return &Pawn{NewBasePiece(c)}
}

func (p *Pawn) String() string {
    if p.Color() == White {
        return "P"
    } else {
        return "p"
    }
}

func (p *Pawn) IsValidMove(b *Board, m Move) (execute func(), valid bool) {
    s := &b.squares
    fromRow, fromCol, toRow, toCol := m.From.Row, m.From.Col, m.To.Row, m.To.Col

    var dir, startRow int
    if p.Color() == White {
        dir, startRow = 1, 1
    } else {
        dir, startRow = -1, 6
    }

    dCol := toCol - fromCol
    dRow := toRow - fromRow
    target := s[toRow][toCol]

    adjacent, _ := s[fromRow][toCol].(*Pawn)
    enPessantViable := adjacent != nil &&
        adjacent.Color() != p.Color() &&
        b.lastMove.To.Row == fromRow &&
        b.lastMove.To.Col == toCol &&
        b.lastMove.From.Row == fromRow+2*dir &&
        b.lastMove.From.Col == toCol

    switch {
    case dCol == 0 && dRow == dir && target == nil:
        // single push
    case dCol == 0 && dRow == 2*dir && fromRow == startRow &&
        target == nil && s[fromRow+dir][fromCol] == nil:
        // double push from starting rank
    case (dCol == 1 || dCol == -1) && dRow == dir &&
        target != nil && target.Color() != p.Color():
        // diagonal capture
    case (dCol == 1 || dCol == -1) && dRow == dir &&
        target == nil && enPessantViable:
        // en pessant
        return func() {
            s[toRow][toCol] = s[fromRow][fromCol]
            s[fromRow][fromCol] = nil
            s[fromRow][toCol] = nil
        }, true
    default:
        return nil, false
    }

    return func() {
        s[toRow][toCol] = s[fromRow][fromCol]
        s[fromRow][fromCol] = nil
    }, true
}
