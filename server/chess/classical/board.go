package classical

import (
	"fmt"
	"log"
	"strings"

	"github.com/aringq10/TCP/server/chess"
)

type Board struct {
    // First index is row: 0 = rank 1, 7 = rank 8.
    // Second index is column: 0 = file a, 7 = file h.
    squares [8][8]Piece
    lastMove chess.Move
}

func NewBoard() *Board {
    var b Board
    s := &b.squares
    for col := range 8 {
        s[1][col] = NewPawn(White)
        s[6][col] = NewPawn(Black)
    }
    s[0][0], s[0][7] = NewRook(White), NewRook(White)
    s[7][0], s[7][7] = NewRook(Black), NewRook(Black)
    return &b
}

func (b *Board) String() string {
    var sb strings.Builder
    sb.WriteByte('\n')
    for row := 7; row >= 0; row-- {
        sb.WriteString(fmt.Sprintf("%d ", row+1))
        for col := range 8 {
            if p := b.squares[row][col]; p != nil {
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

func (b *Board) MakeMove(m chess.Move) bool {
    piece := b.squares[m.From.Row][m.From.Col]
    if piece == nil || piece.Color() != m.PlayerColor {
        log.Println("piece nil or color don't match")
        return false
    }

    exec, valid := piece.IsValidMove(b, m)
    if !valid || exec == nil {
        log.Println("exec nil or IsValidMove == false")
        return false
    }

    exec()
    b.lastMove = m
    return true
}

func (b *Board) ParseSquare(s string) (sq chess.Square, ok bool) {
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
