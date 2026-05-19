package chess

import (
	"fmt"
	"strings"
)

type Square struct {
    Row int
    Col int
}

func ParseSquare(s string) (sq Square, ok bool) {
    if len(s) != 2 {
        return sq, false
    }

    col, row := int(s[0] - 'a'), int(s[1] - '1')
    if col < 0 || col > 7 || row < 0 || row > 7 {
        return sq, false
    }

    sq.Row = row
    sq.Col = col
    return sq, true
}

type Move struct {
    PlayerColor Color
    From Square
    To Square
}

func NewMove(playerColor Color, from Square, to Square) Move {
    return Move{PlayerColor: playerColor, From: from, To: to}
}

// First index is row: 0 = rank 1, 7 = rank 8.
// Second index is column: 0 = file a, 7 = file h.
type Board [8][8]Piece

func (b *Board) String() string {
    var sb strings.Builder
    sb.WriteByte('\n')
    for row := 7; row >= 0; row-- {
        sb.WriteString(fmt.Sprintf("%d ", row+1))
        for col := range 8 {
            if p := b[row][col]; p != nil {
                sb.WriteString(p.String())
            } else {
                sb.WriteByte('*')
            }
            sb.WriteByte(' ')
        }
        sb.WriteByte('\n')
    }
    sb.WriteString("  a b c d e f g h\n")
    return sb.String()
}

func (b *Board) MakeMove(m Move) bool {
    piece := b[m.From.Row][m.From.Col]
    if piece == nil || piece.Color() != m.PlayerColor {
        return false
    }

    exec, valid := piece.IsValidMove(b, m)
    if !valid || exec == nil {
        return false
    }

    exec()
    piece.IncMovesMade()
    return true
}
