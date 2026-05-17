package chess

import (
	"fmt"
	"strings"
)

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

func NewBoard() *Board {
    var b Board
    for col := range 8 {
        b[1][col] = &Pawn{basePiece{color: WHITE}}
        b[6][col] = &Pawn{basePiece{color: BLACK}}
    }
    return &b
}

func parseSquare(s string) (col, row int) {
    return int(s[0] - 'a'), int(s[1] - '1')
}

func (b *Board) LookupMove(from, to string) (piece Piece, fromCol, fromRow, toCol, toRow int, ok bool) {
    fromCol, fromRow = parseSquare(from)
    toCol, toRow = parseSquare(to)
    if fromCol < 0 || fromCol > 7 || fromRow < 0 || fromRow > 7 ||
        toCol < 0 || toCol > 7 || toRow < 0 || toRow > 7 {
        return nil, 0, 0, 0, 0, false
    }
    return b[fromRow][fromCol], fromCol, fromRow, toCol, toRow, true
}

func (b *Board) MakeMove(playerColor Color, from string, to string) bool {
    piece, fromCol, fromRow, toCol, toRow, ok := b.LookupMove(from, to)
    if !ok || piece == nil || piece.Color() != playerColor {
        return false
    }

    valid, exec := piece.IsValidMove(b, fromCol, fromRow, toCol, toRow)
    if !valid || exec == nil {
        return false
    }

    exec()
    piece.IncMovesMade()
    return true
}
