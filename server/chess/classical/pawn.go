package classical

import (
    "github.com/aringq10/TCP/server/chess"
)

type Pawn struct {
    chess.BasePiece
}

func NewPawn(c chess.Color) *Pawn {
    return &Pawn{chess.NewBasePiece(c)}
}

func (p *Pawn) String() string {
    if p.Color() == chess.WHITE {
        return "P"
    } else {
        return "p"
    }
}

func (p *Pawn) IsValidMove(b *chess.Board, m chess.Move) (execute func(), valid bool) {
    fromRow, fromCol, toRow, toCol := m.From.Row, m.From.Col, m.To.Row, m.To.Col

    var dir, startRow int
    if p.Color() == chess.WHITE {
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
        return func() {
            b[toRow][toCol] = b[fromRow][fromCol]
            b[fromRow][fromCol] = nil
            b[fromRow][toCol] = nil
        }, true
    default:
        return nil, false
    }

    return func() {
        b[toRow][toCol] = b[fromRow][fromCol]
        b[fromRow][fromCol] = nil
    }, true
}
