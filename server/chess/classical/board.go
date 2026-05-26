package classical

import (
	"fmt"
	"strings"

	"github.com/aringq10/TCP/server/chess"
)

type Board struct {
    // First index is row: 0 = rank 1, 7 = rank 8.
    // Second index is column: 0 = file a, 7 = file h.
    squares [8][8]Piece
    lastMove Move
    whoseTurn Color
}

type Square struct {
    Row int
    Col int
}

type Move struct {
    PlayerColor Color
    From Square
    To Square
}

func (Move) IsMove() {}

func NewMove(playerColor Color, from Square, to Square) Move {
    return Move{PlayerColor: playerColor, From: from, To: to}
}

func NewBoard() *Board {
    var b Board
    b.whoseTurn = White
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
    cm, ok := m.(Move)
    if !ok {
        return false
    }

    if cm.PlayerColor != b.whoseTurn {
        return false
    }

    piece := b.squares[cm.From.Row][cm.From.Col]
    if piece == nil || piece.Color() != cm.PlayerColor {
        return false
    }

    exec, valid := piece.IsValidMove(b, cm)
    if !valid || exec == nil {
        return false
    }

    exec()
    b.rotateTurn()
    b.lastMove = cm
    return true
}

func (b *Board) rotateTurn() {
    if b.whoseTurn == White {
        b.whoseTurn = Black
    } else {
        b.whoseTurn = White
    }
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
